package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// UserStorer handles CRUD operations for users.
type UserStorer interface {
	// User CRUD commands

	Create(ctx context.Context, c *models.User) (*models.User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*models.User, error)
	Get(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
}
