package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// RepositoryStorer handles operations for repositories and relationship to commands.
type RepositoryStorer interface {
	// These are basic CRUD operations for repository entries.

	Create(ctx context.Context, c *models.Repository) (*models.Repository, error)
	Get(ctx context.Context, id string) (*models.Repository, error)
	GetByName(ctx context.Context, name string) (*models.Repository, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, c models.Repository) (*models.Repository, error)
	List(ctx context.Context, opt *models.ListOptions) ([]*models.Repository, error)

	// These functions handle operations on rel_repositories_command relationship table.

	// GetRepositoriesForCommand retrieves all repositories using the relationship table entries for the given command ID.
	GetRepositoriesForCommand(ctx context.Context, commandID string) ([]*models.Repository, error)
	// AddRepositoryRelForCommand adds an entry for this command id to the given repositoryID.
	AddRepositoryRelForCommand(ctx context.Context, commandID string, repositoryID string) error
	// DeleteAllRepositoryRelForCommand deletes all relationship entries for this command.
	// I.e.: The command was deleted and now the relationships to all repositories this command
	// was assigned to, must also be removed.
	DeleteAllRepositoryRelForCommand(ctx context.Context, commandID string) error
	// DeleteRepositoryRelForCommand deletes a single entry for a repository and a command.
	// I.e.: The repository was deleted so remove its connection to a command.
	// When viewing the command, that repository must not show up any longer.
	DeleteRepositoryRelForCommand(ctx context.Context, commandID string) error
}
