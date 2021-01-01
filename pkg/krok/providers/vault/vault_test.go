package vault

import (
	"errors"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	kerr "github.com/krok-o/krok/errors"
)

type memoryVault struct {
	data     []byte
	initErr  error
	readErr  error
	writeErr error
}

func (m *memoryVault) Init() error {
	m.data = make([]byte, 0)
	return m.initErr
}

// Read will read bytes from the storage medium and return it to the caller.
func (m *memoryVault) Read() (data []byte, err error) {
	return m.data, m.readErr
}

// Write will store the passed in data. How, is up to the implementor. Syncing
// is up the caller. Otherwise data will be overwritten.
func (m *memoryVault) Write(data []byte) error {
	m.data = data
	return m.writeErr
}

func TestNewKrokVault_Flow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	m := &memoryVault{}
	v, err := NewKrokVault(Config{}, Dependencies{
		Logger: logger,
		Storer: m,
	})
	assert.NoError(t, err)

	err = v.LoadSecrets()
	assert.NoError(t, err)
	v.AddSecret("key", []byte("value"))
	err = v.SaveSecrets()
	assert.NoError(t, err)
	err = v.LoadSecrets()
	assert.NoError(t, err)
	value, err := v.GetSecret("key")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value"), value)
	list := v.ListSecrets()
	assert.Len(t, list, 1)
	assert.Equal(t, "key", list[0])
	v.DeleteSecret("key")
	err = v.SaveSecrets()
	assert.NoError(t, err)
	err = v.LoadSecrets()
	assert.NoError(t, err)
	_, err = v.GetSecret("key")
	assert.True(t, errors.Is(err, kerr.ErrNotFound))
}

func TestNewKrokVault_ReadError(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	m := &memoryVault{
		readErr: errors.New("error"),
	}
	v, err := NewKrokVault(Config{}, Dependencies{
		Logger: logger,
		Storer: m,
	})
	assert.NoError(t, err)

	err = v.LoadSecrets()
	assert.Error(t, err)
}
