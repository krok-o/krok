package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// ApiKeys defines operations that an api key provider must have.
type ApiKeys interface {
	// Create an apikey.
	Create(ctx context.Context, key *models.ApiKey) (*models.ApiKey, error)
	// Delete an apikey.
	Delete(ctx context.Context, id int) error
	// List will list all apikeys for a user.
	List(ctx context.Context, userID int) ([]*models.ApiKey, error)
	// Get an apikey.
	Get(ctx context.Context, id int) (*models.ApiKey, error)
}
