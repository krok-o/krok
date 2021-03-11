package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/krok"
)

// Watcher defines the ability to watch and load in plugins.
type Watcher interface {
	// Run starts the watcher. Watches a folder for new plugins/commands to save.
	// If a file appears in the watched folder, it will be picked up and saved into the commands.
	Run(ctx context.Context)
	// Load will load a plugin from a given location.
	Load(ctx context.Context, location string) (krok.Execute, error)
}
