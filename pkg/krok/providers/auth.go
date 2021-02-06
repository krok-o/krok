package providers

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/krok-o/krok/pkg/models"
)

// RepositoryAuth defines the capabilities of a repository authentication storage framework.
type RepositoryAuth interface {
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

// OAuthAuthenticator handles user authentication via OAuth2.
type OAuthAuthenticator interface {
	GetAuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GenerateState(redirectURL string) (string, error)
	VerifyState(rawToken string) (string, error)
}

// TokenIssuer handles creation of user authentication tokens.
type TokenIssuer interface {
	Create(token *models.User) (*oauth2.Token, error)
	Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error)
}
