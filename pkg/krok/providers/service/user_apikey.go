package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
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

// UserAPIKeyService represents the user api key service.
type UserAPIKeyService struct {
	storer        providers.APIKeysStorer
	authenticator providers.ApiKeysAuthenticator
	uuid          providers.UUIDGenerator
	clock         providers.Clock

	userv1.UnimplementedAPIKeyServiceServer
}

// NewUserAPIKeyService creates a new UserAPIKeyService
func NewUserAPIKeyService(
	storer providers.APIKeysStorer,
	authenticator providers.ApiKeysAuthenticator,
	uuid providers.UUIDGenerator,
	time providers.Clock,
) *UserAPIKeyService {
	return &UserAPIKeyService{
		storer:        storer,
		authenticator: authenticator,
		uuid:          uuid,
		clock:         time,
	}
}

// CreateAPIKey creates a user api key pair.
func (s *UserAPIKeyService) CreateAPIKey(ctx context.Context, request *userv1.CreateAPIKeyRequest) (*userv1.APIKey, error) {
	if request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	name := request.Name
	if name == "" {
		name = defaultAPIKeyName
	}

	keySecret, err := s.uuid.Generate()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate unique key")
	}

	keyID, err := s.generateKeyID()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate unique key id")
	}

	encrypted, err := s.authenticator.Encrypt(ctx, []byte(keySecret))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to encrypt unique key")
	}

	key := &models.APIKey{
		Name:         name,
		UserID:       int(request.UserId.Value),
		APIKeyID:     keyID,
		APIKeySecret: encrypted,
		TTL:          s.clock.Now().Add(defaultAPIKeyTTL),
	}
	created, err := s.storer.Create(ctx, key)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create api key")
	}

	ttl, err := ptypes.TimestampProto(key.TTL)
	if err != nil {
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
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	if request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	if err := s.storer.Delete(ctx, int(request.Id.Value), int(request.UserId.Value)); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete api key")
	}

	return &empty.Empty{}, nil
}

// GetAPIKey gets an API key.
func (s *UserAPIKeyService) GetAPIKey(ctx context.Context, request *userv1.GetAPIKeyRequest) (*userv1.APIKey, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	if request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	key, err := s.storer.Get(ctx, int(request.Id.Value), int(request.UserId.Value))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get api key")
	}

	ttl, err := ptypes.TimestampProto(key.TTL)
	if err != nil {
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
	if request.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	keys, err := s.storer.List(ctx, int(request.UserId.Value))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list api keys")
	}

	list := make([]*userv1.APIKey, len(keys))
	for i := range list {
		list[i] = &userv1.APIKey{
			Id:     int32(keys[i].ID),
			Name:   keys[i].Name,
			UserId: int32(keys[i].UserID),
			KeyId:  keys[i].APIKeyID,
		}

		ttl, err := ptypes.TimestampProto(keys[i].TTL)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to create ttl timestamp")
		}

		list[i].Ttl = ttl
	}

	return &userv1.APIKeys{List: list}, nil
}

func (s *UserAPIKeyService) generateKeyID() (string, error) {
	u, err := s.uuid.Generate()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(u))), nil
}
