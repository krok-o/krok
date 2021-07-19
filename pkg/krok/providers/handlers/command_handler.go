package handlers

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// CommandsHandlerDependencies defines the dependencies for the commands handler provider.
type CommandsHandlerDependencies struct {
	Logger        zerolog.Logger
	CommandStorer providers.CommandStorer
	Plugins       providers.Plugins
}

// CommandsHandler is a handler taking care of commands related api calls.
type CommandsHandler struct {
	CommandsHandlerDependencies
}

var _ providers.CommandHandler = &CommandsHandler{}

// NewCommandsHandler creates a new commands handler.
func NewCommandsHandler(deps CommandsHandlerDependencies) *CommandsHandler {
	return &CommandsHandler{
		CommandsHandlerDependencies: deps,
	}
}

// Delete deletes a command.
// swagger:operation DELETE /command/{id} deleteCommand
// Deletes given command.
// ---
// parameters:
// - name: id
//   in: path
//   description: 'The ID of the command to delete'
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     description: 'OK in case the deletion was successful'
//   '400':
//     description: 'in case of missing user context or invalid ID'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'when the deletion operation failed'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()
		// Get first for the name
		command, err := ch.CommandStorer.Get(ctx, n)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("command not found", http.StatusNotFound, err))
			}
			ch.Logger.Debug().Err(err).Msg("Command Get failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to get command", http.StatusBadRequest, err))
		}

		// Delete from database
		if err := ch.CommandStorer.Delete(ctx, command.ID); err != nil {
			ch.Logger.Debug().Err(err).Msg("Command Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete command", http.StatusBadRequest, err))
		}

		// Delete from storage
		if err := ch.Plugins.Delete(ctx, command.Name); err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to delete file from permanent storage.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to delete file from permanent storage", http.StatusInternalServerError, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// List lists commands.
// swagger:operation POST /commands listCommands
// List commands
// ---
// produces:
// - application/json
// parameters:
// - name: listOptions
//   in: body
//   required: false
//   schema:
//     "$ref": "#/definitions/ListOptions"
// responses:
//   '200':
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Command"
//   '500':
//     description: 'failed to get user context'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		opts := &models.ListOptions{}
		if err := c.Bind(opts); err != nil {
			// if we don't have anything to bind, just ignore opts.
			opts = nil
		}

		ctx := c.Request().Context()

		list, err := ch.CommandStorer.List(ctx, opts)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Command List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list commands", http.StatusBadRequest, err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// Get returns a specific command.
// swagger:operation GET /command/{id} getCommand
// Returns a specific command.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   type: integer
//   format: int
//   required: true
// responses:
//   '200':
//     schema:
//       "$ref": "#/definitions/Command"
//   '400':
//     description: 'invalid command id'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to get user context'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		repo, err := ch.CommandStorer.Get(ctx, n)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("command not found", http.StatusNotFound, err))
			}
			apiError := kerr.APIError("failed to get command", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		return c.JSON(http.StatusOK, repo)
	}
}

// Upload a command. To set up anything for the command, like schedules etc,
// the command has to be edited. We don't support uploading the same thing twice.
// If the command binary needs to be updated, delete the command and upload the
// new binary.
// swagger:operation POST /command uploadCommand
// Upload a command. To set up anything for the command, like schedules etc,
// the command has to be edited. We don't support uploading the same thing twice.
// If the command binary needs to be updated, delete the command and upload the
// new binary.
// ---
// produces:
// - application/json
// responses:
//   '201':
//     description: 'in case of successful file upload'
//     schema:
//       "$ref": "#/definitions/Command"
//   '400':
//     description: 'invalid file format or command already exists'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to upload file, create plugin, create command or copy operations'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) Upload() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		dots := strings.Split(file.Filename, ".")
		if len(dots) == 0 {
			return c.JSON(http.StatusBadRequest, kerr.APIError("file name does not contain a dot", http.StatusBadRequest, err))
		}
		name := dots[0]
		ctx := c.Request().Context()
		// check if name is already taken:
		if _, err := ch.CommandStorer.GetByName(ctx, name); err == nil {
			return c.JSON(http.StatusBadRequest, kerr.APIError("command with name already taken", http.StatusBadRequest, err))
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer func(src multipart.File) {
			if err := src.Close(); err != nil {
				ch.Logger.Debug().Err(err).Msg("Failed to close src.")
			}
		}(src)

		// Destination
		tmp, err := ioutil.TempDir("", "upload_folder")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to create upload folder", http.StatusInternalServerError, err))
		}
		defer func() {
			if err := os.RemoveAll(tmp); err != nil {
				ch.Logger.Debug().Err(err).Msg("Warning, failed to clean up temporary folder. Please do that manually.")
			}
		}()
		dst, err := os.Create(filepath.Join(tmp, file.Filename))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to create upload file", http.StatusInternalServerError, err))
		}
		defer func(dst *os.File) {
			if err := dst.Close(); err != nil {
				ch.Logger.Debug().Err(err).Msg("Failed to close destination file.")
			}
		}(dst)

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to copy over uploaded file", http.StatusInternalServerError, err))
		}

		// Create the file first, then the command.
		f, hash, err := ch.Plugins.Create(ctx, dst.Name())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to create plugin", http.StatusInternalServerError, err))
		}

		command := &models.Command{
			Name:     name,
			Filename: filepath.Base(f),
			Location: filepath.Dir(f),
			Hash:     hash,
			Enabled:  true,
		}

		command, err = ch.CommandStorer.Create(ctx, command)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to create command", http.StatusInternalServerError, err))
		}

		return c.JSON(http.StatusCreated, command)
	}
}

// Update updates a command.
// swagger:operation POST /command/update updateCommand
// Updates a given command.
// ---
// produces:
// - application/json
// parameters:
// - name: command
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/Command"
// responses:
//   '200':
//     description: 'successfully updated command'
//     schema:
//       "$ref": "#/definitions/Command"
//   '400':
//     description: 'binding error'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to update the command'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		command := &models.Command{}
		if err := c.Bind(command); err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to bind command.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind command", http.StatusBadRequest, err))
		}

		ctx := c.Request().Context()

		updated, err := ch.CommandStorer.Update(ctx, command)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("command not found", http.StatusNotFound, err))
			}
			ch.Logger.Debug().Err(err).Msg("Command Update failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to update command", http.StatusInternalServerError, err))
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// AddCommandRelForRepository adds a command relationship to a repository.
// swagger:operation POST /command/add-command-rel-for-repository/{cmdid}/{repoid} addCommandRelForRepositoryCommand
// Add a connection to a repository. This will make this command to be executed for events for that repository.
// ---
// parameters:
// - name: cmdid
//   in: path
//   required: true
//   type: integer
//   format: int
// - name: repoid
//   in: path
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     description: 'successfully added relationship'
//   '400':
//     description: 'invalid ids or repositroy not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to add relationship'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) AddCommandRelForRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		cn, err := GetParamAsInt("cmdid", c)
		if err != nil {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		rn, err := GetParamAsInt("repoid", c)
		if err != nil {
			apiError := kerr.APIError("invalid repo id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		if err := ch.CommandStorer.AddCommandRelForRepository(ctx, cn, rn); err != nil {
			ch.Logger.Debug().Err(err).Msg("AddCommandRelForRepository failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to add command relationship to repository", http.StatusInternalServerError, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// RemoveCommandRelForRepository removes a relationship of a command from a repository.
// swagger:operation POST /command/remove-command-rel-for-repository/{cmdid}/{repoid} removeCommandRelForRepositoryCommand
// Remove a relationship to a repository. This command will no longer be running for that repository events.
// ---
// parameters:
// - name: cmdid
//   in: path
//   required: true
//   type: integer
//   format: int
// - name: repoid
//   in: path
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     description: 'successfully removed relationship'
//   '400':
//     description: 'invalid ids or repositroy not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to add relationship'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) RemoveCommandRelForRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		cn, err := GetParamAsInt("cmdid", c)
		if err != nil {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		rn, err := GetParamAsInt("repoid", c)
		if err != nil {
			apiError := kerr.APIError("invalid repo id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		if err := ch.CommandStorer.RemoveCommandRelForRepository(ctx, cn, rn); err != nil {
			ch.Logger.Debug().Err(err).Msg("RemoveCommandRelForRepository failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to remove command relationship to repository", http.StatusInternalServerError, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// AddCommandRelForPlatform adds a command relationship to a platform.
// swagger:operation POST /command/add-command-rel-for-platform/{cmdid}/{repoid} addCommandRelForPlatformCommand
// Adds a connection to a platform for a command. Defines what platform a command supports. These commands will only be able to run for those platforms.
// ---
// parameters:
// - name: cmdid
//   in: path
//   required: true
//   type: integer
//   format: int
// - name: repoid
//   in: path
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     description: 'successfully added relationship'
//   '400':
//     description: 'invalid ids or platform not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to add command relationship to platform'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) AddCommandRelForPlatform() echo.HandlerFunc {
	return func(c echo.Context) error {
		cn, err := GetParamAsInt("cmdid", c)
		if err != nil {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		pid, err := GetParamAsInt("pid", c)
		if err != nil {
			apiError := kerr.APIError("invalid platform id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		found := false
		for _, p := range models.SupportedPlatforms {
			if p.ID == pid {
				found = true
				break
			}
		}
		if !found {
			apiError := kerr.APIError("patform id not found in supported platforms", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		if err := ch.CommandStorer.AddCommandRelForPlatform(ctx, cn, pid); err != nil {
			ch.Logger.Debug().Err(err).Msg("AddCommandRelForPlatform failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to add command relationship to platform", http.StatusInternalServerError, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// RemoveCommandRelForPlatform removes a relationship of a command from a platform.
// swagger:operation POST /command/remove-command-rel-for-platform/{cmdid}/{repoid} removeCommandRelForPlatformCommand
// Remove a relationship to a platform. This command will no longer be running for that platform events.
// ---
// parameters:
// - name: cmdid
//   in: path
//   required: true
//   type: integer
//   format: int
// - name: repoid
//   in: path
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     description: 'successfully removed relationship'
//   '400':
//     description: 'invalid ids or platform not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to add relationship'
//     schema:
//       "$ref": "#/responses/Message"
func (ch *CommandsHandler) RemoveCommandRelForPlatform() echo.HandlerFunc {
	return func(c echo.Context) error {
		cn, err := GetParamAsInt("cmdid", c)
		if err != nil {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		pid, err := GetParamAsInt("pid", c)
		if err != nil {
			apiError := kerr.APIError("invalid platform id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		if err := ch.CommandStorer.RemoveCommandRelForPlatform(ctx, cn, pid); err != nil {
			ch.Logger.Debug().Err(err).Msg("RemoveCommandRelForPlatform failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to remove command relationship to platform", http.StatusInternalServerError, err))
		}

		return c.NoContent(http.StatusOK)
	}
}
