package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/krok-o/krok/pkg/krok/providers"
)

type AuthHandler struct {
	oauthProvider providers.OAuthProvider
}

func NewAuthHandler(oauthProvider providers.OAuthProvider) *AuthHandler {
	return &AuthHandler{oauthProvider: oauthProvider}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		redirectURL := c.QueryParam("redirect_url")

		state, err := h.oauthProvider.GenerateState(redirectURL)
		if err != nil {
			return c.String(http.StatusUnauthorized, "")
		}

		url := h.oauthProvider.GetAuthCodeURL(state)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *AuthHandler) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		state := c.QueryParam("state")
		code := c.QueryParam("code")

		redirectURL, err := h.oauthProvider.VerifyState(state)
		if err != nil {
			return c.String(http.StatusUnauthorized, "")
		}

		token, err := h.oauthProvider.Exchange(ctx, code)
		if err != nil {
			return c.String(http.StatusUnauthorized, "")
		}

		fmt.Println(token)

		return c.Redirect(http.StatusPermanentRedirect, redirectURL)
	}
}
