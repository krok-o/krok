package providers

import (
	"context"

	"github.com/krok-o/krok/pkg/models"
)

// Auth defines the capabilities of a repository authentication storage framework.
type Auth interface {
	// GetRepositoryAuth returns auth data for a repository.
	GetRepositoryAuth(ctx context.Context, id int) (*models.Auth, error)
	// CreateRepositoryAuth creates auth data for a repository in vault.
	CreateRepositoryAuth(ctx context.Context, repositoryID int, info *models.Auth) error
}
