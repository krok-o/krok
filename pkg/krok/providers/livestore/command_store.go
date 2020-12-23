package livestore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	commandsTable    = "commands"
	commandsRelTable = "rel_command_repositories"
	fileLockTable    = "file_lock"
)

// CommandStore is a postgres based store for commands.
type CommandStore struct {
	CommandDependencies
}

// CommandDependencies command specific dependencies such as, the repository store.
// In order to not repeat some SQL, the command store will require the repository
// store and the repository store will require the command store.
type CommandDependencies struct {
	Dependencies
	Connector *Connector
}

// NewCommandStore creates a new CommandStore
func NewCommandStore(deps CommandDependencies) *CommandStore {
	cs := &CommandStore{CommandDependencies: deps}
	// launch the cleanup routine.
	go cs.lockCleaner(context.TODO())
	return cs
}

var _ providers.CommandStorer = &CommandStore{}

// lockCleaner will periodically delete old / stuck entries.
func (s *CommandStore) lockCleaner(ctx context.Context) {
	log := s.Logger.With().Str("func", "lockCleaner").Logger()
	interval := 10 * time.Minute
	for {

		f := func(tx pgx.Tx) error {
			t := time.Now().Add(-10 * time.Minute)
			if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where lock_start < %s", fileLockTable, t)); err != nil {
				log.Debug().Err(err).Msg("Failed to run cleanup")
				return err
			} else {
				log.Debug().Int64("rows", tags.RowsAffected()).Msg("Cleaner run successfully affecting rows...")
			}
			return nil
		}

		if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
			// just log the error, don't stop the cleaner.
			log.Debug().Err(err).Msg("Failed to run cleanup with transaction.")
		}
		// look for old entries.
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			s.Logger.Debug().Msg("Lock Cleaner cancelled.")
		}
	}
}

// Create creates a command record.
func (s *CommandStore) Create(ctx context.Context, c *models.Command) (*models.Command, error) {
	log := s.Logger.With().Str("name", c.Name).Logger()
	// duplicate key value violates unique constraint
	// id will be generated.

	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(name, schedule, filename, hash, location,"+
			" enabled) values($1, $2, $3, $4, $5, $6)", commandsTable),
			c.Name,
			c.Schedule,
			c.Filename,
			c.Hash,
			c.Location,
			c.Enabled); err != nil {
			log.Debug().Err(err).Msg("Failed to create command.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into commands",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into commands",
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, err
	}

	result, err := s.GetByName(ctx, c.Name)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get created command.")
		return nil, err
	}
	return result, nil
}

// Get returns a command model.
func (s *CommandStore) Get(ctx context.Context, id int) (*models.Command, error) {
	log := s.Logger.With().Int("id", id).Str("func", "GetByID").Logger()
	return s.getByX(ctx, log, "id", id)
}

// GetByName returns a command model by name.
func (s *CommandStore) GetByName(ctx context.Context, name string) (*models.Command, error) {
	log := s.Logger.With().Str("name", name).Str("func", "GetByName").Logger()
	return s.getByX(ctx, log, "name", name)
}

// Get returns a command model.
func (s *CommandStore) getByX(ctx context.Context, log zerolog.Logger, field string, value interface{}) (*models.Command, error) {
	log = s.Logger.With().Str("field", field).Interface("value", value).Logger()

	var (
		name      string
		commandID int
		schedule  string
		filename  string
		location  string
		hash      string
		enabled   bool
	)
	f := func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, fmt.Sprintf("select name, id, schedule, filename, location, hash, enabled from %s where %s = $1", commandsTable, field), value).
			Scan(&name, &commandID, &schedule, &filename, &location, &hash, &enabled); err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Query: "select ",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query row.")
			return &kerr.QueryError{
				Query: "select by X",
				Err:   err,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("failed to run in transactions")
		return nil, err
	}

	repositories, err := s.getRepositoriesForCommand(ctx, commandID)
	if err != nil {
		log.Debug().Err(err).Msg("GetRepositoriesForCommand failed")
		return nil, &kerr.QueryError{
			Query: "select id",
			Err:   err,
		}
	}

	return &models.Command{
		Name:         name,
		ID:           commandID,
		Schedule:     schedule,
		Repositories: repositories,
		Filename:     filename,
		Location:     location,
		Hash:         hash,
		Enabled:      enabled,
	}, nil
}

func (s *CommandStore) getRepositoriesForCommand(ctx context.Context, id int) ([]*models.Repository, error) {
	log := s.Logger.With().Int("id", id).Logger()

	// Select the related repositories.
	result := make([]*models.Repository, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select (r.id, name, url) from %s r inner join %s rel"+
			" on r.id = rel.command_id where r.id = $1", repositoriesTable, repositoryRelTable), id)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Query: "select id",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query rel_repositories_command.")
			return &kerr.QueryError{
				Query: "select id",
				Err:   fmt.Errorf("failed to query rel table: %w", err),
			}
		}

		// Repo data here construct, individual repos.
		for rows.Next() {
			var (
				repoID int
				name   string
				url    string
			)
			if err := rows.Scan(&repoID, &name, &url); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select id",
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
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute GetRepositoriesForCommand: %w", err)
	}
	return result, nil
}

// Delete will remove a command.
func (s *CommandStore) Delete(ctx context.Context, id int) error {
	log := s.Logger.With().Int("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if commandTags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", commandsTable), id); err != nil {
			log.Debug().Err(err).Msg("Failed to delete command.")
			return &kerr.QueryError{
				Query: "delete id",
				Err:   fmt.Errorf("failed delete command: %w", err),
			}
		} else if commandTags.RowsAffected() > 0 {
			// Make sure to only delete the relationship if the delete was successful.
			if err := s.deleteAllRepositoryRelForCommand(ctx, id); err != nil {
				log.Debug().Err(err).Msg("Failed to delete repository relationship for command.")
				return &kerr.QueryError{
					Query: "delete id",
					Err:   fmt.Errorf("failed to delete repository relationship for command: %w", err),
				}
			}
		}
		return nil
	}

	return s.Connector.ExecuteWithTransaction(ctx, log, f)
}

func (s *CommandStore) deleteAllRepositoryRelForCommand(ctx context.Context, id int) error {
	log := s.Logger.With().Str("func", "DeleteAllRepositoryRelForCommand").Int("command_id", id).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where command_id = $1", repositoryRelTable),
			id); err != nil {
			log.Debug().Err(err).Msg("Failed to delete relationship between command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from " + repositoryRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from " + repositoryRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to delete from " + repositoryRelTable)
		return err
	}
	return nil
}

// Update modifies a command record.
func (s *CommandStore) Update(ctx context.Context, c *models.Command) (*models.Command, error) {
	log := s.Logger.With().Int("id", c.ID).Str("name", c.Name).Logger()
	var result *models.Command
	f := func(tx pgx.Tx) error {
		// Prevent updating the ID and the creation timestamp.
		// construct update statement:
		commandTags, err := tx.Exec(ctx, fmt.Sprintf("update %s set name = $1, enabled = $2, schedule = $3, filename = $4, location = $5, hash = $6 where id = $7", commandsTable),
			c.Name, c.Enabled, c.Schedule, c.Filename, c.Location, c.Hash, c.ID)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   fmt.Errorf("failed to update: %w", err),
			}
		}
		if commandTags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   kerr.ErrNoRowsAffected,
			}
		}
		result, err = s.Get(ctx, c.ID)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   errors.New("failed to get updated command"),
			}
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, fmt.Errorf("failed to execute update in transaction: %w", err)
	}
	return result, nil
}

// List gets all the command records.
func (s *CommandStore) List(ctx context.Context, opts *models.ListOptions) ([]*models.Command, error) {
	log := s.Logger.With().Str("func", "List").Logger()
	// Select all commands.
	result := make([]*models.Command, 0)
	f := func(tx pgx.Tx) error {
		sql := fmt.Sprintf("select id, name, schedule, filename, hash, location, enabled from %s", commandsTable)
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
					Query: "select all commands",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query commands.")
			return &kerr.QueryError{
				Query: "select all commands",
				Err:   fmt.Errorf("failed to list all commands: %w", err),
			}
		}

		for rows.Next() {
			var (
				id       int
				name     string
				schedule string
				filename string
				hash     string
				location string
				enabled  bool
			)
			if err := rows.Scan(&id, &name, &schedule, &filename, &hash, &location, &enabled); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all commands",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			command := &models.Command{}
			result = append(result, command)
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all commands: %w", err)
	}
	return result, nil
}

// AcquireLock acquires a lock on a file so no other process deals with the same file.
func (s *CommandStore) AcquireLock(ctx context.Context, name string) error {
	log := s.Logger.With().Str("func", "AcquireLock").Str("name", name).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(name, lock_start) values($1, $2)", fileLockTable),
			name, time.Now()); err != nil {
			log.Debug().Err(err).Msg("Failed to acquire lock on file.")
			return &kerr.QueryError{
				Err:   fmt.Errorf("failed to acquire lock: %w", err),
				Query: err.Error(),
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into file_lock",
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to acquire lock")
		return err
	}
	return nil
}

// ReleaseLock releases a lock.
func (s *CommandStore) ReleaseLock(ctx context.Context, name string) error {
	log := s.Logger.With().Str("func", "ReleaseLock").Str("name", name).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where name = $1", fileLockTable),
			name); err != nil {
			log.Debug().Err(err).Msg("Failed to release lock on file.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from file_lock",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from file_lock",
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to release lock")
		return err
	}
	return nil
}

// GetCommandsForRepository returns a list of commands for a repository ID.
func (s *CommandStore) GetCommandsForRepository(ctx context.Context, id int) ([]*models.Command, error) {
	log := s.Logger.With().Int("id", id).Logger()

	// Select the related commands.
	result := make([]*models.Command, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select id, name, schedule, filename, hash, location, enabled from %s as c inner join %s as relc"+
			" on c.repository_id = relc.repository_id where c.repository_id = $1", commandsTable, commandsRelTable), id)
		if err != nil {
			if err.Error() == "no rows in result set" {
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
				storedID string
				name     string
				schedule string
				fileName string
				hash     string
				location string
				enabled  bool
			)
			if err := rows.Scan(&storedID, &name, &schedule, &fileName, &hash, &location, &enabled); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select id",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			command := &models.Command{
				Name:     name,
				ID:       id,
				Schedule: schedule,
				Filename: fileName,
				Location: location,
				Hash:     hash,
				Enabled:  enabled,
			}
			result = append(result, command)
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute GetCommandsForRepository: %w", err)
	}
	return result, nil
}

// AddCommandRelForRepository add an assignment for a command to a repository.
func (s *CommandStore) AddCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error {
	log := s.Logger.With().Str("func", "AddCommandRelForRepository").Int("command_id", commandID).Int("repository_id", repositoryID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(command_id, repository_id) values($1, $2)", commandsRelTable),
			commandID, repositoryID); err != nil {
			log.Debug().Err(err).Msg("Failed to create relationship between command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into " + commandsRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into " + commandsRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to insert into " + commandsRelTable)
		return err
	}
	return nil
}

// DeleteCommandRelForRepository deletes entries for a command.
// I.e.: The command was deleted so remove its connection to any repository which had this command.
func (s *CommandStore) DeleteCommandRelForRepository(ctx context.Context, commandID int) error {
	log := s.Logger.With().Str("func", "DeleteCommandRelForRepository").Int("command_id", commandID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where command_id = $1", commandsRelTable),
			commandID); err != nil {
			log.Debug().Err(err).Msg("Failed to delete relationship between command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from " + commandsRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from " + commandsRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to delete from " + commandsRelTable)
		return err
	}
	return nil
}

// DeleteAllCommandRelForRepository deletes all relationship entries for this repository.
// I.e.: The repository was deleted and now the relationships are gone to all commands.
func (s *CommandStore) DeleteAllCommandRelForRepository(ctx context.Context, repositoryID int) error {
	log := s.Logger.With().Str("func", "DeleteAllCommandRelForRepository").Int("repository_id", repositoryID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where repositroy_id = $1", commandsRelTable),
			repositoryID); err != nil {
			log.Debug().Err(err).Msg("Failed to delete relationship between command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from " + commandsRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from " + commandsRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to delete from " + commandsRelTable)
		return err
	}
	return nil
}
