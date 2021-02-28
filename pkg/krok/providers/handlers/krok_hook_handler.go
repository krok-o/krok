package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// HookDependencies defines the dependencies of this server.
type HookDependencies struct {
	// Load all the hook providers here and decide to which to delegate to.
	Logger          zerolog.Logger
	RepositoryStore providers.RepositoryStorer
	CommandStore    providers.CommandStorer
	GithubPlatform  providers.Platform
	Watcher         providers.Watcher
}

// KrokHookHandler is the main hook handler.
type KrokHookHandler struct {
	HookDependencies
	Config
}

// NewHookHandler creates a new Krok server to listen for all hook related events.
func NewHookHandler(cfg Config, deps HookDependencies) *KrokHookHandler {
	return &KrokHookHandler{
		Config:           cfg,
		HookDependencies: deps,
	}
}

// HandleHooks creates a hook handler.
func (k *KrokHookHandler) HandleHooks() echo.HandlerFunc {
	return func(c echo.Context) error {
		getId := func(id string) (int, error) {
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
		rid, err := getId("rid")
		if err != nil {
			apiError := kerr.APIError("invalid repository id", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		vid, err := getId("vid")
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
		switch vid {
		case models.GITHUB:
			if err := k.GithubPlatform.ValidateRequest(ctx, c.Request(), repo.ID); err != nil {
				apiError := kerr.APIError("failed to validate hook request", http.StatusBadRequest, err)
				return c.JSON(http.StatusBadRequest, apiError)
			}
		}

		// all good, run all attached commands in separate go routines.
		execute := make([]krok.Execute, 0)
		for _, c := range repo.Commands {
			log.Debug().Str("name", c.Name).Msg("Loading command...")
			p, err := k.Watcher.Load(ctx, c.Location)
			if err != nil {
				log.Debug().Str("name", c.Name).Msg("Failed to load command... ignoring.")
				continue
			}
			execute = append(execute, p)
		}

		// Run all commands.
		errs := make([]error, 0)
		errChan := make(chan error, len(execute))
		var wg sync.WaitGroup
		payload, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			apiError := kerr.APIError("failed to get payload", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		defer c.Request().Body.Close()

		for _, e := range execute {
			wg.Add(1)
			go func(ex krok.Execute, errChan chan error) {
				defer wg.Done()
				if _, _, err := ex(string(payload)); err != nil {
					errChan <- err
				}
			}(e, errChan)
		}

		wg.Wait()

		if len(errs) != 0 {
			log.Error().Errs("command_errors", errs).Msg("The following errors happened while executing...")
			apiError := kerr.APIError("failed to validate hook request", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		return c.String(http.StatusOK, "successfully processed event")
	}
}
