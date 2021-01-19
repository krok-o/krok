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
			TTL:          ttl.Add(defaultAPIKeyTTL),
		}).Return(&models.APIKey{ID: 1234}, nil).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				UUID:          uuid,
				Authenticator: authenticator,
				Clock:         clock,
				Storer:        storer,
			},
		}
		key, err := svc.CreateAPIKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
			Name:   "Key-1",
		})
		uuid.AssertExpectations(t)
		authenticator.AssertExpectations(t)
		storer.AssertExpectations(t)
		clock.AssertExpectations(t)
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
			TTL:          ttl.Add(defaultAPIKeyTTL),
		}).Return(&models.APIKey{ID: 1234}, nil).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				UUID:          uuid,
				Authenticator: authenticator,
				Clock:         clock,
				Storer:        storer,
			},
		}
		key, err := svc.CreateAPIKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		uuid.AssertExpectations(t)
		authenticator.AssertExpectations(t)
		storer.AssertExpectations(t)
		clock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, int32(1234), key.Id)
		assert.Equal(t, "My API Key", key.Name)
		assert.Equal(t, int32(1), key.UserId)
		assert.Equal(t, "54aa4e0ce8487fbd3a2ffa79103da00a", key.KeyId)
		assert.Equal(t, "2c038884-8076-4352-8f68-19ef9ea6a584", key.KeySecret)
		assert.Equal(t, int64(1579964400), key.Ttl.GetSeconds())
	})

	t.Run("missing user_id in request", func(t *testing.T) {
		svc := &UserAPIKeyService{}
		_, err := svc.CreateAPIKey(context.Background(), &userv1.CreateAPIKeyRequest{})
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing user_id")
	})

	t.Run("unique key generation error", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("", errors.New("err")).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				UUID: uuid,
			},
		}
		_, err := svc.CreateAPIKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		uuid.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to generate unique key")
	})

	t.Run("unique key id generation error", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("2c038884-8076-4352-8f68-19ef9ea6a584", nil).Once()
		uuid.On("Generate").Return("", errors.New("err")).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				UUID: uuid,
			},
		}
		_, err := svc.CreateAPIKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		uuid.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to generate unique key id")
	})

	t.Run("encrypt creation error", func(t *testing.T) {
		uuid := &mocks.UUIDGenerator{}
		uuid.On("Generate").Return("2c038884-8076-4352-8f68-19ef9ea6a584", nil).Once()
		uuid.On("Generate").Return("64be95eb-1508-4357-99fd-d913cd703e40", nil).Once()

		authenticator := &mocks.ApiKeysAuthenticator{}
		authenticator.On("Encrypt", mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				UUID:          uuid,
				Authenticator: authenticator,
			},
		}
		_, err := svc.CreateAPIKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		uuid.AssertExpectations(t)
		authenticator.AssertExpectations(t)
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

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				UUID:          uuid,
				Authenticator: authenticator,
				Clock:         clock,
				Storer:        storer,
			},
		}
		_, err := svc.CreateAPIKey(context.Background(), &userv1.CreateAPIKeyRequest{
			UserId: wrapperspb.Int32(1),
		})
		uuid.AssertExpectations(t)
		authenticator.AssertExpectations(t)
		storer.AssertExpectations(t)
		clock.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to create api key")
	})
}

func TestUserAPIKeyService_DeleteAPIKey(t *testing.T) {
	t.Run("missing id in request", func(t *testing.T) {
		svc := &UserAPIKeyService{}
		_, err := svc.DeleteAPIKey(context.Background(), &userv1.DeleteAPIKeyRequest{})
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing id")
	})

	t.Run("missing user_id in request", func(t *testing.T) {
		svc := &UserAPIKeyService{}
		_, err := svc.DeleteAPIKey(context.Background(), &userv1.DeleteAPIKeyRequest{
			Id: wrapperspb.Int32(1),
		})
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing user_id")
	})

	t.Run("successfully delete api key", func(t *testing.T) {
		storer := &mocks.APIKeysStorer{}
		storer.On("Delete", context.Background(), 1, 1234).Return(nil).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				Storer: storer,
			},
		}
		_, err := svc.DeleteAPIKey(context.Background(), &userv1.DeleteAPIKeyRequest{
			Id:     wrapperspb.Int32(1),
			UserId: wrapperspb.Int32(1234),
		})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("store delete error", func(t *testing.T) {
		storer := &mocks.APIKeysStorer{}
		storer.On("Delete", context.Background(), 1, 1234).Return(errors.New("err")).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				Storer: storer,
			},
		}
		_, err := svc.DeleteAPIKey(context.Background(), &userv1.DeleteAPIKeyRequest{
			Id:     wrapperspb.Int32(1),
			UserId: wrapperspb.Int32(1234),
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to delete api key")
	})
}

func TestAPIKeyService_GetAPIKey(t *testing.T) {
	t.Run("missing id in request", func(t *testing.T) {
		svc := &UserAPIKeyService{}
		_, err := svc.GetAPIKey(context.Background(), &userv1.GetAPIKeyRequest{})
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing id")
	})

	t.Run("missing user_id in request", func(t *testing.T) {
		svc := &UserAPIKeyService{}
		_, err := svc.GetAPIKey(context.Background(), &userv1.GetAPIKeyRequest{
			Id: wrapperspb.Int32(1),
		})
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing user_id")
	})

	t.Run("successfully get api key", func(t *testing.T) {
		now := time.Now()
		storer := &mocks.APIKeysStorer{}
		storer.On("Get", context.Background(), 1, 1234).
			Return(&models.APIKey{
				ID:       1,
				UserID:   1234,
				Name:     "name",
				APIKeyID: "id",
				TTL:      now,
			}, nil).
			Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				Storer: storer,
			},
		}
		key, err := svc.GetAPIKey(context.Background(), &userv1.GetAPIKeyRequest{
			Id:     wrapperspb.Int32(1),
			UserId: wrapperspb.Int32(1234),
		})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), key.Id)
		assert.Equal(t, "name", key.Name)
		assert.Equal(t, int32(1234), key.UserId)
		assert.Equal(t, "id", key.KeyId)
		assert.Equal(t, now.Unix(), key.Ttl.GetSeconds())
	})

	t.Run("store get api key error", func(t *testing.T) {
		storer := &mocks.APIKeysStorer{}
		storer.On("Get", context.Background(), 1, 1234).Return(nil, errors.New("err")).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				Storer: storer,
			},
		}
		_, err := svc.GetAPIKey(context.Background(), &userv1.GetAPIKeyRequest{
			Id:     wrapperspb.Int32(1),
			UserId: wrapperspb.Int32(1234),
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to get api key")
	})
}

func TestAPIKeyService_ListAPIKeys(t *testing.T) {
	t.Run("missing user_id in request", func(t *testing.T) {
		svc := &UserAPIKeyService{}
		_, err := svc.ListAPIKeys(context.Background(), &userv1.ListAPIKeyRequest{})
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing user_id")
	})

	t.Run("successfully list api keys", func(t *testing.T) {
		storer := &mocks.APIKeysStorer{}
		now := time.Now()
		storer.On("List", context.Background(), 1234).
			Return([]*models.APIKey{
				{
					ID:       1,
					UserID:   1234,
					Name:     "name",
					APIKeyID: "id",
					TTL:      now,
				},
				{
					ID:       2,
					UserID:   1234,
					Name:     "name-2",
					APIKeyID: "id-2",
					TTL:      now,
				},
			}, nil).
			Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				Storer: storer,
			},
		}
		keys, err := svc.ListAPIKeys(context.Background(), &userv1.ListAPIKeyRequest{
			UserId: wrapperspb.Int32(1234),
		})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Len(t, keys.Items, 2)
		assert.Equal(t, int32(1), keys.Items[0].Id)
		assert.Equal(t, "name", keys.Items[0].Name)
		assert.Equal(t, int32(1234), keys.Items[0].UserId)
		assert.Equal(t, "id", keys.Items[0].KeyId)
		assert.Equal(t, now.Unix(), keys.Items[0].Ttl.GetSeconds())
		assert.Equal(t, int32(2), keys.Items[1].Id)
		assert.Equal(t, "name-2", keys.Items[1].Name)
		assert.Equal(t, int32(1234), keys.Items[1].UserId)
		assert.Equal(t, "id-2", keys.Items[1].KeyId)
		assert.Equal(t, now.Unix(), keys.Items[1].Ttl.GetSeconds())
	})

	t.Run("store list api keys error", func(t *testing.T) {
		storer := &mocks.APIKeysStorer{}
		storer.On("List", context.Background(), 1234).Return(nil, errors.New("err")).Once()

		svc := &UserAPIKeyService{
			UserAPIKeyServiceDependencies: UserAPIKeyServiceDependencies{
				Storer: storer,
			},
		}
		_, err := svc.ListAPIKeys(context.Background(), &userv1.ListAPIKeyRequest{
			UserId: wrapperspb.Int32(1234),
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to list api keys")
	})
}
