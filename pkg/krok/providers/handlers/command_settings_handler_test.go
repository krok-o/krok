package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func TestCommandSettingsHandler_BasicFlow(t *testing.T) {
	cs := &mocks.CommandStorer{}
	logger := zerolog.New(os.Stderr)
	csh := NewCommandSettingsHandler(CommandSettingsHandlerDependencies{
		Logger:        logger,
		CommandStorer: cs,
	})
	t.Run("create normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		cs.On("CreateSetting", mock.Anything, &models.CommandSetting{
			CommandID: 1,
			Key:       "key",
			Value:     "value",
			InVault:   false,
		}).Return(&models.CommandSetting{
			ID:        1,
			CommandID: 1,
			Key:       "key",
			Value:     "value",
			InVault:   false,
		}, nil)

		commandSettingsPost := `{"command_id" : 1, "key" : "key", "value": "value", "in_vault": false}`
		commandSettingsReturned := `{"id":1,"command_id":1,"key":"key","value":"value","in_vault":false}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/commands/settings", strings.NewReader(commandSettingsPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = csh.Create()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusCreated, rec.Code)
		assert.Equal(tt, commandSettingsReturned, rec.Body.String())
	})
	t.Run("delete normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		cs.On("DeleteSetting", mock.Anything, 1).Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/command/settings/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err = csh.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})
	t.Run("update normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		cs.On("UpdateSetting", mock.Anything, &models.CommandSetting{
			ID:        1,
			CommandID: 1,
			Key:       "key",
			Value:     "value",
			InVault:   false,
		}).Return(nil)

		commandSettingsPost := `{"id": 1,"command_id" : 1, "key" : "key", "value": "value", "in_vault": false}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/command/settings/update", strings.NewReader(commandSettingsPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = csh.Update()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})
	t.Run("get normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		cs.On("GetSetting", mock.Anything, 1).Return(&models.CommandSetting{
			ID:        1,
			CommandID: 1,
			Key:       "key",
			Value:     "value",
			InVault:   false,
		}, nil)

		settingsExpected := `{"id":1,"command_id":1,"key":"key","value":"value","in_vault":false}
`

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/command/settings/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err = csh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, settingsExpected, rec.Body.String())
	})
	t.Run("list normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		cs.On("ListSettings", mock.Anything, 1).Return([]*models.CommandSetting{{
			CommandID: 1,
			Key:       "key",
			Value:     "value",
			InVault:   false,
		},
		}, nil)

		settingsExpected := `[{"id":0,"command_id":1,"key":"key","value":"value","in_vault":false}]
`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/command/:id/settings")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err = csh.List()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, settingsExpected, rec.Body.String())
	})
}
