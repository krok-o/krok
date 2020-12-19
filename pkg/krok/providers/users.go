package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// UserStorer handles CRUD operations for users.
type UserStorer interface {
	Create(ctx context.Context, c *models.User) (*models.User, error)
	Get(ctx context.Context, id string) (*models.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*models.User, error)
}
