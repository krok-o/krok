package livestore

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers/auth"
	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/vault"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestRepositoryStore_Flow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestRepositoryStore_Create")
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
	a, err := auth.NewKrokAuth(auth.AuthConfig{}, auth.AuthDependencies{
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
	repo, err := rp.Create(ctx, &models.Repository{
		Name: "TestRepo_Create_No_Auth",
		URL:  "https://github.com/krok-o/test",
		VCS:  models.GITHUB,
	})
	assert.NoError(t, err)
	assert.True(t, repo.ID > 0)

	// Get the repo.
	getRepo, err := rp.Get(ctx, repo.ID)
	assert.NoError(t, err)
	assert.Equal(t, repo, getRepo)

	// List repos
	repos, err := rp.List(ctx, &models.ListOptions{})
	assert.NoError(t, err)
	assert.True(t, len(repos) > 0)

	// Update repos
	getRepo.Name = "UpdatedName"
	updatedR, err := rp.Update(ctx, getRepo)
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedName", updatedR.Name)

	// Delete repo
	err = rp.Delete(ctx, getRepo.ID)
	assert.NoError(t, err)

	// Try getting the deleted command should result in NotFound
	_, err = rp.Get(ctx, getRepo.ID)
	assert.True(t, errors.Is(err, kerr.ErrNotFound))
}

func TestRepositoryStore_ListByFilter(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestRepositoryStore_ListByFilter")
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
	a, err := auth.NewKrokAuth(auth.AuthConfig{}, auth.AuthDependencies{
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
	_, err = rp.Create(ctx, &models.Repository{
		Name: "TestRepo_ListByName-1",
		URL:  "https://github.com/krok-o/test",
		VCS:  models.GITHUB,
	})
	assert.NoError(t, err)
	_, err = rp.Create(ctx, &models.Repository{
		Name: "TestRepo_ListByName-2",
		URL:  "https://github.com/krok-o/test",
		VCS:  models.GITHUB,
	})
	assert.NoError(t, err)
	_, err = rp.Create(ctx, &models.Repository{
		Name: "TestRepo_ListByVCS",
		URL:  "https://github.com/krok-o/test",
		VCS:  models.GITEA,
	})
	assert.NoError(t, err)

	allRepos, err := rp.List(ctx, &models.ListOptions{})
	assert.NoError(t, err)
	assert.True(t, len(allRepos) > 2)

	onlyName, err := rp.List(ctx, &models.ListOptions{Name: "TestRepo_ListByName-"})
	assert.NoError(t, err)
	assert.Len(t, onlyName, 2)

	onlyVcs, err := rp.List(ctx, &models.ListOptions{VCS: models.GITEA})
	assert.NoError(t, err)
	assert.Len(t, onlyVcs, 1)
}

func TestRepositoryStore_Create_Unique(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestRepositoryStore_Create_Unique")
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
	a, err := auth.NewKrokAuth(auth.AuthConfig{}, auth.AuthDependencies{
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
	rp := livestore.NewRepositoryStore(livestore.RepositoryDependencies{
		Dependencies: livestore.Dependencies{
			Converter: env,
			Logger:    logger,
		},
		Connector: connector,
		Vault:     v,
		Auth:      a,
	})
	cp := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: connector,
	})
	ctx := context.Background()
	c, err := cp.Create(ctx, &models.Command{
		Name:     "CommandConnectionName",
		Schedule: "Schedule100",
		Filename: "Create_Filename",
		Location: "Location-10",
		Hash:     "Hash-10",
		Enabled:  false,
	})
	assert.NoError(t, err)
	repo, err := rp.Create(ctx, &models.Repository{
		Name: "TestRepo_Create_Unique",
		URL:  "https://github.com/krok-o/test",
		VCS:  models.GITHUB,
	})
	assert.NoError(t, err)
	assert.True(t, repo.ID > 0)

	err = cp.AddCommandRelForRepository(ctx, c.ID, repo.ID)
	assert.NoError(t, err)

	// get the repository to retrieve commands.
	repo, err = rp.Get(ctx, repo.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, repo.Commands)

	// delete the relationship and see if the command was removed
	err = cp.RemoveCommandRelForRepository(ctx, c.ID, repo.ID)
	assert.NoError(t, err)

	// get the repository again
	repo, err = rp.Get(ctx, repo.ID)
	assert.NoError(t, err)
	assert.Empty(t, repo.Commands)
}

func TestRepositoryStore_Create_WithCommands(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestRepositoryStore_Create_WithCommands")
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
	a, err := auth.NewKrokAuth(auth.AuthConfig{}, auth.AuthDependencies{
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
	repo, err := rp.Create(ctx, &models.Repository{
		Name: "TestRepo_Create_WithCommands",
		URL:  "https://github.com/krok-o/test",
		VCS:  models.BITBUCKET,
	})
	assert.NoError(t, err)
	assert.True(t, repo.ID > 0)
	_, err = rp.Create(ctx, &models.Repository{
		Name: "TestRepo_Create_WithCommands",
		URL:  "https://github.com/krok-o/test",
	})
	assert.Error(t, err)
}
