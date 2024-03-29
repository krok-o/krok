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

// EventsStore is a postgres based store for eventStorer.
type EventsStore struct {
	EventsStoreDependencies
}

// EventsStoreDependencies eventsStoreStore specific dependencies.
type EventsStoreDependencies struct {
	Dependencies
	Connector *Connector
}

// NewEventsStorer creates a new eventsStore
func NewEventsStorer(deps EventsStoreDependencies) *EventsStore {
	return &EventsStore{EventsStoreDependencies: deps}
}

var _ providers.EventsStorer = &EventsStore{}

// Create an event.
func (e *EventsStore) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	log := e.Logger.With().Str("event_id", event.EventID).Int("repository_id", event.RepositoryID).Logger()
	var returnID int
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("insert into %s(event_id, created_at, repository_id, payload, vcs, event_type) values($1, $2, $3, $4, $5, $6) returning id", eventsStoreTable)
		row := tx.QueryRow(ctx, query,
			event.EventID,
			event.CreateAt,
			event.RepositoryID,
			event.Payload,
			event.VCS,
			event.EventType)
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
func (e *EventsStore) ListEventsForRepository(ctx context.Context, repoID int, options *models.ListOptions) ([]*models.Event, error) {
	log := e.Logger.With().Str("func", "List").Int("repo_id", repoID).Logger()
	if options == nil {
		options = &models.ListOptions{
			PageSize: defaultPageSize,
		}
	}
	if options.PageSize == 0 {
		options.PageSize = defaultPageSize
	}
	// Select all commands.
	result := make([]*models.Event, 0)
	f := func(tx pgx.Tx) error {
		sql := fmt.Sprintf("select id, event_id, repository_id, created_at, vcs, event_type from %s where repository_id = $1", eventsStoreTable)
		args := []interface{}{
			repoID,
		}
		if options.StartingDate != nil && options.EndDate != nil {
			sql += " and created_at between $2 and $3 limit $4 offset $5"
			args = append(args, *options.StartingDate, *options.EndDate, options.PageSize, options.PageSize*options.Page)
		} else {
			sql += " limit $2 offset $3"
			args = append(args, options.PageSize, options.PageSize*options.Page)
		}
		rows, err := tx.Query(ctx, sql, args...)
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
				storedVCS          int
				storedEventType    string
			)
			if err := rows.Scan(&storedID, &storedEventID, &storedRepositoryID, &storedCreatedAt, &storedVCS, &storedEventType); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: sql,
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			// todo: should add a list of command run event ids...
			event := &models.Event{
				ID:           storedID,
				EventID:      storedEventID,
				CreateAt:     storedCreatedAt,
				RepositoryID: storedRepositoryID,
				VCS:          storedVCS,
				EventType:    storedEventType,
			}
			result = append(result, event)
		}
		return nil
	}
	// TODO: Add command runs
	if err := e.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all events: %w", err)
	}
	return result, nil
}

// GetEvent retrieves details about an event.
func (e *EventsStore) GetEvent(ctx context.Context, id int) (*models.Event, error) {
	// Select all commands from a run belonging to this event and get the command.
	log := e.Logger.With().Int("id", id).Logger()
	// Get all data from the repository table.
	result := &models.Event{}
	f := func(tx pgx.Tx) error {
		var (
			storedID, repoID, vcs       int
			eventID, payload, eventType string
			createdAt                   time.Time
		)
		if err := tx.QueryRow(ctx, fmt.Sprintf("select id, event_id, created_at, repository_id, payload, vcs, event_type from %s where id=$1", eventsStoreTable), id).Scan(&storedID, &eventID, &createdAt, &repoID, &payload, &vcs, &eventType); err != nil {
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
		result.ID = storedID
		result.RepositoryID = repoID
		result.EventID = eventID
		result.CreateAt = createdAt
		result.Payload = payload
		result.VCS = vcs
		result.EventType = eventType
		return nil
	}
	if err := e.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute Get: %w", err)
	}

	commands, err := e.getCommandRunsForEvent(ctx, result.ID)
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("Get failed to get event command runs.")
		return nil, err
	}
	result.CommandRuns = commands
	return result, nil
}

// getCommandRunsForEvent returns a list of command runs for an event.
func (e *EventsStore) getCommandRunsForEvent(ctx context.Context, id int) ([]*models.CommandRun, error) {
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
		return nil, fmt.Errorf("failed to execute getCommandRunsForEvent: %w", err)
	}
	return result, nil
}
