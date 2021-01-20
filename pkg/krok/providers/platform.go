package providers

import (
	"context"
	"net/http"
)

// Platform defines what a platform should be able to do in order for it to
// work with hooks. Once a provider is selected when creating a repository
// given the right authorization the platform provider will create the hook on this
// repository.
type Platform interface {
	// CreateHook creates a hook for the respective platform.
	CreateHook(ctx context.Context) error
	// ValidateRequest will take a hook and verify it being a valid hook request according to
	// platform rules.
	ValidateRequest(ctx context.Context, r *http.Request) error
}
