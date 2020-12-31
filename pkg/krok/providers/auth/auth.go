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
)

// AuthConfig has the configuration options for the vault.
type AuthConfig struct {
}

// AuthDependencies defines the dependencies for the auth provider.
type AuthDependencies struct {
	Logger zerolog.Logger
	Vault  providers.Vault
}

// KrokAuth is the authentication provider for Krok.
type KrokAuth struct {
	AuthConfig
	AuthDependencies
}

// NewKrokAuth creates a new Krok authentication provider.
func NewKrokAuth(cfg AuthConfig, deps AuthDependencies) (*KrokAuth, error) {
	return &KrokAuth{
		AuthConfig:       cfg,
		AuthDependencies: deps,
	}, nil
}

var _ providers.Auth = &KrokAuth{}

// CreateRepositoryAuth creates auth data for a repository in vault.
func (a *KrokAuth) CreateRepositoryAuth(ctx context.Context, repositoryID int, info *models.Auth) error {
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

	return nil
}

// GetRepositoryAuth returns auth data for a repository. Returns ErrNotFound if there is no
// auth info for a repository.
func (a *KrokAuth) GetRepositoryAuth(ctx context.Context, id int) (*models.Auth, error) {
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

	if username == nil && password == nil && sshKey == nil {
		log.Debug().Msg("No auth information for given id.")
		return nil, nil
	}
	result := &models.Auth{
		SSH:      string(sshKey),
		Username: string(username),
		Password: string(password),
	}
	return result, nil
}
