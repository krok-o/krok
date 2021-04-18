package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventHandler_List(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	es := &mocks.EventsStorer{}
	es.On("ListEventsForRepository", mock.Anything, 1, &models.ListOptions{}).Return([]*models.Event{
		{
			ID:           1,
			EventID:      "uuid",
			CreateAt:     time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
			RepositoryID: 1,
			CommandRuns: []*models.CommandRun{
				{
					ID:          1,
					EventID:     1,
					CommandName: "echo",
					Status:      "success",
					Outcome:     "echo this",
					CreateAt:    time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC),
				},
			},
			Payload: `{"id": "uuid", "other": "stuff"}`,
			VCS:     models.GITHUB,
		},
	}, nil)
	eh := NewEventHandler(EventHandlerDependencies{
		Logger:       logger,
		EventsStorer: es,
	})

	t.Run("can list events", func(tt *testing.T) {
		token, err := generateTestToken("test@email.com")
		assert.NoError(tt, err)

		repositoryExpected := `[{"id":1,"event_id":"uuid","create_at":"1981-01-01T01:01:01.000000001Z","repository_id":1,"command_runs":[{"id":1,"event_id":1,"command_name":"echo","status":"success","outcome":"echo this","create_at":"1981-01-01T01:01:01.000000001Z"}],"payload":"{\"id\": \"uuid\", \"other\": \"stuff\"}","vcs":1}]
`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		c := e.NewContext(req, rec)
		c.SetPath("/events/:repoid")
		c.SetParamNames("repoid")
		c.SetParamValues("1")
		err = eh.List()(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rec.Code)
		assert.Equal(tt, repositoryExpected, rec.Body.String())
	})
}
