package handlers

import (
	"errors"
	"net/http"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// CommandRunHandlerDependencies defines the dependencies for the commands handler provider.
type CommandRunHandlerDependencies struct {
	Logger           zerolog.Logger
	CommandRunStorer providers.CommandRunStorer
}

// CommandRunHandler is a handler taking care of commands related api calls.
type CommandRunHandler struct {
	CommandRunHandlerDependencies
}

var _ providers.CommandRunHandler = &CommandRunHandler{}

// NewCommandRunHandler creates a new command run handler.
func NewCommandRunHandler(deps CommandRunHandlerDependencies) *CommandRunHandler {
	return &CommandRunHandler{
		CommandRunHandlerDependencies: deps,
	}
}

// GetCommandRun returns details about a command run.
func (cm *CommandRunHandler) GetCommandRun() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			kapiErr := kerr.APIError("failed to parse parameter", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, kapiErr)
		}
		cr, err := cm.CommandRunStorer.Get(c.Request().Context(), n)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.NoContent(http.StatusNotFound)
			}
			kapiErr := kerr.APIError("failed to get command run", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, kapiErr)
		}
		return c.JSON(http.StatusOK, cr)
	}
}
