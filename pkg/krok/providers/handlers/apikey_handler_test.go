package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type mockApiKeysStore struct {
	providers.APIKeysStorer
	getKey    *models.APIKey
	deleteErr error
	keyList   []*models.APIKey
	id        int
}

func (m *mockApiKeysStore) Create(ctx context.Context, key *models.APIKey) (*models.APIKey, error) {
	key.ID = m.id
	m.id++
	return key, nil
}

func (m *mockApiKeysStore) Delete(ctx context.Context, id int, uid int) error {
	return m.deleteErr
}

func TestApiKeysHandler_CreateApiKeyPair(t *testing.T) {
	maka := &mockApiKeyAuth{}
	mus := &mockUserStorer{}
	aks := &mockApiKeysStore{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		ApiKeyAuth: maka,
	}
	cfg := Config{
		Hostname:       "https://testHost",
		GlobalTokenKey: "secret",
	}
	tp, err := NewTokenProvider(cfg, deps)
	assert.NoError(t, err)
	akh, err := NewApiKeysHandler(cfg, ApiKeysHandlerDependencies{
		Dependencies:  deps,
		APIKeysStore:  aks,
		TokenProvider: tp,
	})
	assert.NoError(t, err)

	t.Run("create happy path", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/generate/:name")
		c.SetParamNames("uid", "name")
		c.SetParamValues("0", "test-key")
		err = akh.CreateApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		fmt.Println(rec.Body.String())
		var key models.APIKey
		err = json.Unmarshal(rec.Body.Bytes(), &key)
		assert.NoError(tt, err)
		assert.NotEmpty(tt, key.APIKeySecret)
		assert.NotEmpty(tt, key.APIKeyID)
		assert.True(tt, key.TTL.After(time.Now()))
		assert.Equal(tt, "test-key", key.Name)
	})
	t.Run("create happy path without name", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/generate")
		c.SetParamNames("uid")
		c.SetParamValues("0")
		err = akh.CreateApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		fmt.Println(rec.Body.String())
		var key models.APIKey
		err = json.Unmarshal(rec.Body.Bytes(), &key)
		assert.NoError(tt, err)
		assert.NotEmpty(tt, key.APIKeySecret)
		assert.NotEmpty(tt, key.APIKeyID)
		assert.True(tt, key.TTL.After(time.Now()))
		assert.Equal(tt, "My Api Key", key.Name)
	})
	t.Run("create no token", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/generate/:name")
		c.SetParamNames("uid")
		c.SetParamValues("0")
		err = akh.CreateApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusUnauthorized, rec.Code)
	})
	t.Run("create invalid user ID", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/generate/:name")
		c.SetParamNames("uid")
		c.SetParamValues("invalid")
		err = akh.CreateApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestApiKeysHandler_DeleteApiKeyPair(t *testing.T) {
	maka := &mockApiKeyAuth{}
	mus := &mockUserStorer{}
	aks := &mockApiKeysStore{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		ApiKeyAuth: maka,
	}
	cfg := Config{
		Hostname:       "https://testHost",
		GlobalTokenKey: "secret",
	}
	tp, err := NewTokenProvider(cfg, deps)
	assert.NoError(t, err)
	akh, err := NewApiKeysHandler(cfg, ApiKeysHandlerDependencies{
		Dependencies:  deps,
		APIKeysStore:  aks,
		TokenProvider: tp,
	})
	assert.NoError(t, err)

	t.Run("delete happy path", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/delete/:keyid")
		c.SetParamNames("uid", "keyid")
		c.SetParamValues("0", "0")
		err = akh.DeleteApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})

	t.Run("delete no token", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		c.SetParamNames("uid", "keyid")
		c.SetParamValues("0", "0")
		err = akh.DeleteApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusUnauthorized, rec.Code)
	})

	t.Run("delete invalid id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		c.SetParamNames("uid", "keyid")
		c.SetParamValues("invalid", "0")
		err = akh.DeleteApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("delete empty id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		err = akh.DeleteApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}
