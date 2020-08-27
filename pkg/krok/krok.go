package krok

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rs/zerolog"
)

// Config defines configuration which this bot needs to Run.
type Config struct {
}

// Dependencies defines the dependencies of this server.
type Dependencies struct {
	// Load all the hook providers here and decide to which to delegate to.
	Logger zerolog.Logger
}

// HookHandler is the server's main handler.
type HookHandler struct {
	Dependencies
	Config
}

// Handler represents what the krok server is capable off.
type Handler interface {
	HandleHooks(ctx context.Context) echo.HandlerFunc
}

// NewHookHandler creates a new krok server to listen for all hook related events.
func NewHookHandler(cfg Config, deps Dependencies) *HookHandler {
	return &HookHandler{
		Config:       cfg,
		Dependencies: deps,
	}
}

// HandleHooks creates a hook handler.
func (k *HookHandler) HandleHooks(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.NoContent(http.StatusBadRequest)
		}
		// Get the right type based on the saved data.
		// Fire off the handler for this event.
		// Return result.
		// switch based on Model Type and use the appropriate provider
		return c.String(http.StatusOK, "successfully processed event")
	}
}
