package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

type Auth interface {
	// GetRepositoryAuth returns auth data for a repository.
	GetRepositoryAuth(ctx context.Context, id string) (*models.Auth, error)
	// CreateRepositoryAuth creates auth data for a repository in vault.
	CreateRepositoryAuth(ctx context.Context, repositoryID string, info *models.Auth) error
}
