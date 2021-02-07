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
	usernameFormat = prefixFormat + "_REPO_USERNAME"
	passwordFormat = prefixFormat + "_REPO_PASSWORD"
	sshKeyFormat   = prefixFormat + "_REPO_SSH_KEY"
	secretFormat   = prefixFormat + "_REPO_SECRET"
)

// RepositoryAuthConfig has the configuration options for the repository auth.
type RepositoryAuthConfig struct {
}

// RepositoryAuthDependencies defines the dependencies for the repository auth provider.
type RepositoryAuthDependencies struct {
	Logger zerolog.Logger
	Vault  providers.Vault
}

// RepoAuth is the authentication provider for Krok repositories.
type RepoAuth struct {
	RepositoryAuthConfig
	RepositoryAuthDependencies
}

// NewRepositoryAuth creates a new repository authentication provider.
func NewRepositoryAuth(cfg RepositoryAuthConfig, deps RepositoryAuthDependencies) (*RepoAuth, error) {
	return &RepoAuth{
		RepositoryAuthConfig:       cfg,
		RepositoryAuthDependencies: deps,
	}, nil
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
	if info.Password != "" {
		log.Debug().Msg("Store password")
		a.Vault.AddSecret(fmt.Sprintf(passwordFormat, repositoryID), []byte(info.Password))
	}
	if info.Username != "" {
		log.Debug().Msg("Store username")
		a.Vault.AddSecret(fmt.Sprintf(usernameFormat, repositoryID), []byte(info.Username))
	}
	if info.SSH != "" {
		log.Debug().Msg("Store ssh key")
		a.Vault.AddSecret(fmt.Sprintf(sshKeyFormat, repositoryID), []byte(info.SSH))
	}
	if info.Secret != "" {
		log.Debug().Msg("Store hook secret")
		a.Vault.AddSecret(fmt.Sprintf(secretFormat, repositoryID), []byte(info.Secret))
	}
	if err := a.Vault.SaveSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to save secrets")
		return fmt.Errorf("failed to save secrets: %w", err)
	}
	return nil
}

// GetRepositoryAuth returns auth data for a repository. Returns ErrNotFound if there is no
// auth info for a repository.
func (a *RepoAuth) GetRepositoryAuth(ctx context.Context, id int) (*models.Auth, error) {
	log := a.Logger.With().Int("id", id).Logger()
	if err := a.Vault.LoadSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to load secrets")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}
	username, err := a.Vault.GetSecret(fmt.Sprintf(usernameFormat, id))
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("GetSecret failed for username")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}

	password, err := a.Vault.GetSecret(fmt.Sprintf(passwordFormat, id))
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("GetSecret failed for password")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}

	sshKey, err := a.Vault.GetSecret(fmt.Sprintf(sshKeyFormat, id))
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("GetSecret failed sshKey")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}

	secret, err := a.Vault.GetSecret(fmt.Sprintf(secretFormat, id))
	if err != nil && !errors.Is(err, kerr.ErrNotFound) {
		log.Debug().Err(err).Msg("GetSecret failed secret")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
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
