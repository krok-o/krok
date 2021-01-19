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

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type mockApiKeyAuth struct {
	providers.ApiKeysAuthenticator
}

func (maka *mockApiKeyAuth) Match(ctx context.Context, key *models.APIKey) error {
	return nil
}

func (maka *mockApiKeyAuth) Encrypt(ctx context.Context, secret []byte) ([]byte, error) {
	return nil, nil
}

func generateTestToken(email string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email // from context
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, nil
}

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

func (m *mockApiKeysStore) Get(ctx context.Context, id int, userID int) (*models.APIKey, error) {
	return m.getKey, nil
}

func (m *mockApiKeysStore) List(ctx context.Context, userID int) ([]*models.APIKey, error) {
	return m.keyList, nil
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

func TestApiKeysHandler_GetApiKeyPair(t *testing.T) {
	maka := &mockApiKeyAuth{}
	mus := &mockUserStorer{}
	aks := &mockApiKeysStore{
		getKey: &models.APIKey{
			ID:           0,
			Name:         "test-key",
			UserID:       0,
			APIKeyID:     "api-key-id",
			APIKeySecret: []byte("secret"),
			TTL:          time.Now().Add(10 * time.Minute),
		},
	}
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

	t.Run("get apikey happy path", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		ekey := &models.APIKey{
			ID:           0,
			Name:         "test-key",
			UserID:       0,
			APIKeyID:     "api-key-id",
			APIKeySecret: []byte("secret"),
		}
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		c.SetParamNames("uid", "keyid")
		c.SetParamValues("0", "0")
		err = akh.GetApiKeyPair()(c)
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

	t.Run("no token", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		c.SetParamNames("uid", "keyid")
		c.SetParamValues("0", "0")
		err = akh.GetApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusUnauthorized, rec.Code)
	})
	t.Run("get invalid id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		c.SetParamNames("uid", "keyid")
		c.SetParamValues("invalid", "0")
		err = akh.GetApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
	t.Run("get invalid user id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		c.SetParamNames("uid", "keyid")
		c.SetParamValues("0", "invalid")
		err = akh.GetApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("get empty id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikey/:keyid")
		err = akh.GetApiKeyPair()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestApiKeysHandler_ListApiKeyPairs(t *testing.T) {
	maka := &mockApiKeyAuth{}
	mus := &mockUserStorer{}
	aks := &mockApiKeysStore{
		keyList: []*models.APIKey{
			{
				ID:           0,
				Name:         "test-key-1",
				UserID:       0,
				APIKeyID:     "test-key-id-1",
				APIKeySecret: []byte("secret1"),
				TTL:          time.Now().Add(10 * time.Minute),
			},
			{
				ID:           1,
				Name:         "test-key-2",
				UserID:       1,
				APIKeyID:     "test-key-id-2",
				APIKeySecret: []byte("secret2"),
				TTL:          time.Now().Add(10 * time.Minute),
			},
		},
	}
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

	t.Run("list apikey happy path", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/user/:uid/apikeys")
		c.SetParamNames("uid")
		c.SetParamValues("0")
		err = akh.ListApiKeyPairs()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		var gotKey []*models.APIKey
		err = json.Unmarshal(rec.Body.Bytes(), &gotKey)
		assert.NoError(tt, err)
		assert.Len(tt, gotKey, 2)
	})

}
