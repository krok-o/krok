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
	APIKeyAuth  providers.APIKeysAuthenticator
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

// TokenHandler creates a JWT token for a given api key pair.
// swagger:operation POST /get-token getToken
// Creates a JWT token for a given api key pair.
// ---
// deprecated: true
// produces:
// - application/json
// responses:
//   '200':
//     description: 'the generated JWT token'
//     schema:
//       "$ref": "#/responses/TokenResponse"
//   '500':
//     description: 'when there was a problem with matching the email, or the api key or generating the token'
func (p *TokenHandler) TokenHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &models.APIKeyAuthRequest{}
		err := c.Bind(request)
		if err != nil {
			p.Logger.Error().Err(err).Msg("Failed to bind request")
			return err
		}
		log := p.Logger.With().Str("email", request.Email).Logger()

		ctx := c.Request().Context()

		// Assert API Key, then Get the request if the api key has matched successfully.
		if err := p.APIKeyAuth.Match(ctx, &models.APIKey{
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

		tr := &models.TokenResponse{
			Token: t.AccessToken,
		}
		return c.JSON(http.StatusOK, tr)
	}
}
