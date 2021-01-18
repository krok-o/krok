package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	userv1 "github.com/krok-o/krok/proto/user/v1"
)

func TestUserApiKeyService_CreateApiKey(t *testing.T) {
	ttl, err := time.Parse(time.RFC3339, "2020-01-18T15:00:00Z")
	assert.NoError(t, err)

	t.Run("successful creation of api key", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("2c038884-8076-4352-8f68-19ef9ea6a584", nil).Once()
		uuid.On("Generate").Return("64be95eb-1508-4357-99fd-d913cd703e40", nil).Once()

		authenticator := &mocks.ApiKeysAuthenticator{}
		authenticator.On("Encrypt", context.Background(), []byte("2c038884-8076-4352-8f68-19ef9ea6a584")).Return([]byte("encvalue"), nil).Once()

		clock := &mocks.Clock{}
		clock.On("Now").Return(ttl).Once()

		storer := &mocks.APIKeysStorer{}
		storer.On("Create", context.Background(), &models.APIKey{
			Name:         "Key-1",
			UserID:       1,
			APIKeyID:     "54aa4e0ce8487fbd3a2ffa79103da00a",
			APIKeySecret: []byte("encvalue"),
			TTL:          ttl.Add(defaultApiKeyTTL),
		}).Return(&models.APIKey{ID: 1234}, nil).Once()

		svc := NewUserAPIKeyService(storer, authenticator, uuid, clock)

		key, err := svc.CreateApiKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
			Name:   "Key-1",
		})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, int32(1234), key.Id)
		assert.Equal(t, "Key-1", key.Name)
		assert.Equal(t, int32(1), key.UserId)
		assert.Equal(t, "54aa4e0ce8487fbd3a2ffa79103da00a", key.KeyId)
		assert.Equal(t, "2c038884-8076-4352-8f68-19ef9ea6a584", key.KeySecret)
		assert.Equal(t, int64(1579964400), key.Ttl.GetSeconds())
	})

	t.Run("successful creation of api key with default name", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("2c038884-8076-4352-8f68-19ef9ea6a584", nil).Once()
		uuid.On("Generate").Return("64be95eb-1508-4357-99fd-d913cd703e40", nil).Once()

		authenticator := &mocks.ApiKeysAuthenticator{}
		authenticator.On("Encrypt", mock.Anything, mock.Anything).Return([]byte("encvalue"), nil).Once()

		clock := &mocks.Clock{}
		clock.On("Now").Return(ttl).Once()

		storer := &mocks.APIKeysStorer{}
		storer.On("Create", context.Background(), &models.APIKey{
			Name:         "My API Key",
			UserID:       1,
			APIKeyID:     "54aa4e0ce8487fbd3a2ffa79103da00a",
			APIKeySecret: []byte("encvalue"),
			TTL:          ttl.Add(defaultApiKeyTTL),
		}).Return(&models.APIKey{ID: 1234}, nil).Once()

		svc := NewUserAPIKeyService(storer, authenticator, uuid, clock)

		key, err := svc.CreateApiKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, int32(1234), key.Id)
		assert.Equal(t, "My API Key", key.Name)
		assert.Equal(t, int32(1), key.UserId)
		assert.Equal(t, "54aa4e0ce8487fbd3a2ffa79103da00a", key.KeyId)
		assert.Equal(t, "2c038884-8076-4352-8f68-19ef9ea6a584", key.KeySecret)
		assert.Equal(t, int64(1579964400), key.Ttl.GetSeconds())
	})

	t.Run("unique key generation error", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("", errors.New("err")).Once()

		authenticator := &mocks.ApiKeysAuthenticator{}
		clock := &mocks.Clock{}
		storer := &mocks.APIKeysStorer{}

		svc := NewUserAPIKeyService(storer, authenticator, uuid, clock)

		_, err := svc.CreateApiKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to generate unique key")
	})

	t.Run("unique key id generation error", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("2c038884-8076-4352-8f68-19ef9ea6a584", nil).Once()
		uuid.On("Generate").Return("", errors.New("err")).Once()

		authenticator := &mocks.ApiKeysAuthenticator{}
		clock := &mocks.Clock{}
		storer := &mocks.APIKeysStorer{}

		svc := NewUserAPIKeyService(storer, authenticator, uuid, clock)

		_, err := svc.CreateApiKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to generate unique key id")
	})

	t.Run("encrypt creation error", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("2c038884-8076-4352-8f68-19ef9ea6a584", nil).Once()
		uuid.On("Generate").Return("64be95eb-1508-4357-99fd-d913cd703e40", nil).Once()

		authenticator := &mocks.ApiKeysAuthenticator{}
		authenticator.On("Encrypt", mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()

		clock := &mocks.Clock{}
		clock.On("Now").Return(ttl).Once()

		storer := &mocks.APIKeysStorer{}

		svc := NewUserAPIKeyService(storer, authenticator, uuid, clock)

		_, err := svc.CreateApiKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to encrypt unique key")
	})

	t.Run("store creation error", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("2c038884-8076-4352-8f68-19ef9ea6a584", nil).Once()
		uuid.On("Generate").Return("64be95eb-1508-4357-99fd-d913cd703e40", nil).Once()

		authenticator := &mocks.ApiKeysAuthenticator{}
		authenticator.On("Encrypt", mock.Anything, mock.Anything).Return([]byte("encvalue"), nil).Once()

		clock := &mocks.Clock{}
		clock.On("Now").Return(ttl).Once()

		storer := &mocks.APIKeysStorer{}
		storer.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()

		svc := NewUserAPIKeyService(storer, authenticator, uuid, clock)

		_, err := svc.CreateApiKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to create key")
	})
}
