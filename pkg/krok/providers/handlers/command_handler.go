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

// CommandsHandlerDependencies defines the dependencies for the commands handler provider.
type CommandsHandlerDependencies struct {
	Logger        zerolog.Logger
	CommandStorer providers.CommandStorer
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
func (ch *CommandsHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		if err := ch.CommandStorer.Delete(ctx, n); err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("command not found", http.StatusNotFound, err))
			}
			ch.Logger.Debug().Err(err).Msg("Command Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete command", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// List lists commands.
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

// Update updates a command.
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
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to update command", http.StatusBadRequest, err))
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// AddCommandRelForRepository adds a command relationship to a repository.
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
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to add command relationship to repository", http.StatusBadRequest, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// RemoveCommandRelForRepository removes a relationship of a command from a repository.
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
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to remove command relationship to repository", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// AddCommandRelForPlatform adds a command relationship to a platform.
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
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to add command relationship to platform", http.StatusBadRequest, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// RemoveCommandRelForPlatform removes a relationship of a command from a platform.
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
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to remove command relationship to platform", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusOK)
	}
}
