package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/krok"
)

// Loader defines the ability to load in plugins.
type Loader interface {
	// Watch a folder for new plugins/commands to save.
	// If a file appears in the watched folder, it will be picked up and saved into the commands.
	Watch(ctx context.Context) error
	// Load will load a plugin from a given location.
	Load(ctx context.Context, location string) (krok.Plugin, error)
}
