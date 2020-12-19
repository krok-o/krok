package livestore

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	Connector       *Connector
}

// NewCommandStore creates a new CommandStore
func NewCommandStore(cfg Config, deps CommandDependencies) *CommandStore {
	cs := &CommandStore{Config: cfg, CommandDependencies: deps}
	// launch the cleanup routine.
	go cs.lockCleaner(context.TODO())
	return cs
}

func (s *CommandStore) lockCleaner(ctx context.Context) {
	interval := 1 * time.Minute
	for {
		// look for old entries.
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			s.Logger.Debug().Msg("Lock Cleaner cancelled.")
		}
	}
}

// Create creates a command record.
func (s *CommandStore) Create(ctx context.Context, c *models.Command) error {

	// duplicate key value violates unique constraint
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

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
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

	return s.Connector.ExecuteWithTransaction(ctx, log, f)
}

// Update modifies a command record.
func (s *CommandStore) Update(ctx context.Context, c *models.Command) (*models.Command, error) {
	log := s.Logger.With().Str("id", c.ID).Str("name", c.Name).Logger()
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
func (s *CommandStore) List(ctx context.Context) (*[]models.Command, error) {
	return nil, nil
}
