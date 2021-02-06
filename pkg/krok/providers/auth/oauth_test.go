package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"gopkg.in/h2non/gock.v1"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func TestOAuthAuthenticator_Exchange(t *testing.T) {
	t.Run("exchange token for existing user", func(t *testing.T) {
		setupGock()
		defer gock.Off()

		mockUserStore := &mocks.UserStorer{}
		mockUserStore.On("GetByEmail", mock.Anything, "test@test.com").Return(&models.User{ID: 1, Email: "test@test.com"}, nil)

		mockTokenIssuer := &mocks.TokenIssuer{}
		mockTokenIssuer.On("Create", &models.User{ID: 1, Email: "test@test.com"}).Return(&oauth2.Token{}, nil)

		auth := NewOAuthAuthenticator(OAuthAuthenticatorConfig{
			BaseURL: "https://test.com",
		}, OAuthAuthenticatorDependencies{
			Issuer:    mockTokenIssuer,
			UserStore: mockUserStore,
		})

		token, err := auth.Exchange(context.Background(), "1234")
		assert.NoError(t, err)
		assert.Equal(t, &oauth2.Token{}, token)
		mockUserStore.AssertExpectations(t)
		mockTokenIssuer.AssertExpectations(t)
	})

	t.Run("exchange token for a new user", func(t *testing.T) {
		setupGock()
		defer gock.Off()

		mockUserStore := &mocks.UserStorer{}
		qerr := &kerr.QueryError{Err: kerr.ErrNotFound}
		mockUserStore.On("GetByEmail", mock.Anything, "test@test.com").Return(nil, qerr)
		mockUserStore.On("Create", mock.Anything, &models.User{
			Email:       "test@test.com",
			DisplayName: "Test User",
		}).Return(&models.User{
			ID:          1,
			Email:       "test@test.com",
			DisplayName: "Test User",
		}, nil)

		mockTokenIssuer := &mocks.TokenIssuer{}
		mockTokenIssuer.On("Create", &models.User{ID: 1, Email: "test@test.com", DisplayName: "Test User"}).Return(&oauth2.Token{}, nil)

		auth := NewOAuthAuthenticator(OAuthAuthenticatorConfig{
			BaseURL: "https://test.com",
		}, OAuthAuthenticatorDependencies{
			Issuer:    mockTokenIssuer,
			UserStore: mockUserStore,
		})

		token, err := auth.Exchange(context.Background(), "1234")
		assert.NoError(t, err)
		assert.Equal(t, &oauth2.Token{}, token)
		mockUserStore.AssertExpectations(t)
		mockTokenIssuer.AssertExpectations(t)
	})
}

func setupGock() {
	gock.New("https://oauth2.googleapis.com").
		Post("/token").
		MatchHeader("Content-Type", "application/x-www-form-urlencoded").
		AddMatcher(func(r *http.Request, _ *gock.Request) (bool, error) {
			switch {
			case r.FormValue("code") != "1234":
				return false, fmt.Errorf("unexpected code %s", r.FormValue("code"))
			case r.FormValue("grant_type") != "authorization_code":
				return false, fmt.Errorf("unexpected grant_type %s", r.FormValue("grant_type"))
			case r.FormValue("redirect_uri") != "https://test.com/auth/callback":
				return false, fmt.Errorf("unexpected redirect_uri %s", r.FormValue("redirect_uri"))
			default:
				return true, nil
			}
		}).
		Reply(200).
		JSON(&oauth2.Token{
			AccessToken:  "aaaaaaa",
			RefreshToken: "rrrrrrr",
		})

	gock.New("https://www.googleapis.com").
		Get("/oauth2/v2/userinfo").
		MatchParam("access_token", "aaaaaaa").
		Reply(200).
		JSON(googleUser{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@test.com",
		})
}

func TestOAuthAuthenticator_GenerateAndVerifyState(t *testing.T) {
	t.Run("generate uuid error", func(t *testing.T) {
		mockUUID := &mocks.UUIDGenerator{}
		mockUUID.On("Generate").Return("", errors.New("uuid err"))

		auth := NewOAuthAuthenticator(OAuthAuthenticatorConfig{
			GlobalTokenKey: "test",
		}, OAuthAuthenticatorDependencies{
			UUID: mockUUID,
		})

		_, err := auth.GenerateState("https://test.com")
		require.EqualError(t, err, "uuid generate: uuid err")

		mockUUID.AssertExpectations(t)
	})

	t.Run("generate and verify valid state", func(t *testing.T) {
		mockUUID := &mocks.UUIDGenerator{}
		mockUUID.On("Generate").Return("3d29d8ca-a836-48de-be74-469268660a34", nil)

		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(time.Now())

		auth := NewOAuthAuthenticator(OAuthAuthenticatorConfig{
			GlobalTokenKey: "test",
		}, OAuthAuthenticatorDependencies{
			UUID:  mockUUID,
			Clock: mockClock,
		})

		state, err := auth.GenerateState("https://test.com")
		require.NoError(t, err)

		redirectURL, err := auth.VerifyState(state)
		require.NoError(t, err)
		require.Equal(t, "https://test.com", redirectURL)

		mockUUID.AssertExpectations(t)
		mockClock.AssertExpectations(t)
	})

	t.Run("generate and verify expired token state", func(t *testing.T) {
		mockUUID := &mocks.UUIDGenerator{}
		mockUUID.On("Generate").Return("3d29d8ca-a836-48de-be74-469268660a34", nil)

		expired, err := time.Parse(time.RFC3339, "2019-01-31T15:00:00Z")
		require.NoError(t, err)
		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(expired)

		auth := NewOAuthAuthenticator(OAuthAuthenticatorConfig{
			GlobalTokenKey: "test",
		}, OAuthAuthenticatorDependencies{
			UUID:  mockUUID,
			Clock: mockClock,
		})

		state, err := auth.GenerateState("https://test.com")
		require.NoError(t, err)

		redirectURL, err := auth.VerifyState(state)
		require.Error(t, err)
		require.Contains(t, err.Error(), "parse token: token is expired")
		require.Equal(t, "", redirectURL)

		mockUUID.AssertExpectations(t)
		mockClock.AssertExpectations(t)
	})
}
