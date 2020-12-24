package livestore

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestCommandStore_Create(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandStore_Create")
	env := environment.NewDockerConverter(environment.Config{}, environment.Dependencies{Logger: logger})
	cp := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: dbaccess.Hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
	})
	// Create the first command.
	c, err := cp.Create(context.Background(), &models.Command{
		Name:         "Test_Create",
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename-create",
		Location:     location,
		Hash:         "hash1",
		Enabled:      false,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, c.ID)
}

func TestCommandStore_Create_NameIsUnique(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandStore_Create_NameIsUnique")
	env := environment.NewDockerConverter(environment.Config{}, environment.Dependencies{Logger: logger})
	cp := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: dbaccess.Hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
	})
	// Create the first command.
	c, err := cp.Create(context.Background(), &models.Command{
		Name:         "Test_Create_Error",
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename-create-error",
		Location:     location,
		Hash:         "hash2",
		Enabled:      false,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, c.ID)

	// Create the second command with the same name.
	_, err = cp.Create(context.Background(), &models.Command{
		Name:         "Test_Create_Error",
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename-create-error-2",
		Location:     location,
		Hash:         "hash3",
		Enabled:      false,
	})
	assert.Error(t, err)
}
