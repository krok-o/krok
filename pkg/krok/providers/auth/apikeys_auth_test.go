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
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

type mockAPIKeysStore struct {
	providers.APIKeysStorer

	key *models.APIKey
	err error
}

func (mak *mockAPIKeysStore) GetByAPIKeyID(ctx context.Context, id string) (*models.APIKey, error) {
	if id == mak.key.APIKeyID {
		return mak.key, mak.err
	}
	return nil, mak.err
}

func TestMatch(t *testing.T) {
	mt := &mocks.Clock{}
	mt.On("Now").Return(time.Now())
	mak := &mockAPIKeysStore{
		key: &models.APIKey{
			ID:           0,
			Name:         "test-key",
			UserID:       1,
			APIKeyID:     "api-key-id",
			APIKeySecret: "secret",
			TTL:          "10m",
			CreateAt:     time.Now(),
		},
	}

	p := NewAPIKeysProvider(APIKeysDependencies{
		APIKeysStore: mak,
		Clock:        mt,
		Logger:       zerolog.New(os.Stderr),
	})
	secret := "secret"
	encrypted, err := p.Encrypt(context.Background(), []byte(secret))
	assert.NoError(t, err)
	mak.key.APIKeySecret = string(encrypted)
	err = p.Match(context.Background(), &models.APIKey{
		ID:           0,
		Name:         "test-key",
		UserID:       1,
		APIKeyID:     "api-key-id",
		APIKeySecret: secret,
		TTL:          "10m",
		CreateAt:     time.Now(),
	})
	assert.NoError(t, err)
	err = p.Match(context.Background(), &models.APIKey{
		ID:           0,
		Name:         "test-key",
		UserID:       1,
		APIKeyID:     "api-key-id",
		APIKeySecret: "secret2",
		TTL:          "10m",
		CreateAt:     time.Now(),
	})
	assert.Error(t, err)
	mak.err = kerr.ErrNotFound
	err = p.Match(context.Background(), &models.APIKey{
		ID:           0,
		Name:         "test-key",
		UserID:       1,
		APIKeyID:     "nope",
		APIKeySecret: secret,
		TTL:          "10m",
		CreateAt:     time.Now(),
	})
	assert.Error(t, err)
	mak = &mockAPIKeysStore{
		key: &models.APIKey{
			ID:           0,
			Name:         "test-key",
			UserID:       1,
			APIKeyID:     "api-key-id",
			APIKeySecret: "secret",
			TTL:          "10m",
			CreateAt:     time.Now().Add(-15 * time.Minute),
		},
	}
	p = NewAPIKeysProvider(APIKeysDependencies{
		APIKeysStore: mak,
		Clock:        mt,
		Logger:       zerolog.New(os.Stderr),
	})
	err = p.Match(context.Background(), &models.APIKey{
		ID:           0,
		Name:         "test-key",
		UserID:       1,
		APIKeyID:     "api-key-id",
		APIKeySecret: secret,
		TTL:          "10m",
		CreateAt:     time.Now().Add(-15 * time.Minute),
	})
	assert.Error(t, err)
}
