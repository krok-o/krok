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
	// Events define the events this hook subscribes to. Since we don't want all hooks
	// to subscribe to all events all the time, we provide the option to the user
	// to select the events.
	CreateHook(ctx context.Context, repo *models.Repository) error
	// ValidateRequest will take a hook and verify it being a valid hook request according to
	// platform rules.
	ValidateRequest(ctx context.Context, r *http.Request, repoID int) error
	// GetEventID Based on the platform, retrieve the ID of the event.
	GetEventID(ctx context.Context, r *http.Request) (string, error)
	// GetEventType Based on the platform, retrieve the Type of the event.
	GetEventType(ctx context.Context, r *http.Request) (string, error)
}

// PlatformTokenProvider defines the operations a token provider must perform.
// A single platform will manage a single token for now. Later maybe we'll provider the
// ability to handle multiple tokens.
type PlatformTokenProvider interface {
	GetTokenForPlatform(vcs int) (string, error)
	SaveTokenForPlatform(token string, vcs int) error
	// for now, people can manually delete the secret from the vault directly.
	//DeleteTokenForPlatform(vcs int) error
}
