package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type mockApiKeysStore struct {
	providers.APIKeys

	key *models.APIKey
	err error
}

func (mak *mockApiKeysStore) GetByApiKeyID(ctx context.Context, id string) (*models.APIKey, error) {
	if id == mak.key.APIKeyID {
		return mak.key, mak.err
	}
	return nil, mak.err
}

func TestMatch(t *testing.T) {
	mak := &mockApiKeysStore{
		key: &models.APIKey{
			ID:           0,
			Name:         "test-key",
			UserID:       1,
			APIKeyID:     "api-key-id",
			APIKeySecret: []byte("secret"),
			TTL:          time.Now(),
		},
	}

	p, err := NewApiKeysProvider(ApiKeysConfig{}, ApiKeysDependencies{
		ApiKeysStore: mak,
		Logger:       zerolog.New(os.Stderr),
	})
	assert.NoError(t, err)
	secret := []byte("secret")
	encrypted, err := p.Encrypt(context.Background(), secret)
	assert.NoError(t, err)
	mak.key.APIKeySecret = encrypted
	err = p.Match(context.Background(), &models.APIKey{
		ID:           0,
		Name:         "test-key",
		UserID:       1,
		APIKeyID:     "api-key-id",
		APIKeySecret: secret,
	})
	assert.NoError(t, err)
	err = p.Match(context.Background(), &models.APIKey{
		ID:           0,
		Name:         "test-key",
		UserID:       1,
		APIKeyID:     "api-key-id",
		APIKeySecret: []byte("secret2"),
	})
	assert.Error(t, err)
	mak.err = kerr.ErrNotFound
	err = p.Match(context.Background(), &models.APIKey{
		ID:           0,
		Name:         "test-key",
		UserID:       1,
		APIKeyID:     "nope",
		APIKeySecret: secret,
	})
	assert.Error(t, err)
}
