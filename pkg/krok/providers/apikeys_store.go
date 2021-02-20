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
	Delete(ctx context.Context, id int, userID int) error
	// List will list all apikeys for a user.
	List(ctx context.Context, userID int) ([]*models.APIKey, error)
	// Get an apikey.
	Get(ctx context.Context, id int, userID int) (*models.APIKey, error)
	// GetByAPIKeyID an apikey by apikeyid.
	GetByAPIKeyID(ctx context.Context, id string) (*models.APIKey, error)
}
