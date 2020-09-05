package livestore

import (
	"context"
	"fmt"

	"github.com/krok-o/krok/pkg/models"

	"github.com/jackc/pgx/v4"
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

// Create
func (s *CommandStore) Create(c *models.Command) error {
	return nil
}

// Get
func (s *CommandStore) Get(ctx context.Context, id string) (*models.Command, error) {
	conn, err := s.connect()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeoutForTransactions)
	defer cancel()

	defer conn.Close(ctx)
	var (
		storedHandle string
		commands     string
	)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

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
	return nil, nil
}

// Delete
func (s *CommandStore) Delete(id string) error {
	return nil
}

// Update
func (s *CommandStore) Update(c *models.Command) (*models.Command, error) {
	return nil, nil
}

// loader contains the error which will be shared by loadValue.
type loader struct {
	s   *CommandStore
	err error
}

// loadValues takes a value, tries to load its value from file and
// returns an error. If there was an error previously, this is a no-op.
func (l *loader) loadValues(v string) string {
	if l.err != nil {
		return ""
	}
	value, err := l.s.Converter.LoadValueFromFile(v)
	l.err = err
	return value
}

func (s *CommandStore) connect() (*pgx.Conn, error) {
	l := &loader{
		s:   s,
		err: nil,
	}
	hostname := l.loadValues(s.Hostname)
	database := l.loadValues(s.Database)
	username := l.loadValues(s.Username)
	password := l.loadValues(s.Password)
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
