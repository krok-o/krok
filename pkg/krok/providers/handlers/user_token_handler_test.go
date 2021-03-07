package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/pkg/server/middleware"
)

func TestNewUserTokenHandler(t *testing.T) {
	e := echo.New()

	t.Run("missing user in context returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := NewUserTokenHandler(UserTokenHandlerDeps{})
		err := handler.Generate()(c)
		assert.NoError(t, err)
		assert.Equal(t, "{\"code\":500,\"message\":\"Failed to get the user context.\",\"error\":\"user not found in context\"}\n", rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("user store error returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})

		mockUserStore := &mocks.UserStorer{}
		mockUserStore.On("Get", mock.Anything, 1).Return(nil, errors.New("err"))

		handler := NewUserTokenHandler(UserTokenHandlerDeps{UserStore: mockUserStore})
		err := handler.Generate()(c)
		assert.NoError(t, err)
		assert.Equal(t, "{\"code\":500,\"message\":\"Failed to get the user.\",\"error\":\"err\"}\n", rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("update error returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})

		mockUserStore := &mocks.UserStorer{}
		mockUserStore.On("Get", mock.Anything, 1).Return(&models.User{}, nil)
		mockUserStore.On("Update", mock.Anything, &models.User{Token: "1234"}).Return(nil, errors.New("err"))

		mockUTG := &mocks.UserTokenGenerator{}
		mockUTG.On("Generate").Return("1234", nil)

		handler := NewUserTokenHandler(UserTokenHandlerDeps{
			UserStore:          mockUserStore,
			UserTokenGenerator: mockUTG,
		})
		err := handler.Generate()(c)
		assert.NoError(t, err)
		assert.Equal(t, "{\"code\":500,\"message\":\"Failed to update the user.\",\"error\":\"err\"}\n", rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("successfully generate token returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})

		mockUserStore := &mocks.UserStorer{}
		mockUserStore.On("Get", mock.Anything, 1).Return(&models.User{}, nil)
		mockUserStore.On("Update", mock.Anything, &models.User{Token: "1234"}).Return(&models.User{Token: "1234"}, nil)

		mockUTG := &mocks.UserTokenGenerator{}
		mockUTG.On("Generate").Return("1234", nil)

		handler := NewUserTokenHandler(UserTokenHandlerDeps{
			UserStore:          mockUserStore,
			UserTokenGenerator: mockUTG,
		})
		err := handler.Generate()(c)
		assert.NoError(t, err)
		assert.Equal(t, "{\"token\":\"1234\"}\n", rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
