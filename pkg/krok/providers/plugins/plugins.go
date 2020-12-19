package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
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
	failureTry := time.Second * 15
	for {
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			return p.Watch(ctx)
		})
		if err := g.Wait(); err != nil {
			p.Logger.
				Error().
				Err(err).
				Msgf("Failed to start the watcher or watcher encountered an error. Try again in %s.",
					failureTry.String())
		}
		select {
		case <-time.After(failureTry):
			// try starting the watcher again
		case <-ctx.Done():
			return
		}
	}
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
			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				if err := p.handleCreateEvent(ctx, event, log); err != nil {
					log.Error().Err(err).Msg("Failed to handle create event.")
				}
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				if err := p.handleRemoveEvent(ctx, event, log); err != nil {
					log.Error().Err(err).Msg("Failed to handle remove event.")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("errors channel closed")
			}
			log.Warn().Err(err).Msg("Error from file watcher.")
		}
	}
}

// handleCreateEvent will handle a create event from the file system. Generally
// these are non-blocking events and can be re-tried by doing the same steps again.
func (p *GoPlugins) handleCreateEvent(ctx context.Context, event fsnotify.Event, log zerolog.Logger) error {
	file := event.Name
	log = log.With().Str("file", file).Logger()

	log.Debug().Msg("New file added.")
	hash, err := p.generateHash(file)
	if err != nil || hash == "" {
		log.Debug().Err(err).Str("hash", hash).Msg("Failed to generate hash for the file.")
		return err
	}
	// TODO: Check if exists and disabled. If hash==thisnewhash enable the plugin.
	if _, err := p.Store.Create(ctx, &models.Command{
		Name:     filepath.Base(file),
		Filename: file,
		Location: p.Location,
		Hash:     hash,
		Enabled:  true,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to add new command.")
	}
	return nil
}

// handleRemoveEvent will handle a remove event from the file system.
func (p *GoPlugins) handleRemoveEvent(ctx context.Context, event fsnotify.Event, log zerolog.Logger) error {
	file := event.Name
	log = log.With().Str("file", file).Logger()

	log.Debug().Msg("File deleted. Disabling plugin.")
	hash, err := p.generateHash(file)
	if err != nil || hash == "" {
		log.Debug().Err(err).Str("hash", hash).Msg("Failed to generate hash for the file.")
		return err
	}
	name := path.Base(file)
	command, err := p.Store.GetByName(ctx, name)
	if errors.Is(err, kerr.NotFound) {
		// no command with this name, nothing to do.
		return nil
	} else if err != nil {
		log.Debug().Err(err).Msg("GetByName failed")
		return err
	}

	command.Enabled = false
	if _, err := p.Store.Update(ctx, command); err != nil {
		log.Debug().Err(err).Msg("Update command failed")
		return err
	}
	return nil
}

// generateHash generates a hash for a file.
func (p *GoPlugins) generateHash(file string) (string, error) {
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
