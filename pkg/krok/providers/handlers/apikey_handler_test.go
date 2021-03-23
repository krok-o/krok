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
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/pkg/server/middleware"
)

type mockAPIKeysStore struct {
	providers.APIKeysStorer
	getKey    *models.APIKey
	deleteErr error
	keyList   []*models.APIKey
	id        int
}

func (m *mockAPIKeysStore) Create(ctx context.Context, key *models.APIKey) (*models.APIKey, error) {
	key.ID = m.id
	m.id++
	return key, nil
}

func (m *mockAPIKeysStore) Delete(ctx context.Context, id int, uid int) error {
	return m.deleteErr
}

func (m *mockAPIKeysStore) Get(ctx context.Context, id int, userID int) (*models.APIKey, error) {
	return m.getKey, nil
}

func (m *mockAPIKeysStore) List(ctx context.Context, userID int) ([]*models.APIKey, error) {
	return m.keyList, nil
}

type mockAPIKeyAuth struct {
	providers.APIKeysAuthenticator
}

func (maka *mockAPIKeyAuth) Match(ctx context.Context, key *models.APIKey) error {
	return nil
}

func (maka *mockAPIKeyAuth) Encrypt(ctx context.Context, secret []byte) ([]byte, error) {
	return nil, nil
}

func TestAPIKeysHandler_CreateAPIKeyPair(t *testing.T) {
	maka := &mocks.APIKeysAuthenticator{}
	mus := &mockUserStorer{}
	aks := &mockAPIKeysStore{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	akh := NewAPIKeysHandler(APIKeysHandlerDependencies{
		Dependencies:  deps,
		APIKeysStore:  aks,
		TokenProvider: tp,
	})

	t.Run("create happy path", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/generate/:name")
		c.SetParamNames("name")
		c.SetParamValues("test-key")
		err = akh.Create()(c)
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
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/generate")
		err = akh.Create()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		fmt.Println(rec.Body.String())
		var key models.APIKey
		err = json.Unmarshal(rec.Body.Bytes(), &key)
		assert.NoError(tt, err)
		assert.NotEmpty(tt, key.APIKeySecret)
		assert.NotEmpty(tt, key.APIKeyID)
		assert.True(tt, key.TTL.After(time.Now()))
		assert.Equal(tt, "My API Key", key.Name)
	})
	t.Run("create no user context", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/apikey/generate/:name")
		err = akh.Create()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusInternalServerError, rec.Code)
	})
}

func TestAPIKeysHandler_DeleteAPIKeyPair(t *testing.T) {
	maka := &mockAPIKeyAuth{}
	mus := &mockUserStorer{}
	aks := &mockAPIKeysStore{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	akh := NewAPIKeysHandler(APIKeysHandlerDependencies{
		Dependencies:  deps,
		APIKeysStore:  aks,
		TokenProvider: tp,
	})

	t.Run("delete happy path", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/delete/:keyid")
		c.SetParamNames("keyid")
		c.SetParamValues("0")
		err = akh.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})

	t.Run("delete no user context", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/apikey/:keyid")
		c.SetParamNames("keyid")
		c.SetParamValues("0")
		err = akh.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusInternalServerError, rec.Code)
	})

	t.Run("delete invalid id", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/:keyid")
		err = akh.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("delete empty id", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/:keyid")
		err = akh.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestAPIKeysHandler_GetAPIKeyPair(t *testing.T) {
	maka := &mockAPIKeyAuth{}
	mus := &mockUserStorer{}
	aks := &mockAPIKeysStore{
		getKey: &models.APIKey{
			ID:           0,
			Name:         "test-key",
			UserID:       0,
			APIKeyID:     "api-key-id",
			APIKeySecret: "secret",
			TTL:          time.Now().Add(10 * time.Minute),
		},
	}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	akh := NewAPIKeysHandler(APIKeysHandlerDependencies{
		Dependencies:  deps,
		APIKeysStore:  aks,
		TokenProvider: tp,
	})

	t.Run("get apikey happy path", func(tt *testing.T) {
		ekey := &models.APIKey{
			ID:           0,
			Name:         "test-key",
			UserID:       0,
			APIKeyID:     "api-key-id",
			APIKeySecret: "secret",
		}
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/:keyid")
		c.SetParamNames("keyid")
		c.SetParamValues("0")
		err = akh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		var gotKey models.APIKey
		err = json.Unmarshal(rec.Body.Bytes(), &gotKey)
		assert.NoError(tt, err)
		assert.Equal(tt, ekey.UserID, gotKey.UserID)
		assert.Equal(tt, ekey.APIKeyID, gotKey.APIKeyID)
		assert.Equal(tt, ekey.APIKeySecret, gotKey.APIKeySecret)
		assert.Equal(tt, ekey.Name, gotKey.Name)
		assert.Equal(tt, ekey.ID, gotKey.ID)
	})

	t.Run("no user context", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/apikey/:keyid")
		err = akh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusInternalServerError, rec.Code)
	})

	t.Run("get invalid id", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/:keyid")
		err = akh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
	t.Run("get invalid user id", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/:keyid")
		err = akh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("get empty id", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey/:keyid")
		err = akh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestAPIKeysHandler_ListAPIKeyPairs(t *testing.T) {
	maka := &mockAPIKeyAuth{}
	mus := &mockUserStorer{}
	aks := &mockAPIKeysStore{
		keyList: []*models.APIKey{
			{
				ID:           0,
				Name:         "test-key-1",
				UserID:       0,
				APIKeyID:     "test-key-id-1",
				APIKeySecret: "secret1",
				TTL:          time.Now().Add(10 * time.Minute),
			},
			{
				ID:           1,
				Name:         "test-key-2",
				UserID:       1,
				APIKeyID:     "test-key-id-2",
				APIKeySecret: "secret2",
				TTL:          time.Now().Add(10 * time.Minute),
			},
		},
	}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	akh := NewAPIKeysHandler(APIKeysHandlerDependencies{
		Dependencies:  deps,
		APIKeysStore:  aks,
		TokenProvider: tp,
	})

	t.Run("list apikey happy path", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", &middleware.UserContext{UserID: 1})
		c.SetPath("/user/apikey")
		err = akh.List()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		var gotKey []*models.APIKey
		err = json.Unmarshal(rec.Body.Bytes(), &gotKey)
		assert.NoError(tt, err)
		assert.Len(tt, gotKey, 2)
	})
}
