package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// CommandStorer handles CRUD operations for commands.
type CommandStorer interface {
	// These are basic CRUD operations for command entries.

	Create(ctx context.Context, c *models.Command) (*models.Command, error)
	Get(ctx context.Context, id string) (*models.Command, error)
	GetByName(ctx context.Context, name string) (*models.Command, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, c *models.Command) (*models.Command, error)
	List(ctx context.Context) ([]*models.Command, error)

	// These functions handle operations on rel_command_repositories relationship table.

	// GetCommandsForRepository retrieves all commands using the relationship table entries for the given repository ID.
	GetCommandsForRepository(ctx context.Context, repoID string) ([]*models.Command, error)
	// AddCommandRelForRepository adds an entry for this command id to the given repositoryID.
	AddCommandRelForRepository(ctx context.Context, commandID string, repositoryID string) error
	// DeleteAllCommandRelForRepository deletes all relationship entries for this repository.
	// I.e.: The repository was deleted and now the relationship is gone.
	DeleteAllCommandRelForRepository(ctx context.Context, repositoryID string) error
	// DeleteCommandRelForRepository deletes a single entry for a command and a repository.
	// I.e.: The command was deleted so remove its connection to a repository.
	DeleteCommandRelForRepository(ctx context.Context, repositoryID string) error
}
