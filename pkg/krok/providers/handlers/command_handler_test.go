package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	if opts == nil {
		return mcs.commandList, nil
	}
	result := make([]*models.Command, 0)
	for _, c := range mcs.commandList {
		if opts.Name != "" {
			if strings.Contains(c.Name, opts.Name) {
				result = append(result, c)
			}
		} else {
			result = append(result, c)
		}
	}
	return result, nil
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
	maka := &mockAPIKeyAuth{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	ch := NewCommandsHandler(CommandsHandlerDependencies{
		Logger:        logger,
		CommandStorer: mcs,
		TokenProvider: tp,
	})

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
		err = ch.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
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
		err = ch.Delete()(c)
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
		err = ch.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestCommandsHandler_GetCommand(t *testing.T) {
	mus := &mockUserStorer{}
	mcs := &mockCommandStorer{
		getCommand: &models.Command{
			Name:     "test-command",
			ID:       0,
			Schedule: "* * * * *",
			Repositories: []*models.Repository{
				{
					Name: "test-repo",
					ID:   0,
					URL:  "https://google.com",
					VCS:  1,
				},
			},
			Filename: "filename",
			Location: "location",
			Hash:     "hash",
			Enabled:  true,
		},
	}
	maka := &mockAPIKeyAuth{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	ch := NewCommandsHandler(CommandsHandlerDependencies{
		Logger:        logger,
		CommandStorer: mcs,
		TokenProvider: tp,
	})

	t.Run("get normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		commandExpected := `{"name":"test-command","id":0,"schedule":"* * * * *","repositories":[{"name":"test-repo","id":0,"url":"https://google.com","vcs":1}],"filename":"filename","location":"location","hash":"hash","enabled":true}
`

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = ch.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, commandExpected, rec.Body.String())
	})

	t.Run("get invalid id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")
		err = ch.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("get empty id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/:id")
		err = ch.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestCommandsHandler_ListCommands(t *testing.T) {
	mus := &mockUserStorer{}
	mcs := &mockCommandStorer{
		commandList: []*models.Command{
			{
				Name:     "test-command1",
				ID:       0,
				Schedule: "10 * * * *",
				Filename: "filename1",
				Location: "location1",
				Hash:     "hash1",
				Enabled:  true,
			},
			{
				Name:     "test-command2",
				ID:       1,
				Schedule: "15 * * * *",
				Filename: "filename2",
				Location: "location2",
				Hash:     "hash2",
				Enabled:  true,
			},
		},
	}
	maka := &mockAPIKeyAuth{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	ch := NewCommandsHandler(CommandsHandlerDependencies{
		Logger:        logger,
		CommandStorer: mcs,
		TokenProvider: tp,
	})

	t.Run("list normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		expectedCommandsResponse := `[{"name":"test-command1","id":0,"schedule":"10 * * * *","filename":"filename1","location":"location1","hash":"hash1","enabled":true},{"name":"test-command2","id":1,"schedule":"15 * * * *","filename":"filename2","location":"location2","hash":"hash2","enabled":true}]
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/commands")
		err = ch.List()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, expectedCommandsResponse, rec.Body.String())
	})

	t.Run("list normal flow with filters", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		listOpts := `{"name": "1"}`
		expectedCommandsResponse := `[{"name":"test-command1","id":0,"schedule":"10 * * * *","filename":"filename1","location":"location1","hash":"hash1","enabled":true}]
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(listOpts))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/commands")
		err = ch.List()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, expectedCommandsResponse, rec.Body.String())
	})

}

func TestCommandsHandler_UpdateCommand(t *testing.T) {
	mus := &mockUserStorer{}
	mcs := &mockCommandStorer{}
	maka := &mockAPIKeyAuth{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	ch := NewCommandsHandler(CommandsHandlerDependencies{
		Logger:        logger,
		CommandStorer: mcs,
		TokenProvider: tp,
	})

	t.Run("update normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		commandPost := `{"name":"test-command1","id":0,"schedule":"10 * * * *","filename":"filename1","location":"location1","hash":"hash1","enabled":true}`
		commandExpected := `{"name":"test-command1","id":0,"schedule":"10 * * * *","filename":"filename1","location":"location1","hash":"hash1","enabled":true}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/command/update", strings.NewReader(commandPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = ch.Update()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, commandExpected, rec.Body.String())
	})

	t.Run("update invalid syntax on body", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		commandPost := `<xml>`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/repository/update", strings.NewReader(commandPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = ch.Update()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestCommandsHandler_AddCommandRelForRepository(t *testing.T) {
	mus := &mockUserStorer{}
	mcs := &mockCommandStorer{
		getCommand: &models.Command{
			Name:     "test-command",
			ID:       0,
			Schedule: "* * * * *",
			Repositories: []*models.Repository{
				{
					Name: "test-repo",
					ID:   0,
					URL:  "https://google.com",
					VCS:  1,
				},
			},
			Filename: "filename",
			Location: "location",
			Hash:     "hash",
			Enabled:  true,
		},
	}
	maka := &mockAPIKeyAuth{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	ch := NewCommandsHandler(CommandsHandlerDependencies{
		Logger:        logger,
		CommandStorer: mcs,
		TokenProvider: tp,
	})

	t.Run("add relation happy path", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/add-command-rel-for-repository/:cmdid/:repoid")
		c.SetParamNames("cmdid", "repoid")
		c.SetParamValues("0", "0")
		err = ch.AddCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})

	t.Run("add relation invalid command id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/add-command-rel-for-repository/:cmdid/:repoid")
		c.SetParamNames("cmdid", "repoid")
		c.SetParamValues("invalid", "0")
		err = ch.AddCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("add relation invalid repo id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/add-command-rel-for-repository/:cmdid/:repoid")
		c.SetParamNames("cmdid", "repoid")
		c.SetParamValues("0", "invalid")
		err = ch.AddCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("add relation empty id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/add-command-rel-for-repository/:cmdid/:repoid")
		err = ch.AddCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestCommandsHandler_RemoveCommandRelForRepository(t *testing.T) {
	mus := &mockUserStorer{}
	mcs := &mockCommandStorer{
		getCommand: &models.Command{
			Name:     "test-command",
			ID:       0,
			Schedule: "* * * * *",
			Repositories: []*models.Repository{
				{
					Name: "test-repo",
					ID:   0,
					URL:  "https://google.com",
					VCS:  1,
				},
			},
			Filename: "filename",
			Location: "location",
			Hash:     "hash",
			Enabled:  true,
		},
	}
	maka := &mockAPIKeyAuth{}
	logger := zerolog.New(os.Stderr)
	deps := Dependencies{
		Logger:     logger,
		UserStore:  mus,
		APIKeyAuth: maka,
	}
	tp, err := NewTokenHandler(deps)
	assert.NoError(t, err)
	ch := NewCommandsHandler(CommandsHandlerDependencies{
		Logger:        logger,
		CommandStorer: mcs,
		TokenProvider: tp,
	})

	t.Run("remove relation happy path", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/remove-command-rel-for-repository/:cmdid/:repoid")
		c.SetParamNames("cmdid", "repoid")
		c.SetParamValues("0", "0")
		err = ch.RemoveCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})

	t.Run("remove relation invalid command id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/remove-command-rel-for-repository/:cmdid/:repoid")
		c.SetParamNames("cmdid", "repoid")
		c.SetParamValues("invalid", "0")
		err = ch.RemoveCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("remove relation invalid repo id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/remove-command-rel-for-repository/:cmdid/:repoid")
		c.SetParamNames("cmdid", "repoid")
		c.SetParamValues("0", "invalid")
		err = ch.RemoveCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("remove relation empty id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/command/remove-command-rel-for-repository/:cmdid/:repoid")
		err = ch.RemoveCommandRelForRepository()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}
