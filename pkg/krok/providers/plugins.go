package providers

import (
	"context"
)

// Plugins handles create and delete events on the file system for concrete plugins / commands.
type Plugins interface {
	Create(ctx context.Context, src string) (string, error)
	Delete(ctx context.Context, name string) error
}
