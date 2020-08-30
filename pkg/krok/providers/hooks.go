package providers

import "context"

// Hooks defines what the hooks can do, which is mostly just Execute.
// Gets the raw payload
// return outcome, success, error
type Hooks interface {
	Execute(ctx context.Context, raw string) (string, bool, error)
}
