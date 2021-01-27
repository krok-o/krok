package krok

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
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

// Handler represents what the Krok server is capable off.
type Handler interface {
	HandleHooks(ctx context.Context) echo.HandlerFunc
}

// NewHookHandler creates a new Krok server to listen for all hook related events.
func NewHookHandler(cfg Config, deps Dependencies) *HookHandler {
	return &HookHandler{
		Config:       cfg,
		Dependencies: deps,
	}
}

// HandleHooks creates a hook handler.
func (k *HookHandler) HandleHooks(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get what repository that links to
		// Get the VCS type for that repository
		// Instantiate the right platform providers
		// Validate the request
		// Send the data to all linked commands
		return c.String(http.StatusOK, "successfully processed event")
	}
}
