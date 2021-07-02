package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const (
	// AccessTokenCookie is the name of the access token cookie.
	AccessTokenCookie = "_a_token_"
	// RefreshTokenCookie is the name of the refresh token cookie.
	RefreshTokenCookie = "_r_token_"
)

// UserAuthHandlerDeps contains the UserAuthHandler dependencies.
type UserAuthHandlerDeps struct {
	Logger      zerolog.Logger
	OAuth       providers.OAuthAuthenticator
	TokenIssuer providers.TokenIssuer
}

// UserAuthHandler handles user authentication.
type UserAuthHandler struct {
	UserAuthHandlerDeps
}

// NewUserAuthHandler creates a new UserAuthHandler.
func NewUserAuthHandler(deps UserAuthHandlerDeps) *UserAuthHandler {
	return &UserAuthHandler{UserAuthHandlerDeps: deps}
}

// OAuthLogin handles a user login.
// swagger:operation GET /auth/login userLogin
//
// User login.
// ---
// parameters:
// - name: redirect_url
//   in: query
//   description: the redirect URL coming from Google to redirect login to
//   required: true
//   type: string
// responses:
//   '307':
//     description: 'the redirect url to the login'
//   '404':
//     description: 'error invalid redirect_url'
//   '401':
//     description: 'error generating state'
func (h *UserAuthHandler) OAuthLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectURL := c.QueryParam("redirect_url")
		if redirectURL == "" {
			h.Logger.Warn().Msg("missing redirect url")
			return c.String(http.StatusBadRequest, "error invalid redirect_url")
		}

		log := h.Logger.With().Str("redirect_url", redirectURL).Logger()

		state, err := h.OAuth.GenerateState(redirectURL)
		if err != nil {
			log.Debug().Err(err).Msg("failed to generate state")
			return c.String(http.StatusUnauthorized, "error generating state")
		}

		url := h.OAuth.GetAuthCodeURL(state)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// OAuthCallback handles the user login callback.
// swagger:operation GET /auth/callback userCallback
//
// This is the url to which Google calls back after a successful login.
// Creates a cookie which will hold the authenticated user.
// ---
// parameters:
// - name: state
//   in: query
//   description: the state variable defined by Google
//   required: true
//   type: string
// - name: code
//   in: query
//   description: the authorization code provided by Google
//   required: true
//   type: string
// responses:
//   '308':
//     description: 'the permanent redirect url'
//   '404':
//     description: 'error invalid state|code'
//   '401':
//     description: 'error verifying state | error during token exchange'
func (h *UserAuthHandler) OAuthCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := h.Logger.With().Logger()

		state := c.QueryParam("state")
		if state == "" {
			log.Warn().Msg("error verifying state")
			return c.String(http.StatusBadRequest, "error invalid state")
		}

		code := c.QueryParam("code")
		if code == "" {
			log.Warn().Msg("error verifying state")
			return c.String(http.StatusBadRequest, "error invalid code")
		}

		redirectURL, err := h.OAuth.VerifyState(state)
		if err != nil {
			log.Error().Err(err).Msg("error verifying state")
			return c.String(http.StatusUnauthorized, "error verifying state")
		}

		token, err := h.OAuth.Exchange(ctx, code)
		if err != nil {
			log.Error().Err(err).Msg("error during token exchange")
			return c.String(http.StatusUnauthorized, "error during token exchange")
		}
		setCookies(c, token)

		return c.Redirect(http.StatusPermanentRedirect, redirectURL)
	}
}

// Refresh handles user token refreshing.
// swagger:operation POST /auth/refresh refreshToken
//
// Refresh the authentication token.
//
// ---
// responses:
//  '200':
//     description: 'Status OK'
//  '401':
//     description: 'refresh token cookie not found|error refreshing token'
func (h *UserAuthHandler) Refresh() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := h.Logger.With().Logger()

		rtoken, err := c.Cookie(RefreshTokenCookie)
		if err != nil {
			log.Error().Err(err).Msg("refresh token cookie not found")
			return c.String(http.StatusUnauthorized, "error getting refresh token")
		}

		token, err := h.TokenIssuer.Refresh(ctx, rtoken.Value)
		if err != nil {
			log.Error().Err(err).Msg("error refreshing token")
			return c.String(http.StatusUnauthorized, "error refreshing token")
		}
		setCookies(c, token)

		return c.NoContent(http.StatusOK)
	}
}

func setCookies(c echo.Context, token *oauth2.Token) {
	c.SetCookie(&http.Cookie{
		Path:     "/",
		Name:     AccessTokenCookie,
		Value:    token.AccessToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	c.SetCookie(&http.Cookie{
		Path:     "/",
		Name:     RefreshTokenCookie,
		Value:    token.RefreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
