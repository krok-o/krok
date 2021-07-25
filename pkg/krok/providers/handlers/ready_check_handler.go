package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
)

// ReadyCheckHandlerDependencies defines the dependencies for the ready check handler provider.
type ReadyCheckHandlerDependencies struct {
	Logger  zerolog.Logger
	Checker providers.Ready
}

// ReadyCheckHandler is a handler taking care of ready check related api calls.
type ReadyCheckHandler struct {
	ReadyCheckHandlerDependencies
}

var _ providers.ReadyHandler = &ReadyCheckHandler{}

// NewReadyCheckHandler creates a new ready check handler.
func NewReadyCheckHandler(deps ReadyCheckHandlerDependencies) *ReadyCheckHandler {
	return &ReadyCheckHandler{
		ReadyCheckHandlerDependencies: deps,
	}
}

// Ready check if Krok is ready to handle requests.
func (r *ReadyCheckHandler) Ready() echo.HandlerFunc {
	return func(c echo.Context) error {
		if ready := r.Checker.Ready(c.Request().Context()); !ready {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}
