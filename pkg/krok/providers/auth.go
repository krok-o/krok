package providers

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/krok-o/krok/pkg/models"
)

// Auth defines the capabilities of a repository authentication storage framework.
type Auth interface {
	// GetRepositoryAuth returns auth data for a repository.
	GetRepositoryAuth(ctx context.Context, id int) (*models.Auth, error)
	// CreateRepositoryAuth creates auth data for a repository in vault.
	CreateRepositoryAuth(ctx context.Context, repositoryID int, info *models.Auth) error
}

// ApiKeysAuthenticator deals with authenticating api keys.
type ApiKeysAuthenticator interface {
	// Match matches a given user's api keys with the stored ones.
	Match(ctx context.Context, key *models.APIKey) error
	// Encrypt takes an api key secret and encrypts it for storage.
	Encrypt(ctx context.Context, secret []byte) ([]byte, error)
}

type OAuthProvider interface {
	GetAuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GenerateState() (string, error)
	VerifyState(rawToken string) error
}
