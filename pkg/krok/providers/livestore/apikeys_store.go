package livestore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	apiKeysTable = "apikeys"
)

// APIKeysStore is a postgres based store for APIKeys.
type APIKeysStore struct {
	APIKeysDependencies
	Config
}

// APIKeysDependencies APIKeys specific dependencies.
type APIKeysDependencies struct {
	Dependencies
	Connector *Connector
}

// NewAPIKeysStore creates a new APIKeysStore
func NewAPIKeysStore(cfg Config, deps APIKeysDependencies) *APIKeysStore {
	return &APIKeysStore{Config: cfg, APIKeysDependencies: deps}
}

var _ providers.APIKeys = &APIKeysStore{}

// Create an apikey.
func (a *APIKeysStore) Create(ctx context.Context, key *models.APIKey) (*models.APIKey, error) {
	log := a.Logger.With().Str("name", key.Name).Str("id", key.APIKeyID).Logger()
	var returnID int
	f := func(tx pgx.Tx) error {
		if row, err := tx.Query(ctx, fmt.Sprintf("insert into %s(name, api_key_id, api_key_secret, user_id, ttl) values($1, $2, $3, $4, $5) returning id", apiKeysTable),
			key.Name,
			key.APIKeyID,
			key.APIKeySecret,
			key.UserID,
			key.TTL); err != nil {
			log.Debug().Err(err).Msg("Failed to create api key.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into apikeys",
			}
		} else if row.CommandTag().RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into apikeys",
			}
		} else {
			if err := row.Scan(&returnID); err != nil {
				return &kerr.QueryError{
					Err:   err,
					Query: "scanning row",
				}
			}
		}
		return nil
	}

	if err := a.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, err
	}
	result, err := a.Get(ctx, returnID)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get created api key")
		return nil, err
	}
	return result, nil
}

// Delete an apikey.
func (a *APIKeysStore) Delete(ctx context.Context, id int) error {
	log := a.Logger.With().Int("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", apiKeysTable),
			id); err != nil {
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from apikeys",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from apikeys",
			}
		}
		return nil
	}

	return a.Connector.ExecuteWithTransaction(ctx, log, f)
}

// List will list all apikeys for a user.
func (a *APIKeysStore) List(ctx context.Context, userID int) ([]*models.APIKey, error) {
	log := a.Logger.With().Str("func", "ListApiKeys").Logger()
	// Select all users.
	result := make([]*models.APIKey, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select id, name, api_key_id, ttl from %s "+
			"where user_id = $1", apiKeysTable), userID)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Query: "select all apikeys",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query apikeys.")
			return &kerr.QueryError{
				Query: "select all apikeys",
				Err:   fmt.Errorf("failed to list all apikeys: %w", err),
			}
		}

		for rows.Next() {
			var (
				id       int
				name     string
				apiKeyID string
				ttl      time.Time
			)
			if err := rows.Scan(&id, &name, &apiKeyID, &ttl); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all users",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			key := &models.APIKey{
				ID:       id,
				Name:     name,
				TTL:      ttl,
				APIKeyID: apiKeyID,
			}
			result = append(result, key)
		}
		return nil
	}
	if err := a.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all api keys: %w", err)
	}
	return result, nil
}

// Get an apikey.
func (a *APIKeysStore) Get(ctx context.Context, id int) (*models.APIKey, error) {
	log := a.Logger.With().Int("id", id).Logger()
	var (
		storedID       int
		storedName     string
		storedAPIKeyID string
		storedUserID   string
		storedTTL      time.Time
	)
	f := func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "select id, name, api_key_id, user_id, ttl from %s where id = $1", id).
			Scan(&storedID, &storedName, &storedAPIKeyID, &storedUserID, &storedTTL)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Err:   kerr.ErrNotFound,
					Query: "select apikey",
				}
			}
			return err
		}
		return nil
	}
	if err := a.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to run transaction for get api key.")
		return nil, err
	}

	return &models.APIKey{
		ID:       storedID,
		Name:     storedName,
		UserID:   storedUserID,
		APIKeyID: storedAPIKeyID,
		TTL:      storedTTL,
	}, nil
}
