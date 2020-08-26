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

// KrokServer is the server's main handler.
type KrokServer struct {
	Dependencies
	Config
}

// Krok represents what the krok server is capable off.
type Krok interface {
	HandleHooks(ctx context.Context) echo.HandlerFunc
}

// NewKrok creates a new krok server to listen for all hook related events.
func NewKrok(cfg Config, deps Dependencies) *KrokServer {
	return &KrokServer{
		Config:       cfg,
		Dependencies: deps,
	}
}

// HandleHooks creates a hook handler.
func (k *KrokServer) HandleHooks(ctx context.Context) echo.HandlerFunc {
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
