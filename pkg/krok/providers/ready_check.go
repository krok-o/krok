package providers

import "context"

// Ready provides a ready check for Krok.
type Ready interface {
	Ready(ctx context.Context) bool
}
