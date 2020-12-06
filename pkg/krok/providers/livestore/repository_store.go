package livestore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/models"
)

const (
	repositoriesTable  = "repositories"
	repositoryRelTable = "rel_repositories_command"
)

// RepositoryStore is a postgres based store for repositories.
type RepositoryStore struct {
	Config
	RepositoryDependencies
}

// RepositoryDependencies repository specific dependencies such as, the command store.
type RepositoryDependencies struct {
	Dependencies
	Connector *Connector
}

// NewRepositoryStore creates a new RepositoryStore
func NewRepositoryStore(cfg Config, deps RepositoryDependencies) *RepositoryStore {
	return &RepositoryStore{Config: cfg, RepositoryDependencies: deps}
}

// GetRepositoriesForCommand returns a list of repositories for a command ID.
func (r *RepositoryStore) GetRepositoriesForCommand(ctx context.Context, id string) ([]*models.Repository, error) {
	// Select the related repositories.
	result := make([]*models.Repository, 0)
	log := r.Logger.With().Str("id", id).Logger()
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select repository_id from %s where command_id = $1", repositoryRelTable), id)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Query: "select id: " + id,
					Err:   kerr.NotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query rel_repositories_command.")
			return &kerr.QueryError{
				Query: "select id: " + id,
				Err:   fmt.Errorf("failed to query rel table: %w", err),
			}
		}
		for rows.Next() {
			var (
				repoID string
			)
			if err := rows.Scan(&repoID); err != nil {
				log.Debug().Err(err).Msg("Failed to scan repoID.")
				return &kerr.QueryError{
					Query: "select id: " + id,
					Err:   fmt.Errorf("failed to scan repoID: %w", err),
				}
			}
			repo, err := r.Get(ctx, repoID)
			if err != nil {
				return &kerr.QueryError{
					Query: "get id: " + repoID,
					Err:   fmt.Errorf("failed to get repoID: %w", err),
				}
			}
			result = append(result, repo)
		}
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute GetRepositoriesForCommand: %w", err)
	}
	return result, nil
}

// Get fetches a repository by ID.
func (r *RepositoryStore) Get(ctx context.Context, id string) (*models.Repository, error) {
	return nil, nil
}
