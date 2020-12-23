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

// ApiKeysStore is a postgres based store for ApiKeys.
type ApiKeysStore struct {
	ApiKeysDependencies
	Config
}

// ApiKeysDependencies ApiKeys specific dependencies.
type ApiKeysDependencies struct {
	Dependencies
	Connector *Connector
}

// NewApiKeysStore creates a new ApiKeysStore
func NewApiKeysStore(cfg Config, deps ApiKeysDependencies) *ApiKeysStore {
	return &ApiKeysStore{Config: cfg, ApiKeysDependencies: deps}
}

var _ providers.ApiKeys = &ApiKeysStore{}

// Create an apikey.
func (a *ApiKeysStore) Create(ctx context.Context, key *models.ApiKey) (*models.ApiKey, error) {
	log := a.Logger.With().Str("name", key.Name).Str("id", key.ApiKeyID).Logger()
	var returnId int
	f := func(tx pgx.Tx) error {
		if row, err := tx.Query(ctx, fmt.Sprintf("insert into %s(name, api_key_id, api_key_secret, user_id, ttl) values($1, $2, $3, $4, $5) returning id", apiKeysTable),
			key.Name,
			key.ApiKeyID,
			key.ApiKeySecret,
			key.UserID,
			key.TTL); err != nil {
			log.Debug().Err(err).Msg("Failed to create api key.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into apikeys",
			}
		} else if row.CommandTag().RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.NoRowsAffected,
				Query: "insert into apikeys",
			}
		} else {
			if err := row.Scan(&returnId); err != nil {
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
	result, err := a.Get(ctx, returnId)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get created api key")
		return nil, err
	}
	return result, nil
}

// Delete an apikey.
func (a *ApiKeysStore) Delete(ctx context.Context, id int) error {
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
				Err:   kerr.NoRowsAffected,
				Query: "delete from apikeys",
			}
		}
		return nil
	}

	return a.Connector.ExecuteWithTransaction(ctx, log, f)
}

// List will list all apikeys for a user.
func (a *ApiKeysStore) List(ctx context.Context, userID int) ([]*models.ApiKey, error) {
	log := a.Logger.With().Str("func", "ListApiKeys").Logger()
	// Select all users.
	result := make([]*models.ApiKey, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select id, name, api_key_id, ttl from %s "+
			"where user_id = $1", apiKeysTable), userID)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Query: "select all apikeys",
					Err:   kerr.NotFound,
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
				apiKeyId string
				ttl      time.Time
			)
			if err := rows.Scan(&id, &name, &apiKeyId, &ttl); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all users",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			key := &models.ApiKey{
				ID:       id,
				Name:     name,
				TTL:      ttl,
				ApiKeyID: apiKeyId,
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
func (a *ApiKeysStore) Get(ctx context.Context, id int) (*models.ApiKey, error) {
	log := a.Logger.With().Int("id", id).Logger()
	var (
		storedID       int
		storedName     string
		storedApiKeyID string
		storedUserID   string
		storedTTL      time.Time
	)
	f := func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "select id, name, api_key_id, user_id, ttl from %s where id = $1", id).
			Scan(&storedID, &storedName, &storedApiKeyID, &storedUserID, &storedTTL)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Err:   kerr.NotFound,
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

	return &models.ApiKey{
		ID:       storedID,
		Name:     storedName,
		UserID:   storedUserID,
		ApiKeyID: storedApiKeyID,
		TTL:      storedTTL,
	}, nil
}
