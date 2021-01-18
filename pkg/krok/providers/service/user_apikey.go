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
	defaultApiKeyName = "My API Key"
	defaultApiKeyTTL  = 7 * 24 * time.Hour
)

// UserAPIKeyService represents the user api key service.
type UserAPIKeyService struct {
	storer        providers.APIKeysStorer
	authenticator providers.ApiKeysAuthenticator
	uuid          providers.UUIDGenerator
	clock         providers.Clock

	userv1.UnimplementedApiKeyServiceServer
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

// CreateApiKey creates a user api key pair.
func (s *UserAPIKeyService) CreateApiKey(ctx context.Context, request *userv1.CreateAPIKeyRequest) (*userv1.ApiKey, error) {
	if request.UserId == nil {
		return nil, status.Error(codes.Internal, "missing user_id")
	}

	name := request.Name
	if name == "" {
		name = defaultApiKeyName
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
		TTL:          s.clock.Now().Add(defaultApiKeyTTL),
	}
	created, err := s.storer.Create(ctx, key)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create key")
	}

	ttl, err := ptypes.TimestampProto(key.TTL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create ttl timestamp")
	}

	return &userv1.ApiKey{
		Id:        int32(created.ID),
		Name:      key.Name,
		UserId:    int32(key.UserID),
		KeyId:     keyID,
		KeySecret: keySecret,
		Ttl:       ttl,
	}, nil
}

// DeleteApiKey deletes a user api key pair.
func (s *UserAPIKeyService) DeleteApiKey(ctx context.Context, request *userv1.DeleteAPIKeyRequest) (*empty.Empty, error) {
	if request.GetId() == nil {
		return nil, status.Error(codes.Internal, "missing id")
	}

	if request.GetUserId() == nil {
		return nil, status.Error(codes.Internal, "missing user_id")
	}

	if err := s.storer.Delete(ctx, int(request.Id.Value), int(request.UserId.Value)); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete api key")
	}

	return &empty.Empty{}, nil
}

func (s *UserAPIKeyService) generateKeyID() (string, error) {
	u, err := s.uuid.Generate()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(u))), nil
}
