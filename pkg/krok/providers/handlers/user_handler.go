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

// UserHandler .
type UserHandler struct {
	UserHandlerDependencies
}

// UserHandlerDependencies .
type UserHandlerDependencies struct {
	Logger     zerolog.Logger
	UserStore  providers.UserStorer
	APIKeyAuth providers.APIKeysAuthenticator
}

// NewUserHandler .
func NewUserHandler(deps UserHandlerDependencies) *UserHandler {
	return &UserHandler{
		UserHandlerDependencies: deps,
	}
}

var _ providers.UserHandler = &UserHandler{}

// GetUser returns a user.
// swagger:operation GET /user/{id} getUser
// Gets the user with the corresponding ID.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   type: integer
//   format: int
//   required: true
// responses:
//   '200':
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: 'invalid user id'
//     schema:
//       "$ref": "#/responses/Message"
//   '404':
//     description: 'user not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to get user'
//     schema:
//       "$ref": "#/responses/Message"
func (u *UserHandler) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		// Get the user from store.
		user, err := u.UserStore.Get(ctx, n)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				return c.JSON(http.StatusNotFound, kerr.APIError("user not found", http.StatusNotFound, err))
			}
			apiError := kerr.APIError("failed to get user", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}

		return c.JSON(http.StatusOK, user)
	}
}

// ListUsers lists all users.
// swagger:operation POST /users listUsers
// List users
// ---
// produces:
// - application/json
// responses:
//   '200':
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/User"
//   '500':
//     description: 'failed to list user'
//     schema:
//       "$ref": "#/responses/Message"
func (u *UserHandler) ListUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		// Get the users from store.
		users, err := u.UserStore.List(ctx)
		if err != nil {
			apiError := kerr.APIError("failed to list users", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}

		return c.JSON(http.StatusOK, users)
	}
}

// DeleteUser deletes a user.
// swagger:operation DELETE /user/{id} deleteUser
// Deletes the given user.
// ---
// parameters:
// - name: id
//   in: path
//   description: 'The ID of the user to delete'
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
//   '404':
//     description: 'in case of user not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'when the deletion operation failed'
//     schema:
//       "$ref": "#/responses/Message"
func (u *UserHandler) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		n, err := GetParamAsInt("id", c)
		if err != nil {
			apiError := kerr.APIError("invalid id", http.StatusBadRequest, nil)
			return c.JSON(http.StatusBadRequest, apiError)
		}
		ctx := c.Request().Context()

		if err := u.UserStore.Delete(ctx, n); err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				apiError := kerr.APIError("user not found", http.StatusNotFound, err)
				return c.JSON(http.StatusNotFound, apiError)
			}
			apiError := kerr.APIError("failed to delete users", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}

		return c.NoContent(http.StatusOK)
	}
}

// UpdateUser update a specific user.
// swagger:operation POST /user/update updateUser
// Updates an existing user.
// ---
// produces:
// - application/json
// consumes:
// - application/json
// parameters:
// - name: user
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// responses:
//   '200':
//     description: 'user successfully updated'
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: 'invalid json payload'
//     schema:
//       "$ref": "#/responses/Message"
//   '404':
//     description: 'user not found'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to update user'
//     schema:
//       "$ref": "#/responses/Message"
func (u *UserHandler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var update *models.User
		if err := c.Bind(&update); err != nil {
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind user", http.StatusBadRequest, err))
		}
		result, err := u.UserStore.Update(c.Request().Context(), update)
		if err != nil {
			if errors.Is(err, kerr.ErrNotFound) {
				apiError := kerr.APIError("user not found", http.StatusNotFound, err)
				return c.JSON(http.StatusNotFound, apiError)
			}
			apiError := kerr.APIError("failed to update user", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}

		return c.JSON(http.StatusOK, result)
	}
}

// CreateUser creates a new user.
// swagger:operation POST /user createUser
// Creates a new user
// ---
// produces:
// - application/json
// consumes:
// - application/json
// parameters:
// - name: user
//   in: body
//   required: true
//   schema:
//     "$ref": "#/definitions/User"
// responses:
//   '200':
//     description: 'the created user'
//     schema:
//       "$ref": "#/definitions/User"
//   '400':
//     description: 'invalid json payload'
//     schema:
//       "$ref": "#/responses/Message"
//   '500':
//     description: 'failed to create user or generating a new api key'
//     schema:
//       "$ref": "#/responses/Message"
func (u *UserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var create *models.User
		if err := c.Bind(&create); err != nil {
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind user", http.StatusBadRequest, err))
		}
		// Create the user
		result, err := u.UserStore.Create(c.Request().Context(), create)
		if err != nil {
			apiError := kerr.APIError("failed to create user", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}

		// Create initial API key
		key, err := u.APIKeyAuth.Generate(c.Request().Context(), "New API Key", result.ID)
		if err != nil {
			apiError := kerr.APIError("failed to create new api keys for user", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}
		result.APIKeys = append(result.APIKeys, key)
		return c.JSON(http.StatusCreated, result)
	}
}
