package livestore

import (
	"context"
	"errors"
	"fmt"
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
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
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
		Enabled:      true,
		Image:        "krokhook/slack-notification:v0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	setting, err := cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key",
		Value:     "value",
		InVault:   false,
	})
	assert.NoError(t, err)
	assert.Equal(t, &models.CommandSetting{
		ID:        setting.ID,
		CommandID: c.ID,
		Key:       "key",
		Value:     "value",
		InVault:   false,
	}, setting)
	list, err := cp.ListSettings(ctx, c.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	setting = list[0]
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
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	fileStore := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	err := fileStore.Init()
	assert.NoError(t, err)
	v := vault.NewKrokVault(vault.Dependencies{Logger: logger, Storer: fileStore})
	connector := livestore.NewDatabaseConnector(livestore.Config{
		Hostname: hostname,
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
		Enabled:  false,
		Image:    "krokhook/slack-notification:v0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	// put setting into vault
	setting, err := cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key",
		Value:     "confidential_value",
		InVault:   true,
	})
	assert.NoError(t, err)
	assert.True(t, setting.ID > 0)

	err = v.LoadSecrets()
	assert.NoError(t, err)

	list, err := cp.ListSettings(ctx, c.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	setting = list[0]
	assert.Equal(t, "confidential_value", setting.Value)

	vKey := fmt.Sprintf("command_setting_%d_%s", c.ID, setting.Key)
	value, err := v.GetSecret(vKey)
	assert.NoError(t, err)
	assert.Equal(t, string(value), "confidential_value")

	getSetting, err := cp.GetSetting(ctx, setting.ID)
	assert.NoError(t, err)
	assert.Equal(t, "confidential_value", getSetting.Value)
}

func TestCommandSettings_CascadingDelete(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
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
		Name:         "Test_CascadeDelete_Setting_1",
		Schedule:     "test-schedule-setting-1",
		Repositories: nil,
		Enabled:      true,
		Image:        "krokhook/slack-notification:v0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	setting, err := cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key-5",
		Value:     "value",
		InVault:   false,
	})
	assert.NoError(t, err)
	assert.True(t, setting.ID > 0)
	list, err := cp.ListSettings(ctx, c.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	setting = list[0]

	err = cp.Delete(ctx, c.ID)
	assert.NoError(t, err)

	// Try getting the setting should result in NotFound
	_, err = cp.GetSetting(ctx, setting.ID)
	assert.True(t, errors.Is(err, kerr.ErrNotFound))
}

func TestCommandSettings_CantCreateSameKeyAndCommandCombination(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
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
		Name:         "Test_CreateError_Setting_1",
		Schedule:     "test-schedule-setting-1",
		Repositories: nil,
		Enabled:      true,
		Image:        "krokhook/slack-notification:v0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	_, err = cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key-5",
		Value:     "value",
		InVault:   false,
	})
	assert.NoError(t, err)
	_, err = cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key-5",
		Value:     "value",
		InVault:   false,
	})
	assert.Error(t, err)
	_, err = cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: 999,
		Key:       "key-5",
		Value:     "value",
		InVault:   false,
	})
	assert.Error(t, err)
}

func TestCommandSettings_UpdateInVault(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandSettings_UpdateInVault")
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	fileStore := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	err := fileStore.Init()
	assert.NoError(t, err)
	v := vault.NewKrokVault(vault.Dependencies{Logger: logger, Storer: fileStore})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
		Vault: v,
	})
	assert.NoError(t, err)
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:         "Test_UpdateVault_Setting_1",
		Schedule:     "test-schedule-setting-1",
		Repositories: nil,
		Enabled:      true,
		Image:        "krokhook/slack-notification:v0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	_, err = cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key",
		Value:     "value",
		InVault:   true,
	})
	assert.NoError(t, err)
	list, err := cp.ListSettings(ctx, c.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	setting := list[0]
	// Update setting
	setting.Value = "new_value"
	err = cp.UpdateSetting(ctx, setting)
	assert.NoError(t, err)
	vKey := fmt.Sprintf("command_setting_%d_%s", c.ID, setting.Key)

	err = v.LoadSecrets()
	assert.NoError(t, err)

	secret, err := v.GetSecret(vKey)
	assert.NoError(t, err)
	assert.Equal(t, "new_value", string(secret))
}

func TestCommandSettings_ErrorOnListIfValueDoesntExistsInVault(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestCommandSettings_ErrorOnListIfValueDoesntExistsInVault")
	env := environment.NewDockerConverter(environment.Dependencies{Logger: logger})
	fileStore := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	err := fileStore.Init()
	assert.NoError(t, err)
	v := vault.NewKrokVault(vault.Dependencies{Logger: logger, Storer: fileStore})
	cp, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Connector: livestore.NewDatabaseConnector(livestore.Config{
			Hostname: hostname,
			Database: dbaccess.Db,
			Username: dbaccess.Username,
			Password: dbaccess.Password,
		}, livestore.Dependencies{
			Logger:    logger,
			Converter: env,
		}),
		Vault: v,
	})
	assert.NoError(t, err)
	ctx := context.Background()
	// Create the first command.
	c, err := cp.Create(ctx, &models.Command{
		Name:         "Test_Error_On_No_Value_Setting_1",
		Schedule:     "test-schedule-setting-1",
		Repositories: nil,
		Enabled:      true,
		Image:        "krokhook/slack-notification:v0.0.1",
	})
	assert.NoError(t, err)
	assert.True(t, 0 < c.ID)

	_, err = cp.CreateSetting(ctx, &models.CommandSetting{
		CommandID: c.ID,
		Key:       "key99",
		Value:     "value",
		InVault:   true,
	})
	assert.NoError(t, err)
	err = v.LoadSecrets()
	assert.NoError(t, err)
	vKey := fmt.Sprintf("command_setting_%d_%s", c.ID, "key99")
	v.DeleteSecret(vKey)
	err = v.SaveSecrets()
	assert.NoError(t, err)

	_, err = cp.ListSettings(ctx, c.ID)
	assert.Error(t, err)
}
