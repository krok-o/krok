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

	"github.com/krok-o/krok/pkg/krok/providers"
)

type mockPlatformTokenProvider struct {
	providers.PlatformTokenProvider
}

func (m *mockPlatformTokenProvider) SaveTokenForPlatform(token string, vcs int) error {
	return nil
}

func TestVCSTokenHandler_Create(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	mpt := &mockPlatformTokenProvider{}
	vtp, err := NewVCSTokenHandler(Config{}, VCSTokenHandlerDependencies{
		Logger:        logger,
		TokenProvider: mpt,
	})
	assert.NoError(t, err)
	token, err := generateTestToken("test@email.com")
	assert.NoError(t, err)

	tokenPost := `{"vcs" : 1, "token" : "github_token"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/vcs-token", strings.NewReader(tokenPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = vtp.Create()(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
