package plugins

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/krok/providers/tar"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLock struct {
}

func (ml *mockLock) Close() error {
	return nil
}

func TestPluginProviderFlow(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "TestNewGoPluginsProvider")
	location := filepath.Join(tmp, "location")
	err := os.Mkdir(location, os.ModePerm)
	assert.NoError(t, err)
	logger := zerolog.New(os.Stderr)
	dst := filepath.Join(tmp, "test.tar.gz")
	mcs := &mocks.CommandStorer{}
	mcs.On("AcquireLock", mock.Anything, dst).Return(&mockLock{}, nil)
	tarer := tar.NewTarer(tar.Dependencies{
		Logger: logger,
	})
	p := NewPluginsProvider(Config{
		Location: location,
	}, Dependencies{
		Logger: logger,
		Store:  mcs,
		Tar:    tarer,
	})
	// test the copying of test data to tmp folder
	err = p.Copy(filepath.Join("testdata", "test.tar.gz"), dst)
	assert.NoError(t, err)

	// test create
	hash, err := p.Create(context.Background(), dst)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}
