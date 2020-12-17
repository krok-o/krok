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
	prefixFormat   = "%s_"
	usernameFormat = prefixFormat + "_REPO_USERNAME"
	passwordFormat = prefixFormat + "_REPO_PASSWORD"
	sshKeyFormat   = prefixFormat + "_REPO_SSH_KEY"
)

// Config has the configuration options for the vault.
type Config struct {
}

// Dependencies defines the dependencies for the plugin provider.
type Dependencies struct {
	Logger zerolog.Logger
	Vault  providers.Vault
}

// KrokAuth is the authentication provider for Krok.
type KrokAuth struct {
	Config
	Dependencies
}

// NewKrokAuth creates a new Krok authentication provider.
func NewKrokAuth(cfg Config, deps Dependencies) (*KrokAuth, error) {
	return &KrokAuth{
		Config:       cfg,
		Dependencies: deps,
	}, nil
}

// GetRepositoryAuth returns auth data for a repository. Returns NotFound if there is no
// auth info for a repository.
func (a *KrokAuth) GetRepositoryAuth(ctx context.Context, id string) (*models.Auth, error) {
	log := a.Logger.With().Str("id", id).Logger()
	if err := a.Vault.LoadSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to load secrets")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}
	username, err := a.Vault.GetSecret(fmt.Sprintf(usernameFormat, id))
	if !errors.Is(err, kerr.NotFound) {
		log.Debug().Err(err).Msg("GetSecret failed for username")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}

	password, err := a.Vault.GetSecret(fmt.Sprintf(passwordFormat, id))
	if !errors.Is(err, kerr.NotFound) {
		log.Debug().Err(err).Msg("GetSecret failed for password")
		return nil, fmt.Errorf("failed to get repository auth: %w", err)
	}

	sshKey, err := a.Vault.GetSecret(fmt.Sprintf(sshKeyFormat, id))
	if !errors.Is(err, kerr.NotFound) {
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
