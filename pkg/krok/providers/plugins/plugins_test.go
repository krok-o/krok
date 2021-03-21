package plugins

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type mockCommandStorer struct {
	providers.CommandStorer
	getCommand *models.Command
	getError   error
	id         int
}

func (mcs *mockCommandStorer) Update(ctx context.Context, command *models.Command) (*models.Command, error) {
	return command, nil
}

func (mcs *mockCommandStorer) Create(ctx context.Context, command *models.Command) (*models.Command, error) {
	command.ID = mcs.id
	mcs.id++
	mcs.getCommand = command
	return command, nil
}

func (mcs *mockCommandStorer) Get(ctx context.Context, id int) (*models.Command, error) {
	return mcs.getCommand, mcs.getError
}

func (mcs *mockCommandStorer) GetByName(ctx context.Context, name string) (*models.Command, error) {
	return mcs.getCommand, mcs.getError
}

type mockLock struct {
}

func (ml *mockLock) Close() error {
	return nil
}

func (mcs *mockCommandStorer) AcquireLock(ctx context.Context, name string) (io.Closer, error) {
	return &mockLock{}, nil
}

// Test the flow of the watcher. Create a location to watch and copy the command from
// testdata to this location and wait until the watcher picks it up.
func TestPluginProviderFlow(t *testing.T) {
	location, _ := ioutil.TempDir("", "TestNewGoPluginsProvider")
	logger := zerolog.New(os.Stderr)
	mcs := &mockCommandStorer{
		getError: kerr.ErrNotFound,
	}
	pp, err := NewGoPluginsProvider(Config{
		Location: location,
	}, Dependencies{
		Logger: logger,
		Store:  mcs,
	})

	assert.NoError(t, err)
	go pp.Run(context.Background())

	// Wait for the watcher to start up...
	time.Sleep(1 * time.Second)
	//err = copyTestPlugin(location)
	file, err := ioutil.TempFile(location, "test")
	assert.NoError(t, err)

	// Wait for the watcher to pick up the new file and call create.
	time.Sleep(1 * time.Second)
	assert.Equal(t, filepath.Base(file.Name()), mcs.getCommand.Name)
	assert.Equal(t, 0, mcs.getCommand.ID)
	assert.Equal(t, filepath.Base(file.Name()), mcs.getCommand.Filename)
	assert.Equal(t, location, mcs.getCommand.Location)
	assert.NotEqual(t, "", mcs.getCommand.Hash)
}
