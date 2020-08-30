package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// CommandStorer handles CRUD operations for commands.
type CommandStorer interface {
	Create(ctx context.Context, c *models.Command) error
	Get(ctx context.Context, id string) (*models.Command, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, c *models.Command) (*models.Command, error)
}
