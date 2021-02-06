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
	TokenIssuer providers.UserTokenIssuer
}

// UserAuthHandler handles user authentication.
type UserAuthHandler struct {
	UserAuthHandlerDeps
}

// NewUserAuthHandler creates a new UserAuthHandler.
func NewUserAuthHandler(deps UserAuthHandlerDeps) *UserAuthHandler {
	return &UserAuthHandler{UserAuthHandlerDeps: deps}
}

// Login handles a user login.
func (h *UserAuthHandler) Login() echo.HandlerFunc {
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

// Callback handles the user login callback.
func (h *UserAuthHandler) Callback() echo.HandlerFunc {
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

		return c.String(http.StatusOK, "")
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
