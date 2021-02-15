package filevault

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

const (
	vaultFileName = "krok_vault.dat"
	keySize       = 32
)

// Config has the configuration options for the file vault.
type Config struct {
	// Location is the location of the saved vault file.
	Location string
	// Key is the key to use to unlock the vault file. It must be 32 bit.
	// If it isn't, it will be padded if it's larger, it will be cropped.
	Key string
}

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger zerolog.Logger
}

// FileStorer is a vault backed by an encrypted file.
type FileStorer struct {
	Config
	Dependencies

	pathToFile string
	counter    uint64
	key        []byte
	lock       sync.RWMutex
}

// NewFileStorer creates a vault which contains secrets.
// The format is:
// KEY=VALUE
// KEY2=VALUE2
func NewFileStorer(cfg Config, deps Dependencies) *FileStorer {
	return &FileStorer{
		Config:       cfg,
		Dependencies: deps,
	}
}

// Init initializes the vault file.
func (v *FileStorer) Init() error {
	path := filepath.Join(v.Location, vaultFileName)
	log := v.Logger.With().Str("location", path).Logger()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Debug().Msg("Vault file doesn't exist. Creating...")
		if _, err := os.Create(path); err != nil {
			log.Debug().Err(err).Msg("Failed to create vault file.")
			return err
		}
	} else if err != nil {
		log.Debug().Err(err).Msg("Failed to check if vault file is accessible.")
		return err
	}

	v.pathToFile = path
	v.key = v.padOrCropKey()
	return nil
}

// padOrCropKey checks if the given key is more or less than the required key size and
// crops or extends it accordingly. Padding is done with 0s.
func (v *FileStorer) padOrCropKey() []byte {
	if len(v.Key) < keySize {
		diff := keySize - len(v.Key)
		for i := 0; i < diff; i++ {
			v.Key += "0"
		}
	} else if len(v.Key) > keySize {
		v.Key = v.Key[:keySize]
	}
	return []byte(v.Key)
}

// Read defines a read for the FileVaultStorer. It decrypts the file to get to the content.
func (v *FileStorer) Read() ([]byte, error) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	log := v.Logger.With().Str("location", v.pathToFile).Logger()
	data, err := ioutil.ReadFile(v.pathToFile)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read file.")
		return nil, err
	}
	// decrypt the data so plain data is returned upon reading.
	return v.decrypt(data)
}

// Write will store the passed in data. How, is up to the implementor. Syncing
// is up the caller. Otherwise data will be overwritten. The file vault storer will encrypt the file
// upon writing.
func (v *FileStorer) Write(data []byte) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	log := v.Logger.With().Str("location", v.pathToFile).Logger()
	encrypted, err := v.encrypt(data)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to encrypt data to be written out.")
		return err
	}
	return ioutil.WriteFile(v.pathToFile, []byte(encrypted), 0400)
}

// encrypt uses an aes cipher provided by the certificate file for encryption.
// We don't store the password anywhere. An error will be thrown in case the encryption
// operation encounters a problem. Gaia uses AES GCM to encrypt the vault file. For Nonce it's
// using a constantly increasing number which is stored with the file. GCM allows for better
// password verification in which case we don't have to guess what was wrong any longer.
// In the end we encrypt the whole thing to Base64 for ease of saving an handling.
func (v *FileStorer) encrypt(data []byte) (string, error) {
	if len(data) < 1 {
		// User has deleted all the secrets. the file will be empty.
		return "", nil
	}
	key := v.key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	v.counter++
	nonce := make([]byte, 12)
	binary.LittleEndian.PutUint64(nonce, v.counter)
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	hexNonce := hex.EncodeToString(nonce)
	hexChiperText := hex.EncodeToString(ciphertext)
	content := fmt.Sprintf("%s||%s", hexNonce, hexChiperText)
	finalMsg := hex.EncodeToString([]byte(content))
	return finalMsg, nil
}

func (v *FileStorer) decrypt(encodedData []byte) ([]byte, error) {
	if len(encodedData) < 1 {
		v.Logger.Debug().Msg("the vault is empty")
		return nil, nil
	}
	key := v.key
	decodedMsg, err := hex.DecodeString(string(encodedData))
	if err != nil {
		v.Logger.Debug().Err(err).Msg("Failed to decode encoded data.")
		return nil, err
	}
	split := strings.Split(string(decodedMsg), "||")
	if len(split) < 2 {
		v.Logger.Error().Strs("split", split).Msg("Invalid data format.")
		return nil, errors.New("invalid number of splits")
	}
	nonce, err := hex.DecodeString(split[0])
	if err != nil {
		v.Logger.Debug().Err(err).Msg("Failed to decode nonce")
		return nil, err
	}
	data, err := hex.DecodeString(split[1])
	if err != nil {
		v.Logger.Debug().Err(err).Msg("Failed to decode data")
		return nil, err
	}
	v.counter = binary.LittleEndian.Uint64(nonce)
	block, err := aes.NewCipher(key)
	if err != nil {
		v.Logger.Debug().Err(err).Msg("Failed to create new Cipher")
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		v.Logger.Debug().Err(err).Msg("Failed to create new GCM block")
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		v.Logger.Debug().Err(err).Msg("Failed to Open data.")
		return nil, err
	}
	return plaintext, nil
}
