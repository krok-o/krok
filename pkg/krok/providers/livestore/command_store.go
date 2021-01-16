package livestore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cirello.io/pglock"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	commandsTable        = "commands"
	commandSettingsTable = "command_settings"
)

// CommandStore is a postgres based store for commands.
type CommandStore struct {
	CommandDependencies
}

// CommandDependencies command specific dependencies such as, the repository store.
// In order to not repeat some SQL, the command store will require the repository
// store and the repository store will require the command store.
type CommandDependencies struct {
	Dependencies
	Connector *Connector
}

// NewCommandStore creates a new CommandStore
func NewCommandStore(deps CommandDependencies) (*CommandStore, error) {
	db, err := deps.Connector.GetDB()
	if err != nil {
		deps.Logger.Debug().Err(err).Msg("Failed to get DB for locking.")
		return nil, err
	}
	c, err := pglock.New(db,
		pglock.WithLeaseDuration(10*time.Second),
		pglock.WithHeartbeatFrequency(1*time.Second),
	)
	if err != nil {
		deps.Logger.Debug().Err(err).Msg("Cannot create lock client.")
		return nil, err
	}
	if err := c.CreateTable(); err != nil && !strings.Contains(err.Error(), "relation \"locks\" already exists") {
		deps.Logger.Debug().Err(err).Msg("Failed to create lock table.")
		return nil, err
	}
	cs := &CommandStore{CommandDependencies: deps}
	return cs, nil
}

var _ providers.CommandStorer = &CommandStore{}

// Create creates a command record.
func (s *CommandStore) Create(ctx context.Context, c *models.Command) (*models.Command, error) {
	log := s.Logger.With().Str("name", c.Name).Logger()
	// duplicate key value violates unique constraint
	// id will be generated.

	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(name, schedule, filename, hash, location,"+
			" enabled) values($1, $2, $3, $4, $5, $6)", commandsTable),
			c.Name,
			c.Schedule,
			c.Filename,
			c.Hash,
			c.Location,
			c.Enabled); err != nil {
			log.Debug().Err(err).Msg("Failed to create command.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into commands",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into commands",
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, err
	}

	result, err := s.GetByName(ctx, c.Name)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get created command.")
		return nil, err
	}
	return result, nil
}

// Get returns a command model.
func (s *CommandStore) Get(ctx context.Context, id int) (*models.Command, error) {
	log := s.Logger.With().Int("id", id).Str("func", "GetByID").Logger()
	return s.getByX(ctx, log, "id", id)
}

// GetByName returns a command model by name.
func (s *CommandStore) GetByName(ctx context.Context, name string) (*models.Command, error) {
	log := s.Logger.With().Str("name", name).Str("func", "GetByName").Logger()
	return s.getByX(ctx, log, "name", name)
}

// Get returns a command model.
func (s *CommandStore) getByX(ctx context.Context, log zerolog.Logger, field string, value interface{}) (*models.Command, error) {
	log = s.Logger.With().Str("field", field).Interface("value", value).Logger()

	var (
		name      string
		commandID int
		schedule  string
		filename  string
		location  string
		hash      string
		enabled   bool
	)
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("select name, id, schedule, filename, location, hash, enabled from %s where %s = $1", commandsTable, field)
		if err := tx.QueryRow(ctx, query, value).
			Scan(&name, &commandID, &schedule, &filename, &location, &hash, &enabled); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: query,
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query row.")
			return &kerr.QueryError{
				Query: query,
				Err:   err,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("failed to run in transactions")
		return nil, err
	}

	repositories, err := s.getRepositoriesForCommand(ctx, commandID)
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("GetRepositoriesForCommand failed")
		return nil, &kerr.QueryError{
			Query: "select id",
			Err:   err,
		}
	}

	return &models.Command{
		Name:         name,
		ID:           commandID,
		Schedule:     schedule,
		Repositories: repositories,
		Filename:     filename,
		Location:     location,
		Hash:         hash,
		Enabled:      enabled,
	}, nil
}

func (s *CommandStore) getRepositoriesForCommand(ctx context.Context, id int) ([]*models.Repository, error) {
	log := s.Logger.With().Int("id", id).Logger()

	// Select the related repositories.
	result := make([]*models.Repository, 0)
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("select r.id, name, url, vcs from %s r inner join %s rel"+
			" on r.id = rel.repository_id where rel.command_id = $1", repositoriesTable, commandsRepositoriesRelTable)
		rows, err := tx.Query(ctx, query, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Debug().Err(err).Str("query", query).Msg("no repositories found for command.")
				return &kerr.QueryError{
					Query: query,
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query rel_repositories_command.")
			return &kerr.QueryError{
				Query: query,
				Err:   fmt.Errorf("failed to query rel table: %w", err),
			}
		}

		// Repo data here construct, individual repos.
		for rows.Next() {
			var (
				repoID int
				name   string
				url    string
				vcs    int
			)
			if err := rows.Scan(&repoID, &name, &url, &vcs); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select id",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			repo := &models.Repository{
				Name: name,
				ID:   id,
				URL:  url,
				VCS:  vcs,
			}
			result = append(result, repo)
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute GetRepositoriesForCommand: %w", err)
	}
	return result, nil
}

// Delete will remove a command.
func (s *CommandStore) Delete(ctx context.Context, id int) error {
	log := s.Logger.With().Int("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", commandsTable), id); err != nil {
			log.Debug().Err(err).Msg("Failed to delete command.")
			return &kerr.QueryError{
				Query: "delete id",
				Err:   fmt.Errorf("failed delete command: %w", err),
			}
		}
		return nil
	}

	return s.Connector.ExecuteWithTransaction(ctx, log, f)
}

// Update modifies a command record.
func (s *CommandStore) Update(ctx context.Context, c *models.Command) (*models.Command, error) {
	log := s.Logger.With().Int("id", c.ID).Str("name", c.Name).Logger()
	f := func(tx pgx.Tx) error {
		// Prevent updating the ID and the creation timestamp.
		// construct update statement:
		commandTags, err := tx.Exec(ctx, fmt.Sprintf("update %s set name = $1, enabled = $2, schedule = $3, filename = $4, location = $5, hash = $6 where id = $7", commandsTable),
			c.Name, c.Enabled, c.Schedule, c.Filename, c.Location, c.Hash, c.ID)
		if err != nil {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   fmt.Errorf("failed to update: %w", err),
			}
		}
		if commandTags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Query: "update :" + c.Name,
				Err:   kerr.ErrNoRowsAffected,
			}
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return nil, fmt.Errorf("failed to execute update in transaction: %w", err)
	}
	result, err := s.Get(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List gets all the command records.
func (s *CommandStore) List(ctx context.Context, opts *models.ListOptions) ([]*models.Command, error) {
	log := s.Logger.With().Str("func", "List").Logger()
	// Select all commands.
	result := make([]*models.Command, 0)
	f := func(tx pgx.Tx) error {
		sql := fmt.Sprintf("select id, name, schedule, filename, hash, location, enabled from %s", commandsTable)
		where := " where "
		filters := make([]string, 0)
		if opts.Name != "" {
			filters = append(filters, "name = %"+opts.Name+"%")
		}
		filter := strings.Join(filters, " AND ")
		if filter != "" {
			sql += where + filter
		}
		rows, err := tx.Query(ctx, sql)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select all commands",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query commands.")
			return &kerr.QueryError{
				Query: "select all commands",
				Err:   fmt.Errorf("failed to list all commands: %w", err),
			}
		}

		for rows.Next() {
			var (
				id       int
				name     string
				schedule string
				filename string
				hash     string
				location string
				enabled  bool
			)
			if err := rows.Scan(&id, &name, &schedule, &filename, &hash, &location, &enabled); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all commands",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			command := &models.Command{
				Name:     name,
				ID:       id,
				Schedule: schedule,
				Filename: filename,
				Location: location,
				Hash:     hash,
				Enabled:  enabled,
			}
			result = append(result, command)
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all commands: %w", err)
	}
	return result, nil
}

// AcquireLock acquires a lock on a file so no other process deals with the same file.
func (s *CommandStore) AcquireLock(ctx context.Context, name string) (*pglock.Lock, error) {
	log := s.Logger.With().Str("func", "AcquireLock").Str("name", name).Logger()
	db, err := s.Connector.GetDB()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get DB for locking.")
		return nil, err
	}
	c, err := pglock.New(db,
		pglock.WithLeaseDuration(3*time.Second),
		pglock.WithHeartbeatFrequency(1*time.Second),
	)
	if err != nil {
		log.Debug().Err(err).Msg("Cannot create lock client.")
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	l, err := c.AcquireContext(ctx, name)
	if err != nil {
		log.Debug().Err(err).Msg("unexpected error while acquiring 1st lock")
		return nil, err
	}
	return l, nil
}

// AddCommandRelForRepository add an assignment for a command to a repository.
func (s *CommandStore) AddCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error {
	log := s.Logger.With().Str("func", "AddCommandRelForRepository").Int("command_id", commandID).Int("repository_id", repositoryID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(command_id, repository_id) values($1, $2)", commandsRepositoriesRelTable),
			commandID, repositoryID); err != nil {
			log.Debug().Err(err).Msg("Failed to create relationship between command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into " + commandsRepositoriesRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into " + commandsRepositoriesRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to insert into " + commandsRepositoriesRelTable)
		return err
	}
	return nil
}

// RemoveCommandRelForRepository remove a relation to a repository for a command.
func (s *CommandStore) RemoveCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error {
	log := s.Logger.With().Str("func", "RemoveCommandRelForRepository").Int("command_id", commandID).Int("repository_id", repositoryID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where command_id = $1 and repository_id = $2", commandsRepositoriesRelTable),
			commandID, repositoryID); err != nil {
			log.Debug().Err(err).Msg("Failed to remove relationship for command and repository.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from " + commandsRepositoriesRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from " + commandsRepositoriesRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to delete from " + commandsRepositoriesRelTable)
		return err
	}
	return nil
}

// CreateSetting will create a setting will take a list of settings and save it for a command.
// Since settings are usually submitted in a batch, this takes a list of settings by default.
func (s *CommandStore) CreateSetting(ctx context.Context, settings []*models.CommandSetting) error {
	log := s.Logger.With().Str("func", "CreateSetting").Int("count", len(settings)).Logger()
	f := func(tx pgx.Tx) error {
		// construct a batch and send it.
		batch := &pgx.Batch{}
		for _, setting := range settings {
			// TODO : If InVault, store in vault and generate a unique key for the vault and the
			// value of the actual setting here will be the key of the generated vault setting name as a reference.
			batch.Queue(fmt.Sprintf("insert into %s(command_id, key, value, in_vault) values($1, $2, $3, $4)", commandSettingsTable),
				setting.CommandID,
				setting.Key,
				setting.Value,
				setting.InVault)
		}

		br := tx.SendBatch(ctx, batch)
		if tags, err := br.Exec(); err != nil {
			log.Debug().Err(err).Msg("Failed to batch create settings.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into command_settings",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into command_settings",
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		return err
	}
	return nil
}

// DeleteSetting takes a
func (s *CommandStore) DeleteSetting(ctx context.Context, id int) error {
	log := s.Logger.With().Int("id", id).Logger()
	f := func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", commandSettingsTable), id); err != nil {
			log.Debug().Err(err).Msg("Failed to delete command setting.")
			return &kerr.QueryError{
				Query: "delete id",
				Err:   fmt.Errorf("failed delete command setting: %w", err),
			}
		}
		return nil
	}

	return s.Connector.ExecuteWithTransaction(ctx, log, f)
}

func (s *CommandStore) ListSettings(ctx context.Context, commandID int) ([]*models.CommandSetting, error) {
	panic("implement me")
}

// GetSetting returns a single setting for an ID.
func (s *CommandStore) GetSetting(ctx context.Context, id int) (*models.CommandSetting, error) {
	log := s.Logger.With().Int("id", id).Logger()

	var (
		storedID  int
		commandID int
		key       string
		value     string
		inVault   bool
	)
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("select id, command_id, key, value, in_vault from %s where id = $1", commandSettingsTable)
		if err := tx.QueryRow(ctx, query, id).
			Scan(&storedID, &commandID, &key, &value, &inVault); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: query,
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query row.")
			return &kerr.QueryError{
				Query: query,
				Err:   err,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("failed to run in transactions")
		return nil, err
	}

	if inVault {
		value = "****************"
	}

	return &models.CommandSetting{
		ID:        storedID,
		CommandID: commandID,
		Key:       key,
		Value:     value,
		InVault:   inVault,
	}, nil
}

func (s *CommandStore) UpdateSetting(ctx context.Context, setting *models.CommandSetting) error {
	panic("implement me")
}
