package livestore

import (
	"context"
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
	Config
}

// UserDependencies user specific dependencies.
type UserDependencies struct {
	Dependencies
	Connector *Connector
	ApiKeys   providers.ApiKeys
}

// NewUserStore creates a new UserStore
func NewUserStore(cfg Config, deps UserDependencies) *UserStore {
	return &UserStore{Config: cfg, UserDependencies: deps}
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
				Err:   err,
				Query: "insert into users",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.NoRowsAffected,
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
func (s *UserStore) Delete(ctx context.Context, id string) error {
	log := s.Logger.With().Str("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, "delete from users where id = $1",
			id); err != nil {
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from users",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.NoRowsAffected,
				Query: "delete from users",
			}
		}
		return nil
	}

	return s.Connector.ExecuteWithTransaction(ctx, log, f)
}

// Get retrieves a user.
func (s *UserStore) Get(ctx context.Context, id string) (*models.User, error) {
	log := s.Logger.With().Str("func", "GetByID").Logger()
	return s.getByX(ctx, log, "id", id)
}

// GetByEmail retrieves a user by its email address.
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	log := s.Logger.With().Str("func", "GetByEmail").Logger()
	return s.getByX(ctx, log, "email", email)
}

func (s *UserStore) getByX(ctx context.Context, log zerolog.Logger, field string, value interface{}) (*models.User, error) {
	log = log.With().Str("fields", field).Interface("value", value).Logger()
	var (
		storedEmail       string
		storedDisplayName string
		storedID          string
		storedLastLogin   time.Time
	)
	f := func(tx pgx.Tx) error {
		withWhere := fmt.Sprintf("select id, email, display_name last_login from users where %s = $1", field)
		err := tx.QueryRow(ctx, withWhere, value).
			Scan(&storedID, &storedEmail, &storedDisplayName, &storedLastLogin)
		if err != nil {
			if err.Error() == "no rows in result set" {
				return &kerr.QueryError{
					Err:   kerr.NotFound,
					Query: "select user",
				}
			}
			return err
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to run transaction for GetByField.")
		return nil, err
	}

	apiKeys, err := s.ApiKeys.List(ctx, storedID)
	// if we didn't find any, that's fine.
	if !errors.Is(err, kerr.NotFound) {
		log.Debug().Err(err).Msg("Failed to get api keys for user.")
		return nil, fmt.Errorf("failed to get api keys for user: %w", err)
	}

	return &models.User{
		Email:       storedEmail,
		DisplayName: storedDisplayName,
		ID:          storedID,
		ApiKeys:     apiKeys,
		LastLogin:   storedLastLogin}, nil
}

// Update updates a user with a given email address.
//func (s *UserStore) Update(ctx context.Context, email string, newUser models.User) error {
//	if _, err = tx.Exec(ctx, "update users set email=$1, password=$2, confirm_code=$3, max_staples=$4 where email=$5",
//		newUser.Email,
//		newUser.Password,
//		newUser.ConfirmCode,
//		newUser.MaxStaples,
//		email); err != nil {
//		return err
//	}
//	return err
//}
