package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// Config has the configuration options for the repository handler.
type Config struct {
	Proto              string
	Hostname           string
	GlobalTokenKey     string
	GoogleClientID     string
	GoogleClientSecret string
}

// Dependencies defines the dependencies for the repository handler provider.
type Dependencies struct {
	Logger      zerolog.Logger
	UserStore   providers.UserStorer
	ApiKeyAuth  providers.ApiKeysAuthenticator
	TokenIssuer providers.TokenIssuer
}

// TokenHandler is a token provider for the handlers.
type TokenHandler struct {
	Dependencies
}

// NewTokenHandler creates a new token handler which deals with generating and handling tokens.
func NewTokenHandler(deps Dependencies) (*TokenHandler, error) {
	tp := &TokenHandler{
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
func (p *TokenHandler) TokenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &ApiKeyAuthRequest{}
		err := c.Bind(request)
		if err != nil {
			p.Logger.Error().Err(err).Msg("Failed to bind request")
			return err
		}
		log := p.Logger.With().Str("email", request.Email).Logger()

		ctx := c.Request().Context()

		// Assert Api Key, then Get the request if the api key has matched successfully.
		if err := p.ApiKeyAuth.Match(ctx, &models.APIKey{
			APIKeyID:     request.APIKeyID,
			APIKeySecret: request.APIKeySecret,
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to match api keys.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("Failed to match api keys", http.StatusInternalServerError, err))
		}

		u, err := p.UserStore.GetByEmail(ctx, request.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("Failed to get user", http.StatusInternalServerError, err))
		}

		t, err := p.TokenIssuer.Create(u)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to generate token", http.StatusInternalServerError, err))
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t.AccessToken,
		})
	}
}
