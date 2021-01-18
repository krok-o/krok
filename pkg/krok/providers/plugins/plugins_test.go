package plugins

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type mockCommandStorer struct {
	providers.CommandStorer
	getCommand  *models.Command
	deleteErr   error
	commandList []*models.Command
}

func (mcs *mockCommandStorer) Update(ctx context.Context, command *models.Command) (*models.Command, error) {
	return command, nil
}

func (mcs *mockCommandStorer) Create(ctx context.Context, command *models.Command) (*models.Command, error) {
	return command, nil
}

// Test the flow of the watcher. Create a location to watch and copy the command from
// testdata to this location and wait until the watcher picks it up.
func TestPluginProviderFlow(t *testing.T) {
	location, _ := ioutil.TempDir("", "TestNewGoPluginsProvider")
	logger := zerolog.New(os.Stderr)
	mcs := &mockCommandStorer{}
	_, err := NewGoPluginsProvider(context.Background(), Config{
		Location: location,
	}, Dependencies{
		Logger: logger,
		Store:  mcs,
	})

	assert.NoError(t, err)
	err = copyTestPlugin(location)
	assert.NoError(t, err)
}

// copyTestPlugin copies over the test plugin.
func copyTestPlugin(location string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	src := filepath.Join(cwd, "testdata")
	sourceFile, err := os.Open(filepath.Join(src, "test"))
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create(filepath.Join(location, "test"))
	if err != nil {
		return err
	}
	defer newFile.Close()

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		return err
	}
	log.Printf("Copied %d bytes.", bytesCopied)
	return nil
}
