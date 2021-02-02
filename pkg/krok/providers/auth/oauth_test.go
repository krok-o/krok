package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
)

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
