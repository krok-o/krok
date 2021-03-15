package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func TestHandleHooks(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	mrs := &mocks.RepositoryStorer{}
	mrs.On("Get", mock.Anything, 1).Return(&models.Repository{ID: 1}, nil)
	mgp := &mocks.Platform{}
	mgp.On("ValidateRequest", mock.Anything, mock.Anything, 1).Return(nil)
	mgp.On("GetEventID", mock.Anything, mock.Anything).Return("id", nil)
	platformProviders := make(map[int]providers.Platform)
	platformProviders[models.GITHUB] = mgp
	es := &mocks.EventsStorer{}
	es.On("Create", mock.Anything, mock.Anything).Return(&models.Event{
		ID:           1,
		EventID:      "id",
		RepositoryID: 1,
		Commands:     make([]*models.Command, 0),
		Payload:      "",
	}, nil)
	ex := &mocks.Executor{}
	ex.On("CreateRun", mock.Anything, mock.Anything).Return(nil)
	deps := HookDependencies{
		Logger:            logger,
		RepositoryStore:   mrs,
		PlatformProviders: platformProviders,
		EventsStorer:      es,
		Executer:          ex,
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
