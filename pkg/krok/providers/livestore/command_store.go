package livestore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/krok-o/krok/pkg/models"
)

// CommandStore is a postgres based store for commands.
type CommandStore struct {
	Config
	Dependencies
}

// NewCommandStore creates a new CommandStore
func NewCommandStore(cfg Config, deps Dependencies) *CommandStore {
	return &CommandStore{Config: cfg, Dependencies: deps}
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
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeoutForTransactions)
	defer cancel()

	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Debug().Err(err).Msg("Failed to close connection.")
		}
	}()
	var (
		storedHandle string
		commands     string
	)
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to begin transaction.")
		return nil, err
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
		// Relationship manager will get all the repositories which belong to this command.
		repositories []models.Repository
		filename     string
		location     string
		hash         string
		enabled      bool
	)
	err = tx.QueryRow(ctx, "select name, id, schedule, filename, location, hash, enabled from commands where id = $1", id).Scan(&storedHandle, &commands)
	if err != nil {
		if err.Error() == "no rows in result set" {

			return nil, nil
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		log.Debug().Err(err).Msg("Failed to commit transaction.")
		return nil, err
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
		return nil, l.err
	}
	url := fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s", hostname, database, username, password)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to connect to the database")
		return nil, err
	}
	return conn, nil
}
