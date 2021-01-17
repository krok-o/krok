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
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestApiKeys_Flow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Config{}, environment.Dependencies{Logger: logger})
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
	ctx := context.Background()
	apiKey, err := ap.Create(ctx, &models.APIKey{
		Name:         "Main",
		UserID:       1,
		APIKeyID:     "keyid",
		APIKeySecret: []byte("secret"),
		TTL:          time.Now(),
	})
	assert.NoError(t, err)
	assert.True(t, apiKey.ID > 0)

	// Get the apiKey.
	getKey, err := ap.Get(ctx, apiKey.ID, apiKey.UserID)
	assert.NoError(t, err)
	assert.Equal(t, apiKey, getKey)

	// Get the apiKey by api key id.
	getKey, err = ap.GetByApiKeyID(ctx, apiKey.APIKeyID)
	assert.NoError(t, err)
	assert.Equal(t, apiKey, getKey)

	// List keys
	keys, err := ap.List(ctx, 1)
	assert.NoError(t, err)
	assert.True(t, len(keys) > 0)

	// Delete apiKey
	err = ap.Delete(ctx, apiKey.ID, 1)
	assert.NoError(t, err)

	// Try getting the deleted command should result in NotFound
	_, err = ap.Get(ctx, apiKey.ID, apiKey.UserID)
	assert.True(t, errors.Is(err, kerr.ErrNotFound))
}
