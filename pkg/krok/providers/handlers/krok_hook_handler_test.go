package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleHooks(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	mrs := &mocks.RepositoryStorer{}
	mrs.On("Get", mock.Anything, 1).Return(&models.Repository{ID: 1}, nil)
	mgp := &mocks.Platform{}
	mgp.On("ValidateRequest", mock.Anything, mock.Anything, 1).Return(nil)
	platformProviders := make(map[int]providers.Platform)
	platformProviders[models.GITHUB] = mgp
	deps := HookDependencies{
		Logger:            logger,
		RepositoryStore:   mrs,
		PlatformProviders: platformProviders,
	}

	hh := NewHookHandler(deps)
	t.Run("handle hook event", func(tt *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/hooks/:rid/:vid/callback")
		c.SetParamNames("rid", "vid")
		c.SetParamValues("1", "1")
		err := hh.HandleHooks()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
	})
}