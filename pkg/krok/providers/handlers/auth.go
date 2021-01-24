package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/cache"
	"github.com/krok-o/krok/pkg/models"
)

// Config has the configuration options for the repository handler.
type Config struct {
	Hostname       string
	GlobalTokenKey string
}

// Dependencies defines the dependencies for the repository handler provider.
type Dependencies struct {
	Logger     zerolog.Logger
	UserStore  providers.UserStorer
	ApiKeyAuth providers.ApiKeysAuthenticator
	UserCache  *cache.UserCache
}

// TokenProvider is a token provider for the handlers.
type TokenProvider struct {
	Config
	Dependencies

	// A cache to track authenticated users.
	cache *cache.UserCache
}

// NewTokenProvider creates a new token provider which deals with generating and handling tokens.
func NewTokenProvider(cfg Config, deps Dependencies) (*TokenProvider, error) {
	tp := &TokenProvider{
		Config:       cfg,
		Dependencies: deps,
	}
	return tp, nil
}

// ApiKeyAuthRequest contains a user email and their api key.
type ApiKeyAuthRequest struct {
	Email        string `json:"email"`
	APIKeyID     string `json:"api_key_id"`
	APIKeySecret string `json:"api_key_secret"`
}

// TokenHandler creates a JWT token for a given api key pair.
func (p *TokenProvider) TokenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &ApiKeyAuthRequest{}
		err := c.Bind(request)
		if err != nil {
			p.Logger.Error().Err(err).Msg("Failed to bind request")
			return err
		}
		log := p.Logger.With().Str("email", request.Email).Logger()
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()

		// Assert Api Key, then Get the request if the api key has matched successfully.
		if err := p.ApiKeyAuth.Match(ctx, &models.APIKey{
			APIKeyID:     request.APIKeyID,
			APIKeySecret: []byte(request.APIKeySecret),
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to match api keys.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("Failed to match api keys", http.StatusInternalServerError, err))
		}

		u, err := p.UserStore.GetByEmail(ctx, request.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("Failed to get user", http.StatusInternalServerError, err))
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = u.Email // from context
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(p.Config.GlobalTokenKey))
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate token.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("Failed to generate token", http.StatusInternalServerError, err))
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}

// GetToken gets the JWT token from the echo context
func (p *TokenProvider) GetToken(c echo.Context) (*jwt.Token, error) {
	// Get the token
	raw := c.Request().Header.Get("Authorization")
	return p.GetTokenRaw(raw)
}

// GetTokenRaw gets the JWT token from a raw token string.
func (p *TokenProvider) GetTokenRaw(raw string) (*jwt.Token, error) {
	split := strings.Split(raw, " ")
	if len(split) != 2 {
		return nil, errors.New("unauthorized")
	}
	raw = split[1]

	token, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		signingMethodError := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		switch token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return []byte(p.Config.GlobalTokenKey), nil
		default:
			return nil, signingMethodError
		}
	})
	if err != nil {
		p.Logger.Error().Err(err).Msg("Failed to get token")
		return nil, err
	}

	// Check and auth the user here. If the user doesn't exist, throw an error.
	// This error can be then used to say that the user needs to register.
	claims := token.Claims.(jwt.MapClaims)
	iEmail, ok := claims["email"]
	if !ok {
		p.Logger.Error().Msg("No email found in token claim.")
		return nil, errors.New("invalid token signature")
	}

	email := iEmail.(string)
	if v, ok := p.cache.Has(email); ok && !v.Expired() {
		return token, nil
	}

	user, err := p.UserStore.GetByEmail(context.Background(), email)
	if err != nil {
		p.Logger.Err(err).Msg("Failed to get user by email.")
		return nil, err
	}

	// cache user.
	// if the email happens to already exist, the cache will be refreshed.
	p.cache.Add(user.Email, user.ID)

	return token, nil
}
