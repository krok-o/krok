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

// EventHandlerDependencies defines the dependencies for the vcs token handler provider.
type EventHandlerDependencies struct {
	Logger       zerolog.Logger
	EventsStorer providers.EventsStorer
}

// EventHandler is a handler taking care of vcs token related api calls.
type EventHandler struct {
	EventHandlerDependencies
}

var _ providers.EventHandler = &EventHandler{}

// NewEventHandler creates a new event handler.
func NewEventHandler(deps EventHandlerDependencies) *EventHandler {
	return &EventHandler{
		EventHandlerDependencies: deps,
	}
}

// List handles the list rest event.
func (r *EventHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		opts := &models.ListOptions{}
		if err := c.Bind(opts); err != nil {
			// if we don't have anything to bind, just ignore opts.
			opts = nil
		}
		n, err := GetParamAsInt("repoid", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()

		list, err := r.EventsStorer.ListEventsForRepository(ctx, n, opts)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Event List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list events", http.StatusBadRequest, err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// Get retrieves a repository and displays the unique URL for which this repo is responsible for.
func (r *EventHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		// Get the event from store.
		event, err := r.EventsStorer.GetEvent(ctx, n)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("event not found", http.StatusNotFound, err))
			}
			apiError := kerr.APIError("failed to get event", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		return c.JSON(http.StatusOK, event)
	}
}
