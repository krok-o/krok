package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// RepositoryStorer handles CRUD operations for repositories.
type RepositoryStorer interface {
	Create(ctx context.Context, c *models.Repository) error
	Get(ctx context.Context, id string) (*models.Repository, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, c models.Repository) (*models.Repository, error)
}
