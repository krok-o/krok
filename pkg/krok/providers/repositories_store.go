package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// RepositoryStorer handles operations for repositories and relationship to commands.
type RepositoryStorer interface {
	// These are basic CRUD operations for repository entries.

	Create(ctx context.Context, c *models.Repository) (*models.Repository, error)
	Get(ctx context.Context, id int) (*models.Repository, error)
	GetByName(ctx context.Context, name string) (*models.Repository, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, c *models.Repository) (*models.Repository, error)
	List(ctx context.Context, opt *models.ListOptions) ([]*models.Repository, error)
}
