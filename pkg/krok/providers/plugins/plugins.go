package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cirello.io/pglock"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config has the configuration options for the plugins.
type Config struct {
	// Location is the folder to put the plugin into.
	Location string
}

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger zerolog.Logger
	Store  providers.CommandStorer
	Tar    providers.Tar
}

// Plugins is a plugin handler which uses basic Go plugins and the plugins package.
type Plugins struct {
	Config
	Dependencies
}

// NewPluginsProvider creates a new Go based plugin provider.
// Starts the folder watcher.
func NewPluginsProvider(cfg Config, deps Dependencies) *Plugins {
	return &Plugins{Config: cfg, Dependencies: deps}
}

// Create will handle creating a plugin, including un-taring and copying the plugin into the right location.
// It returns the hash of the command. Creating the command is not its responsibility.
func (p *Plugins) Create(ctx context.Context, file string) (string, error) {
	log := p.Logger.With().Str("file", file).Logger()

	l, err := p.Store.AcquireLock(ctx, file)
	if err != nil {
		if errors.Is(err, pglock.ErrNotAcquired) {
			log.Debug().Msg("Some other process is already handling this file's create event.")
			return "", nil
		}
		return "", err
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Debug().Err(err).Msg("Failed to release lock...")
		}
	}()

	// un-tar
	dir := filepath.Dir(file)
	base := filepath.Base(file)
	f, err := os.Open(file)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to open archive.")
		return "", err
	}
	if err := p.Tar.Untar(dir, f); err != nil {
		log.Debug().Err(err).Msg("Failed to un-tar archive.")
		return "", err
	}
	dots := strings.Split(base, ".")
	if len(dots) < 1 {
		return "", errors.New("no extensions found in filename")
	}
	extractedFile := filepath.Join(dir, dots[0])
	dst := filepath.Join(p.Location, dots[0])

	// Copy file to permanent storage
	if err := p.Copy(extractedFile, dst); err != nil {
		log.Debug().Err(err).Str("extracted_file", extractedFile).Msg("Failed to copy file to permanent storage.")
		return "", err
	}

	log.Debug().Msg("New file added.")
	hash, err := p.generateHash(dst)
	if err != nil || hash == "" {
		log.Debug().Err(err).Str("hash", hash).Msg("Failed to generate hash for the file.")
		return "", err
	}
	log.Debug().Msg("New command successfully created.")
	return hash, nil
}

// Copy the src file to dst. Any existing file will be overwritten and will not
func (p *Plugins) Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		p.Logger.Debug().Err(err).Str("src", src).Msg("Failed to open source.")
		return err
	}
	defer func(in *os.File) {
		if err := in.Close(); err != nil {
			p.Logger.Debug().Err(err).Str("in", in.Name()).Msg("Failed to close source file.")
		}
	}(in)

	out, err := os.Create(dst)
	if err != nil {
		p.Logger.Debug().Err(err).Str("dst", dst).Msg("Failed to create destination.")
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// Delete will handle deleting a plugin from permanent storage.
func (p *Plugins) Delete(ctx context.Context, name string) error {
	// bail early if file doesn't exist so we don't acquire the lock.
	stat, err := os.Stat(filepath.Join(p.Location, name))
	if err != nil {
		p.Logger.Debug().Err(err).Str("name", name).Msg("Failed to find file with name in permanent storage.")
		return err
	}
	file := stat.Name()
	log := log.With().Str("file", file).Logger()
	l, err := p.Store.AcquireLock(ctx, file)
	if err != nil {
		if errors.Is(err, pglock.ErrNotAcquired) {
			log.Debug().Msg("Some other process is already handling this file's delete event.")
			return nil
		}
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Debug().Err(err).Msg("Failed to release lock...")
		}
	}()

	if err := os.Remove(file); err != nil {
		log.Debug().Err(err).Msg("Failed to remove file.")
		return err
	}
	log.Debug().Msg("File deleted.")
	return nil
}

// generateHash generates a hash for a file.
func (p *Plugins) generateHash(file string) (string, error) {
	log := p.Logger.With().Str("file", file).Logger()
	hasher := sha256.New()
	f, err := os.Open(file)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to open file.")
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Debug().Err(err).Msg("Failed to close file descriptor.")
		}
	}()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Debug().Err(err).Msg("Failed to hash file.")
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), err
}
