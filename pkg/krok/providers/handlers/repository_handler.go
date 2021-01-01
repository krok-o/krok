package handlers

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// RepoHandlerDependencies defines the dependencies for the repository handler provider.
type RepoHandlerDependencies struct {
	Dependencies
	RepositoryStorer providers.RepositoryStorer
	TokenProvider    *TokenProvider
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
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to get token", http.StatusBadRequest, err))
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
		return c.JSON(http.StatusCreated, created)
	}
}

// DeleteRepository handles the Delete rest event.
func (r *RepoHandler) DeleteRepository() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := r.TokenProvider.GetToken(c)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to get token", http.StatusBadRequest, err))
		}
		id := c.Param("id")
		if id == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
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
			return err
		}
		id := c.Param("id")
		if id == "" {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		n, err := strconv.Atoi(id)
		if err != nil {
			apiError := kerr.APIError("failed to convert id to number", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		repo, err := r.RepositoryStorer.Get(ctx, n)
		if err != nil {
			return err
		}
		uurl, err := r.generateUniqueCallBackURL(repo)
		if err != nil {
			return err
		}
		repo.UniqueURL = uurl
		var result = struct {
			Repository models.Repository `json:"repository"`
		}{
			Repository: *repo,
		}
		return c.JSON(http.StatusOK, result)
	}
}

// ListRepositories handles the List rest event.
func (r *RepoHandler) ListRepositories() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := r.TokenProvider.GetToken(c)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to get Token.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to get token", http.StatusBadRequest, err))
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
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to get token", http.StatusBadRequest, err))
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
