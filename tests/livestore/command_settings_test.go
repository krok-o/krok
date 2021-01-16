package livestore

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/vault"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/tests/dbaccess"
)

func TestCommandSettings_Flow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandSettings_Create")
	env := environment.NewDockerConverter(environment.Config{}, environment.Dependencies{Logger: logger})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: dbaccess.Hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
	})
	assert.NoError(t, err)
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:         "Test_Create_Setting_1",
		Schedule:     "test-schedule-setting-1",
		Repositories: nil,
		Filename:     "test-filename-setting-1",
		Location:     location,
		Hash:         "settings-hash1",
		Enabled:      true,
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	err = cp.CreateSetting(ctx, &models.CommandSetting{
		ID:        0,
		CommandID: c.ID,
		Key:       "key",
		Value:     "value",
		InVault:   false,
	})
	assert.NoError(t, err)
	list, err := cp.ListSettings(ctx, c.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	setting := list[0]
	// Get the setting.
	getSetting, err := cp.GetSetting(ctx, setting.ID)
	assert.NoError(t, err)
	assert.Equal(t, &models.CommandSetting{
		ID:        setting.ID,
		CommandID: c.ID,
		Key:       "key",
		Value:     "value",
		InVault:   false,
	}, getSetting)

	// Update setting
	setting.Value = "new_value"
	err = cp.UpdateSetting(ctx, setting)
	assert.NoError(t, err)
	updatedSetting, err := cp.GetSetting(ctx, setting.ID)
	assert.NoError(t, err)
	assert.Equal(t, "new_value", updatedSetting.Value)
	// Delete setting
	err = cp.DeleteSetting(ctx, setting.ID)
	assert.NoError(t, err)

	// Try getting the setting should result in NotFound
	_, err = cp.GetSetting(ctx, setting.ID)
	assert.True(t, errors.Is(err, kerr.ErrNotFound))
}

func TestCommandSettings_Vault(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandSettings_Vault")
	env := environment.NewDockerConverter(environment.Config{}, environment.Dependencies{Logger: logger})
	fileStore, err := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	assert.NoError(t, err)
	err = fileStore.Init()
	assert.NoError(t, err)
	v, err := vault.NewKrokVault(vault.Config{}, vault.Dependencies{Logger: logger, Storer: fileStore})
	assert.NoError(t, err)
	connector := livestore.NewDatabaseConnector(livestore.Config{
		Hostname: dbaccess.Hostname,
		Database: dbaccess.Db,
		Username: dbaccess.Username,
		Password: dbaccess.Password,
	}, livestore.Dependencies{
		Logger:    logger,
		Converter: env,
	})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: connector,
		Vault:     v,
	})
	assert.NoError(t, err)
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:     "Test_Relationship_Vault",
		Schedule: "Test_Relationship_Vault-test-schedule",
		Filename: "Test_Relationship_Vault-test-filename-create",
		Location: location,
		Hash:     "Test_Relationship_Vault-hash1",
		Enabled:  false,
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	// put setting into vault
	err = cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key",
		Value:     "confidential_value",
		InVault:   true,
	})
	assert.NoError(t, err)

	err = v.LoadSecrets()
	assert.NoError(t, err)
	value, err := v.GetSecret("command_setting_1_key")
	assert.NoError(t, err)
	assert.Equal(t, string(value), "confidential_value")

	list, err := cp.ListSettings(ctx, c.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	setting := list[0]

	assert.Equal(t, "***********", setting.Value)

	getSetting, err := cp.GetSetting(ctx, setting.ID)
	assert.NoError(t, err)
	assert.Equal(t, "confidential_value", getSetting.Value)
}
