package livestore

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestEventsStore_Create(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	// location, _ := ioutil.TempDir("", "TestEventsStore_Create")
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	es := livestore.NewEventsStorer(livestore.EventsStoreDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
	})
	ctx := context.Background()
	event, err := es.Create(ctx, &models.Event{
		EventID:      "uuid1",
		CreateAt:     time.Now(),
		RepositoryID: 1,
		CommandRuns:  make([]*models.CommandRun, 0),
		Payload:      "{}",
	})
	assert.NoError(t, err)
	assert.NotEqual(t, 0, event.ID, "Event ID should have been a sequence and increased to above 0.")
}

func TestEventsStore_GetWithRuns(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	es := livestore.NewEventsStorer(livestore.EventsStoreDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
	})
	crs := livestore.NewCommandRunStore(livestore.CommandRunDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
	})

	ctx := context.Background()
	event, err := es.Create(ctx, &models.Event{
		EventID:      "uuid2",
		CreateAt:     time.Now(),
		RepositoryID: 1,
		Payload:      "{}",
	})
	assert.NoError(t, err)
	assert.NotEqual(t, 0, event.ID, "Event ID should have been a sequence and increased to above 0.")

	run := &models.CommandRun{
		ID:          1,
		EventID:     event.ID,
		CommandName: "test-command",
		Status:      "failed",
		Outcome:     "file not found",
		CreateAt:    time.Now(),
	}
	r, err := crs.CreateRun(ctx, run)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, r.ID, "Expected id to not equal 0 because it's an automatic sequencer.")

	// get the event and see if the command run is assigned to it
	e, err := es.GetEvent(ctx, r.EventID)
	assert.NoError(t, err)
	assert.Equal(t, event.ID, e.ID)
	assert.Equal(t, event.EventID, e.EventID)
	assert.Equal(t, event.Payload, e.Payload)
	assert.Equal(t, run.CommandName, e.CommandRuns[0].CommandName)
	assert.Equal(t, run.EventID, e.CommandRuns[0].EventID)
	assert.Equal(t, run.ID, e.CommandRuns[0].ID)
	assert.Equal(t, run.Outcome, e.CommandRuns[0].Outcome)
	assert.Equal(t, run.Status, e.CommandRuns[0].Status)
}

func TestEventsStore_List(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	es := livestore.NewEventsStorer(livestore.EventsStoreDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
	})

	t.Run("basic list function", func(tt *testing.T) {
		ctx := context.Background()
		event, err := es.Create(ctx, &models.Event{
			EventID:      "uuid3",
			CreateAt:     time.Now(),
			RepositoryID: 1,
			Payload:      "{}",
		})
		assert.NoError(t, err)
		assert.NotEqual(t, 0, event.ID, "Event ID should have been a sequence and increased to above 0.")

		events, err := es.ListEventsForRepository(ctx, 1, models.ListOptions{})
		assert.NoError(tt, err)
		assert.NotZero(tt, len(events), "events list should not have come back as empty")
	})

}
