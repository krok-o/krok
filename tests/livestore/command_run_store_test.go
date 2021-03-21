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

func TestCommandRun_Create(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
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
	run := &models.CommandRun{
		EventID:     1,
		CommandName: "test-command",
		Status:      "failed",
		Outcome:     "file not found",
		CreateAt:    time.Now(),
	}
	ctx := context.Background()
	r, err := crs.CreateRun(ctx, run)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, r.ID, "Expected id to not equal 0 because it's an automatic sequencer.")
}

func TestCommandRun_UpdateRunStatus(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
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
	run := &models.CommandRun{
		EventID:     1,
		CommandName: "test-command",
		Status:      "running",
		Outcome:     "",
		CreateAt:    time.Now(),
	}
	ctx := context.Background()
	r, err := crs.CreateRun(ctx, run)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, r.ID, "Expected id to not equal 0 because it's an automatic sequencer.")

	err = crs.UpdateRunStatus(ctx, r.ID, "success", "all good")
	assert.NoError(t, err)

	r, err = crs.Get(ctx, r.ID)
	assert.NoError(t, err)
	assert.Equal(t, "success", r.Status)
	assert.Equal(t, "all good", r.Outcome)

	err = crs.UpdateRunStatus(ctx, 999, "success", "all good")
	assert.Error(t, err)

	_, err = crs.Get(ctx, 999)
	assert.Error(t, err)
}
