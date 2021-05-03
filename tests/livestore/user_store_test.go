package livestore

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestUserStore_Flow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	clock := &mocks.Clock{}
	clock.On("Now").Return(time.Now())
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
	ap := livestore.NewAPIKeysStore(livestore.APIKeysDependencies{
		Connector: connector,
	})
	up := livestore.NewUserStore(livestore.UserDependencies{
		Connector: connector,
		APIKeys:   ap,
		Time:      clock,
	})
	ctx := context.Background()
	user, err := up.Create(ctx, &models.User{
		DisplayName: "DisplayName",
		Email:       "valid-1@email.com",
	})
	assert.NoError(t, err)
	assert.True(t, user.ID > 0)

	// Get the user.
	getUser, err := up.Get(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, getUser)

	// List users
	users, err := up.List(ctx)
	assert.NoError(t, err)
	assert.True(t, len(users) > 0)

	// Update users
	getUser.DisplayName = "UpdatedName"
	getUser.Token = "UpdatedToken"
	updatedU, err := up.Update(ctx, getUser)
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedName", updatedU.DisplayName)
	assert.Equal(t, "UpdatedToken", updatedU.Token)

	// Delete user
	err = up.Delete(ctx, getUser.ID)
	assert.NoError(t, err)

	// Try getting the deleted command should result in NotFound
	_, err = up.Get(ctx, getUser.ID)
	assert.True(t, errors.Is(err, kerr.ErrNotFound))
}

func TestUserStore_Create_Unique(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	clock := &mocks.Clock{}
	clock.On("Now").Return(time.Now())
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
	ap := livestore.NewAPIKeysStore(livestore.APIKeysDependencies{
		Connector: connector,
	})
	up := livestore.NewUserStore(livestore.UserDependencies{
		Connector: connector,
		APIKeys:   ap,
		Time:      clock,
	})
	ctx := context.Background()
	_, err := up.Create(ctx, &models.User{
		DisplayName: "DisplayName",
		Email:       "valid-2@email.com",
		LastLogin:   time.Now(),
	})
	assert.NoError(t, err)
	_, err = up.Create(ctx, &models.User{
		DisplayName: "DisplayName",
		Email:       "valid-2@email.com",
		LastLogin:   time.Now(),
	})
	assert.Error(t, err)
}
