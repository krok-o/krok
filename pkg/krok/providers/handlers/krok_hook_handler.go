package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// HookDependencies defines the dependencies of this server.
type HookDependencies struct {
	// Load all the hook providers here and decide to which to delegate to.
	Logger            zerolog.Logger
	RepositoryStore   providers.RepositoryStorer
	PlatformProviders map[int]providers.Platform
	Executer          providers.Executor
	EventsStorer      providers.EventsStorer
}

// KrokHookHandler is the main hook handler.
type KrokHookHandler struct {
	HookDependencies
}

// NewHookHandler creates a new Krok server to listen for all hook related events.
func NewHookHandler(deps HookDependencies) *KrokHookHandler {
	return &KrokHookHandler{
		HookDependencies: deps,
	}
}

// HandleHooks creates a hook handler.
func (k *KrokHookHandler) HandleHooks() echo.HandlerFunc {
	return func(c echo.Context) error {
		getID := func(id string) (int, error) {
			i := c.Param(id)
			if i == "" {
				return 0, errors.New("id is empty")
			}

			n, err := strconv.Atoi(i)
			if err != nil {
				return 0, errors.New("parameter is not a valid integer")
			}
			return n, nil
		}
		rid, err := getID("rid")
		if err != nil {
			apiError := kerr.APIError("invalid repository id", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		vid, err := getID("vid")
		if err != nil {
			apiError := kerr.APIError("invalid platform id", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		log := k.Logger.With().Int("rid", rid).Int("vid", vid).Logger()
		ctx := c.Request().Context()
		repo, err := k.RepositoryStore.Get(ctx, rid)
		if err != nil {
			apiError := kerr.APIError("can't find repository", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		log.Debug().Str("name", repo.Name).Msg("Found repository...")
		// Validate the request that it's a valid and subscribed to event.
		provider, ok := k.PlatformProviders[vid]
		if !ok {
			err := fmt.Errorf("vcs provider with id %d is not supported", vid)
			return c.JSON(http.StatusBadRequest, kerr.APIError("unable to find vcs provider", http.StatusBadRequest, err))
		}
		//if err := provider.ValidateRequest(ctx, c.Request(), repo.ID); err != nil {
		//	apiError := kerr.APIError("failed to validate hook request", http.StatusBadRequest, err)
		//	return c.JSON(http.StatusBadRequest, apiError)
		//}
		payload, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			apiError := kerr.APIError("failed to get payload", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		id, err := provider.GetEventID(ctx, c.Request())
		if err != nil {
			apiError := kerr.APIError("failed to get event ID from provider", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		event := &models.Event{
			RepositoryID: rid,
			CreateAt:     time.Now(), // TODO: replace this with the timer thingy.
			EventID:      id,
			Payload:      string(payload),
		}
		// Create an ID for this event from the database.
		storedEvent, err := k.EventsStorer.Create(ctx, event)
		if err != nil {
			apiError := kerr.APIError("failed to store event", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		// Create a run which runs the commands attached to this event.
		if err := k.Executer.CreateRun(ctx, storedEvent, repo.Commands); err != nil {
			apiError := kerr.APIError("failed to start run for event", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		return c.String(http.StatusOK, "successfully processed event")
	}
}
