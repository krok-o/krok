package auth

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// ApiKeysConfig has the configuration options for the api keys auth provider.
type ApiKeysConfig struct {
}

// ApiKeysDependencies defines the dependencies for the apikeys provider.
type ApiKeysDependencies struct {
	Logger       zerolog.Logger
	ApiKeysStore providers.APIKeysStorer
}

// ApiKeysProvider is the authentication provider for api keys.
type ApiKeysProvider struct {
	ApiKeysConfig
	ApiKeysDependencies
}

// NewApiKeysProvider creates a new authentication provider for api keys.
func NewApiKeysProvider(cfg ApiKeysConfig, deps ApiKeysDependencies) (*ApiKeysProvider, error) {
	return &ApiKeysProvider{
		ApiKeysConfig:       cfg,
		ApiKeysDependencies: deps,
	}, nil
}

var _ providers.ApiKeysAuthenticator = &ApiKeysProvider{}

// Match matches a given user's api keys with the stored ones.
func (a *ApiKeysProvider) Match(ctx context.Context, key *models.APIKey) error {
	// It doesn't matter who the api keys belong to at this stage.
	log := a.Logger.With().Str("id", key.APIKeyID).Str("name", key.Name).Logger()
	storedKey, err := a.ApiKeysStore.GetByApiKeyID(ctx, key.APIKeyID)
	if err != nil {
		log.Debug().Err(err).Msg("ApiKeys Get failed.")
		return fmt.Errorf("failed to get api key: %w", err)
	}
	return bcrypt.CompareHashAndPassword([]byte(storedKey.APIKeySecret), []byte(key.APIKeySecret))
}

// Encrypt takes an api key secret and encrypts it for storage.
func (a *ApiKeysProvider) Encrypt(ctx context.Context, secret []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
}
