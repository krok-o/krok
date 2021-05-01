package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
)

func TestVaultHandler_Create(t *testing.T) {
	logger := zerolog.New(os.Stderr)

	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)

	t.Run("create a vault secret", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("AddSecret", "key", []byte("value")).Return(nil)
		vp.On("SaveSecrets").Return(nil)
		vaultSettingPost := `{"key": "key", "value": "value"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/vault/secret", strings.NewReader(vaultSettingPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.CreateSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusCreated, rec.Code)
	})

	t.Run("create a vault secret -- load fails", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(errors.New("nope"))
		vaultSettingPost := `{"key": "key", "value": "value"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/vault/secret", strings.NewReader(vaultSettingPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.CreateSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusInternalServerError, rec.Code)
	})

	t.Run("create a vault secret -- save fails", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("AddSecret", "key", []byte("value")).Return(nil)
		vp.On("SaveSecrets").Return(errors.New("nope"))
		vaultSettingPost := `{"key": "key", "value": "value"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/vault/secret", strings.NewReader(vaultSettingPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.CreateSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusInternalServerError, rec.Code)
	})

	t.Run("create a vault secret -- invalid body", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		e := echo.New()
		body := `yaml: content`
		req := httptest.NewRequest(http.MethodPost, "/vault/secret", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.CreateSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

}

func TestVaultHandler_Update(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)

	t.Run("update a vault secret", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("GetSecret", "key").Return([]byte("value"), nil)
		vp.On("AddSecret", "key", []byte("value1")).Return(nil)
		vp.On("SaveSecrets").Return(nil)
		vaultSettingPost := `{"key": "key", "value": "value1"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/vault/secret/update", strings.NewReader(vaultSettingPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.UpdateSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})

	t.Run("update a vault secret which does not exist", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("GetSecret", "key").Return(nil, errors.New("nope"))
		vaultSettingPost := `{"key": "key", "value": "value"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/vault/secret/update", strings.NewReader(vaultSettingPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.UpdateSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusNotFound, rec.Code)
	})
}

func TestVaultHandler_Get(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)

	t.Run("get a vault secret", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("GetSecret", "key").Return([]byte("value"), nil)
		vaultSettingResponse := `{"key":"key","value":"value"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/vault/setting/:name")
		c.SetParamNames("name")
		c.SetParamValues("key")
		err = vh.GetSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		body, err := ioutil.ReadAll(rec.Result().Body)
		assert.NoError(tt, err)
		assert.Equal(tt, vaultSettingResponse, string(body))
	})
	t.Run("get a vault secret that does not exist", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("GetSecret", "key").Return(nil, errors.New("nope"))
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/vault/setting/:name")
		c.SetParamNames("name")
		c.SetParamValues("key")
		err = vh.GetSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusNotFound, rec.Code)
	})
	t.Run("get a vault secret without secret name", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/vault/secret/:name")
		err = vh.GetSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestVaultHandler_Delete(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)

	t.Run("delete a vault secret", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("GetSecret", "key").Return(nil, nil)
		vp.On("DeleteSecret", "key")
		vp.On("SaveSecrets").Return(nil)
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/vault/secret/:name")
		c.SetParamNames("name")
		c.SetParamValues("key")
		err = vh.DeleteSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})
	t.Run("delete a secret which does not exist", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("GetSecret", "key").Return(nil, errors.New("nope"))
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/vault/secret/:name")
		c.SetParamNames("name")
		c.SetParamValues("key")
		err = vh.DeleteSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusNotFound, rec.Code)
	})
	t.Run("delete a secret without secret name", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/vault/secret/:name")
		err = vh.GetSecret()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestVaultHandler_List(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)

	t.Run("list vault secrets", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(nil)
		vp.On("ListSecrets").Return([]string{"key1", "key2"}, nil)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/vault/secrets", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.ListSecrets()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		expected := `["key1","key2"]
`
		body, err := ioutil.ReadAll(rec.Body)
		assert.NoError(tt, err)
		assert.Equal(tt, expected, string(body))
	})
	t.Run("list vault secrets with error", func(tt *testing.T) {
		vp := &mocks.Vault{}
		vh := NewVaultHandler(VaultHandlerDependencies{
			Logger: logger,
			Vault:  vp,
		})
		vp.On("LoadSecrets").Return(errors.New("nope"))
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/vault/secrets", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = vh.ListSecrets()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusInternalServerError, rec.Code)
	})
}
