package livestore

import (
	"context"

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
		if _, err := tx.Exec(ctx, "insert into users(id, password, confirm_code, max_staples) values($1, $2, $3, $4)",
			email,
			password,
			"",
			DefaultMaxStaples); err != nil {
			log.Debug().Err(err).Msg("Failed to create user.")
			return err
		}
		return nil
	}

	return s.Connector.ExecuteWithTransaction(ctx, s.Logger, f)
}

// Delete deletes a user from the db.
func (s *UserStore) Delete(ctx context.Context, email string) error {
	conn, err := s.connect()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutForTransactions)
	defer cancel()

	defer conn.Close(ctx)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, "delete from users where email = $1",
		email); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get retrieves a user.
func (s *UserStore) Get(ctx context.Context, email string) (*models.User, error) {
	conn, err := s.connect()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutForTransactions)
	defer cancel()

	defer conn.Close(ctx)
	var (
		storedEmail string
		password    []byte
		confirmCode string
		maxStaples  int
	)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "select email, password, confirm_code, max_staples from users where email = $1", email).Scan(&storedEmail, &password, &confirmCode, &maxStaples)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
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
	conn, err := s.connect()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutForTransactions)
	defer cancel()

	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // this is safe to call even if commit is called first.

	if _, err = tx.Exec(ctx, "update users set email=$1, password=$2, confirm_code=$3, max_staples=$4 where email=$5",
		newUser.Email,
		newUser.Password,
		newUser.ConfirmCode,
		newUser.MaxStaples,
		email); err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}
