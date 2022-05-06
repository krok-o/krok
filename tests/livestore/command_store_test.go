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
	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/vault"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestCommandStore_Flow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
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
	assert.NoError(t, err)
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:          "Test_Create",
		Schedule:      "test-schedule",
		Repositories:  nil,
		Enabled:       false,
		Image:         "krokhook/slack-notification:v0.0.1",
		RequiresClone: true,
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	// Get the command.
	cGet, err := cp.Get(ctx, c.ID)
	assert.NoError(t, err)
	assert.Equal(t, &models.Command{
		Name:          "Test_Create",
		ID:            c.ID,
		Schedule:      "test-schedule",
		Repositories:  []*models.Repository{},
		Enabled:       false,
		Image:         "krokhook/slack-notification:v0.0.1",
		RequiresClone: true,
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
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandStore_RelationshipFlow")
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	fileStore := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	err := fileStore.Init()
	assert.NoError(t, err)
	v := vault.NewKrokVault(vault.Dependencies{Logger: logger, Storer: fileStore})
	assert.NoError(t, err)
	connector := livestore.NewDatabaseConnector(livestore.Config{
		Hostname: hostname,
		Database: dbaccess.Db,
		Username: dbaccess.Username,
		Password: dbaccess.Password,
	}, livestore.Dependencies{
		Logger:    logger,
		Converter: env,
	})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: connector,
	})
	assert.NoError(t, err)
	rp := livestore.NewRepositoryStore(livestore.RepositoryDependencies{
		Dependencies: livestore.Dependencies{
			Converter: env,
			Logger:    logger,
		},
		Connector: connector,
		Vault:     v,
	})
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:          "Test_Relationship_Flow",
		Schedule:      "Test_Relationship_Flow-test-schedule",
		Enabled:       false,
		Image:         "krokhook/slack-notification:v0.0.1",
		RequiresClone: true,
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
	// Create the second command.
	c2, err := cp.Create(ctx, &models.Command{
		Name:          "Test_Relationship_Flow-2",
		Schedule:      "Test_Relationship_Flow-test-schedule-2",
		Enabled:       false,
		Image:         "krokhook/slack-notification:v0.0.1",
		RequiresClone: true,
	})
	assert.NoError(t, err)

	// add repository relationship
	err = cp.AddCommandRelForRepository(ctx, c2.ID, repo.ID)
	assert.NoError(t, err)

	// Get and check the repository connection
	c2, err = cp.Get(ctx, c2.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, c2.Repositories)

	// Remove the repository
	err = rp.Delete(ctx, repo.ID)
	assert.NoError(t, err)

	// get again to get repositories
	c2, err = cp.Get(ctx, c2.ID)
	assert.NoError(t, err)
	assert.Empty(t, c2.Repositories)
}

func TestCommandStore_Create_Unique(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
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
	assert.NoError(t, err)
	// Create the first command.
	c, err := cp.Create(context.Background(), &models.Command{
		Name:          "Test_Create_Error",
		Schedule:      "test-schedule",
		Repositories:  nil,
		Image:         "krokhook/slack-notification:v0.0.1",
		Enabled:       false,
		RequiresClone: false,
	})
	require.NoError(t, err)
	assert.True(t, 0 < c.ID)

	// Create the second command with the same name.
	_, err = cp.Create(context.Background(), &models.Command{
		Name:         "Test_Create_Error",
		Schedule:     "test-schedule",
		Repositories: nil,
		Enabled:      false,
		Image:        "krokhook/slack-notification:v0.0.1",
	})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unique constraint \"commands_name_key\""))
}

func TestCommandStore_PlatformRelationshipFlow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	connector := livestore.NewDatabaseConnector(livestore.Config{
		Hostname: hostname,
		Database: dbaccess.Db,
		Username: dbaccess.Username,
		Password: dbaccess.Password,
	}, livestore.Dependencies{
		Logger:    logger,
		Converter: env,
	})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: connector,
	})
	assert.NoError(t, err)
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:    "Test_Relationship_Flow_Platform",
		Enabled: true,
		Image:   "krokhook/slack-notification:v0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)
	err = cp.AddCommandRelForPlatform(ctx, c.ID, models.GITHUB)
	assert.NoError(t, err)

	supported, err := cp.IsPlatformSupported(ctx, c.ID, models.GITHUB)
	assert.NoError(t, err)
	assert.True(t, supported)

	supported, err = cp.IsPlatformSupported(ctx, c.ID, 999)
	assert.NoError(t, err)
	assert.False(t, supported)

	// Get command, platforms should be in platform list.
	command, err := cp.Get(ctx, c.ID)
	assert.NoError(t, err)
	assert.Contains(t, command.Platforms, models.SupportedPlatforms[models.GITHUB], "Github not found in the supported platforms list.")

	// remove the relation
	err = cp.RemoveCommandRelForPlatform(ctx, c.ID, models.GITHUB)
	assert.NoError(t, err)
	// platform list should be empty.
	command, err = cp.Get(ctx, c.ID)
	assert.NoError(t, err)
	assert.Empty(t, command.Platforms, "Github not found in the supported platforms list.")

	supported, err = cp.IsPlatformSupported(ctx, c.ID, models.GITHUB)
	assert.NoError(t, err)
	assert.False(t, supported)
}

func TestCommandStore_Update(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
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
	assert.NoError(t, err)
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:          "Test_Update",
		Schedule:      "test-schedule",
		Enabled:       false,
		Image:         "krokhook/slack-notification:v0.0.1",
		RequiresClone: true,
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	// Update command
	updatedC, err := cp.Update(ctx, &models.Command{ID: c.ID, Name: "UpdatedName2"})
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedName2", updatedC.Name)
	// Make sure nothing else changed.
	assert.Equal(t, c.Schedule, updatedC.Schedule)
	assert.Equal(t, c.Enabled, updatedC.Enabled)
	assert.Equal(t, c.Image, updatedC.Image)
}
