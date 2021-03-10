package livestore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// UserStore is a postgres based store for users.
type UserStore struct {
	UserDependencies
}

// UserDependencies user specific dependencies.
type UserDependencies struct {
	Dependencies
	Connector *Connector
	APIKeys   providers.APIKeysStorer
}

// NewUserStore creates a new UserStore
func NewUserStore(deps UserDependencies) *UserStore {
	return &UserStore{UserDependencies: deps}
}

// Create saves a user in the db.
func (s *UserStore) Create(ctx context.Context, user *models.User) (*models.User, error) {
	log := s.Logger.With().Str("display_name", user.DisplayName).Str("email", user.Email).Logger()

	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, "insert into users(email, last_login, display_name) values($1, $2, $3)",
			user.Email,
			time.Now(),
			user.DisplayName); err != nil {
			log.Debug().Err(err).Msg("Failed to create user.")
			return &kerr.QueryError{
				Err:   fmt.Errorf("failed create user: %w", err),
				Query: "insert into users",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into users",
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, s.Logger, f); err != nil {
		log.Debug().Err(err).Msg("Failed to create user.")
		return nil, err
	}

	// Get the newly created user and return it.
	user, err := s.GetByEmail(ctx, user.Email)
	if err != nil {
		log.Debug().Err(err).Msg("Created user not found.")
		return nil, err
	}

	return user, nil
}

// Delete deletes a user from the db.
func (s *UserStore) Delete(ctx context.Context, id int) error {
	log := s.Logger.With().Int("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, "delete from users where id = $1",
			id); err != nil {
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from users",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from users",
			}
		}
		return nil
	}

	return s.Connector.ExecuteWithTransaction(ctx, log, f)
}

// Get retrieves a user.
func (s *UserStore) Get(ctx context.Context, id int) (*models.User, error) {
	log := s.Logger.With().Str("func", "GetByID").Logger()
	return s.getByX(ctx, log, "id", id)
}

// GetByEmail retrieves a user by its email address.
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	log := s.Logger.With().Str("func", "GetByEmail").Logger()
	return s.getByX(ctx, log, "email", email)
}

// getByX abstracts the ability to define concrete fields to retrieve users by.
// i.e: email, id, lastLogin...
func (s *UserStore) getByX(ctx context.Context, log zerolog.Logger, field string, value interface{}) (*models.User, error) {
	log = log.With().Str("fields", field).Interface("value", value).Logger()
	var (
		storedEmail       string
		storedDisplayName string
		storedID          int
		storedLastLogin   time.Time
		storedToken       sql.NullString
	)
	f := func(tx pgx.Tx) error {
		withWhere := fmt.Sprintf("select id, email, display_name, last_login, token from users where %s = $1", field)
		err := tx.QueryRow(ctx, withWhere, value).
			Scan(&storedID, &storedEmail, &storedDisplayName, &storedLastLogin, &storedToken)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Err:   kerr.ErrNotFound,
					Query: withWhere,
				}
			}
			log.Debug().Err(err).Msg("Failed to run select for users.")
			return err
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to run transaction for GetByField.")
		return nil, fmt.Errorf("failed to run user transaction: %w", err)
	}

	apiKeys, err := s.APIKeys.List(ctx, storedID)
	// if we didn't find any, that's fine.
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("Failed to get api keys for user.")
		return nil, fmt.Errorf("failed to get api keys for user: %w", err)
	}

	return &models.User{
		Email:       storedEmail,
		DisplayName: storedDisplayName,
		ID:          storedID,
		APIKeys:     apiKeys,
		LastLogin:   storedLastLogin,
		Token:       storedToken.String,
	}, nil
}

// Update updates a user with a given email address.
func (s *UserStore) Update(ctx context.Context, user *models.User) (*models.User, error) {
	log := s.Logger.With().Int("id", user.ID).Str("email", user.Email).Logger()

	f := func(tx pgx.Tx) error {
		query := "update users set display_name=$1, token=$2 where id=$3"
		if tags, err := tx.Exec(ctx, query, user.DisplayName, user.Token, user.ID); err != nil {
			log.Debug().Err(err).Msg("Failed to update user.")
			return &kerr.QueryError{
				Err:   err,
				Query: "update users",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "update users",
			}
		}
		return nil
	}

	// update, get, return
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, err
	}
	newUser, err := s.Get(ctx, user.ID)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get updated user.")
	}

	return newUser, err
}

// List all users. This will not return api keys. For those we need an explicit get.
func (s *UserStore) List(ctx context.Context) ([]*models.User, error) {
	log := s.Logger.With().Str("func", "List").Logger()
	// Select all users.
	result := make([]*models.User, 0)
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, "select id, email, display_name, last_login from users")
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select all users",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query users.")
			return &kerr.QueryError{
				Query: "select all users",
				Err:   fmt.Errorf("failed to list all users: %w", err),
			}
		}

		for rows.Next() {
			var (
				id          int
				email       string
				displayName string
				lastLogin   time.Time
			)
			if err := rows.Scan(&id, &email, &displayName, &lastLogin); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all users",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			user := &models.User{
				DisplayName: displayName,
				ID:          id,
				Email:       email,
				LastLogin:   lastLogin,
			}
			result = append(result, user)
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all users: %w", err)
	}
	return result, nil
}

// GetByToken retrieves a user by personal access token.
func (s *UserStore) GetByToken(ctx context.Context, token string) (*models.User, error) {
	return s.getByX(ctx, s.Logger, "token", token)
}
