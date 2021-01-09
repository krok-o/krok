package providers

import "context"

// Hooks defines what the hooks can do, which is mostly just Execute.
// Gets the raw payload
// return outcome, success, error
type Hooks interface {
	// Execute will be called for the command which can be executed.
	// opts is a variable number of arguments which can be given to a hook.
	// exp.: Environment properties, auth information, tokens, etc.
	Execute(ctx context.Context, raw string, opts ...interface{}) (string, bool, error)
}
