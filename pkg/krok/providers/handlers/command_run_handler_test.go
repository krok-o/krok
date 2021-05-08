package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
)

func TestCommandRunHandler_GetCommandRun(t *testing.T) {
	logger := zerolog.New(os.Stderr)

	t.Run("when there is a request to return a command run that exists", func(tt *testing.T) {
		mcrs := &mocks.CommandRunStorer{}
		mcrs.On("Get", mock.Anything, 1).Return(&models.CommandRun{
			ID:          1,
			EventID:     1,
			CommandName: "echo",
			Status:      "success",
			Outcome:     "this is echoed",
			CreateAt:    time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
		}, nil)
		ch := NewCommandRunHandler(CommandRunHandlerDependencies{
			Logger:           logger,
			CommandRunStorer: mcrs,
		})
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		commandRunExpected := `{"id":1,"event_id":1,"command_name":"echo","status":"success","outcome":"this is echoed","create_at":"1981-01-01T01:01:01.000000001Z"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/run/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err = ch.GetCommandRun()(c)
		assert.NoError(tt, err)
		content, err := ioutil.ReadAll(rec.Body)
		assert.NoError(tt, err)
		assert.Equal(tt, commandRunExpected, string(content))
	})

	t.Run("when there is an error getting a command run", func(tt *testing.T) {
		mcrs := &mocks.CommandRunStorer{}
		mcrs.On("Get", mock.Anything, 1).Return(nil, errors.New("nope"))
		ch := NewCommandRunHandler(CommandRunHandlerDependencies{
			Logger:           logger,
			CommandRunStorer: mcrs,
		})
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/run/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err = ch.GetCommandRun()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusInternalServerError, rec.Code)
	})
	t.Run("when the command does not exist", func(tt *testing.T) {
		mcrs := &mocks.CommandRunStorer{}
		mcrs.On("Get", mock.Anything, 1).Return(nil, kerr.ErrNotFound)
		ch := NewCommandRunHandler(CommandRunHandlerDependencies{
			Logger:           logger,
			CommandRunStorer: mcrs,
		})
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/run/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err = ch.GetCommandRun()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusNotFound, rec.Code)
	})
	t.Run("when the command id is invalid", func(tt *testing.T) {
		mcrs := &mocks.CommandRunStorer{}
		ch := NewCommandRunHandler(CommandRunHandlerDependencies{
			Logger:           logger,
			CommandRunStorer: mcrs,
		})
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/run/:id")
		c.SetParamNames("id")
		c.SetParamValues("asdf")
		err = ch.GetCommandRun()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
		mcrs.AssertNotCalled(tt, "Get", mock.Anything, 1)
	})
}
