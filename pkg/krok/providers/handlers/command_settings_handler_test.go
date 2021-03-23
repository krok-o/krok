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

func TestCommandSettingsHandler_Create(t *testing.T) {
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
		}).Return(nil)

		commandSettingsPost := `{"command_id" : 1, "key" : "key", "value": "value", "in_vault": false}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/commands/settings", strings.NewReader(commandSettingsPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = csh.Create()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusCreated, rec.Code)
	})
}
