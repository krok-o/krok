package livestore

import (
	"context"
	"fmt"
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

var testID = 0

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
		Name:         fmt.Sprintf("Test%d", testID),
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename",
		Location:     location,
		Hash:         "hash",
		Enabled:      false,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, c.ID)
	testID++
}

func TestCommandStore_Create_NameIsUnique(t *testing.T) {
	// increment the test id at least once
	defer func() { testID++ }()
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
	name := fmt.Sprintf("Test%d", testID)
	c, err := cp.Create(context.Background(), &models.Command{
		Name:         name,
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename",
		Location:     location,
		Hash:         "hash",
		Enabled:      false,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, c.ID)

	// Create the second command with the same name.
	_, err = cp.Create(context.Background(), &models.Command{
		Name:         name,
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename",
		Location:     location,
		Hash:         "hash",
		Enabled:      false,
	})
	assert.Error(t, err)
}
