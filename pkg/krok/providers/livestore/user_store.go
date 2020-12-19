package livestore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

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
}

// NewUserStore creates a new UserStore
func NewUserStore(cfg Config, deps UserDependencies) *UserStore {
	return &UserStore{Config: cfg, UserDependencies: deps}
}

// Create saves a user in the db.
func (s *UserStore) Create(ctx context.Context, user *models.User) error {
	log := s.Logger.With().Str("id", user.ID).Logger()

	f := func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "insert into users(email, last_login, display_name) values($1, $2, $3)",
			user.Email,
			time.Now(),
			user.DisplayName); err != nil {
			log.Debug().Err(err).Msg("Failed to create user.")
			return err
		}
		return nil
	}

	return s.Connector.ExecuteWithTransaction(ctx, s.Logger, f)
}

// Delete deletes a user from the db.
func (s *UserStore) Delete(ctx context.Context, email string) error {
	if _, err = tx.Exec(ctx, "delete from users where email = $1",
		email); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get retrieves a user.
func (s *UserStore) Get(ctx context.Context, email string) (*models.User, error) {
	err = tx.QueryRow(ctx, "select email, password, confirm_code, max_staples from users where email = $1", email).Scan(&storedEmail, &password, &confirmCode, &maxStaples)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &models.User{
		Email:       storedEmail,
		Password:    string(password),
		ConfirmCode: confirmCode,
		MaxStaples:  maxStaples}, nil
}

// Update updates a user with a given email address.
func (s *UserStore) Update(ctx context.Context, email string, newUser models.User) error {
	if _, err = tx.Exec(ctx, "update users set email=$1, password=$2, confirm_code=$3, max_staples=$4 where email=$5",
		newUser.Email,
		newUser.Password,
		newUser.ConfirmCode,
		newUser.MaxStaples,
		email); err != nil {
		return err
	}
	return err
}
