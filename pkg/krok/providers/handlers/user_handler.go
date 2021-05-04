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
	Logger    zerolog.Logger
	UserStore providers.UserStorer
}

// NewUserHandler .
func NewUserHandler(deps UserHandlerDependencies) *UserHandler {
	return &UserHandler{
		UserHandlerDependencies: deps,
	}
}

var _ providers.UserHandler = &UserHandler{}

// GetUser returns a user.
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

// ListUsers .
func (u *UserHandler) ListUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		// Get the users from store.
		users, err := u.UserStore.List(ctx)
		if err != nil {
			apiError := kerr.APIError("failed to list users", http.StatusBadRequest, err)
			return c.JSON(http.StatusBadRequest, apiError)
		}

		return c.JSON(http.StatusOK, users)
	}
}

// DeleteUser .
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

// UpdateUser .
func (u *UserHandler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var update *models.User
		if err := c.Bind(&update); err != nil {
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind user", http.StatusBadRequest, err))
		}
		// we have this on `-` but let's make sure it's not set.
		update.Token = nil
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

// CreateUser .
func (u *UserHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var create *models.User
		if err := c.Bind(&create); err != nil {
			return c.JSON(http.StatusBadRequest, kerr.APIError("failed to bind user", http.StatusBadRequest, err))
		}
		result, err := u.UserStore.Create(c.Request().Context(), create)
		if err != nil {
			apiError := kerr.APIError("failed to create user", http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, apiError)
		}

		return c.JSON(http.StatusCreated, result)
	}
}
