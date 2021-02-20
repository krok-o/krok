package auth

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const (
	tokenFormat = prefixFormat + "_VCS_TOKEN"
)

// TokenProviderDependencies defines the dependencies for the token provider.
type TokenProviderDependencies struct {
	Logger zerolog.Logger
	Vault  providers.Vault
}

// TokenProvider is the provider which saves and manages tokens for the various platforms.
type TokenProvider struct {
	TokenProviderDependencies
}

// NewPlatformTokenProvider creates a new Token provider for the platforms.
func NewPlatformTokenProvider(deps TokenProviderDependencies) *TokenProvider {
	return &TokenProvider{
		TokenProviderDependencies: deps,
	}
}

var _ providers.PlatformTokenProvider = &TokenProvider{}

// GetTokenForPlatform will retrieve the token for this VCS.
func (t *TokenProvider) GetTokenForPlatform(vcs int) (string, error) {
	log := t.Logger.With().Int("vcs", vcs).Logger()
	if err := t.Vault.LoadSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to load secrets")
		return "", fmt.Errorf("failed to get secrets: %w", err)
	}
	token, err := t.Vault.GetSecret(fmt.Sprintf(tokenFormat, vcs))
	if err != nil {
		log.Debug().Err(err).Msg("GetSecret failed for token")
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	return string(token), nil
}

// SaveTokenForPlatform will save the token for this VCS.
func (t *TokenProvider) SaveTokenForPlatform(token string, vcs int) error {
	log := t.Logger.With().Int("vcs", vcs).Logger()
	if token == "" {
		return errors.New("token is empty")
	}
	if err := t.Vault.LoadSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to load secrets")
		return fmt.Errorf("failed to get repository auth: %w", err)
	}
	log.Debug().Msg("Store token")
	t.Vault.AddSecret(fmt.Sprintf(tokenFormat, vcs), []byte(token))

	if err := t.Vault.SaveSecrets(); err != nil {
		log.Debug().Err(err).Msg("Failed to save secrets")
		return fmt.Errorf("failed to save secrets: %w", err)
	}
	return nil
}
