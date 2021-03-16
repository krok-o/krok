package providers

import (
	"context"
)

// Watcher defines the ability to watch and load in plugins.
type Watcher interface {
	// Run starts the watcher. Watches a folder for new plugins/commands to save.
	// If a file appears in the watched folder, it will be picked up and saved into the commands.
	Run(ctx context.Context)
}
