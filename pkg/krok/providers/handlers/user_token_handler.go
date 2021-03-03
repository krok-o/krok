package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	krokmiddleware "github.com/krok-o/krok/pkg/server/middleware"
)

// UserTokenHandlerDeps represents the UserTokenHandler dependencies.
type UserTokenHandlerDeps struct {
	Logger     zerolog.Logger
	UserStore  providers.UserStorer
	APIKeyAuth providers.APIKeysAuthenticator
	UUID       providers.UUIDGenerator
}

// UserTokenHandler represents the user personal token handler.
type UserTokenHandler struct {
	UserTokenHandlerDeps
}

// NewUserTokenHandler creates a new UserTokenHandler.
func NewUserTokenHandler(deps UserTokenHandlerDeps) *UserTokenHandler {
	return &UserTokenHandler{UserTokenHandlerDeps: deps}
}

// Generate generates users new personal access token.
func (h *UserTokenHandler) Generate() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		uc, err := krokmiddleware.GetUserContext(c)
		if err != nil {
			h.Logger.Warn().Err(err).Msg("error getting user context")
			apiErr := kerr.APIError("Failed to get the user context.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiErr)
		}

		user, err := h.UserStore.Get(ctx, uc.UserID)
		if err != nil {
			h.Logger.Warn().Int("user_id", uc.UserID).Err(err).Msg("failed to get the user")
			apiErr := kerr.APIError("Failed to get the user.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiErr)
		}

		secret, err := h.UUID.Generate()
		if err != nil {
			h.Logger.Error().Int("user_id", uc.UserID).Err(err).Msg("failed to generate the uuid secret")
			apiErr := kerr.APIError("Failed to generate the secret.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiErr)
		}

		token, err := h.APIKeyAuth.Encrypt(ctx, []byte(secret))
		if err != nil {
			h.Logger.Error().Int("user_id", uc.UserID).Err(err).Msg("failed to encrypt the uuid secret")
			apiErr := kerr.APIError("Failed to encrypt the secret.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiErr)
		}

		user.Token = string(token)
		updated, err := h.UserStore.Update(ctx, user)
		if err != nil {
			h.Logger.Error().Int("user_id", uc.UserID).Err(err).Msg("failed to update the user")
			apiErr := kerr.APIError("Failed to update the user.", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiErr)
		}

		h.Logger.Debug().Int("user_id", updated.ID).Msg("successfully generated a new token")
		return c.JSON(http.StatusOK, map[string]string{"token": updated.Token})
	}
}
