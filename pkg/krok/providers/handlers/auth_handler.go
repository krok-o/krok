package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const (
	AccessTokenCookie  = "_a_token_"
	RefreshTokenCookie = "_r_token_"
)

type AuthHandler struct {
	OAuthProvider providers.OAuthProvider
	TokenIssuer   providers.TokenIssuer
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectURL := c.QueryParam("redirect_url")
		if redirectURL == "" {
			return c.String(http.StatusBadRequest, "error invalid redirect_url")
		}

		state, err := h.OAuthProvider.GenerateState(redirectURL)
		if err != nil {
			return c.String(http.StatusUnauthorized, "")
		}

		url := h.OAuthProvider.GetAuthCodeURL(state)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *AuthHandler) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		state := c.QueryParam("state")
		if state == "" {
			return c.String(http.StatusBadRequest, "error invalid state")
		}

		code := c.QueryParam("code")
		if code == "" {
			return c.String(http.StatusBadRequest, "error invalid code")
		}

		redirectURL, err := h.OAuthProvider.VerifyState(state)
		if err != nil {
			return c.String(http.StatusUnauthorized, "error verifying state")
		}

		token, err := h.OAuthProvider.Exchange(ctx, code)
		if err != nil {
			return c.String(http.StatusUnauthorized, "error during token exchange")
		}
		setCookies(c, token)

		return c.Redirect(http.StatusPermanentRedirect, redirectURL)
	}
}

func (h *AuthHandler) Refresh() echo.HandlerFunc {
	return func(c echo.Context) error {
		rtoken, err := c.Cookie(RefreshTokenCookie)
		if err != nil {
			return c.String(http.StatusUnauthorized, "error getting refresh token")
		}

		token, err := h.TokenIssuer.Refresh(rtoken.Value)
		if err != nil {
			return c.String(http.StatusUnauthorized, "error refreshing tokens")
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
