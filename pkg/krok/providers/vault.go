package providers

// Vault defines the capabilities of the vault storage.
type Vault interface {
	// LoadSecrets unlocks the vault and loads in all secrets.
	LoadSecrets() error
	// ListSecrets lists all secret names. Not the values.
	ListSecrets() []string
	// SaveSecrets saves all the secrets to the vault. Persisting new values.
	SaveSecrets() error
	// AddSecret adds a value to the vault.
	AddSecret(key string, value []byte)
	// DeleteSecret deletes a secret from the vault.
	DeleteSecret(key string)
	// GetSecret returns a single secret's value from the vault.
	GetSecret(key string) ([]byte, error)
}

// VaultStorer defines the interface for storing things in the vault.
// This can be any kind of store which supports these operations.
type VaultStorer interface {
	// Init initializes the medium by creating the file, or bootstrapping the
	// db or simply setting up an in-memory mock storage device. The Init
	// function of a storage medium should be idempotent. Meaning it should
	// be callable multiple times without changing the underlying medium.
	Init() error
	// Read will read bytes from the storage medium and return it to the caller.
	Read() (data []byte, err error)
	// Write will store the passed in data. How, is up to the implementor. Syncing
	// is up the caller. Otherwise data will be overwritten.
	Write(data []byte) error
}
