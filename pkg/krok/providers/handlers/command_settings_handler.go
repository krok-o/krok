package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// CommandSettingsHandlerDependencies defines the dependencies for the command settings handler provider.
type CommandSettingsHandlerDependencies struct {
	Logger        zerolog.Logger
	CommandStorer providers.CommandStorer
}

// CommandSettingsHandler is a handler taking care of command settings related api calls.
type CommandSettingsHandler struct {
	CommandSettingsHandlerDependencies
}

var _ providers.CommandSettingsHandler = &CommandSettingsHandler{}

// NewCommandSettingsHandler creates a new command settings handler.
func NewCommandSettingsHandler(deps CommandSettingsHandlerDependencies) *CommandSettingsHandler {
	return &CommandSettingsHandler{
		CommandSettingsHandlerDependencies: deps,
	}
}

// Delete deletes a setting.
func (ch *CommandSettingsHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		if err := ch.CommandStorer.DeleteSetting(ctx, n); err != nil {
			ch.Logger.Debug().Err(err).Msg("Command Setting Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete command setting", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// List lists commands.
func (ch *CommandSettingsHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		list, err := ch.CommandStorer.ListSettings(ctx, n)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Command List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list commands", http.StatusBadRequest, err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// Get returns a specific setting.
func (ch *CommandSettingsHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		repo, err := ch.CommandStorer.GetSetting(ctx, n)
		if err != nil {
			apiError := kerr.APIError("failed to get command setting", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		return c.JSON(http.StatusOK, repo)
	}
}

// Update updates a setting.
func (ch *CommandSettingsHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		setting := &models.CommandSetting{}
		if err := c.Bind(setting); err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to bind command.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind command", http.StatusBadRequest, err))
		}

		ctx := c.Request().Context()
		if err := ch.CommandStorer.UpdateSetting(ctx, setting); err != nil {
			ch.Logger.Debug().Err(err).Msg("Command setting update failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to update command setting", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Create creates a command setting.
func (ch *CommandSettingsHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		setting := &models.CommandSetting{}
		if err := c.Bind(setting); err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to bind command setting.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind command setting", http.StatusBadRequest, err))
		}

		ctx := c.Request().Context()
		if err := ch.CommandStorer.CreateSetting(ctx, setting); err != nil {
			ch.Logger.Debug().Err(err).Msg("Command setting create failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to create command setting", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusCreated)
	}
}
