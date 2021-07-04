package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

type mockUserStorer struct {
	providers.UserStorer
}

func (mus *mockUserStorer) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return &models.User{
		DisplayName: "testUser",
		Email:       email,
		ID:          0,
		LastLogin:   time.Now(),
		APIKeys: []*models.APIKey{
			{
				ID:           0,
				Name:         "test",
				UserID:       0,
				APIKeyID:     "apikeyid",
				APIKeySecret: "secret",
				TTL:          "10m",
				CreateAt:     time.Now(),
			},
		},
	}, nil
}

type mockRepositoryStorer struct {
	providers.RepositoryStorer
	id        int
	getRepo   *models.Repository
	deleteErr error
	listRepo  []*models.Repository
}

func (mrs *mockRepositoryStorer) Create(ctx context.Context, repo *models.Repository) (*models.Repository, error) {
	repo.ID = mrs.id
	mrs.id++
	return repo, nil
}

func (mrs *mockRepositoryStorer) Update(ctx context.Context, repo *models.Repository) (*models.Repository, error) {
	return repo, nil
}

func (mrs *mockRepositoryStorer) Get(ctx context.Context, id int) (*models.Repository, error) {
	return mrs.getRepo, nil
}

func (mrs *mockRepositoryStorer) List(ctx context.Context, opts *models.ListOptions) ([]*models.Repository, error) {
	return mrs.listRepo, nil
}

func (mrs *mockRepositoryStorer) Delete(ctx context.Context, id int) error {
	return mrs.deleteErr
}

type mockGithubPlatformProvider struct {
	providers.Platform
}

func (g *mockGithubPlatformProvider) CreateHook(ctx context.Context, repo *models.Repository) error {
	return nil
}

func TestRepoHandler_CreateRepository(t *testing.T) {
	mrs := &mocks.RepositoryStorer{}
	mars := &mocks.RepositoryAuth{}
	mars.On("CreateRepositoryAuth", mock.Anything, mock.Anything, &models.Auth{Secret: "secret"}).Return(nil)
	mg := &mockGithubPlatformProvider{}
	logger := zerolog.New(os.Stderr)
	cfg := RepoConfig{
		Protocol: "http",
		HookBase: "hookbase",
	}
	t.Run("positive flow of create", func(tt *testing.T) {
		mrs = &mocks.RepositoryStorer{}
		mrs.On("Create", mock.Anything, &models.Repository{
			Name: "test-name",
			URL:  "https://github.com/Skarlso/test",
			VCS:  1,
			Auth: &models.Auth{
				Secret: "secret",
			},
		}).Return(&models.Repository{
			Name: "test-name",
			URL:  "https://github.com/Skarlso/test",
			ID:   1,
			VCS:  1,
			Auth: &models.Auth{
				Secret: "secret",
			}}, nil)
		rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
			Logger:           logger,
			RepositoryStorer: mrs,
			PlatformProviders: map[int]providers.Platform{
				models.GITHUB: mg,
			},
			Auth: mars,
		})
		assert.NoError(t, err)
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryPost := `{"name" : "test-name", "url" : "https://github.com/Skarlso/test", "vcs" : 1, "auth": {"secret": "secret"}}`
		repositoryExpected := `{"name":"test-name","id":1,"url":"https://github.com/Skarlso/test","vcs":1,"auth":{"secret":"secret"},"unique_url":"http://hookbase/rest/api/1/hooks/1/1/callback"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/repository", strings.NewReader(repositoryPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = rh.Create()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusCreated, rec.Code)
		assert.Equal(tt, repositoryExpected, rec.Body.String())
	})

	t.Run("positive flow of create with project id", func(tt *testing.T) {
		mrs = &mocks.RepositoryStorer{}
		mrs.On("Create", mock.Anything, &models.Repository{
			Name:   "test-name",
			URL:    "https://github.com/Skarlso/test",
			VCS:    1,
			GitLab: &models.GitLab{ProjectID: 10},
			Auth: &models.Auth{
				Secret: "secret",
			},
		}).Return(&models.Repository{
			Name:   "test-name",
			URL:    "https://github.com/Skarlso/test",
			ID:     1,
			VCS:    1,
			GitLab: &models.GitLab{ProjectID: 10},
			Auth: &models.Auth{
				Secret: "secret",
			}}, nil)
		rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
			Logger:           logger,
			RepositoryStorer: mrs,
			PlatformProviders: map[int]providers.Platform{
				models.GITHUB: mg,
			},
			Auth: mars,
		})
		assert.NoError(t, err)
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryPost := `{"name" : "test-name", "url" : "https://github.com/Skarlso/test", "vcs" : 1, "git_lab":{"project_id": 10}, "auth": {"secret": "secret"}}`
		repositoryExpected := `{"name":"test-name","id":1,"url":"https://github.com/Skarlso/test","vcs":1,"git_lab":{"project_id":10},"auth":{"secret":"secret"},"unique_url":"http://hookbase/rest/api/1/hooks/1/1/callback"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/repository", strings.NewReader(repositoryPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = rh.Create()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusCreated, rec.Code)
		assert.Equal(tt, repositoryExpected, rec.Body.String())
	})

	t.Run("invalid post data", func(tt *testing.T) {
		mrs = &mocks.RepositoryStorer{}
		rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
			Logger:           logger,
			RepositoryStorer: mrs,
			PlatformProviders: map[int]providers.Platform{
				models.GITHUB: mg,
			},
			Auth: mars,
		})
		assert.NoError(t, err)
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryPost := `<xml>`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/repository", strings.NewReader(repositoryPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = rh.Create()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestRepoHandler_UpdateRepository(t *testing.T) {
	mrs := &mockRepositoryStorer{}
	mars := &mocks.RepositoryAuth{}
	logger := zerolog.New(os.Stderr)
	cfg := RepoConfig{
		Protocol: "http",
		HookBase: "hookbase",
	}
	rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
		Logger:           logger,
		RepositoryStorer: mrs,
		Auth:             mars,
	})
	assert.NoError(t, err)

	t.Run("update normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryPost := `{"name":"updated-name","id":0,"url":"https://github.com/Skarlso/test","vcs":1}`
		repositoryExpected := `{"name":"updated-name","id":0,"url":"https://github.com/Skarlso/test","vcs":1,"unique_url":"http://hookbase/rest/api/1/hooks/0/1/callback"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/repository/update", strings.NewReader(repositoryPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = rh.Update()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, repositoryExpected, rec.Body.String())
	})

	t.Run("update invalid syntax on body", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryPost := `<xml>`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/repository/update", strings.NewReader(repositoryPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = rh.Update()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestRepoHandler_GetRepository(t *testing.T) {
	mrs := &mockRepositoryStorer{
		getRepo: &models.Repository{
			Name: "test-name",
			ID:   0,
			URL:  "https://github.com/Skarlso/test",
			VCS:  1,
		},
	}
	mars := &mocks.RepositoryAuth{}
	mars.On("GetRepositoryAuth", mock.Anything, 0).Return(&models.Auth{
		Secret: "secret",
	}, nil)
	logger := zerolog.New(os.Stderr)
	cfg := RepoConfig{
		Protocol: "http",
		HookBase: "hookbase",
	}
	rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
		Logger:           logger,
		RepositoryStorer: mrs,
		Auth:             mars,
	})
	assert.NoError(t, err)

	t.Run("get normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryExpected := `{"name":"test-name","id":0,"url":"https://github.com/Skarlso/test","vcs":1,"auth":{"secret":"secret"},"unique_url":"http://hookbase/rest/api/1/hooks/0/1/callback"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/repository/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = rh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, repositoryExpected, rec.Body.String())
	})

	t.Run("get normal flow with project id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		mrs.getRepo.GitLab = &models.GitLab{ProjectID: 10}
		repositoryExpected := `{"name":"test-name","id":0,"url":"https://github.com/Skarlso/test","vcs":1,"git_lab":{"project_id":10},"auth":{"secret":"secret"},"unique_url":"http://hookbase/rest/api/1/hooks/0/1/callback"}
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/repository/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = rh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, repositoryExpected, rec.Body.String())
	})

	t.Run("get invalid id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/repository/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")
		err = rh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})

	t.Run("empty id", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/repository/:id")
		err = rh.Get()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
}

func TestRepoHandler_ListRepositories(t *testing.T) {
	mrs := &mockRepositoryStorer{
		listRepo: []*models.Repository{
			{
				Name: "test-name",
				ID:   0,
				URL:  "https://github.com/Skarlso/test",
				VCS:  1,
			},
			{
				Name: "test-name2",
				ID:   1,
				URL:  "https://github.com/Skarlso/test2",
				VCS:  0,
				GitLab: &models.GitLab{
					ProjectID: 10,
				},
			},
		},
	}
	logger := zerolog.New(os.Stderr)
	cfg := RepoConfig{
		Protocol: "http",
		HookBase: "hookbase",
	}
	rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
		Logger:           logger,
		RepositoryStorer: mrs,
	})
	assert.NoError(t, err)

	t.Run("list normal flow", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryExpected := `[{"name":"test-name","id":0,"url":"https://github.com/Skarlso/test","vcs":1},{"name":"test-name2","id":1,"url":"https://github.com/Skarlso/test2","vcs":0,"git_lab":{"project_id":10}}]
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/repositories")
		err = rh.List()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, repositoryExpected, rec.Body.String())
	})
}

func TestRepoHandler_DeleteRepository(t *testing.T) {
	mrs := &mockRepositoryStorer{}
	logger := zerolog.New(os.Stderr)
	cfg := RepoConfig{
		Protocol: "http",
		HookBase: "hookbase",
	}
	rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
		Logger:           logger,
		RepositoryStorer: mrs,
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
		c.SetPath("/repository/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")
		err = rh.Delete()(c)
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
		c.SetPath("/repository/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")
		err = rh.Delete()(c)
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
		c.SetPath("/repository/:id")
		err = rh.Delete()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusBadRequest, rec.Code)
	})
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
