package handlers

import (
	"errors"
	"net/http"
	"strconv"

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
	TokenProvider *TokenHandler
}

// CommandsHandler is a handler taking care of commands related api calls.
type CommandsHandler struct {
	Config
	CommandsHandlerDependencies
}

var _ providers.CommandHandler = &CommandsHandler{}

// NewCommandsHandler creates a new commands handler.
func NewCommandsHandler(cfg Config, deps CommandsHandlerDependencies) (*CommandsHandler, error) {
	return &CommandsHandler{
		Config:                      cfg,
		CommandsHandlerDependencies: deps,
	}, nil
}

// Delete deletes a command.
func (ch *CommandsHandler) Delete() echo.HandlerFunc {
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

		if err := ch.CommandStorer.Delete(ctx, n); err != nil {
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

		repo, err := ch.CommandStorer.Get(ctx, n)
		if err != nil {
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
			ch.Logger.Debug().Err(err).Msg("Command Update failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to update command", http.StatusBadRequest, err))
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// Create is unimplemented.
func (ch *CommandsHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, kerr.APIError("unimplemented", http.StatusInternalServerError, errors.New("unimplemented")))
	}
}

// AddCommandRelForRepository adds a command relationship to a repository.
func (ch *CommandsHandler) AddCommandRelForRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		cmdID := c.Param("cmdid")
		if cmdID == "" {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		repoID := c.Param("repoid")
		if repoID == "" {
			apiError := kerr.APIError("invalid repository id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		cn, err := strconv.Atoi(cmdID)
		if err != nil {
			apiError := kerr.APIError("failed to convert command id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		rn, err := strconv.Atoi(repoID)
		if err != nil {
			apiError := kerr.APIError("failed to convert repository id to number", http.StatusBadRequest, err)
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
		cmdID := c.Param("cmdid")
		if cmdID == "" {
			apiError := kerr.APIError("invalid command id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		repoID := c.Param("repoid")
		if repoID == "" {
			apiError := kerr.APIError("invalid repository id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		cn, err := strconv.Atoi(cmdID)
		if err != nil {
			apiError := kerr.APIError("failed to convert command id to number", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		rn, err := strconv.Atoi(repoID)
		if err != nil {
			apiError := kerr.APIError("failed to convert repository id to number", http.StatusBadRequest, err)
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
