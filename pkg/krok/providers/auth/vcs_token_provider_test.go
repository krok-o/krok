package auth

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/vault"
	"github.com/krok-o/krok/pkg/models"
)

func TestTokenProvider_SaveTokenForPlatform(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	location, _ := ioutil.TempDir("", "TestTokenProvider_SaveTokenForPlatform")
	fileStore := filevault.NewFileStorer(filevault.Config{
		Location: location,
		Key:      "password123",
	}, filevault.Dependencies{Logger: logger})
	err := fileStore.Init()
	assert.NoError(t, err)
	v := vault.NewKrokVault(vault.Dependencies{Logger: logger, Storer: fileStore})
	tp := NewPlatformTokenProvider(TokenProviderDependencies{
		Logger: logger,
		Vault:  v,
	})
	_, err = tp.GetTokenForPlatform(99)
	assert.Error(t, err)
	err = tp.SaveTokenForPlatform("github_token", models.GITHUB)
	assert.NoError(t, err)
	token, err := tp.GetTokenForPlatform(models.GITHUB)
	assert.NoError(t, err)
	assert.Equal(t, "github_token", token)
}
