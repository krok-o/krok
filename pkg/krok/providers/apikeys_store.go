package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// APIKeysStorer defines operations that an api key provider must have
type APIKeysStorer interface {
	// Create an apikey.
	Create(ctx context.Context, key *models.APIKey) (*models.APIKey, error)
	// Delete an apikey.
	Delete(ctx context.Context, id int) error
	// List will list all apikeys for a user.
	List(ctx context.Context, userID int) ([]*models.APIKey, error)
	// Get an apikey.
	Get(ctx context.Context, id int) (*models.APIKey, error)
	// GetByApiKeyID an apikey by apikeyid.
	GetByApiKeyID(ctx context.Context, id string) (*models.APIKey, error)
}
