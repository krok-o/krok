package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
	krokmiddleware "github.com/krok-o/krok/pkg/server/middleware"
)

type UserTokenHandlerDeps struct {
	Logger     zerolog.Logger
	UserStore  providers.UserStorer
	APIKeyAuth providers.ApiKeysAuthenticator
}

type UserTokenHandler struct {
	UserTokenHandlerDeps
}

func NewUserTokenHandler(deps UserTokenHandlerDeps) *UserTokenHandler {
	return &UserTokenHandler{UserTokenHandlerDeps: deps}
}

func (h *UserTokenHandler) Generate() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		uc, err := krokmiddleware.GetUserContext(c)
		if err != nil {
			h.Logger.Debug().Err(err).Msg("error getting user context")
			return c.String(http.StatusInternalServerError, "failed to get user context")
		}

		user, err := h.UserStore.Get(ctx, uc.UserID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		secret := uuid.New()
		token, err := h.APIKeyAuth.Encrypt(ctx, []byte(secret.String()))
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		user.Token = string(token)
		updated, err := h.UserStore.Update(ctx, user)
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		return c.JSON(http.StatusOK, map[string]string{"token": updated.Token})
	}
}

func (h *UserTokenHandler) Revoke() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		uc, err := krokmiddleware.GetUserContext(c)
		if err != nil {
			h.Logger.Debug().Err(err).Msg("error getting user context")
			return c.String(http.StatusInternalServerError, "failed to get user context")
		}

		user, err := h.UserStore.Get(ctx, uc.UserID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		user.Token = ""
		if _, err := h.UserStore.Update(ctx, user); err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		return c.String(http.StatusOK, "")
	}
}
