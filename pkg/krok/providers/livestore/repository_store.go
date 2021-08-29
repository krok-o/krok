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
	repositoriesTable = "repositories"
)

// RepositoryStore is a postgres based store for repositories.
type RepositoryStore struct {
	RepositoryDependencies
}

// RepositoryDependencies repository specific dependencies such as, the command store.
type RepositoryDependencies struct {
	Dependencies
	Connector *Connector
	Vault     providers.Vault
}

// NewRepositoryStore creates a new RepositoryStore
func NewRepositoryStore(deps RepositoryDependencies) *RepositoryStore {
	return &RepositoryStore{RepositoryDependencies: deps}
}

var _ providers.RepositoryStorer = &RepositoryStore{}

// Create creates a repository. Upon creating we don't assign any commands yet. So we don't save those here.
// We do save auth information into the vault.
func (r *RepositoryStore) Create(ctx context.Context, c *models.Repository) (*models.Repository, error) {
	log := r.Logger.With().Str("name", c.Name).Logger()
	// duplicate key value violates unique constraint
	// id will be generated.

	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(name, url, vcs, project_id) values($1, $2, $3, $4)", repositoriesTable),
			c.Name,
			c.URL,
			c.VCS,
			c.GitLab.GetProjectID()); err != nil {
			log.Debug().Err(err).Msg("Failed to create repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into repository",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
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
		log.Debug().Err(err).Msg("Failed to get created repository.")
		return nil, err
	}
	result.Events = c.Events
	return result, nil
}

// Delete removes a repository and all.
func (r *RepositoryStore) Delete(ctx context.Context, id int) error {
	log := r.Logger.With().Int("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if tag, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", repositoriesTable), id); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select id",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to delete repository.")
			return &kerr.QueryError{
				Query: "delete id",
				Err:   fmt.Errorf("failed to delete repository: %w", err),
			}
		} else if tag.RowsAffected() == 0 {
			return kerr.ErrNoRowsAffected
		}
		return nil
	}

	return r.Connector.ExecuteWithTransaction(ctx, log, f)
}

// Update can only update the name of the repository. If auth information is updated for the repository,
// it has to be re-created. Since auth is stored elsewhere.
func (r *RepositoryStore) Update(ctx context.Context, c *models.Repository) (*models.Repository, error) {
	log := r.Logger.With().Int("id", c.ID).Str("name", c.Name).Logger()
	f := func(tx pgx.Tx) error {
		// Prevent updating the ID and the creation timestamp.
		// construct update statement:
		tags, err := tx.Exec(ctx, fmt.Sprintf("update %s set name = $1 where id = $2", repositoriesTable),
			c.Name, c.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			return &kerr.QueryError{
				Query: "select id",
				Err:   kerr.ErrNotFound,
			}
		}
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   fmt.Errorf("failed to update: %w", err),
			}
		}
		if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   kerr.ErrNoRowsAffected,
			}
		}
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, fmt.Errorf("failed to execute update in transaction: %w", err)
	}
	result, err := r.Get(ctx, c.ID)
	if err != nil {
		return nil, &kerr.QueryError{
			Query: "update :" + c.Name,
			Err:   errors.New("failed to get updated repository"),
		}
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
		sql := fmt.Sprintf("select id, name, url, vcs, project_id from %s", repositoriesTable)
		where := " where "
		filters := make([]string, 0)
		if opts.Name != "" {
			filters = append(filters, "name LIKE '%"+opts.Name+"%'")
		}
		if opts.VCS != 0 {
			filters = append(filters, fmt.Sprintf("vcs = %d", opts.VCS))
		}
		filter := strings.Join(filters, " AND ")
		if filter != "" {
			sql += where + filter
		}
		rows, err := tx.Query(ctx, sql)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select all repositories",
					Err:   kerr.ErrNotFound,
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
				id        int
				name      string
				url       string
				vcs       int
				projectID int // this field needs to be a pointer because it can be nil which will result in a nil value.
			)
			if err := rows.Scan(&id, &name, &url, &vcs, &projectID); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all repositories",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			repository := &models.Repository{
				Name: name,
				ID:   id,
				URL:  url,
				VCS:  vcs,
				GitLab: &models.GitLab{
					ProjectID: projectID,
				},
			}
			result = append(result, repository)
		}
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all repositories: %w", err)
	}
	return result, nil
}

// Get retrieves a single repository using its ID.
func (r *RepositoryStore) Get(ctx context.Context, id int) (*models.Repository, error) {
	log := r.Logger.With().Str("func", "Get").Logger()
	return r.getByX(ctx, log, "id", id)
}

// GetByName retrieves a single repository using its name.
func (r *RepositoryStore) GetByName(ctx context.Context, name string) (*models.Repository, error) {
	log := r.Logger.With().Str("func", "GetByName").Logger()
	return r.getByX(ctx, log, "name", name)
}

// Get fetches a repository by ID.
// Also returns Auth information for the repository.
func (r *RepositoryStore) getByX(ctx context.Context, log zerolog.Logger, field string, value interface{}) (*models.Repository, error) {
	log = r.Logger.With().Str("field", field).Interface("value", value).Logger()
	// Get all data from the repository table.
	result := &models.Repository{}
	f := func(tx pgx.Tx) error {
		var (
			id, vcs   int
			name, url string
			projectID int
		)
		if err := tx.QueryRow(ctx, fmt.Sprintf("select id, name, url, vcs, project_id from %s where %s=$1", repositoriesTable, field), value).Scan(&id, &name, &url, &vcs, &projectID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select id",
					Err:   kerr.ErrNotFound,
				}
			}
			return &kerr.QueryError{
				Query: "select id",
				Err:   fmt.Errorf("failed to scan: %w", err),
			}
		}
		result.ID = id
		result.Name = name
		result.URL = url
		result.VCS = vcs
		result.GitLab = &models.GitLab{ProjectID: projectID}
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute Get: %w", err)
	}

	// Get all commands from the rel table.
	commands, err := r.getCommandsForRepository(ctx, result.ID)
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("Get failed to get repository commands.")
		return nil, err
	}
	result.Commands = commands
	return result, nil
}

// getCommandsForRepository returns a list of commands for a repository ID.
func (r *RepositoryStore) getCommandsForRepository(ctx context.Context, id int) ([]*models.Command, error) {
	log := r.Logger.With().Int("id", id).Logger()

	// Select the related commands.
	result := make([]*models.Command, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select c.id, name, schedule, enabled, image from %s as c inner join %s as relc"+
			" on c.id = relc.command_id where relc.repository_id = $1", commandsTable, commandsRepositoriesRelTable), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select id",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query relationship.")
			return &kerr.QueryError{
				Query: "select commands for repository",
				Err:   fmt.Errorf("failed to query rel table: %w", err),
			}
		}

		for rows.Next() {
			var (
				storedID int
				name     string
				schedule string
				enabled  bool
				image    string
			)
			if err := rows.Scan(&storedID, &name, &schedule, &enabled, &image); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select id",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			command := &models.Command{
				Name:     name,
				ID:       storedID,
				Schedule: schedule,
				Enabled:  enabled,
				Image:    image,
			}
			result = append(result, command)
		}
		return nil
	}
	if err := r.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute GetCommandsForRepository: %w", err)
	}
	return result, nil
}
