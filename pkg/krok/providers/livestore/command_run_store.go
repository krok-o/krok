package livestore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	commandRunTable = "command_run"
)

// CommandRunStore is a postgres based store for CommandRunStore.
type CommandRunStore struct {
	CommandRunDependencies
}

// CommandRunDependencies CommandRunStore specific dependencies.
type CommandRunDependencies struct {
	Dependencies
	Connector *Connector
}

// NewCommandRunStore creates a new CommandRunStore
func NewCommandRunStore(deps CommandRunDependencies) *CommandRunStore {
	return &CommandRunStore{CommandRunDependencies: deps}
}

var _ providers.CommandRunStorer = &CommandRunStore{}

// CreateRun creates a new Command run entry.
func (a *CommandRunStore) CreateRun(ctx context.Context, cmdRun *models.CommandRun) (*models.CommandRun, error) {
	log := a.Logger.With().Int("event_id", cmdRun.EventID).Logger()
	var returnID int
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("insert into %s(event_id, command_name, status, outcome, created_at) values($1, $2, $3, $4, $5) returning id", commandRunTable)
		row := tx.QueryRow(ctx, query,
			cmdRun.EventID,
			cmdRun.CommandName,
			cmdRun.Status,
			cmdRun.Outcome,
			cmdRun.CreateAt)
		if err := row.Scan(&returnID); err != nil {
			log.Debug().Err(err).Str("query", query).Msg("Failed to scan row.")
			return &kerr.QueryError{
				Err:   err,
				Query: query,
			}
		}
		return nil
	}
	if err := a.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, err
	}
	cmdRun.ID = returnID
	return cmdRun, nil
}

// UpdateRunStatus takes an id a status and an outcome and updates a run with it. This is a convenient method around
// update which breaks the normal update flow so it's easy to call by the providers.
func (a *CommandRunStore) UpdateRunStatus(ctx context.Context, id int, status string, outcome string) error {
	log := a.Logger.With().Int("id", id).Str("status", status).Logger()
	f := func(tx pgx.Tx) error {
		// Prevent updating the ID and the creation timestamp.
		// construct update statement:
		tags, err := tx.Exec(ctx, fmt.Sprintf("update %s set status = $1, outcome = $2 where id = $3", commandRunTable),
			status, outcome, id)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + status,
				Err:   fmt.Errorf("failed to update: %w", err),
			}
		}
		if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Query: "update :" + status,
				Err:   kerr.ErrNoRowsAffected,
			}
		}
		return nil
	}
	if err := a.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return fmt.Errorf("failed to execute update in transaction: %w", err)
	}
	return nil
}
