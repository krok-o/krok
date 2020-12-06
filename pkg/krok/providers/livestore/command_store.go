package livestore

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const commandsTable = "commands"

// CommandStore is a postgres based store for commands.
type CommandStore struct {
	Config
	CommandDependencies
}

// CommandDependencies command specific dependencies such as, the repository store.
// In order to not repeat some SQL, the command store will require the repository
// store and the repository store will require the command store.
type CommandDependencies struct {
	Dependencies
	RepositoryStore providers.RepositoryStorer
}

// NewCommandStore creates a new CommandStore
func NewCommandStore(cfg Config, deps CommandDependencies) *CommandStore {
	return &CommandStore{Config: cfg, CommandDependencies: deps}
}

// Create creates a command record.
func (s *CommandStore) Create(ctx context.Context, c *models.Command) error {
	return nil
}

// Get returns a command model.
func (s *CommandStore) Get(ctx context.Context, id string) (*models.Command, error) {
	log := s.Logger.With().Str("id", id).Logger()

	var (
		name      string
		commandID string
		schedule  string
		filename  string
		location  string
		hash      string
		enabled   bool
	)
	f := func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, fmt.Sprintf("select name, id, schedule, filename, location, hash, enabled from %s where id = $1", commandsTable), id).
			Scan(&name, &commandID, &schedule, &filename, &location, &hash, &enabled); err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Query: "select id: " + id,
					Err:   kerr.NotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query row.")
			return &kerr.QueryError{
				Query: "select id: " + id,
				Err:   err,
			}
		}
		return nil
	}

	if err := s.executeWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("failed to run in transactions")
		return nil, err
	}

	repositories, err := s.RepositoryStore.GetRepositoriesForCommand(ctx, id)
	if err != nil {
		log.Debug().Err(err).Msg("GetRepositoriesForCommand failed")
		return nil, &kerr.QueryError{
			Query: "select id: " + id,
			Err:   err,
		}
	}
	//// Select the related repositories.
	//rows, err := tx.Query(ctx, "select repository_id from rel_repositories_command where command_id = $1", id)
	//if err != nil {
	//	if err.Error() == "no rows in result set" {
	//		return nil, nil
	//	}
	//	log.Debug().Err(err).Msg("Failed to query rel_repositories_command.")
	//	return nil, err
	//}
	//for rows.Next() {
	//	var (
	//		repoID string
	//	)
	//	if err := rows.Scan(&repoID); err != nil {
	//		log.Debug().Err(err).Msg("Failed to scan repoID.")
	//		return nil, err
	//	}
	//	// Fetch the repository details.
	//	repo, err := s.RepositoryStore.Get(ctx, id)
	//	if err != nil {
	//		log.Debug().Err(err).Msg("Failed to get repository details.")
	//		return nil, err
	//	}
	//	repositories = append(repositories, repo)
	//}

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

// Delete will remove a command.
func (s *CommandStore) Delete(ctx context.Context, id string) error {
	log := s.Logger.With().Str("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if _, err := s.Get(ctx, id); errors.Is(err, kerr.NotFound) {
			log.Debug().Err(err).Msg("Command not found")
			return fmt.Errorf("command to be deleted not found: %w", err)
		} else if err != nil {
			log.Debug().Err(err).Msg("Failed to get command.")
			return fmt.Errorf("failed get command: %w", err)
		}

		if commandTags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", commandsTable), id); err != nil {
			log.Debug().Err(err).Msg("Failed to delete command.")
			return &kerr.QueryError{
				Query: "delete id: " + id,
				Err:   fmt.Errorf("failed get command: %w", err),
			}
		} else if commandTags.RowsAffected() > 0 {
			// Make sure to only delete the relationship if the delete was successful.
			if err := s.RepositoryStore.DeleteAllRepositoryRelForCommand(ctx, id); err != nil {
				log.Debug().Err(err).Msg("Failed to delete repository relationship for command.")
				return &kerr.QueryError{
					Query: "delete id: " + id,
					Err:   fmt.Errorf("failed to delete repository relationship for command: %w", err),
				}
			}
		}
		return nil
	}

	return s.executeWithTransaction(ctx, log, f)
}

// Update modifies a command record.
func (s *CommandStore) Update(ctx context.Context, c *models.Command) (*models.Command, error) {
	var result *models.Command
	f := func(tx pgx.Tx) error {
		// Prevent updating the ID and the creation timestamp.
		// construct update statement:
		commandTags, err := tx.Exec(ctx, fmt.Sprintf("update %s set name = $1, enabled = $2, schedule = $3, filename = $4, location = $5, hash = $6", commandsTable),
			c.Name, c.Enabled, c.Schedule, c.Filename, c.Location, c.Hash)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   fmt.Errorf("failed to update: %w", err),
			}
		}
		if commandTags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   errors.New("no rows were affected"),
			}
		}
		result, err = s.Get(ctx, c.ID)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   errors.New("no rows were affected"),
			}
		return nil
	}
	return nil, nil
}

// List gets all the command records.
func (s *CommandStore) List(ctx context.Context) (*[]models.Command, error) {
	return nil, nil
}

// Takes a query and executes it inside a transaction.
func (s *CommandStore) executeWithTransaction(ctx context.Context, log zerolog.Logger, f func(tx pgx.Tx) error) error {
	conn, err := s.connect()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to connect to database.")
		return fmt.Errorf("database connection error: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, timeoutForTransactions)
	defer cancel()

	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Debug().Err(err).Msg("Failed to close connection.")
		}
	}()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to begin transaction.")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := f(tx); err != nil {
		log.Debug().Err(err).Msg("Failed to call method for the transaction.")
		return fmt.Errorf("failed to execute method: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction.")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// loader contains the error which will be shared by loadValue.
type loader struct {
	s   *CommandStore
	err error
}

// loadValue takes a value, tries to load its value from file and
// returns an error. If there was an error previously, this is a no-op.
func (l *loader) loadValue(v string) string {
	if l.err != nil {
		return ""
	}
	value, err := l.s.Converter.LoadValueFromFile(v)
	l.err = err
	return value
}

// connect will load all necessary values from secret and try to connect
// to a database.
func (s *CommandStore) connect() (*pgx.Conn, error) {
	l := &loader{
		s:   s,
		err: nil,
	}
	hostname := l.loadValue(s.Hostname)
	database := l.loadValue(s.Database)
	username := l.loadValue(s.Username)
	password := l.loadValue(s.Password)
	if l.err != nil {
		s.Logger.Error().Err(l.err).Msg("Failed to load database credentials.")
		return nil, fmt.Errorf("failed to load database credentials: %w", l.err)
	}
	url := fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s", hostname, database, username, password)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to connect to the database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return conn, nil
}
