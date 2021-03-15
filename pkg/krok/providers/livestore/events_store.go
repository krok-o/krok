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

const (
	eventsStoreTable = "events"
	defaultPageSize  = 10
)

// eventsStore is a postgres based store for eventStorer.
type eventsStore struct {
	EventsStoreDependencies
}

// EventsStoreDependencies eventsStoreStore specific dependencies.
type EventsStoreDependencies struct {
	Dependencies
	Connector *Connector
}

// NewEventsStorer creates a new eventsStore
func NewEventsStorer(deps EventsStoreDependencies) *eventsStore {
	return &eventsStore{EventsStoreDependencies: deps}
}

var _ providers.EventsStorer = &eventsStore{}

// Create an event.
func (e *eventsStore) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	log := e.Logger.With().Str("event_id", event.EventID).Int("repository_id", event.RepositoryID).Logger()
	var returnID int
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("insert into %s(event_id, created_at, repository_id, payload) values($1, $2, $3, $4) returning id", eventsStoreTable)
		row := tx.QueryRow(ctx, query,
			event.EventID,
			event.CreateAt,
			event.RepositoryID,
			event.Payload)
		if err := row.Scan(&returnID); err != nil {
			log.Debug().Err(err).Str("query", query).Msg("Failed to scan row.")
			return &kerr.QueryError{
				Err:   err,
				Query: query,
			}
		}
		return nil
	}
	if err := e.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, err
	}
	event.ID = returnID
	return event, nil
}

// ListEventsForRepository gets paginated list of events for a repository.
// It does not return the command runs and the payloads to prevent potentially big chunks of transfer data.
// To get those, one must do a Get.
func (e *eventsStore) ListEventsForRepository(ctx context.Context, repoID int, options models.ListOptions) ([]*models.Event, error) {
	log := e.Logger.With().Str("func", "List").Int("repo_id", repoID).Logger()
	if options.PageSize == 0 {
		options.PageSize = defaultPageSize
	}
	// Select all commands.
	result := make([]*models.Event, 0)
	f := func(tx pgx.Tx) error {
		sql := fmt.Sprintf("select id, event_id, repository_id, created_at from %s limit %d offset %d where repository_id = $1", eventsStoreTable, options.PageSize, options.PageSize*options.Page)
		args := []interface{}{
			repoID,
		}
		if options.StartingDate != nil && options.EndDate != nil {
			sql += " where created_at between $1 and $2"
			args = append(args, options.StartingDate, options.EndDate)
		}
		rows, err := tx.Query(ctx, sql, args)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select all events",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query events.")
			return &kerr.QueryError{
				Query: "select all events",
				Err:   fmt.Errorf("failed to list all events: %w", err),
			}
		}

		for rows.Next() {
			var (
				storedID           int
				storedEventID      string
				storedRepositoryID int
				storedCreatedAt    time.Time
			)
			if err := rows.Scan(&storedID, &storedEventID, &storedRepositoryID, &storedCreatedAt); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all events",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			event := &models.Event{
				ID:           storedID,
				EventID:      storedEventID,
				CreateAt:     storedCreatedAt,
				RepositoryID: storedRepositoryID,
			}
			result = append(result, event)
		}
		return nil
	}
	if err := e.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all events: %w", err)
	}
	return result, nil
}

// GetEvent retrieves details about an event.
func (e *eventsStore) GetEvent(ctx context.Context, id int) (*models.Event, error) {
	// Select all commands from a run belonging to this event and get the command.
	log := e.Logger.With().Int("id", id).Logger()
	// Get all data from the repository table.
	result := &models.Event{}
	f := func(tx pgx.Tx) error {
		var (
			storedId, repoID int
			eventID, payload string
			createdAt        time.Time
		)
		if err := tx.QueryRow(ctx, fmt.Sprintf("select id, event_id, created_at, repository_id, payload from %s where id=$1", eventsStoreTable), id).Scan(&storedId, &eventID, &createdAt, &repoID, &payload); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select id in events",
					Err:   kerr.ErrNotFound,
				}
			}
			return &kerr.QueryError{
				Query: "select id",
				Err:   fmt.Errorf("failed to scan: %w", err),
			}
		}
		result.ID = storedId
		result.RepositoryID = repoID
		result.EventID = eventID
		result.CreateAt = createdAt
		result.Payload = payload
		return nil
	}
	if err := e.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute Get: %w", err)
	}

	commands, err := e.getCommandRunForEvent(ctx, result.ID)
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("Get failed to get event command runs.")
		return nil, err
	}
	result.CommandRuns = commands
	return result, nil
}

// getCommandRunForEvent returns a list of command runs for an event.
func (e *eventsStore) getCommandRunForEvent(ctx context.Context, id int) ([]*models.CommandRun, error) {
	log := e.Logger.With().Int("id", id).Logger()

	// Select the related commands.
	result := make([]*models.CommandRun, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select id, command_name, event_id, status, outcome, created_at from %s where event_id = $1", commandRunTable), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select id",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query for command runs.")
			return &kerr.QueryError{
				Query: "select command runs for event",
				Err:   fmt.Errorf("failed to query command run table: %w", err),
			}
		}

		for rows.Next() {
			var (
				storedID          int
				storedCommandName string
				storedEventID     int
				storedStatus      string
				storedOutcome     string
				storedCreatedAt   time.Time
			)
			if err := rows.Scan(&storedID, &storedCommandName, &storedEventID, &storedStatus, &storedOutcome, &storedCreatedAt); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select id from command_runs",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			command := &models.CommandRun{
				ID:          storedID,
				EventID:     storedEventID,
				CommandName: storedCommandName,
				Status:      storedStatus,
				Outcome:     storedOutcome,
				CreateAt:    storedCreatedAt,
			}
			result = append(result, command)
		}
		return nil
	}
	if err := e.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute getCommandRunForEvent: %w", err)
	}
	return result, nil
}
