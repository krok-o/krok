package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func TestTokenIssuer_Create(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-01-31T15:00:00Z")

	userInput := &models.User{ID: 1}
	expectedTokenResponse := &oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODA0ODM3MDAsImlhdCI6MTU4MDQ4MjgwMCwic3ViIjoiMSJ9.apom8FiBl_QEfRYVkp-PDETLFzAdEFzVZLVMqrkj6Uc`,
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODc2MDAsImlhdCI6MTU4MDQ4MjgwMCwic3ViIjoiMSJ9.U3ocf3xQv8r5bzbr3l9IwAnCpqMDkOfdsNxkUktINSU`,
		Expiry:       time.Unix(1580483700, 0),
	}

	t.Run("create token success with valid input", func(t *testing.T) {
		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(now)

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			Clock: mockClock,
		})

		token, err := ti.Create(userInput)
		mockClock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, token, expectedTokenResponse)
	})
}

func TestTokenIssuer_Refresh(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-01-31T15:00:00Z")
	expectedTokenResponse := &oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODA0ODM3MDAsImlhdCI6MTU4MDQ4MjgwMCwic3ViIjoiMSJ9.apom8FiBl_QEfRYVkp-PDETLFzAdEFzVZLVMqrkj6Uc`,
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODc2MDAsImlhdCI6MTU4MDQ4MjgwMCwic3ViIjoiMSJ9.U3ocf3xQv8r5bzbr3l9IwAnCpqMDkOfdsNxkUktINSU`,
		Expiry:       time.Unix(1580483700, 0),
	}

	t.Run("refresh expired token returns error", func(t *testing.T) {
		expiredToken := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODA0ODM3MDAsImlhdCI6MTU4MDQ4MjgwMCwic3ViIjoiMSJ9.apom8FiBl_QEfRYVkp-PDETLFzAdEFzVZLVMqrkj6Uc`

		cfg := TokenIssuerConfig{
			GlobalTokenKey: "test",
		}
		deps := TokenIssuerDependencies{}
		ti := NewTokenIssuer(cfg, deps)

		token, err := ti.Refresh(context.Background(), expiredToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired by")
		assert.Nil(t, token)
	})

	t.Run("get user from store error", func(t *testing.T) {
		// 50 year expiry
		validToken := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTMyMDYxOTAsImlhdCI6MTYxMjYwMTM5MCwic3ViIjoiMSJ9.fWt4SleEugX_9kCc2aNQeCkkG0_4SZhww0Z9u3ASwe0`

		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("Get", mock.Anything, 1).Return(nil, errors.New("err"))

		cfg := TokenIssuerConfig{
			GlobalTokenKey: "test",
		}
		deps := TokenIssuerDependencies{
			Clock:     providers.NewClock(),
			UserStore: mockUserStorer,
		}
		ti := NewTokenIssuer(cfg, deps)

		token, err := ti.Refresh(context.Background(), validToken)
		assert.EqualError(t, err, "user store get: err")
		assert.Nil(t, token)
		mockUserStorer.AssertExpectations(t)
	})

	t.Run("get user from store error", func(t *testing.T) {
		// 50 year expiry
		validToken := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTMyMDYxOTAsImlhdCI6MTYxMjYwMTM5MCwic3ViIjoiMSJ9.fWt4SleEugX_9kCc2aNQeCkkG0_4SZhww0Z9u3ASwe0`

		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("Get", mock.Anything, 1).Return(&models.User{ID: 1}, nil)

		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(now)

		cfg := TokenIssuerConfig{
			GlobalTokenKey: "test",
		}
		deps := TokenIssuerDependencies{
			Clock:     mockClock,
			UserStore: mockUserStorer,
		}
		ti := NewTokenIssuer(cfg, deps)

		token, err := ti.Refresh(context.Background(), validToken)
		mockClock.AssertExpectations(t)
		mockUserStorer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, token, expectedTokenResponse)
	})
}
