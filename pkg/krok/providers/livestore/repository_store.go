package livestore

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"

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

var _ providers.RepositoryStorer = &RepositoryStore{}

// Create creates a repository. Upon creating we don't assign any commands yet. So we don't save those here.
// We do save auth information into the vault.
func (r *RepositoryStore) Create(ctx context.Context, c *models.Repository) (*models.Repository, error) {
	log := r.Logger.With().Str("name", c.Name).Logger()
	// duplicate key value violates unique constraint
	// id will be generated.

	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(name, url) values($1, $2)", repositoriesTable),
			c.Name,
			c.URL); err != nil {
			log.Debug().Err(err).Msg("Failed to create repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into repository",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.NoRowsAffected,
				Query: "insert into repository",
			}
		}
		return nil
	}

	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, err
	}

	result, err := r.GetByName(ctx, c.Name)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get created command.")
		return nil, err
	}

	if err := r.Auth.CreateRepositoryAuth(ctx, result.ID, c.Auth); err != nil {
		log.Debug().Err(err).Msg("Failed to store auth information.")
		return nil, err
	}

	return result, nil
}

func (r *RepositoryStore) Delete(ctx context.Context, id string) error {
	log := r.Logger.With().Str("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if commandTags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", repositoriesTable), id); err != nil {
			log.Debug().Err(err).Msg("Failed to delete repository.")
			return &kerr.QueryError{
				Query: "delete id: " + id,
				Err:   fmt.Errorf("failed to delete repository: %w", err),
			}
		} else if commandTags.RowsAffected() > 0 {
			// Make sure to only delete the relationship if the delete was successful.
			if err := r.CommandStore.DeleteAllCommandRelForRepository(ctx, id); err != nil {
				log.Debug().Err(err).Msg("Failed to delete repository relationship for repository.")
				return &kerr.QueryError{
					Query: "delete id: " + id,
					Err:   fmt.Errorf("failed to delete repository relationship for command: %w", err),
				}
			}
		}
		return nil
	}

	return r.Connector.ExecuteWithTransaction(ctx, log, f)
}

// Update can only update the name of the repository. If auth information is updated for the repository,
// it has to be re-created. Since auth is stored elsewhere.
func (r *RepositoryStore) Update(ctx context.Context, c models.Repository) (*models.Repository, error) {
	log := r.Logger.With().Str("id", c.ID).Str("name", c.Name).Logger()
	var result *models.Repository
	f := func(tx pgx.Tx) error {
		// Prevent updating the ID and the creation timestamp.
		// construct update statement:
		tags, err := tx.Exec(ctx, fmt.Sprintf("update %s set name = $1", commandsTable),
			c.Name)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   fmt.Errorf("failed to update: %w", err),
			}
		}
		if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   kerr.NoRowsAffected,
			}
		}
		result, err = r.Get(ctx, c.ID)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   errors.New("failed to get updated repository"),
			}
		}
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, fmt.Errorf("failed to execute update in transaction: %w", err)
	}
	return result, nil
}

// List all repositories or the ones specified by the filter opts.
// We are ignoring auth information here.
func (r *RepositoryStore) List(ctx context.Context, opts *models.ListOptions) ([]*models.Repository, error) {
	log := r.Logger.With().Str("func", "List").Logger()
	// Select all repositories.
	result := make([]*models.Repository, 0)
	f := func(tx pgx.Tx) error {
		sql := fmt.Sprintf("select id, name, url from %s", repositoriesTable)
		where := " where "
		filters := make([]string, 0)
		if opts.Name != "" {
			filters = append(filters, "name = %"+opts.Name+"%")
		}
		filter := strings.Join(filters, " AND ")
		if filter != "" {
			sql += where + filter
		}
		rows, err := tx.Query(ctx, sql)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Query: "select all repositories",
					Err:   kerr.NotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query repositories.")
			return &kerr.QueryError{
				Query: "select all repositories",
				Err:   fmt.Errorf("failed to list all repositories: %w", err),
			}
		}

		for rows.Next() {
			var (
				id   string
				name string
				url  string
			)
			if err := rows.Scan(&id, &name, &url); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all repositories",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			repository := &models.Repository{}
			result = append(result, repository)
		}
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all repositories: %w", err)
	}
	return result, nil
}

func (r *RepositoryStore) AddRepositoryRelForCommand(ctx context.Context, commandID string, repositoryID string) error {
	log := r.Logger.With().Str("func", "AddRepositoryRelForCommand").Str("command_id", commandID).Str("repository_id", repositoryID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(command_id, repository_id) values($1, $2)", repositoryRelTable),
			commandID, repositoryID); err != nil {
			log.Debug().Err(err).Msg("Failed to create relationship between repository and command.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into " + repositoryRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.NoRowsAffected,
				Query: "insert into " + repositoryRelTable,
			}
		}
		return nil
	}

	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to insert into " + repositoryRelTable)
		return err
	}
	return nil
}

func (r *RepositoryStore) DeleteAllRepositoryRelForCommand(ctx context.Context, commandID string) error {
	log := r.Logger.With().Str("func", "DeleteAllRepositoryRelForCommand").Str("command_id", commandID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where command_id = $1", repositoryRelTable),
			commandID); err != nil {
			log.Debug().Err(err).Msg("Failed to delete relationship between command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from " + repositoryRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.NoRowsAffected,
				Query: "delete from " + repositoryRelTable,
			}
		}
		return nil
	}

	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to delete from " + repositoryRelTable)
		return err
	}
	return nil
}

func (r *RepositoryStore) DeleteRepositoryRelForCommand(ctx context.Context, repositoryID string) error {
	log := r.Logger.With().Str("func", "DeleteRepositoryRelForCommand").Str("repository_id", repositoryID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where repository_id = $1", repositoryRelTable),
			repositoryID); err != nil {
			log.Debug().Err(err).Msg("Failed to delete relationship between command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from " + repositoryRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.NoRowsAffected,
				Query: "delete from " + repositoryRelTable,
			}
		}
		return nil
	}

	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to delete from " + repositoryRelTable)
		return err
	}
	return nil
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

func (r *RepositoryStore) Get(ctx context.Context, id string) (*models.Repository, error) {
	log := r.Logger.With().Str("func", "Get").Logger()
	return r.getByX(ctx, log, "id", id)
}

func (r *RepositoryStore) GetByName(ctx context.Context, name string) (*models.Repository, error) {
	log := r.Logger.With().Str("func", "GetByName").Logger()
	return r.getByX(ctx, log, "name", name)
}

// Get fetches a repository by ID.
// Also returns Auth information for the repository.
func (r *RepositoryStore) getByX(ctx context.Context, log zerolog.Logger, field string, value interface{}) (*models.Repository, error) {
	log = r.Logger.With().Str("field", field).Interface("value", value).Logger()
	// Get all data from the repository table.
	var result *models.Repository
	f := func(tx pgx.Tx) error {
		var (
			id, name, url string
		)
		if err := tx.QueryRow(ctx, fmt.Sprintf("select id, name, url from widgets where %s=$1", field), value).Scan(&id, &name, &url); err != nil {
			return &kerr.QueryError{
				Query: "select id",
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
	commands, err := r.CommandStore.GetCommandsForRepository(ctx, result.ID)
	if !errors.Is(err, kerr.NotFound) {
		log.Debug().Err(err).Msg("Get failed to get repository commands.")
		return nil, err
	}
	result.Commands = commands

	// Get auth info
	auth, err := r.Auth.GetRepositoryAuth(ctx, result.ID)
	if !errors.Is(err, kerr.NotFound) {
		log.Debug().Err(err).Msg("GetRepositoryAuth failed.")
		return nil, err
	}
	result.Auth = auth
	return result, nil
}
