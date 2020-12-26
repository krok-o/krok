package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// CommandStorer handles CRUD operations for commands.
type CommandStorer interface {
	// These are basic CRUD operations for command entries.

	Create(ctx context.Context, c *models.Command) (*models.Command, error)
	Get(ctx context.Context, id int) (*models.Command, error)
	GetByName(ctx context.Context, name string) (*models.Command, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, c *models.Command) (*models.Command, error)
	List(ctx context.Context, opts *models.ListOptions) ([]*models.Command, error)

	// These functions handle operations on rel_command_repositories relationship table.

	// AddCommandRelForRepository adds an entry for this command id to the given repositoryID.
	AddCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error
	// RemoveCommandRelForRepository remove a relation to a repository for a command.
	RemoveCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error

	// Lock file functionality to prevent processing multiple commands simultaneously.

	// AcquireLock creates a lock entry for a name
	AcquireLock(ctx context.Context, name string) error
	// ReleaseLock deletes a lock entry for a name
	ReleaseLock(ctx context.Context, name string) error
}
