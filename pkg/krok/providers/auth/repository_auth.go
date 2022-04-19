package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	prefixFormat   = "%d_"
	usernameFormat = prefixFormat + "REPO_USERNAME"
	passwordFormat = prefixFormat + "REPO_PASSWORD"
	sshKeyFormat   = prefixFormat + "REPO_SSH_KEY"
	secretFormat   = prefixFormat + "REPO_SECRET"
)

// RepositoryAuthDependencies defines the dependencies for the repository auth provider.
type RepositoryAuthDependencies struct {
	Logger zerolog.Logger
	Vault  providers.Vault
}

// RepoAuth is the authentication provider for Krok repositories.
type RepoAuth struct {
	RepositoryAuthDependencies
}

// NewRepositoryAuth creates a new repository authentication provider.
func NewRepositoryAuth(deps RepositoryAuthDependencies) *RepoAuth {
	return &RepoAuth{
		RepositoryAuthDependencies: deps,
	}
}

var _ providers.RepositoryAuth = &RepoAuth{}

// CreateRepositoryAuth creates auth data for a repository in vault.
func (a *RepoAuth) CreateRepositoryAuth(ctx context.Context, repositoryID int, info *models.Auth) error {
	log := a.Logger.With().Str("func", "CreateRepositoryAuth").Int("repository_id", repositoryID).Logger()
	if info == nil {
		log.Debug().Msg("No auth information for repository. Skip storing anything.")
		return nil
	}

	if err := a.Vault.LoadSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to load secrets")
		return fmt.Errorf("failed to get repository auth: %w", err)
	}

	addIfNotEmpty := func(k, v string) {
		if v != "" {
			a.Vault.AddSecret(k, []byte(v))
		}
	}
	addIfNotEmpty(fmt.Sprintf(passwordFormat, repositoryID), info.Password)
	addIfNotEmpty(fmt.Sprintf(usernameFormat, repositoryID), info.Username)
	addIfNotEmpty(fmt.Sprintf(sshKeyFormat, repositoryID), info.SSH)
	addIfNotEmpty(fmt.Sprintf(secretFormat, repositoryID), info.Secret)

	if err := a.Vault.SaveSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to save secrets")
		return fmt.Errorf("failed to save secrets: %w", err)
	}
	return nil
}

type secretGetter struct {
	err error
}

// getSecret returns the value of a secret. If there was an error, this is a no-op.
func (s *secretGetter) getSecret(a *RepoAuth, log zerolog.Logger, secret string) []byte {
	if s.err != nil {
		return nil
	}
	value, err := a.Vault.GetSecret(secret)
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("GetSecret failed")
		s.err = fmt.Errorf("failed to get repository auth: %w", err)
	}
	return value
}

// GetRepositoryAuth returns auth data for a repository. Returns ErrNotFound if there is no
// auth info for a repository.
func (a *RepoAuth) GetRepositoryAuth(ctx context.Context, id int) (*models.Auth, error) {
	log := a.Logger.With().Int("id", id).Logger()
	if err := a.Vault.LoadSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to load secrets")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}

	getter := &secretGetter{}
	username := getter.getSecret(a, log, fmt.Sprintf(usernameFormat, id))
	password := getter.getSecret(a, log, fmt.Sprintf(passwordFormat, id))
	sshKey := getter.getSecret(a, log, fmt.Sprintf(sshKeyFormat, id))
	secret := getter.getSecret(a, log, fmt.Sprintf(secretFormat, id))
	if getter.err != nil {
		return nil, getter.err
	}

	if username == nil && password == nil && sshKey == nil && secret == nil {
		log.Debug().Msg("No auth information for given id.")
		return nil, nil
	}
	result := &models.Auth{
		SSH:      string(sshKey),
		Username: string(username),
		Password: string(password),
		Secret:   string(secret),
	}
	return result, nil
}
