package github

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type mockAuthProvider struct {
	providers.RepositoryAuth
	getAuth *models.Auth
}

func (mp *mockAuthProvider) GetRepositoryAuth(ctx context.Context, id int) (*models.Auth, error) {
	return mp.getAuth, nil
}

type mockPlatformTokenProvider struct {
	providers.PlatformTokenProvider
}

func (mptp *mockPlatformTokenProvider) GetTokenForPlatform(vcs int) (string, error) {
	return "token", nil
}

type mockGithubRepositoryService struct {
	Hook     *github.Hook
	Response *github.Response
	Error    error
	Owner    string
	Repo     string
}

func (mgc *mockGithubRepositoryService) CreateHook(ctx context.Context, owner, repo string, hook *github.Hook) (*github.Hook, *github.Response, error) {
	if owner != mgc.Owner {
		return nil, nil, errors.New("owner did not equal expected owner: was: " + owner)
	}
	if repo != mgc.Repo {
		return nil, nil, errors.New("repo did not equal expected repo: was: " + repo)
	}
	return mgc.Hook, mgc.Response, mgc.Error
}

func TestGithub_CreateHook(t *testing.T) {
	cliLogger := zerolog.New(os.Stderr)
	mp := &mockAuthProvider{
		getAuth: &models.Auth{
			SSH:      "ssh",
			Username: "username",
			Password: "password",
			Secret:   "secret",
		},
	}
	mptp := &mockPlatformTokenProvider{}
	npp := NewGithubPlatformProvider(Dependencies{
		Logger:                cliLogger,
		PlatformTokenProvider: mptp,
		AuthProvider:          mp,
	})
	mock := &mockGithubRepositoryService{}
	mock.Hook = &github.Hook{
		Name: github.String("test hook"),
		URL:  github.String("https://api.github.com/repos/krok-o/krok/hooks/44321286/test"),
	}
	mock.Response = &github.Response{
		Response: &http.Response{
			Status: "Ok",
		},
	}
	mock.Owner = "krok-o"
	mock.Repo = "krok"
	npp.repoMock = mock
	err := npp.CreateHook(context.Background(), &models.Repository{
		Name:      "test",
		ID:        0,
		URL:       "https://github.com/krok-o/krok",
		VCS:       models.GITHUB,
		Auth:      &models.Auth{Secret: "secret"},
		UniqueURL: "https://krok.com/hooks/0/0/callback",
		Events:    []string{"push"},
	})
	assert.NoError(t, err)
}

func TestGithub_CreateHook_InvalidURL(t *testing.T) {
	cliLogger := zerolog.New(os.Stderr)
	mp := &mockAuthProvider{
		getAuth: &models.Auth{
			SSH:      "ssh",
			Username: "username",
			Password: "password",
			Secret:   "secret",
		},
	}
	mptp := &mockPlatformTokenProvider{}
	npp := NewGithubPlatformProvider(Dependencies{
		Logger:                cliLogger,
		PlatformTokenProvider: mptp,
		AuthProvider:          mp,
	})
	mock := &mockGithubRepositoryService{}
	mock.Hook = &github.Hook{
		Name: github.String("test hook"),
		URL:  github.String("https://api.github.com/repos/krok-o/krok/hooks/44321286/test"),
	}
	mock.Response = &github.Response{
		Response: &http.Response{
			Status: "Ok",
		},
	}
	mock.Owner = "krok-o"
	mock.Repo = "krok"
	npp.repoMock = mock
	err := npp.CreateHook(context.Background(), &models.Repository{
		Name:      "test",
		ID:        0,
		URL:       "https://github.com/krok-o",
		VCS:       models.GITHUB,
		Auth:      &models.Auth{Secret: "secret"},
		UniqueURL: "https://krok.com/hooks/0/0/callback",
		Events:    []string{"push"},
	})
	assert.Error(t, err)
}

func TestGithub_GetEventID(t *testing.T) {
	cliLogger := zerolog.New(os.Stderr)
	mp := &mockAuthProvider{
		getAuth: &models.Auth{
			SSH:      "ssh",
			Username: "username",
			Password: "password",
			Secret:   "secret",
		},
	}
	mptp := &mockPlatformTokenProvider{}
	npp := NewGithubPlatformProvider(Dependencies{
		Logger:                cliLogger,
		PlatformTokenProvider: mptp,
		AuthProvider:          mp,
	})
	header := http.Header{}
	header.Add("X-GitHub-Delivery", "ID")
	id, err := npp.GetEventID(context.Background(), &http.Request{
		Header: header,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ID", id)

	_, err = npp.GetEventID(context.Background(), &http.Request{})
	assert.Errorf(t, err, "event id not found for request")
}
