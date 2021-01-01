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

	"github.com/krok-o/krok/pkg/krok/providers"
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
				APIKeySecret: []byte("secret"),
				TTL:          time.Now().Add(10 * time.Minute),
			},
		},
	}, nil
}

type mockRepositoryStorer struct {
	providers.RepositoryStorer
	id int
}

func (mrs *mockRepositoryStorer) Create(ctx context.Context, repo *models.Repository) (*models.Repository, error) {
	repo.ID = mrs.id
	mrs.id++
	return repo, nil
}

type mockApiKeyAuth struct {
	providers.ApiKeysAuthenticator
}

func (maka *mockApiKeyAuth) Match(ctx context.Context, key *models.APIKey) error {
	return nil
}

func TestRepoHandler_CreateRepository(t *testing.T) {
	mus := &mockUserStorer{}
	mrs := &mockRepositoryStorer{}
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
	rh, err := NewRepositoryHandler(cfg, RepoHandlerDependencies{
		Dependencies:     deps,
		RepositoryStorer: mrs,
		TokenProvider:    tp,
	})

	assert.NoError(t, err)

	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)

	repositoryPost := `{"name" : "test-name", "url" : "https://github.com/Skarlso/test", "vcs" : 1}`
	repositoryExpected := `{"name":"test-name","id":0,"url":"https://github.com/Skarlso/test","vcs":1,"unique_url":"http://testHost/0/1/callback"}
`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/repository", strings.NewReader(repositoryPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = rh.CreateRepository()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, repositoryExpected, rec.Body.String())
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
