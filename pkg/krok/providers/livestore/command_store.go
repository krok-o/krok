package livestore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

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
	conn, err := s.connect()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to connect to database.")
		return nil, fmt.Errorf("database connection error: %w", err)
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
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Debug().Err(err).Msg("Failed to rollback transaction.")
		}
	}()

	var (
		name      string
		commandID string
		schedule  string
		filename  string
		location  string
		hash      string
		enabled   bool
	)
	err = tx.QueryRow(ctx, "select name, id, schedule, filename, location, hash, enabled from commands where id = $1", id).
		Scan(&name, &commandID, &schedule, &filename, &location, &hash, &enabled)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, &errors.QueryError{
				Query: "select id: " + id,
				Err:   errors.NotFound,
			}
		}
		log.Debug().Err(err).Msg("Failed to query row.")
		return nil, &errors.QueryError{
			Query: "select id: " + id,
			Err:   err,
		}
	}

	repositories, err := s.RepositoryStore.GetRepositoriesForCommand(ctx, id)
	if err != nil {
		log.Debug().Err(err).Msg("GetRepositoriesForCommand failed")
		return nil, &errors.QueryError{
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

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction.")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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

// Delete will remove a command.
func (s *CommandStore) Delete(ctx context.Context, id string) error {
	return nil
}

// Update modifies a command record.
func (s *CommandStore) Update(ctx context.Context, c *models.Command) (*models.Command, error) {
	return nil, nil
}

// List gets all the command records.
func (s *CommandStore) List(ctx context.Context) (*[]models.Command, error) {
	return nil, nil
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
