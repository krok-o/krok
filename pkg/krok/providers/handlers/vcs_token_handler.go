package handlers

import (
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
// swagger:operation POST /vcs-token createVcsToken
// Create a new token for a platform like Github, Gitlab, Gitea...
// ---
// consumes:
// - application/json
// parameters:
// - name: secret
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/VCSToken"
// responses:
//   '200':
//     description: 'OK setting successfully create'
//   '400':
//     description: 'invalid json payload'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to create secret'
//     schema:
//       "$ref": "#/responses/Message"
func (r *VCSTokenHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		vcsToken := &models.VCSToken{}
		if err := c.Bind(vcsToken); err != nil {
			r.Logger.Debug().Err(err).Msg("Failed to bind vcs token.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind vcs token", http.StatusBadRequest, err))
		}

		if err := r.TokenProvider.SaveTokenForPlatform(vcsToken.Token, vcsToken.VCS); err != nil {
			r.Logger.Debug().Err(err).Msg("VCS token creation failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("VCS token creation failed", http.StatusInternalServerError, err))
		}

		return c.NoContent(http.StatusCreated)
	}
}
