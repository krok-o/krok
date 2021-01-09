package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// CommandsHandlerDependencies defines the dependencies for the commands handler provider.
type CommandsHandlerDependencies struct {
	Dependencies
	CommandStorer providers.CommandStorer
	TokenProvider *TokenProvider
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

// DeleteCommand deletes a command.
func (ch *CommandsHandler) DeleteCommand() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ch.TokenProvider.GetToken(c)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
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
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		if err := ch.CommandStorer.Delete(ctx, n); err != nil {
			ch.Logger.Debug().Err(err).Msg("Command Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete command", http.StatusBadRequest, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// ListCommands lists commands.
func (ch *CommandsHandler) ListCommands() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ch.TokenProvider.GetToken(c)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}

		opts := &models.ListOptions{}
		if err := c.Bind(opts); err != nil {
			// if we don't have anything to bind, just ignore opts.
			opts = nil
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		list, err := ch.CommandStorer.List(ctx, opts)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Command List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list commands", http.StatusBadRequest, err))
		}
		return c.JSON(http.StatusOK, list)
	}
}

// GetCommand returns a specific command.
func (ch *CommandsHandler) GetCommand() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ch.TokenProvider.GetToken(c)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
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
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		repo, err := ch.CommandStorer.Get(ctx, n)
		if err != nil {
			apiError := kerr.APIError("failed to get command", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		return c.JSON(http.StatusOK, repo)
	}
}

// UpdateCommand updates a command.
func (ch *CommandsHandler) UpdateCommand() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ch.TokenProvider.GetToken(c)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
		command := &models.Command{}
		err = c.Bind(command)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to bind command.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind command", http.StatusBadRequest, err))
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		updated, err := ch.CommandStorer.Update(ctx, command)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Command Update failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to update command", http.StatusBadRequest, err))
		}
		return c.JSON(http.StatusOK, updated)
	}
}

// AddCommandRelForRepository adds a command relationship to a repository.
func (ch *CommandsHandler) AddCommandRelForRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ch.TokenProvider.GetToken(c)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
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
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
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
		_, err := ch.TokenProvider.GetToken(c)
		if err != nil {
			ch.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
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
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		if err := ch.CommandStorer.RemoveCommandRelForRepository(ctx, cn, rn); err != nil {
			ch.Logger.Debug().Err(err).Msg("RemoveCommandRelForRepository failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to remove command relationship to repository", http.StatusBadRequest, err))
		}
		return c.NoContent(http.StatusOK)
	}
}
