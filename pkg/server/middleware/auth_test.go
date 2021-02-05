package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	hf := UserAuthentication(&UserAuthenticationConfig{
		GlobalTokenKey: "test",
		CookieName:     "_a_token_",
	})(handler)

	t.Run("valid token via header", func(t *testing.T) {
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

	t.Run("valid token via cookie", func(t *testing.T) {
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

	t.Run("invalid token via header returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", "invalid")
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
		assert.Nil(t, c.Get("user"))
	})

	t.Run("invalid token via cookie returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("Cookie", "_a_token_=invalid")
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
		assert.Nil(t, c.Get("user"))
	})

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		c := e.NewContext(req, res)
		err := hf(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
		assert.Nil(t, c.Get("user"))
	})
}
