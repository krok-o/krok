package auth

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/vault"
	"github.com/krok-o/krok/pkg/models"
)

func TestKrokAuth_CreateRepositoryAuth(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestKrokAuth_CreateRepositoryAuth")
	fileStore, err := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	assert.NoError(t, err)
	err = fileStore.Init()
	assert.NoError(t, err)
	v, err := vault.NewKrokVault(vault.Config{}, vault.Dependencies{Logger: logger, Storer: fileStore})
	assert.NoError(t, err)
	auth, err := NewKrokAuth(Config{}, Dependencies{
		Logger: logger,
		Vault:  v,
	})
	assert.NoError(t, err)

	info := &models.Auth{
		SSH:      "testssh",
		Username: "testusername",
		Password: "testpassword",
	}
	ctx := context.Background()
	err = auth.CreateRepositoryAuth(ctx, 1, info)
	assert.NoError(t, err)

	// Get the repository info
	a, err := auth.GetRepositoryAuth(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, info, a)
}

func TestKrokAuth_CreateRepositoryAuthPartialAuth(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestKrokAuth_CreateRepositoryAuthPartialAuth")
	fileStore, err := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	assert.NoError(t, err)
	err = fileStore.Init()
	assert.NoError(t, err)
	v, err := vault.NewKrokVault(vault.Config{}, vault.Dependencies{Logger: logger, Storer: fileStore})
	assert.NoError(t, err)
	auth, err := NewKrokAuth(Config{}, Dependencies{
		Logger: logger,
		Vault:  v,
	})
	assert.NoError(t, err)

	info := &models.Auth{
		SSH: "testssh",
	}
	ctx := context.Background()
	err = auth.CreateRepositoryAuth(ctx, 1, info)
	assert.NoError(t, err)

	// Get the repository info
	a, err := auth.GetRepositoryAuth(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, info, a)
}
