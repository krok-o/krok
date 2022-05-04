package livestore

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	Vault     providers.Vault
}

// NewCommandStore creates a new CommandStore
func NewCommandStore(deps CommandDependencies) (*CommandStore, error) {
	return &CommandStore{CommandDependencies: deps}, nil
}

var _ providers.CommandStorer = &CommandStore{}

// Create creates a command record.
func (s *CommandStore) Create(ctx context.Context, c *models.Command) (*models.Command, error) {
	log := s.Logger.With().Str("name", c.Name).Logger()
	// duplicate key value violates unique constraint
	// id will be generated.

	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(name, schedule, enabled, image, requires_clone) values($1, $2, $3, $4, $5)", commandsTable),
			c.Name,
			c.Schedule,
			c.Enabled,
			c.Image,
			c.RequiresClone); err != nil {
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
		name          string
		commandID     int
		schedule      string
		enabled       bool
		image         string
		requiresClone bool
	)
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("select name, id, schedule, enabled, image, requires_clone from %s where %s = $1", commandsTable, field)
		if err := tx.QueryRow(ctx, query, value).
			Scan(&name, &commandID, &schedule, &enabled, &image, &requiresClone); err != nil {
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

	platforms, err := s.getPlatformsForCommand(ctx, commandID)
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("getPlatformsForCommand failed")
		return nil, &kerr.QueryError{
			Query: "select id",
			Err:   err,
		}
	}

	return &models.Command{
		Name:          name,
		ID:            commandID,
		Schedule:      schedule,
		Repositories:  repositories,
		Enabled:       enabled,
		Image:         image,
		Platforms:     platforms,
		RequiresClone: requiresClone,
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

// getPlatformsForCommand returns a list of platforms which this command supports.
func (s *CommandStore) getPlatformsForCommand(ctx context.Context, id int) ([]models.Platform, error) {
	log := s.Logger.With().Int("id", id).Logger()

	// Select the related platforms.
	var result []models.Platform
	f := func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("select platform_id from %s where command_id = $1", commandsPlatformsRelTable), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select id",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query relationship.")
			return &kerr.QueryError{
				Query: "select commands for platforms",
				Err:   fmt.Errorf("failed to query rel table: %w", err),
			}
		}

		for rows.Next() {
			var (
				platformID int
			)
			if err := rows.Scan(&platformID); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select id",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			result = append(result, models.SupportedPlatforms[platformID])
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute getPlatformsForCommand: %w", err)
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
		args := make([]interface{}, 0)
		sets := make([]string, 0)

		if c.Name != "" {
			args = append(args, c.Name)
			sets = append(sets, "name = $"+strconv.Itoa(len(args)))
		}
		if c.Schedule != "" {
			args = append(args, c.Schedule)
			sets = append(sets, "schedule = $"+strconv.Itoa(len(args)))
		}

		// TODO: change this to a reference type on the enabled to check whether it was supplied or not.
		args = append(args, c.Enabled)
		sets = append(sets, "enabled = $"+strconv.Itoa(len(args)))

		set := strings.Join(sets, ",")
		args = append(args, c.ID)

		commandTags, err := tx.Exec(ctx, fmt.Sprintf("update %s set %s where id = $%d", commandsTable, set, len(args)),
			args...)
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
		sql := fmt.Sprintf("select id, name, schedule, enabled, image, requires_clone from %s", commandsTable)
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
				id            int
				name          string
				schedule      string
				image         string
				enabled       bool
				requiresClone bool
			)
			if err := rows.Scan(&id, &name, &schedule, &enabled, &image, &requiresClone); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all commands",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			command := &models.Command{
				Name:          name,
				ID:            id,
				Schedule:      schedule,
				Enabled:       enabled,
				Image:         image,
				RequiresClone: requiresClone,
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

// CreateSetting will create a setting for a command.
func (s *CommandStore) CreateSetting(ctx context.Context, setting *models.CommandSetting) (*models.CommandSetting, error) {
	log := s.Logger.
		With().
		Str("func", "CreateSetting").
		Str("key", setting.Key).
		Bool("in_vault", setting.InVault).
		Logger()
	rollBackValue := ""
	var returnedID int
	f := func(tx pgx.Tx) error {
		value := setting.Value
		if setting.InVault {
			value = s.generateUniqueVaultID(setting.CommandID, setting.Key)
			if err := s.Vault.LoadSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to load secrets.")
				return err
			}
			s.Vault.AddSecret(value, []byte(setting.Value))
			if err := s.Vault.SaveSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to save secrets.")
				return err
			}
			rollBackValue = value
		}
		query := fmt.Sprintf("insert into %s(command_id, key, value, in_vault) values($1, $2, $3, $4) returning id", commandSettingsTable)
		rows := tx.QueryRow(ctx, query,
			setting.CommandID,
			setting.Key,
			value,
			setting.InVault)
		if err := rows.Scan(&returnedID); err != nil {
			log.Debug().Err(err).Str("query", query).Msg("Failed to scan row.")
			return &kerr.QueryError{
				Err:   err,
				Query: query,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		// delete all the possibly created vault settings
		if rollBackValue != "" {
			if err := s.Vault.LoadSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to load secrets.")
				return nil, err
			}
			s.Vault.DeleteSecret(rollBackValue)
			if err := s.Vault.SaveSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to save secrets.")
				return nil, err
			}
		}
		return nil, err
	}
	setting, err := s.GetSetting(ctx, returnedID)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get created setting")
		return nil, err
	}
	return setting, nil
}

// generateUniqueVaultID generates a unique vault key based on the command id and the key name.
func (s *CommandStore) generateUniqueVaultID(commandID int, key string) string {
	return fmt.Sprintf("command_setting_%d_%s", commandID, key)
}

// DeleteSetting takes a
func (s *CommandStore) DeleteSetting(ctx context.Context, id int) error {
	log := s.Logger.With().Int("id", id).Logger()
	// We only delete the values once they have successfully been removed from the DB.
	// If the vault delete would fail that's less of a problem compared to if we
	// remove the vault value and the database reference remains to it.
	toDeleteVaultValues := make([]string, 0)
	f := func(tx pgx.Tx) error {
		// if setting is in vault, we delete from there as well.
		setting, err := s.GetSetting(ctx, id)
		if err != nil {
			return err
		}
		if setting.InVault {
			toDeleteVaultValues = append(toDeleteVaultValues, setting.Value)
		}
		if _, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where id = $1", commandSettingsTable), id); err != nil {
			log.Debug().Err(err).Msg("Failed to delete command setting.")
			return &kerr.QueryError{
				Query: "delete id",
				Err:   fmt.Errorf("failed delete command setting: %w", err),
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to run transaction.")
		return err
	}

	// remove the vault values if any
	if len(toDeleteVaultValues) > 0 {
		if err := s.Vault.LoadSecrets(); err != nil {
			log.Debug().Err(err).Msg("Failed to load secrets.")
			return err
		}
		for _, v := range toDeleteVaultValues {
			s.Vault.DeleteSecret(v)
		}
		if err := s.Vault.SaveSecrets(); err != nil {
			log.Debug().Err(err).Msg("Failed to save secrets.")
			return err
		}
	}
	return nil
}

// ListSettings lists all settings for a command.
func (s *CommandStore) ListSettings(ctx context.Context, commandID int) ([]*models.CommandSetting, error) {
	log := s.Logger.With().Str("func", "ListSettings").Logger()
	// Select all commands.
	result := make([]*models.CommandSetting, 0)
	f := func(tx pgx.Tx) error {
		sql := fmt.Sprintf("select id, command_id, key, value, in_vault from %s where command_id = $1", commandSettingsTable)
		rows, err := tx.Query(ctx, sql, commandID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &kerr.QueryError{
					Query: "select all command settings",
					Err:   kerr.ErrNotFound,
				}
			}
			log.Debug().Err(err).Msg("Failed to query command settings.")
			return &kerr.QueryError{
				Query: "select all command settings",
				Err:   fmt.Errorf("failed to list all command settings: %w", err),
			}
		}

		for rows.Next() {
			var (
				id              int
				storedCommandID int
				key             string
				value           string
				inVault         bool
			)
			if err := rows.Scan(&id, &storedCommandID, &key, &value, &inVault); err != nil {
				log.Debug().Err(err).Msg("Failed to scan.")
				return &kerr.QueryError{
					Query: "select all command settings",
					Err:   fmt.Errorf("failed to scan: %w", err),
				}
			}
			if inVault {
				if err := s.Vault.LoadSecrets(); err != nil {
					log.Debug().Err(err).Msg("Failed to load secrets.")
					return err
				}
				v, err := s.Vault.GetSecret(value)
				if err != nil {
					log.Debug().Err(err).Msg("Failed to get value for secret from vault.")
					return err
				}
				value = string(v)
			}
			setting := &models.CommandSetting{
				ID:        id,
				CommandID: storedCommandID,
				Key:       key,
				Value:     value,
				InVault:   inVault,
			}
			result = append(result, setting)
		}
		return nil
	}
	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		return nil, fmt.Errorf("failed to execute List all command settings: %w", err)
	}
	return result, nil
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
		if err := s.Vault.LoadSecrets(); err != nil {
			log.Debug().Err(err).Msg("Failed to load secrets.")
			return nil, err
		}
		b, err := s.Vault.GetSecret(value)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}

	return &models.CommandSetting{
		ID:        storedID,
		CommandID: commandID,
		Key:       key,
		Value:     value,
		InVault:   inVault,
	}, nil
}

// UpdateSetting updates the value of a setting. Transferring values is not supported. Aka.:
// If a value was in Vault it must remain in vault. If it was in db it must remain in db.
// Updating the key is also not supported.
// Update: Only the value can be modified.
func (s *CommandStore) UpdateSetting(ctx context.Context, setting *models.CommandSetting) error {
	log := s.Logger.
		With().
		Str("func", "UpdateSetting").
		Str("key", setting.Key).
		Bool("in_vault", setting.InVault).
		Logger()
	var (
		rollBackValue []byte
		rollBackKey   string
	)
	f := func(tx pgx.Tx) error {
		storedSetting, err := s.GetSetting(ctx, setting.ID)
		if err != nil {
			return err
		}

		// If it was in vault, it's easier to just overwrite whatever was in vault.
		if storedSetting.InVault {
			value := s.generateUniqueVaultID(storedSetting.CommandID, storedSetting.Key)
			if err := s.Vault.LoadSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to load secrets.")
				return err
			}
			if rollBackValue, err = s.Vault.GetSecret(value); err != nil {
				log.Debug().Err(err).Msg("Failed to get secret key.")
				return err
			}
			s.Vault.AddSecret(value, []byte(setting.Value))
			if err := s.Vault.SaveSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to save secrets.")
				return err
			}
			rollBackKey = value
			setting.Value = value
		}
		if tags, err := tx.Exec(ctx, fmt.Sprintf("update %s set value = $1 where id = $2", commandSettingsTable),
			setting.Value, storedSetting.ID); err != nil {
			log.Debug().Err(err).Msg("Failed to update setting.")
			return &kerr.QueryError{
				Err:   err,
				Query: "update command_settings",
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "update command setting",
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to execute with transaction.")
		// We set back the vault value to its original value if key is not empty.
		if rollBackKey != "" {
			if err := s.Vault.LoadSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to load secrets.")
				return err
			}
			s.Vault.AddSecret(rollBackKey, rollBackValue)
			if err := s.Vault.SaveSecrets(); err != nil {
				log.Debug().Err(err).Msg("Failed to save secrets.")
				return err
			}
		}
		return err
	}
	return nil
}

// AddCommandRelForPlatform adds a relationship for a platform on a command. This means
// that this command will support this platform. If the relationship doesn't exist
// this command will not run on that platform.
func (s *CommandStore) AddCommandRelForPlatform(ctx context.Context, commandID int, platformID int) error {
	log := s.Logger.With().Str("func", "AddCommandRelForPlatform").Int("command_id", commandID).Int("platform_id", platformID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("insert into %s(command_id, platform_id) values($1, $2)", commandsPlatformsRelTable),
			commandID, platformID); err != nil {
			log.Debug().Err(err).Msg("Failed to create relationship between command and platform.")
			return &kerr.QueryError{
				Err:   err,
				Query: "insert into " + commandsPlatformsRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "insert into " + commandsPlatformsRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to insert into " + commandsPlatformsRelTable)
		return err
	}
	return nil
}

// RemoveCommandRelForPlatform removes the above relationship, disabling this command
// for that platform. Meaning this command will not be executed if that platform is
// detected.
func (s *CommandStore) RemoveCommandRelForPlatform(ctx context.Context, commandID int, platformID int) error {
	log := s.Logger.With().Str("func", "RemoveCommandRelForPlatform").Int("command_id", commandID).Int("platform_id", platformID).Logger()
	f := func(tx pgx.Tx) error {
		if tags, err := tx.Exec(ctx, fmt.Sprintf("delete from %s where command_id = $1 and platform_id = $2", commandsPlatformsRelTable),
			commandID, platformID); err != nil {
			log.Debug().Err(err).Msg("Failed to remove relationship for command and platform.")
			return &kerr.QueryError{
				Err:   err,
				Query: "delete from " + commandsPlatformsRelTable,
			}
		} else if tags.RowsAffected() == 0 {
			return &kerr.QueryError{
				Err:   kerr.ErrNoRowsAffected,
				Query: "delete from " + commandsPlatformsRelTable,
			}
		}
		return nil
	}

	if err := s.Connector.ExecuteWithTransaction(ctx, log, f); err != nil {
		log.Debug().Err(err).Msg("Failed to delete from " + commandsPlatformsRelTable)
		return err
	}
	return nil
}

// IsPlatformSupported returns if a command supports a platform or not.
func (s *CommandStore) IsPlatformSupported(ctx context.Context, commandID, platformID int) (bool, error) {
	log := s.Logger.With().Int("command_id", commandID).Int("platform_id", platformID).Logger()
	var result int
	f := func(tx pgx.Tx) error {
		query := fmt.Sprintf("select count(1) from %s where command_id = $1 and platform_id = $2", commandsPlatformsRelTable)
		if err := tx.QueryRow(ctx, query, commandID, platformID).Scan(&result); err != nil {
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
		log.Debug().Err(err).Msg("Failed to query " + commandsPlatformsRelTable)
		return false, err
	}
	return result == 1, nil
}
