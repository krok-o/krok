package auth

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	keyTTL = 7 * 24 * time.Hour
)

// APIKeysDependencies defines the dependencies for the apikeys provider.
type APIKeysDependencies struct {
	Logger       zerolog.Logger
	APIKeysStore providers.APIKeysStorer
	Clock        providers.Clock
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

// Generate a secret and a key ID pair.
func (a *APIKeysProvider) Generate(ctx context.Context, name string, userID int) (*models.APIKey, error) {
	// generate the key secret
	// this will be displayed once, then never shown again, ever.
	secret, err := a.generateUniqueKey()
	if err != nil {
		return nil, err
	}

	// generate the key id
	// this will be displayed once, then never shown again, ever.
	keyID, err := a.generateKeyID()
	if err != nil {
		return nil, err
	}

	encrypted, err := a.Encrypt(ctx, []byte(secret))
	if err != nil {
		return nil, err
	}

	key := &models.APIKey{
		Name:         name,
		UserID:       userID,
		APIKeyID:     keyID,
		APIKeySecret: string(encrypted),
		TTL:          a.Clock.Now().Add(keyTTL),
	}

	generatedKey, err := a.APIKeysStore.Create(ctx, key)
	if err != nil {
		a.Logger.Debug().Err(err).Msg("Failed to generate a key.")
		return nil, err
	}
	key.ID = generatedKey.ID
	key.APIKeySecret = secret
	return key, nil
}

// Generate a unique api key for a user.
func (a *APIKeysProvider) generateUniqueKey() (string, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}

	return u.String(), nil
}

// Generate a unique api key for a user.
func (a *APIKeysProvider) generateKeyID() (string, error) {
	u, err := a.generateUniqueKey()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(u))), nil
}
