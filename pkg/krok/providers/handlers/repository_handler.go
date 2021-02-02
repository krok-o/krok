package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// RepoHandlerDependencies defines the dependencies for the repository handler provider.
type RepoHandlerDependencies struct {
	RepositoryStorer  providers.RepositoryStorer
	TokenProvider     *TokenProvider
	Logger            zerolog.Logger
	PlatformProviders map[int]providers.Platform
}

// RepoHandler is a handler taking care of repository related api calls.
type RepoHandler struct {
	Config
	RepoHandlerDependencies
}

var _ providers.RepositoryHandler = &RepoHandler{}

// NewRepositoryHandler creates a new repository handler.
func NewRepositoryHandler(cfg Config, deps RepoHandlerDependencies) (*RepoHandler, error) {
	return &RepoHandler{
		Config:                  cfg,
		RepoHandlerDependencies: deps,
	}, nil
}

// CreateRepository handles the Create rest event.
func (r *RepoHandler) CreateRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := r.TokenProvider.GetToken(c)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
		repo := &models.Repository{}
		err = c.Bind(repo)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to bind repository.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind repository", http.StatusBadRequest, err))
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		created, err := r.RepositoryStorer.Create(ctx, repo)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Repository CreateRepository failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to create repository", http.StatusBadRequest, err))
		}
		uurl, err := r.generateUniqueCallBackURL(created)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to generate unique url.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to generate unique call back url", http.StatusBadRequest, err))
		}
		created.UniqueURL = uurl
		// Look for the right providers in the list of providers for the given VCS type.
		// If it's not found, throw an error.
		var (
			provider providers.Platform
			ok       bool
		)
		if provider, ok = r.PlatformProviders[repo.VCS]; !ok {
			err := fmt.Errorf("vcs provider with id %d is not supported", repo.VCS)
			return c.JSON(http.StatusBadRequest, kerr.APIError("unable to find vcs provider", http.StatusBadRequest, err))
		}
		if err := provider.CreateHook(ctx, repo); err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to create Hook")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to create hook", http.StatusInternalServerError, err))
		}
		return c.JSON(http.StatusCreated, created)
	}
}

// DeleteRepository handles the Delete rest event.
// TODO: Delete the hook here as well?
func (r *RepoHandler) DeleteRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := r.TokenProvider.GetToken(c)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to get Token.")
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
		if err := r.RepositoryStorer.Delete(ctx, n); err != nil {
			r.Logger.Debug().Err(err).Msg("Repository Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete repository", http.StatusBadRequest, err))
		}
		return c.NoContent(http.StatusOK)
	}
}

// GetRepository retrieves a repository and displays the unique URL for which this repo is responsible for.
func (r *RepoHandler) GetRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := r.TokenProvider.GetToken(c)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to get Token.")
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
		repo, err := r.RepositoryStorer.Get(ctx, n)
		if err != nil {
			apiError := kerr.APIError("failed to get repository", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		uurl, err := r.generateUniqueCallBackURL(repo)
		if err != nil {
			apiError := kerr.APIError("failed to generate unique callback url for repository", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		repo.UniqueURL = uurl
		return c.JSON(http.StatusOK, repo)
	}
}

// ListRepositories handles the List rest event.
func (r *RepoHandler) ListRepositories() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := r.TokenProvider.GetToken(c)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}

		opts := &models.ListOptions{}
		if err := c.Bind(opts); err != nil {
			// if we don't have anything to bind, just ignore opts.
			opts = nil
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()

		list, err := r.RepositoryStorer.List(ctx, opts)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Repository List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list repository", http.StatusBadRequest, err))
		}
		return c.JSON(http.StatusOK, list)
	}
}

// UpdateRepository handles the update rest event.
func (r *RepoHandler) UpdateRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := r.TokenProvider.GetToken(c)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusUnauthorized, kerr.APIError("failed to get token", http.StatusUnauthorized, err))
		}
		repo := &models.Repository{}
		err = c.Bind(repo)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to bind repository.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind repository", http.StatusBadRequest, err))
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		updated, err := r.RepositoryStorer.Update(ctx, repo)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Repository UpdateRepository failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to update repository", http.StatusBadRequest, err))
		}
		uurl, err := r.generateUniqueCallBackURL(updated)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Repository generateUniqueCallBackURL failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to update repository", http.StatusBadRequest, err))
		}
		updated.UniqueURL = uurl
		return c.JSON(http.StatusOK, updated)
	}
}

// generateUniqueCallBackURL takes a repository and generates a unique URL based on the ID and Type of the repo
// and the configured Krok hostname.
func (r *RepoHandler) generateUniqueCallBackURL(repo *models.Repository) (string, error) {
	u, err := url.Parse(r.Config.Hostname)
	if err != nil {
		r.Logger.Debug().Err(err).Msg("Failed to generate a unique URL for repository.")
		return "", err
	}
	u.Path = path.Join(u.Path, strconv.Itoa(repo.ID), strconv.Itoa(repo.VCS), "callback")
	return u.String(), nil
}
