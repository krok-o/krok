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

	t.Run("user is in the cache", func(t *testing.T) {
		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(now)

		mockUserCache := &mocks.UserCache{}
		mockUserCache.On("Has", testEmail).Return(&providers.UserCacheItem{User: &models.User{ID: 1}}, true)
		mockUserCache.On("Add", testEmail, 1)

		mockUserStorer := &mocks.UserStorer{}

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			Clock:     mockClock,
			UserCache: mockUserCache,
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockUserCache.AssertExpectations(t)
		mockClock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, token, expectedTokenResponse)
	})

	t.Run("user is not in the cache but is in the database", func(t *testing.T) {
		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(now)

		mockUserCache := &mocks.UserCache{}
		mockUserCache.On("Has", testEmail).Return(nil, false)
		mockUserCache.On("Add", testEmail, 1)

		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("GetByEmail", mock.Anything, testEmail).Return(&models.User{ID: 1}, nil)

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			Clock:     mockClock,
			UserCache: mockUserCache,
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockUserCache.AssertExpectations(t)
		mockClock.AssertExpectations(t)
		mockUserStorer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, token, expectedTokenResponse)
	})

	t.Run("user is not in cache or database", func(t *testing.T) {
		mockClock := &mocks.Clock{}
		mockClock.On("Now").Return(now)

		mockUserCache := &mocks.UserCache{}
		mockUserCache.On("Has", testEmail).Return(nil, false)
		mockUserCache.On("Add", testEmail, 1)

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
			UserCache: mockUserCache,
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockUserCache.AssertExpectations(t)
		mockClock.AssertExpectations(t)
		mockUserStorer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, token, expectedTokenResponse)
	})

	t.Run("get user returns an unexpected error", func(t *testing.T) {
		mockUserCache := &mocks.UserCache{}
		mockUserCache.On("Has", testEmail).Return(nil, false)

		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("GetByEmail", mock.Anything, testEmail).Return(nil, errors.New("unexpected error"))

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			UserCache: mockUserCache,
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockUserCache.AssertExpectations(t)
		mockUserStorer.AssertExpectations(t)
		assert.EqualError(t, err, "get user: unexpected error")
		assert.Nil(t, token)
	})

	t.Run("create user returns an unexpected error", func(t *testing.T) {
		mockUserCache := &mocks.UserCache{}
		mockUserCache.On("Has", testEmail).Return(nil, false)

		mockUserStorer := &mocks.UserStorer{}
		mockUserStorer.On("GetByEmail", mock.Anything, testEmail).Return(nil, &kerr.QueryError{
			Err: kerr.ErrNotFound,
		})
		mockUserStorer.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("unexpected error"))

		ti := NewTokenIssuer(TokenIssuerConfig{
			GlobalTokenKey: "test",
		}, TokenIssuerDependencies{
			UserCache: mockUserCache,
			UserStore: mockUserStorer,
		})

		token, err := ti.Create(context.Background(), createInput)
		mockUserCache.AssertExpectations(t)
		mockUserStorer.AssertExpectations(t)
		assert.EqualError(t, err, "create user: unexpected error")
		assert.Nil(t, token)
	})
}
