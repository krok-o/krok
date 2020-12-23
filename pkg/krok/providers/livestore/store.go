package livestore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const timeoutForTransactions = 1 * time.Minute

// Config has the configuration options for the store
type Config struct {
	Hostname string
	Database string
	Username string
	Password string
}

// Dependencies defines the dependencies of this command store
type Dependencies struct {
	Logger    zerolog.Logger
	Converter providers.EnvironmentConverter
}

// Connector defines the connector's structure.
type Connector struct {
	Config
	Dependencies
}

// NewDatabaseConnector defines some common functionality between the database dependent
// providers.
func NewDatabaseConnector(cfg Config, deps Dependencies) *Connector {
	return &Connector{
		Config:       cfg,
		Dependencies: deps,
	}
}

// ExecuteWithTransaction takes a query and executes it inside a transaction.
func (s *Connector) ExecuteWithTransaction(ctx context.Context, log zerolog.Logger, f func(tx pgx.Tx) error) error {
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
	s   *Connector
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
func (s *Connector) connect() (*pgx.Conn, error) {
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
