package handlers

import (
	"errors"
	"net/http"

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
// swagger:operation DELETE /command/settings/{id} deleteCommandSetting
// Deletes a given command setting.
// ---
// parameters:
// - name: id
//   in: path
//   description: 'The ID of the command setting to delete'
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     description: 'OK in case the deletion was successful'
//   '400':
//     description: 'invalid id'
//     schema:
//       "$ref": "#/responses/Message"
//   '404':
//     description: 'command setting not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'when the deletion operation failed'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandSettingsHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		if err := ch.CommandStorer.DeleteSetting(ctx, n); err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("command setting not found", http.StatusNotFound, err))
			}
			ch.Logger.Debug().Err(err).Msg("Command Setting Delete failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to delete command setting", http.StatusInternalServerError, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// List lists command settings.
// swagger:operation POST /command/{id}/settings listCommandSettings
// List settings for a command.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   description: 'The ID of the command to list settings for'
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/CommandSetting"
//   '400':
//     description: 'invalid id'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to list settings'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandSettingsHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		list, err := ch.CommandStorer.ListSettings(ctx, n)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Command List failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to list commands", http.StatusInternalServerError, err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// Get returns a specific setting.
// swagger:operation GET /command/settings/{id} getCommandSetting
// Get a specific setting.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   description: 'The ID of the command setting to retrieve'
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     schema:
//       "$ref": "#/definitions/CommandSetting"
//   '400':
//     description: 'invalid command id'
//     schema:
//       "$ref": "#/responses/Message"
//   '404':
//     description: 'command setting not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to get command setting'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandSettingsHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		repo, err := ch.CommandStorer.GetSetting(ctx, n)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("command setting not found", http.StatusNotFound, err))
			}
			apiError := kerr.APIError("failed to get command setting", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}

		return c.JSON(http.StatusOK, repo)
	}
}

// Update updates a setting.
// swagger:operation POST /command/settings/update updateCommandSetting
// Updates a given command setting.
// ---
// produces:
// - application/json
// parameters:
// - name: setting
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/CommandSetting"
// responses:
//   '200':
//     description: 'successfully updated command setting'
//   '400':
//     description: 'binding error'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to update the command setting'
//     schema:
//       "$ref": "#/responses/Message"
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
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to update command setting", http.StatusInternalServerError, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Create creates a command setting.
// swagger:operation POST /command/settings/update updateCommandSetting
// Create a new command setting.
// ---
// produces:
// - application/json
// parameters:
// - name: setting
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/CommandSetting"
// responses:
//   '200':
//     description: 'successfully created command setting'
//   '400':
//     description: 'binding error'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to create the command setting'
//     schema:
//       "$ref": "#/responses/Message"
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
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to create command setting", http.StatusInternalServerError, err))
		}

		return c.NoContent(http.StatusCreated)
	}
}
