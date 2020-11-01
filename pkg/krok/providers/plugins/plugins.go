package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"os"
	"path/filepath"
	"plugin"

	"github.com/krok-o/krok/pkg/models"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok"
	"github.com/krok-o/krok/pkg/krok/providers"
)

// Config has the configuration options for the plugins.
type Config struct {
	// Location is the folder to watch.
	Location string
}

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger zerolog.Logger
	Store  providers.CommandStorer
}

// GoPlugins is a plugin handler which uses basic Go plugins and the plugins package.
type GoPlugins struct {
	Config
	Dependencies
}

// NewGoPluginsProvider creates a new Go based plugin provider.
// Starts the folder watcher.
func NewGoPluginsProvider(ctx context.Context, cfg Config, deps Dependencies) (*GoPlugins, error) {
	p := &GoPlugins{Config: cfg, Dependencies: deps}
	if _, err := os.Stat(cfg.Location); os.IsNotExist(err) {
		deps.Logger.Err(err).Str("location", cfg.Location).Msg("Location does not exist.")
		return nil, err
	}
	go p.run(ctx)
	return p, nil
}

// run start the watcher and run until context is done.
func (p *GoPlugins) run(ctx context.Context) {
	// start an error group waiter... if there was an error from the watcher, it was
	// fatal and we should try again after a little while to start it.

	<-ctx.Done()
}

// Watch a folder for new plugins/commands to load.
// If a file appears in the watched folder, it will be picked up and saved into the commands.
func (p *GoPlugins) Watch(ctx context.Context) error {
	log := p.Logger.With().Str("location", p.Location).Logger()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to start folder watcher.")
		return err
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Debug().Err(err).Msg("Failed to close watcher.")
		}
	}()
	if err := watcher.Add(p.Location); err != nil {
		log.Debug().Err(err).Msg("Failed to add folder to watcher.")
		return err
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Debug().Err(err).Msg("Events channel closed.")
				return errors.New("events channel closed")
			}
			if event.Op&fsnotify.Create != fsnotify.Create {
				// We only care about the create action.
				break
			}
			file := event.Name
			log.Debug().Str("file", file).Msg("New file added.")
			hasher, err := p.generateHash(file, log)
			id, err := krok.GenerateUUID()
			if err != nil {
				log.Debug().Err(err).Str("file", file).Msg("Failed to generate new ID for resource.")
				break
			}
			if _, err := p.Store.Create(ctx, &models.Command{
				Name:     filepath.Base(file),
				ID:       id,
				Filename: file,
				Location: p.Location,
				Hash:     ,
				Enabled:  true,
			}); err != nil {
				// Log the error but move on.
				log.Error().Err(err).Msg("Failed to add new command.")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("errors channel closed")
			}
			log.Warn().Err(err).Msg("Error from file watcher.")
		}
	}
}

// generateHash generates a hash for a file.
func (p *GoPlugins) generateHash(file string) (string, error) {
	log := p.Logger.With().Str("file", file).Logger()
	hasher := sha256.New()
	f, err := os.Open(file)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to open file.")
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Debug().Err(err).Msg("Failed to hash file.")
		return nil, err
	}
	return hex.EncodeToString(hasher.Sum(nil)), err
}

// Load will load a plugin from a given location.
// This will be called on demand given a location to a plugin when the command is about to be executed.
func (p *GoPlugins) Load(ctx context.Context, location string) (krok.Plugin, error) {
	log := p.Logger.With().Str("location", p.Location).Logger()
	plug, err := plugin.Open(location)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to open Plugin.")
		return nil, err
	}

	symPlugin, err := plug.Lookup("Execute")
	if err != nil {
		log.Warn().Err(err).Msg("Failed to lookup Symbol Execute.")
		return nil, err
	}

	krokPlugin, ok := symPlugin.(krok.Plugin)
	if !ok {
		log.Warn().Err(err).Msg("Loaded plugin is not of type Krok.Plugin.")
		return nil, err
	}
	return krokPlugin, nil
}
