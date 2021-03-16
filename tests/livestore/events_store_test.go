package livestore

import (
	"context"
	"fmt"
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

	_, err = es.Create(ctx, &models.Event{
		EventID:      "uuid1",
		CreateAt:     time.Now(),
		RepositoryID: 1,
		CommandRuns:  make([]*models.CommandRun, 0),
		Payload:      "{}",
	})
	assert.Error(t, err)
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

	_, err = es.GetEvent(ctx, 999)
	assert.Error(t, err)
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

	t.Run("basic list errors", func(tt *testing.T) {
		ctx := context.Background()
		es, err := es.ListEventsForRepository(ctx, 999, models.ListOptions{})
		assert.NoError(tt, err)
		assert.Empty(tt, es)
	})

	t.Run("filter between dates", func(tt *testing.T) {
		ctx := context.Background()
		event1, err := es.Create(ctx, &models.Event{
			EventID:      "uuid4",
			CreateAt:     time.Date(2021, 03, 12, 13, 0, 0, 0, time.UTC),
			RepositoryID: 1,
			Payload:      "{}",
		})
		assert.NoError(t, err)
		assert.NotEqual(t, 0, event1.ID, "Event ID should have been a sequence and increased to above 0.")
		event2, err := es.Create(ctx, &models.Event{
			EventID:      "uuid5",
			CreateAt:     time.Date(2005, 03, 12, 13, 0, 0, 0, time.UTC),
			RepositoryID: 1,
			Payload:      "{}",
		})
		assert.NoError(t, err)
		assert.NotEqual(t, 0, event2.ID, "Event ID should have been a sequence and increased to above 0.")

		from := time.Date(2021, 02, 12, 13, 0, 0, 0, time.UTC)
		to := time.Date(2021, 03, 13, 13, 0, 0, 0, time.UTC)
		events, err := es.ListEventsForRepository(ctx, 1, models.ListOptions{
			StartingDate: &from,
			EndDate:      &to,
		})
		assert.NoError(tt, err)
		assert.Len(tt, events, 1, "events list should not have come back as empty")
		assert.Equal(tt, event1.ID, events[0].ID)

		from = time.Date(2005, 02, 12, 13, 0, 0, 0, time.UTC)
		to = time.Date(2006, 03, 13, 13, 0, 0, 0, time.UTC)
		events, err = es.ListEventsForRepository(ctx, 1, models.ListOptions{
			StartingDate: &from,
			EndDate:      &to,
		})
		assert.NoError(tt, err)
		assert.Len(tt, events, 1, "events list should not have come back as empty")
		assert.Equal(tt, event2.ID, events[0].ID)

		from = time.Date(2005, 02, 12, 13, 0, 0, 0, time.UTC)
		to = time.Date(2021, 03, 13, 13, 0, 0, 0, time.UTC)
		events, err = es.ListEventsForRepository(ctx, 1, models.ListOptions{
			StartingDate: &from,
			EndDate:      &to,
		})
		assert.NoError(tt, err)
		assert.Len(tt, events, 2, "events list should not have come back as empty")
	})

	t.Run("pagination", func(tt *testing.T) {
		ctx := context.Background()
		_, err := es.Create(ctx, &models.Event{
			EventID:      "uuid14",
			CreateAt:     time.Date(2005, 03, 12, 13, 1, 0, 0, time.UTC),
			RepositoryID: 1,
			Payload:      "{}",
		})
		assert.NoError(tt, err)
		_, err = es.Create(ctx, &models.Event{
			EventID:      "uuid15",
			CreateAt:     time.Date(2005, 03, 12, 13, 2, 0, 0, time.UTC),
			RepositoryID: 1,
			Payload:      "{}",
		})
		assert.NoError(tt, err)
		event3, err := es.Create(ctx, &models.Event{
			EventID:      "uuid16",
			CreateAt:     time.Date(2005, 03, 12, 13, 3, 0, 0, time.UTC),
			RepositoryID: 1,
			Payload:      "{}",
		})
		assert.NoError(tt, err)

		events, err := es.ListEventsForRepository(ctx, 1, models.ListOptions{
			PageSize: 1,
		})

		assert.NoError(tt, err)
		assert.Len(tt, events, 1)

		events, err = es.ListEventsForRepository(ctx, 1, models.ListOptions{
			PageSize: 2,
		})

		assert.NoError(tt, err)
		assert.Len(tt, events, 2)

		from := time.Date(2005, 03, 12, 13, 1, 0, 0, time.UTC)
		to := time.Date(2005, 03, 12, 13, 4, 0, 0, time.UTC)
		events, err = es.ListEventsForRepository(ctx, 1, models.ListOptions{
			StartingDate: &from,
			EndDate:      &to,
			PageSize:     1,
			Page:         3,
		})

		assert.NoError(tt, err)
		assert.Len(tt, events, 1)
		fmt.Println(events[0])
		assert.Equal(tt, event3.ID, events[0].ID)
	})
}
