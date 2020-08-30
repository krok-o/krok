package providers

import "context"

// Loader defines the ability to load in plugins.
type Loader interface {
	// Watches a folder for new plugins/commands to load.
	Watch(ctx context.Context) error
	// Load will load a command / plugin.
	Load(ctx context.Context, f string) error
}
