package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
)

func TestTokenIssuer_Create(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-01-31T15:00:00Z")

	testEmail := "test@test.com"
	createInput := &models.UserAuthDetails{Email: testEmail, FirstName: "Test", LastName: "Name"}
	expectedTokenResponse := &oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODA0ODM3MDAsImlhdCI6MTU4MDQ4MjgwMCwic3ViIjoiMSJ9.apom8FiBl_QEfRYVkp-PDETLFzAdEFzVZLVMqrkj6Uc`,
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODc2MDAsImlhdCI6MTU4MDQ4MjgwMCwic3ViIjoiMSJ9.U3ocf3xQv8r5bzbr3l9IwAnCpqMDkOfdsNxkUktINSU`,
		Expiry:       time.Unix(1580483700, 0),
	}

	t.Run("user is in the database", func(t *testing.T) {
		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(now)

		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("GetByEmail", mock.Anything, testEmail).Return(&models.User{ID: 1}, nil)

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			Clock:     mockClock,
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockClock.AssertExpectations(t)
		mockUserStorer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, token, expectedTokenResponse)
	})

	t.Run("user is not in database", func(t *testing.T) {
		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(now)

		mockUserStorer := &mocks.UserStorer{}
		qerr := &kerr.QueryError{Err: kerr.ErrNotFound}
		mockUserStorer.On("GetByEmail", mock.Anything, testEmail).Return(nil, qerr)
		mockUserStorer.On("Create", mock.Anything, &models.User{
			Email:       testEmail,
			DisplayName: "Test Name",
		}).Return(&models.User{
			ID:          1,
			Email:       testEmail,
			DisplayName: "Test Name",
		}, nil)

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			Clock:     mockClock,
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockClock.AssertExpectations(t)
		mockUserStorer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, token, expectedTokenResponse)
	})

	t.Run("get user returns an unexpected error", func(t *testing.T) {
		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("GetByEmail", mock.Anything, testEmail).Return(nil, errors.New("unexpected error"))

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockUserStorer.AssertExpectations(t)
		assert.EqualError(t, err, "get user: unexpected error")
		assert.Nil(t, token)
	})

	t.Run("create user returns an unexpected error", func(t *testing.T) {
		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("GetByEmail", mock.Anything, testEmail).Return(nil, &kerr.QueryError{
			Err: kerr.ErrNotFound,
		})
		mockUserStorer.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("unexpected error"))

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockUserStorer.AssertExpectations(t)
		assert.EqualError(t, err, "create user: unexpected error")
		assert.Nil(t, token)
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
