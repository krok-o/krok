package auth

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// APIKeysDependencies defines the dependencies for the apikeys provider.
type APIKeysDependencies struct {
	Logger       zerolog.Logger
	APIKeysStore providers.APIKeysStorer
}

// APIKeysProvider is the authentication provider for api keys.
type APIKeysProvider struct {
	APIKeysDependencies
}

// NewAPIKeysProvider creates a new authentication provider for api keys.
func NewAPIKeysProvider(deps APIKeysDependencies) *APIKeysProvider {
	return &APIKeysProvider{
		APIKeysDependencies: deps,
	}
}

var _ providers.APIKeysAuthenticator = &APIKeysProvider{}

// Match matches a given user's api keys with the stored ones.
func (a *APIKeysProvider) Match(ctx context.Context, key *models.APIKey) error {
	// It doesn't matter who the api keys belong to at this stage.
	log := a.Logger.With().Str("id", key.APIKeyID).Str("name", key.Name).Logger()
	storedKey, err := a.APIKeysStore.GetByAPIKeyID(ctx, key.APIKeyID)
	if err != nil {
		log.Debug().Err(err).Msg("APIKeys Get failed.")
		return fmt.Errorf("failed to get api key: %w", err)
	}
	return bcrypt.CompareHashAndPassword([]byte(storedKey.APIKeySecret), []byte(key.APIKeySecret))
}

// Encrypt takes an api key secret and encrypts it for storage.
func (a *APIKeysProvider) Encrypt(ctx context.Context, secret []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
}
