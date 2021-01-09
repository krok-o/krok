package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type mockCommandStorer struct {
	providers.CommandStorer
	getCommand  *models.Command
	deleteErr   error
	commandList []*models.Command
}

func (mcs *mockCommandStorer) Update(ctx context.Context, command *models.Command) (*models.Command, error) {
	return command, nil
}

func (mcs *mockCommandStorer) Get(ctx context.Context, id int) (*models.Command, error) {
	return mcs.getCommand, nil
}

func (mcs *mockCommandStorer) List(ctx context.Context, opts *models.ListOptions) ([]*models.Command, error) {
	return mcs.commandList, nil
}

func (mcs *mockCommandStorer) Delete(ctx context.Context, id int) error {
	return mcs.deleteErr
}

func (mcs *mockCommandStorer) AddCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error {
	return nil
}

func (mcs *mockCommandStorer) RemoveCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error {
	return nil
}

func TestCommandsHandler_DeleteCommand(t *testing.T) {
	mus := &mockUserStorer{}
	mcs := &mockCommandStorer{}
	maka := &mockApiKeyAuth{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		ApiKeyAuth: maka,
	}
	cfg := Config{
		Hostname:       "http://testHost",
		GlobalTokenKey: "secret",
	}
	tp, err := NewTokenProvider(cfg, deps)
	assert.NoError(t, err)
	ch, err := NewCommandsHandler(cfg, CommandsHandlerDependencies{
		Dependencies:  deps,
		CommandStorer: mcs,
		TokenProvider: tp,
	})
	assert.NoError(t, err)

	t.Run("delete normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = ch.DeleteCommand()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})

	t.Run("delete no token", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/command/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = ch.DeleteCommand()(c)
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
		c.SetPath("/command/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")
		err = ch.DeleteCommand()(c)
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
		c.SetPath("/command/:id")
		err = ch.DeleteCommand()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}
