package vault

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
)

// Config has the configuration options for the vault.
type Config struct {
}

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger zerolog.Logger
	Storer providers.VaultStorer
}

// KrokVault is the vault used by Krok.
type KrokVault struct {
	Config
	Dependencies

	data map[string][]byte
	sync.RWMutex
}

// NewKrokVault creates a vault which contains secrets.
// The format is:
// KEY=VALUE
// KEY2=VALUE2
func NewKrokVault(cfg Config, deps Dependencies) (*KrokVault, error) {
	return &KrokVault{
		Config:       cfg,
		Dependencies: deps,
	}, nil
}

// ParseToMap will update the Vault data map with values from
// an encrypted file content.
func (v *KrokVault) parseToMap(data []byte) error {
	if len(data) < 1 {
		return nil
	}
	row := bytes.Split(data, []byte("\n"))
	for _, r := range row {
		d := bytes.Split(r, []byte("="))
		v.data[string(d[0])] = d[1]
	}
	return nil
}

// ParseFromMap will create a joined by new line set of key value
// pairs ready to be saved.
func (v *KrokVault) parseFromMap() []byte {
	data := make([][]byte, 0)
	for key, value := range v.data {
		s := fmt.Sprintf("%s=%s", key, value)
		data = append(data, []byte(s))
	}

	return bytes.Join(data, []byte("\n"))
}

// LoadSecrets unlocks the vault and loads in all secrets.
func (v *KrokVault) LoadSecrets() error {
	data, err := v.Storer.Read()
	if err != nil {
		return err
	}
	return v.parseToMap(data)
}

// SaveSecrets saves all the secrets to the vault. Persisting new values.
func (v *KrokVault) SaveSecrets() error {
	data := v.parseFromMap()
	// clear the hash after saving so the system always has a fresh view of the vault.
	v.data = make(map[string][]byte)
	return v.Storer.Write(data)
}

// ListSecrets lists all secret names. Not the values.
func (v *KrokVault) ListSecrets() []string {
	v.RLock()
	defer v.RUnlock()
	m := make([]string, 0)
	for k := range v.data {
		m = append(m, k)
	}
	return m
}

// AddSecret adds a value to the vault.
// Add will overwrite if the key already exists and not warn.
func (v *KrokVault) AddSecret(key string, value []byte) {
	v.Lock()
	defer v.Unlock()
	v.data[key] = value
}

// DeleteSecret deletes a secret from the vault.
// DeleteSecret is a no-op if the data doesn't exist.
func (v *KrokVault) DeleteSecret(key string) {
	v.Lock()
	defer v.Unlock()
	delete(v.data, key)
}

// GetSecret returns a value for a key. This operation is safe to use concurrently.
// Get will return an error if the data doesn't exist.
func (v *KrokVault) GetSecret(key string) ([]byte, error) {
	v.RLock()
	defer v.RUnlock()
	val, ok := v.data[key]
	if !ok {
		return []byte{}, fmt.Errorf("key '%s' not found in vault", key)
	}

	return val, nil
}
