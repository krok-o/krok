package livestore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	platformsTable = "platforms"
)

// PlatformStore is a postgres based store for platforms.
type PlatformStore struct {
	PlatformDependencies
}

// PlatformDependencies platform specific dependencies.
type PlatformDependencies struct {
	Dependencies
	Connector *Connector
}

// NewPlatformStore creates a new PlatformStore
func NewPlatformStore(deps PlatformDependencies) *PlatformStore {
	return &PlatformStore{PlatformDependencies: deps}
}

var _ providers.PlatformStorer = &PlatformStore{}

// Create will add a new platform to Krok.
func (p *PlatformStore) Create(ctx context.Context, platform *models.Platform) (*models.Platform, error) {
	log := p.Logger.With().Str("name", platform.Name).Logger()
	var returnID int
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("insert into %s (name, enabled) values($1, $2) returning id", platformsTable)
		row := tx.QueryRow(ctx, query,
			platform.Name,
			platform.Enabled)

		if err := row.Scan(&returnID); err != nil {
			log.Debug().Err(err).Str("query", query).Msg("Failed to scan row.")
			return &kerr.QueryError{
				Err:   err,
				Query: query,
			}
		}
		return nil
	}

	if err := p.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, err
	}
	platform.ID = returnID
	return platform, nil
}

// Update a platform to set it to enabled or disabled.
// Commands cannot be executed on a disabled platform.
func (p *PlatformStore) Update(ctx context.Context, platform *models.Platform) (*models.Platform, error) {
	log := p.Logger.With().Int("id", platform.ID).Str("name", platform.Name).Bool("enabled", platform.Enabled).Logger()
	f := func(tx pgx.Tx) error {
		// Prevent updating the ID and the creation timestamp.
		// construct update statement:
		tags, err := tx.Exec(ctx, fmt.Sprintf("update %s set enabled = $1", platformsTable), platform.Enabled)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + platform.Name,
				Err:   fmt.Errorf("failed to update: %w", err),
			}
		}
		if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Query: "update :" + platform.Name,
				Err:   kerr.ErrNoRowsAffected,
			}
		}
		return nil
	}
	if err := p.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, fmt.Errorf("failed to execute update in transaction: %w", err)
	}
	result, err := p.Get(ctx, platform.ID)
	if err != nil {
		return nil, &kerr.QueryError{
			Query: "update :" + platform.Name,
			Err:   errors.New("failed to get updated platform"),
		}
	}
	return result, nil
}

// Delete removes a supported platform from krok. Removing support for a platform
// will stop commands from running for that platform.
func (p *PlatformStore) Delete(ctx context.Context, id int) error {
	log := p.Logger.With().Int("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", platformsTable),
			id); err != nil {
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from platforms",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from platforms",
			}
		}
		return nil
	}

	return p.Connector.ExecuteWithTransaction(ctx, log, f)
}

// List will list all supported platforms.
func (p *PlatformStore) List(ctx context.Context, enabled *bool) ([]*models.Platform, error) {
	log := p.Logger.With().Str("func", "List").Logger()
	// Select all platforms.
	result := make([]*models.Platform, 0)
	f := func(tx pgx.Tx) error {

		where := ""
		if enabled != nil {
			// avoid having to $2 it and sql injections for the interpreter.
			if *enabled {
				where += "where enabled = true"
			} else {
				where += "where enabled = false"
			}
		}

		rows, err := tx.Query(ctx, fmt.Sprintf("select id, name, enabled from %s "+where, platformsTable))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select all platforms",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query platform.")
			return &kerr.QueryError{
				Query: "select all platforms",
				Err:   fmt.Errorf("failed to list all platforms: %w", err),
			}
		}

		for rows.Next() {
			var (
				id      int
				name    string
				enabled bool
			)
			if err := rows.Scan(&id, &name, &enabled); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all platorms",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			platform := &models.Platform{
				ID:      id,
				Name:    name,
				Enabled: enabled,
			}
			result = append(result, platform)
		}
		return nil
	}
	if err := p.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all platforms: %w", err)
	}
	return result, nil
}

func (a *PlatformStore) Get(ctx context.Context, id int) (*models.Platform, error) {
	log := a.Logger.With().Int("id", id).Logger()
	var (
		storedID      int
		storedName    string
		storedEnabled bool
	)
	f := func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, fmt.Sprintf("select id, name, enabled from %s where id = $1", platformsTable), id).
			Scan(&storedID, &storedName, &storedEnabled)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Err:   kerr.ErrNotFound,
					Query: "select platform",
				}
			}
			return err
		}
		return nil
	}
	if err := a.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to run transaction for get platform.")
		return nil, err
	}

	return &models.Platform{
		ID:      storedID,
		Name:    storedName,
		Enabled: storedEnabled,
	}, nil
}

func (a *PlatformStore) GetByName(ctx context.Context, name string) (*models.Platform, error) {
	log := a.Logger.With().Str("name", name).Logger()
	var (
		storedID      int
		storedName    string
		storedEnabled bool
	)
	f := func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, fmt.Sprintf("select id, name, enabled from %s where name = $1", platformsTable), name).
			Scan(&storedID, &storedName, &storedEnabled)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Err:   kerr.ErrNotFound,
					Query: "select platform",
				}
			}
			return err
		}
		return nil
	}
	if err := a.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to run transaction for get platform.")
		return nil, err
	}

	return &models.Platform{
		ID:      storedID,
		Name:    storedName,
		Enabled: storedEnabled,
	}, nil
}
