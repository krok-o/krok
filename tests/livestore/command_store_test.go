package livestore

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers/auth"
	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/vault"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestCommandStore_Flow(t *testing.T) {
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
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:         "Test_Create",
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename-create",
		Location:     location,
		Hash:         "hash1",
		Enabled:      false,
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	// Get the command.
	cGet, err := cp.Get(ctx, c.ID)
	assert.NoError(t, err)
	assert.Equal(t, &models.Command{
		Name:         "Test_Create",
		ID:           c.ID,
		Schedule:     "test-schedule",
		Repositories: []*models.Repository{},
		Filename:     "test-filename-create",
		Location:     location,
		Hash:         "hash1",
		Enabled:      false,
	}, cGet)

	// List commands
	commands, err := cp.List(ctx, &models.ListOptions{})
	assert.NoError(t, err)
	assert.True(t, len(commands) > 0)

	// Update command
	cGet.Name = "UpdatedName"
	updatedC, err := cp.Update(ctx, cGet)
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedName", updatedC.Name)

	// Delete commands
	err = cp.Delete(ctx, c.ID)
	assert.NoError(t, err)

	// Try getting the deleted command should result in NotFound
	_, err = cp.Get(ctx, c.ID)
	assert.True(t, errors.Is(err, kerr.ErrNotFound))
}

func TestCommandStore_RelationshipFlow(t *testing.T) {
	// TODO This relationship definitely needs to be refactored.
	// setup repository provider
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandStore_RelationshipFlow")
	env := environment.NewDockerConverter(environment.Config{}, environment.Dependencies{Logger: logger})
	fileStore, err := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	assert.NoError(t, err)
	err = fileStore.Init()
	assert.NoError(t, err)
	v, err := vault.NewKrokVault(vault.Config{}, vault.Dependencies{Logger: logger, Storer: fileStore})
	assert.NoError(t, err)
	a, err := auth.NewKrokAuth(auth.Config{}, auth.Dependencies{
		Logger: logger,
		Vault:  v,
	})
	assert.NoError(t, err)
	connector := livestore.NewDatabaseConnector(livestore.Config{
		Hostname: dbaccess.Hostname,
		Database: dbaccess.Db,
		Username: dbaccess.Username,
		Password: dbaccess.Password,
	}, livestore.Dependencies{
		Logger:    logger,
		Converter: env,
	})
	cp := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: connector,
	})
	rp := livestore.NewRepositoryStore(livestore.RepositoryDependencies{
		Dependencies: livestore.Dependencies{
			Converter: env,
			Logger:    logger,
		},
		Connector: connector,
		Vault:     v,
		Auth:      a,
	})
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:         "Test_Relationship_Flow",
		Schedule:     "Test_Relationship_Flow-test-schedule",
		Repositories: nil,
		Filename:     "Test_Relationship_Flow-test-filename-create",
		Location:     location,
		Hash:         "Test_Relationship_Flow-hash1",
		Enabled:      false,
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)
	// Add repository relation
	repo, err := rp.Create(ctx, &models.Repository{
		Name: "TestRepo1",
		URL:  "https://github.com/Skarlso/test",
		Auth: &models.Auth{
			SSH:      "testSSH",
			Username: "testUsername",
			Password: "testPassword",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	err = rp.AddRepositoryRelForCommand(ctx, c.ID, repo.ID)
	assert.NoError(t, err)
	err = cp.AddCommandRelForRepository(ctx, c.ID, repo.ID)
	assert.NoError(t, err)

	cget, err := cp.Get(ctx, c.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, cget.Repositories)
	assert.Len(t, cget.Repositories, 1)

	repositories := cget.Repositories
	assert.NotEmpty(t, repositories)
	assert.Len(t, repositories, 1)

	// deleting a command removes the relationship from the repository
	err = cp.Delete(ctx, cget.ID)
	assert.NoError(t, err)
	// get again to retrieve repository information
	repo, err = rp.Get(ctx, repo.ID)
	assert.NoError(t, err)
	commands := repo.Commands
	assert.Empty(t, commands)

	// deleting the repository removes the relationship from the command
}

func TestCommandStore_AcquireAndReleaseLock(t *testing.T) {
	logger := zerolog.New(os.Stderr)
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
	ctx := context.Background()
	// Acquire lock
	err := cp.AcquireLock(ctx, "lock-test")
	assert.NoError(t, err)
	// Release lock
	err = cp.ReleaseLock(ctx, "lock-test")
	assert.NoError(t, err)
	// Can acquire again after release
	err = cp.AcquireLock(ctx, "lock-test")
	assert.NoError(t, err)
	// Can't acquire lock again for the same name
	err = cp.AcquireLock(ctx, "lock-test")
	assert.True(t, errors.Is(err, kerr.ErrAcquireLockFailed))
	// Release a none existing lock
	err = cp.ReleaseLock(ctx, "invalid")
	assert.True(t, errors.Is(err, kerr.ErrNoRowsAffected))
}

func TestCommandStore_Create_Unique(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandStore_Create_Unique")
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
	assert.True(t, 0 < c.ID)

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
	assert.True(t, strings.Contains(err.Error(), "unique constraint \"commands_name_key\""))
	// Create the second command with the same filename.
	_, err = cp.Create(context.Background(), &models.Command{
		Name:         "Test_Create_Error-2",
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename-create-error",
		Location:     location,
		Hash:         "hash3",
		Enabled:      false,
	})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unique constraint \"commands_filename_key\""))
	// Create the second command with the same hash.
	_, err = cp.Create(context.Background(), &models.Command{
		Name:         "Test_Create_Error-2",
		Schedule:     "test-schedule",
		Repositories: nil,
		Filename:     "test-filename-create-error=2",
		Location:     location,
		Hash:         "hash2",
		Enabled:      false,
	})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unique constraint \"commands_hash_key\""))
}
