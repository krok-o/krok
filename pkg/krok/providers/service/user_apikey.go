package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
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

// UserApiKeyService represents the user api key service.
type UserApiKeyService struct {
	storer        providers.ApiKeysStorer
	authenticator providers.ApiKeysAuthenticator
	uuid          providers.UUIDGenerator
	clock         providers.Clock

	userv1.UnimplementedApiKeyServiceServer
}

// NewUserApiKeyService creates a new UserApiKeyService
func NewUserApiKeyService(
	storer providers.ApiKeysStorer,
	authenticator providers.ApiKeysAuthenticator,
	uuid providers.UUIDGenerator,
	time providers.Clock,
) *UserApiKeyService {
	return &UserApiKeyService{
		storer:        storer,
		authenticator: authenticator,
		uuid:          uuid,
		clock:         time,
	}
}

// CreateApiKey creates a user api key.
func (s *UserApiKeyService) CreateApiKey(ctx context.Context, request *userv1.CreateApiKeyRequest) (*userv1.ApiKey, error) {
	if request.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "must provide user_id")
	}

	uid, err := strconv.Atoi(request.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert user_id")
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
		UserID:       uid,
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
		UserId:    int32(uid),
		KeyId:     keyID,
		KeySecret: keySecret,
		Ttl:       ttl,
	}, nil
}

func (s *UserApiKeyService) generateKeyID() (string, error) {
	u, err := s.uuid.Generate()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(u))), nil
}
