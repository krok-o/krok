package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
	userv1 "github.com/krok-o/krok/proto/user/v1"
)

const (
	defaultAPIKeyName = "My API Key"
	defaultAPIKeyTTL  = 7 * 24 * time.Hour
)

// UserAPIKeyServiceDependencies represents the services dependencies.
type UserAPIKeyServiceDependencies struct {
	Logger        zerolog.Logger
	Storer        providers.APIKeysStorer
	Authenticator providers.ApiKeysAuthenticator
	UUID          providers.UUIDGenerator
	Clock         providers.Clock
}

// UserAPIKeyService represents the user api key service.
type UserAPIKeyService struct {
	UserAPIKeyServiceDependencies

	userv1.UnimplementedAPIKeyServiceServer
}

// CreateAPIKey creates a user api key pair.
func (s *UserAPIKeyService) CreateAPIKey(ctx context.Context, request *userv1.CreateAPIKeyRequest) (*userv1.APIKey, error) {
	userId := request.GetUserId()
	if request.UserId == nil {
		s.Logger.Info().Msg("missing user_id")
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	name := request.Name
	if name == "" {
		name = defaultAPIKeyName
	}

	log := s.Logger.Log().
		Int32("user_id", userId.GetValue()).
		Str("name", name)

	keySecret, err := s.UUID.Generate()
	if err != nil {
		log.Err(err).Msg("failed to generate unique key")
		return nil, status.Error(codes.Internal, "failed to generate unique key")
	}

	keyID, err := s.generateKeyID()
	if err != nil {
		log.Err(err).Msg("failed to generate unique key id")
		return nil, status.Error(codes.Internal, "failed to generate unique key id")
	}

	encrypted, err := s.Authenticator.Encrypt(ctx, []byte(keySecret))
	if err != nil {
		log.Err(err).Msg("failed to encrypt unique key")
		return nil, status.Error(codes.Internal, "failed to encrypt unique key")
	}

	key := &models.APIKey{
		Name:         name,
		UserID:       int(userId.GetValue()),
		APIKeyID:     keyID,
		APIKeySecret: encrypted,
		TTL:          s.Clock.Now().Add(defaultAPIKeyTTL),
	}
	created, err := s.Storer.Create(ctx, key)
	if err != nil {
		log.Err(err).Msg("error creating api key in store")
		return nil, status.Error(codes.Internal, "failed to create api key")
	}

	ttl, err := ptypes.TimestampProto(key.TTL)
	if err != nil {
		log.Err(err).Msg("failed to create ttl timestamp")
		return nil, status.Error(codes.Internal, "failed to create ttl timestamp")
	}

	return &userv1.APIKey{
		Id:        int32(created.ID),
		Name:      key.Name,
		UserId:    int32(key.UserID),
		KeyId:     keyID,
		KeySecret: keySecret,
		Ttl:       ttl,
	}, nil
}

// DeleteAPIKey deletes a user api key pair.
func (s *UserAPIKeyService) DeleteAPIKey(ctx context.Context, request *userv1.DeleteAPIKeyRequest) (*empty.Empty, error) {
	id := request.GetId()
	if id == nil {
		s.Logger.Info().Msg("missing id")
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	userId := request.GetUserId()
	if userId == nil {
		s.Logger.Info().Msg("missing user_id")
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	log := s.Logger.Log().
		Int32("id", id.GetValue()).
		Int32("user_id", userId.GetValue())

	if err := s.Storer.Delete(ctx, int(id.GetValue()), int(userId.GetValue())); err != nil {
		log.Err(err).Msg("error deleting api key from store")
		return nil, status.Error(codes.Internal, "failed to delete api key")
	}

	return &empty.Empty{}, nil
}

// GetAPIKey gets an API key.
func (s *UserAPIKeyService) GetAPIKey(ctx context.Context, request *userv1.GetAPIKeyRequest) (*userv1.APIKey, error) {
	id := request.GetId()
	if id == nil {
		s.Logger.Info().Msg("missing id")
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	userId := request.GetUserId()
	if userId == nil {
		s.Logger.Info().Msg("missing user_id")
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	log := s.Logger.Log().
		Int32("id", id.GetValue()).
		Int32("user_id", userId.GetValue())

	key, err := s.Storer.Get(ctx, int(id.GetValue()), int(userId.GetValue()))
	if err != nil {
		log.Err(err).Msg("error getting api key from store")
		return nil, status.Error(codes.Internal, "failed to get api key")
	}

	ttl, err := ptypes.TimestampProto(key.TTL)
	if err != nil {
		log.Err(err).Msg("failed to create ttl timestamp")
		return nil, status.Error(codes.Internal, "failed to create ttl timestamp")
	}

	return &userv1.APIKey{
		Id:     int32(key.ID),
		Name:   key.Name,
		UserId: int32(key.UserID),
		KeyId:  key.APIKeyID,
		Ttl:    ttl,
	}, nil
}

// ListAPIKeys lists API keys for a user.
func (s *UserAPIKeyService) ListAPIKeys(ctx context.Context, request *userv1.ListAPIKeyRequest) (*userv1.APIKeys, error) {
	userId := request.GetUserId()
	if userId == nil {
		s.Logger.Info().Msg("missing user_id")
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	log := s.Logger.Log().
		Int32("user_id", userId.GetValue())

	keys, err := s.Storer.List(ctx, int(userId.GetValue()))
	if err != nil {
		log.Err(err).Msg("error listing api keys from store")
		return nil, status.Error(codes.Internal, "failed to list api keys")
	}

	items := make([]*userv1.APIKey, len(keys))
	for i := range items {
		items[i] = &userv1.APIKey{
			Id:     int32(keys[i].ID),
			Name:   keys[i].Name,
			UserId: int32(keys[i].UserID),
			KeyId:  keys[i].APIKeyID,
		}

		ttl, err := ptypes.TimestampProto(keys[i].TTL)
		if err != nil {
			log.Err(err).Msg("failed to create ttl timestamp")
			return nil, status.Error(codes.Internal, "failed to create ttl timestamp")
		}

		items[i].Ttl = ttl
	}

	return &userv1.APIKeys{Items: items}, nil
}

func (s *UserAPIKeyService) generateKeyID() (string, error) {
	u, err := s.UUID.Generate()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(u))), nil
}
