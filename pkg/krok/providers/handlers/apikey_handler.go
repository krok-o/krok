package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	krokmiddleware "github.com/krok-o/krok/pkg/server/middleware"
)

// APIKeysHandlerDependencies defines the dependencies for the api keys handler provider.
type APIKeysHandlerDependencies struct {
	Dependencies
	APIKeysStore  providers.APIKeysStorer
	TokenProvider *TokenHandler
	Clock         providers.Clock
}

// APIKeysHandler is a handler taking care of api keys related api calls.
type APIKeysHandler struct {
	APIKeysHandlerDependencies
}

var _ providers.APIKeysHandler = &APIKeysHandler{}

// NewAPIKeysHandler creates a new api key pair handler.
func NewAPIKeysHandler(deps APIKeysHandlerDependencies) *APIKeysHandler {
	return &APIKeysHandler{
		APIKeysHandlerDependencies: deps,
	}
}

// Create creates an api key pair for a given user.
// swagger:operation POST /user/apikey/generate/{name} createApiKey
// Creates an api key pair for a given user.
// ---
// produces:
// - application/json
// parameters:
// - name: name
//   in: path
//   required: true
//   description: "the name of the key"
//   type: string
// responses:
//   '200':
//     description: 'the generated api key pair'
//     schema:
//       "$ref": "#/definitions/APIKey"
//   '400':
//     description: 'failed to generate unique key or value'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'when failed to get user context'
//     schema:
//       "$ref": "#/responses/Message"
func (a *APIKeysHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		uc, err := krokmiddleware.GetUserContext(c)
		if err != nil {
			apiError := kerr.APIError("failed to get user context", http.StatusInternalServerError, nil)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		name := c.Param("name")
		if name == "" {
			name = "My API Key"
		}
		ctx := c.Request().Context()
		key, err := a.APIKeyAuth.Generate(ctx, name, uc.UserID)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("APIKey Create failed.")
			return c.JSON(http.StatusInternalServerError, kerr.APIError("failed to generate new unique key", http.StatusInternalServerError, err))
		}
		return c.JSON(http.StatusOK, key)
	}
}

// Delete deletes a set of api keys for a given user with a given id.
// swagger:operation DELETE /user/apikey/delete/{keyid} deleteApiKey
// Deletes a set of api keys for a given user with a given id.
// ---
// parameters:
// - name: keyid
//   in: path
//   description: 'The ID of the key to delete'
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     description: 'OK in case the deletion was successful'
//   '400':
//     description: 'in case of missing user context or invalid ID'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'when the deletion operation failed'
//     schema:
//       "$ref": "#/responses/Message"
func (a *APIKeysHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		uc, err := krokmiddleware.GetUserContext(c)
		if err != nil {
			apiError := kerr.APIError("failed to get user context", http.StatusInternalServerError, nil)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		kn, err := GetParamAsInt("keyid", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()
		if err := a.APIKeysStore.Delete(ctx, kn, uc.UserID); err != nil {
			a.Logger.Debug().Err(err).Msg("APIKey Delete failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to delete api key", http.StatusBadRequest, err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// List lists all api keys for a given user.
// swagger:operation POST /user/apikey listApiKeys
// Lists all api keys for a given user.
// ---
// produces:
// - application/json
// responses:
//   '200':
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/APIKey"
//   '500':
//     description: 'failed to get user context'
//     schema:
//       "$ref": "#/responses/Message"
func (a *APIKeysHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		uc, err := krokmiddleware.GetUserContext(c)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("error getting user context")
			return c.String(http.StatusInternalServerError, "failed to get user context")
		}

		ctx := c.Request().Context()
		list, err := a.APIKeysStore.List(ctx, uc.UserID)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("APIKeys List failed.")
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to list api keys", http.StatusBadRequest, err))
		}

		return c.JSON(http.StatusOK, list)
	}
}

// Get returns a given api key.
// swagger:operation GET /user/apikey/{keyid} getApiKeys
// Returns a given api key.
// ---
// produces:
// - application/json
// parameters:
// - name: keyid
//   in: path
//   description: "The ID of the key to return"
//   required: true
//   type: integer
//   format: int
// responses:
//   '200':
//     schema:
//       "$ref": "#/definitions/APIKey"
//   '500':
//     description: 'failed to get user context'
//     schema:
//       "$ref": "#/responses/Message"
func (a *APIKeysHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		uc, err := krokmiddleware.GetUserContext(c)
		if err != nil {
			a.Logger.Debug().Err(err).Msg("error getting user context")
			return c.String(http.StatusInternalServerError, "failed to get user context")
		}

		kn, err := GetParamAsInt("keyid", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		ctx := c.Request().Context()
		key, err := a.APIKeysStore.Get(ctx, kn, uc.UserID)
		if err != nil {
			apiError := kerr.APIError("failed to get api key", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		return c.JSON(http.StatusOK, key)
	}
}
