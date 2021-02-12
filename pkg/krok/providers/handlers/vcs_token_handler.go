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

// VCSTokenHandlerDependencies defines the dependencies for the vcs token handler provider.
type VCSTokenHandlerDependencies struct {
	Logger        zerolog.Logger
	TokenProvider providers.PlatformTokenProvider
}

// VCSTokenHandler is a handler taking care of vcs token related api calls.
type VCSTokenHandler struct {
	VCSTokenHandlerDependencies
}

var _ providers.VCSTokenHandler = &VCSTokenHandler{}

// NewVCSTokenHandler creates a new vcs token handler.
func NewVCSTokenHandler(deps VCSTokenHandlerDependencies) *VCSTokenHandler {
	return &VCSTokenHandler{
		VCSTokenHandlerDependencies: deps,
	}
}

// Create handles the Create rest event.
func (r *VCSTokenHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		vcsToken := &models.VCSToken{}
		if err := c.Bind(vcsToken); err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to bind vcs token.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind vcs token", http.StatusBadRequest, err))
		}

		if err := r.TokenProvider.SaveTokenForPlatform(vcsToken.Token, vcsToken.VCS); err != nil {
			r.Logger.Debug().Err(err).Msg("VCS token creation failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("VCS token creation failed", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusCreated)
	}
}

// Delete handles the Delete rest event.
func (r *VCSTokenHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, kerr.APIError("unimplemented", http.StatusInternalServerError, errors.New("unimplemented")))
	}
}

// Get retrieves a repository and displays the unique URL for which this repo is responsible for.
func (r *VCSTokenHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, kerr.APIError("unimplemented", http.StatusInternalServerError, errors.New("unimplemented")))
	}
}

// List handles the List rest event.
func (r *VCSTokenHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, kerr.APIError("unimplemented", http.StatusInternalServerError, errors.New("unimplemented")))
	}
}

// Update handles the update rest event.
func (r *VCSTokenHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, kerr.APIError("unimplemented", http.StatusInternalServerError, errors.New("unimplemented")))
	}
}
