package livestore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
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
	Connector    *Connector
	CommandStore providers.CommandStorer
	Vault        providers.Vault
	Auth         providers.Auth
}

// NewRepositoryStore creates a new RepositoryStore
func NewRepositoryStore(cfg Config, deps RepositoryDependencies) *RepositoryStore {
	return &RepositoryStore{Config: cfg, RepositoryDependencies: deps}
}

// GetRepositoriesForCommand returns a list of repositories for a command ID.
// This, does not return Auth information.
func (r *RepositoryStore) GetRepositoriesForCommand(ctx context.Context, id string) ([]*models.Repository, error) {
	log := r.Logger.With().Str("id", id).Logger()
	if id == "" {
		return nil, fmt.Errorf("GetRepositoriesForCommand failed with %w", kerr.InvalidArgument)
	}

	// Select the related repositories.
	result := make([]*models.Repository, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select id, name, url from %s as r inner join %s as rel"+
			" on r.command_id = rel.command_id where r.command_id = $1", repositoriesTable, repositoryRelTable), id)
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

		// Repo data here construct, individual repos.
		for rows.Next() {
			var (
				repoID string
				name   string
				url    string
			)
			if err := rows.Scan(&repoID, &name, &url); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select id: " + id,
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			repo := &models.Repository{
				Name: name,
				ID:   id,
				URL:  url,
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
// Also returns Auth information for the repository.
func (r *RepositoryStore) Get(ctx context.Context, rid string) (*models.Repository, error) {
	log := r.Logger.With().Str("id", rid).Logger()
	if rid == "" {
		log.Warn().Msg("ID is empty.")
		return nil, fmt.Errorf("Get failed with %w", kerr.InvalidArgument)
	}
	// Get all data from the repository table.
	var result *models.Repository
	f := func(tx pgx.Tx) error {
		var (
			id, name, url string
		)
		if err := tx.QueryRow(ctx, "select id, name, url from widgets where id=$1", rid).Scan(&id, &name, &url); err != nil {
			return &kerr.QueryError{
				Query: "select id: " + rid,
				Err:   fmt.Errorf("failed to scan: %w", err),
			}
		}
		result.ID = id
		result.Name = name
		result.URL = url
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute Get: %w", err)
	}

	// Get all commands from the rel table.
	commands, err := r.CommandStore.GetCommandsForRepository(ctx, rid)
	if !errors.Is(err, kerr.NotFound) {
		log.Debug().Err(err).Msg("Get failed to get repository commands.")
		return nil, err
	}
	result.Commands = commands

	// Get auth info
	auth, err := r.Auth.GetRepositoryAuth(ctx, rid)
	if !errors.Is(err, kerr.NotFound) {
		log.Debug().Err(err).Msg("GetRepositoryAuth failed.")
		return nil, err
	}
	result.Auth = auth
	return result, nil
}
