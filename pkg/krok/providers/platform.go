package providers

import (
	"context"
	"net/http"

	"github.com/krok-o/krok/pkg/models"
)

// Platform defines what a platform should be able to do in order for it to
// work with hooks. Once a provider is selected when creating a repository
// given the right authorization the platform provider will create the hook on this
// repository.
type Platform interface {
	// CreateHook creates a hook for the respective platform.
	CreateHook(ctx context.Context, repo *models.Repository) error
	// ValidateRequest will take a hook and verify it being a valid hook request according to
	// platform rules.
	ValidateRequest(ctx context.Context, r *http.Request) error
}

// PlatformTokenProvider defines the operations a token provider must perform.
// A single platform will manage a single token for now. Later maybe we'll provider the
// ability to handle multiple tokens.
type PlatformTokenProvider interface {
	// Token related CRUD operations
	GetTokenForPlatform(vcs int) (string, error)
	SaveTokenForPlatform(token string, vcs int) error
	// for now, people can manually delete the secret from the vault directly.
	//DeleteTokenForPlatform(vcs int) error
}
