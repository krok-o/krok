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

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// Config has the configuration options for the repository handler.
type Config struct {
	Hostname       string
	GlobalTokenKey string
}

// Dependencies defines the dependencies for the repository handler provider.
type Dependencies struct {
	Logger       zerolog.Logger
	UserStore    providers.UserStorer
	ApiKeysStore providers.APIKeys
}

// TokenProvider is a token provider for the handlers.
type TokenProvider struct {
	Config
	Dependencies
}

// NewTokenProvider creates a new token provider which deals with generating and handling tokens.
func NewTokenProvider(cfg Config, deps Dependencies) (*TokenProvider, error) {
	return &TokenProvider{Config: cfg, Dependencies: deps}, nil
}

// TokenHandler creates a JWT token for a given user.
func (p *TokenProvider) TokenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {

		// TODO: This either needs to get an email in the body to generate a token for,
		// or check if there is an API key and secret provided and generate the token
		// based on that.

		user := &models.User{}
		err := c.Bind(user)
		if err != nil {
			p.Logger.Error().Err(err).Msg("Failed to bind user")
			return err
		}
		log := p.Logger.With().Str("email", user.Email).Logger()
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()

		u, err := p.UserStore.Get(ctx, user.ID)
		if err != nil {
			return err
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
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}

// GetToken gets the JWT token from the echo context
func (p *TokenProvider) GetToken(c echo.Context) (*jwt.Token, error) {
	// Get the token
	jwtRaw := c.Request().Header.Get("Authorization")
	split := strings.Split(jwtRaw, " ")
	if len(split) != 2 {
		return nil, errors.New("unauthorized")
	}
	jwtString := split[1]
	// Parse token
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
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

	return token, nil
}
