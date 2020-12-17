package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

type Auth interface {
	// GetRepositoryAuth returns auth data for a repository.
	GetRepositoryAuth(ctx context.Context, id string) (*models.Auth, error)
}
