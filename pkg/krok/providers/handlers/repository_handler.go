package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	api = "/rest/api/1"
)

// RepoConfig represents configuration entities that the repository requires.
type RepoConfig struct {
	Protocol string
	HookBase string
}

// RepoHandlerDependencies defines the dependencies for the repository handler provider.
type RepoHandlerDependencies struct {
	Auth              providers.RepositoryAuth
	RepositoryStorer  providers.RepositoryStorer
	Logger            zerolog.Logger
	PlatformProviders map[int]providers.Platform
}

// RepoHandler is a handler taking care of repository related api calls.
type RepoHandler struct {
	RepoConfig
	RepoHandlerDependencies
}

var _ providers.RepositoryHandler = &RepoHandler{}

// NewRepositoryHandler creates a new repository handler.
func NewRepositoryHandler(cfg RepoConfig, deps RepoHandlerDependencies) (*RepoHandler, error) {
	return &RepoHandler{
		RepoConfig:              cfg,
		RepoHandlerDependencies: deps,
	}, nil
}

// Create handles the Create rest event.
func (r *RepoHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		repo := &models.Repository{}
		if err := c.Bind(repo); err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to bind repository.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind repository", http.StatusBadRequest, err))
		}

		if ok, field, err := repo.Validate(); !ok {
			r.Logger.Debug().Err(err).Str("field", field).Msg("Repository validation failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("repository validation failed", http.StatusBadRequest, err))
		}

		ctx := c.Request().Context()
		created, err := r.RepositoryStorer.Create(ctx, repo)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Repository CreateRepository failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to create repository", http.StatusBadRequest, err))
		}

		// Once the creation succeeded, create the auth values
		if err := r.Auth.CreateRepositoryAuth(ctx, created.ID, repo.Auth); err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to store auth information.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to create repository auth information", http.StatusBadRequest, err))
		}
		created.Auth = repo.Auth

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
		if provider, ok = r.PlatformProviders[created.VCS]; !ok {
			err := fmt.Errorf("vcs provider with id %d is not supported", created.VCS)
			return c.JSON(http.StatusBadRequest, kerr.APIError("unable to find vcs provider", http.StatusBadRequest, err))
		}
		if err := provider.CreateHook(ctx, created); err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("token does not exist for platform, please create first.", http.StatusNotFound, err))
			}
			r.Logger.Debug().Err(err).Msg("Failed to create Hook")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to create hook", http.StatusInternalServerError, err))
		}
		return c.JSON(http.StatusCreated, created)
	}
}

// Delete handles the Delete rest event.
// TODO: Delete the hook here as well?
func (r *RepoHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		if err := r.RepositoryStorer.Delete(ctx, n); err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("repository not found", http.StatusNotFound, err))
			}
			r.Logger.Debug().Err(err).Msg("Repository Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete repository", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// Get retrieves a repository and displays the unique URL for which this repo is responsible for.
func (r *RepoHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		// Get the repo from store.
		repo, err := r.RepositoryStorer.Get(ctx, n)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("repository not found", http.StatusNotFound, err))
			}
			apiError := kerr.APIError("failed to get repository", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		// Get the auth information for the repository
		auth, err := r.Auth.GetRepositoryAuth(ctx, repo.ID)
		if err != nil {
			apiError := kerr.APIError("failed to get repository auth information", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		repo.Auth = auth

		uurl, err := r.generateUniqueCallBackURL(repo)
		if err != nil {
			apiError := kerr.APIError("failed to generate unique callback url for repository", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		repo.UniqueURL = uurl
		return c.JSON(http.StatusOK, repo)
	}
}

// List handles the List rest event.
func (r *RepoHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		opts := &models.ListOptions{}
		if err := c.Bind(opts); err != nil {
			// if we don't have anything to bind, just ignore opts.
			opts = nil
		}

		ctx := c.Request().Context()

		list, err := r.RepositoryStorer.List(ctx, opts)
		if err != nil {
			r.Logger.Debug().Err(err).Msg("Repository List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list repository", http.StatusBadRequest, err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// Update handles the update rest event.
func (r *RepoHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		repo := &models.Repository{}
		if err := c.Bind(repo); err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to bind repository.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind repository", http.StatusBadRequest, err))
		}

		ctx := c.Request().Context()

		updated, err := r.RepositoryStorer.Update(ctx, repo)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("repository not found", http.StatusNotFound, err))
			}
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
	u, err := url.Parse(fmt.Sprintf("%s://%s", r.RepoConfig.Protocol, r.RepoConfig.HookBase))
	if err != nil {
		r.Logger.Debug().Err(err).Msg("Failed to generate a unique URL for repository.")
		return "", err
	}
	u.Path = path.Join(u.Path, api, "hooks", strconv.Itoa(repo.ID), strconv.Itoa(repo.VCS), "callback")
	return u.String(), nil
}
