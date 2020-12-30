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
	TokenProvider    TokenProvider
}

// RepoHandler is a handler taking care of repository related api calls.
type RepoHandler struct {
	Config
	RepoHandlerDependencies
}

var _ providers.RepositoriesHandler = &RepoHandler{}

// NewRepositoryHandler creates a new repository handler.
func NewRepositoryHandler(cfg Config, deps RepoHandlerDependencies) (*RepoHandler, error) {
	return &RepoHandler{
		Config:                  cfg,
		RepoHandlerDependencies: deps,
	}, nil
}

// Create handles the Create rest event.
func (r *RepoHandler) Create() echo.HandlerFunc {
	panic("implement me")
}

// Delete handles the Delete rest event.
func (r *RepoHandler) Delete() echo.HandlerFunc {
	panic("implement me")
}

// Get retrieves a repository and displays the unique URL for which this repo is responsible for.
func (r *RepoHandler) Get() echo.HandlerFunc {
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

// List handles the List rest event.
func (r *RepoHandler) List() echo.HandlerFunc {
	panic("implement me")
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
