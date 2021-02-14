package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func generateTestToken(t *testing.T) string {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   "1",
		ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
	}).SignedString([]byte("test"))
	require.NoError(t, err)
	return token
}

func TestUserAuthentication(t *testing.T) {
	e := echo.New()
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	cfg := UserMiddlewareConfig{
		GlobalTokenKey: "test",
		CookieName:     "_a_token_",
	}

	// Invalid token (tampered)
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTMzMzIxOT_INVALID_sInN1YiI6IjEifQ.5zlBJrc4hY9ENDT49OaKtWPk4WG0APj3JS2aETnEtbs"

	t.Run("valid jwt token via header", func(t *testing.T) {
		hf := NewUserMiddleware(cfg, UserMiddlewareDeps{}).JWT()(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", generateTestToken(t))
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, c.Response().Status)
		uc := c.Get("user").(*UserContext)
		assert.Equal(t, &UserContext{UserID: 1}, uc)
	})

	t.Run("valid jwt token via cookie", func(t *testing.T) {
		hf := NewUserMiddleware(cfg, UserMiddlewareDeps{}).JWT()(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Cookie", "_a_token_="+generateTestToken(t))
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, c.Response().Status)
		uc := c.Get("user").(*UserContext)
		assert.Equal(t, &UserContext{UserID: 1}, uc)
	})

	t.Run("invalid jwt token via header returns 401", func(t *testing.T) {
		hf := NewUserMiddleware(cfg, UserMiddlewareDeps{}).JWT()(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", invalidToken)
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
		assert.Nil(t, c.Get("user"))
	})

	t.Run("invalid jwt token via cookie returns 401", func(t *testing.T) {
		hf := NewUserMiddleware(cfg, UserMiddlewareDeps{}).JWT()(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Cookie", "_a_token_="+invalidToken)
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
		assert.Nil(t, c.Get("user"))
	})

	t.Run("valid api token via header", func(t *testing.T) {
		mockUserStore := &mocks.UserStorer{}
		hf := NewUserMiddleware(cfg, UserMiddlewareDeps{UserStore: mockUserStore}).JWT()(handler)

		testToken := "$2a$10$v5Gkd/DL2BpUPbEwrgXOpeMG.T/eU4e7doEY/VcGHQ5dtIn.zTn8G"
		mockUserStore.On("GetByToken", mock.Anything, testToken).Return(&models.User{ID: 1}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", "Bearer "+testToken)
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, c.Response().Status)
		uc := c.Get("user").(*UserContext)
		assert.Equal(t, &UserContext{UserID: 1}, uc)
		mockUserStore.AssertExpectations(t)
	})

	t.Run("invalid api token via header returns 401", func(t *testing.T) {
		mockUserStore := &mocks.UserStorer{}
		hf := NewUserMiddleware(cfg, UserMiddlewareDeps{UserStore: mockUserStore}).JWT()(handler)

		testToken := "$2a$10$v5Gkd/DL2BpUPbEwrgXOpeMG.T/eU4e7doEY/VcGHQ5dtIn.zTn8G"
		mockUserStore.On("GetByToken", mock.Anything, testToken).Return(nil, errors.New("err"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", "Bearer "+testToken)
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
		assert.Nil(t, c.Get("user"))
	})

	t.Run("no token", func(t *testing.T) {
		hf := NewUserMiddleware(cfg, UserMiddlewareDeps{}).JWT()(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
		assert.Nil(t, c.Get("user"))
	})
}
