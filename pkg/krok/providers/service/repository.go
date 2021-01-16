package service

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
	repov1 "github.com/krok-o/krok/proto/repository/v1"
)

// RepositoryServiceConfig represents the RepositoryService config.
type RepositoryServiceConfig struct {
	Hostname string
}

// RepositoryService is the gRPC server implementation for repository interactions.
type RepositoryService struct {
	config RepositoryServiceConfig
	storer providers.RepositoryStorer

	repov1.UnimplementedRepositoryServiceServer
}

// NewRepositoryService creates a new RepositoryService.
func NewRepositoryService(config RepositoryServiceConfig, storer providers.RepositoryStorer) *RepositoryService {
	return &RepositoryService{config: config, storer: storer}
}

// CreateRepository creates a repository.
func (s *RepositoryService) CreateRepository(ctx context.Context, request *repov1.CreateRepositoryRequest) (*repov1.Repository, error) {
	repository, err := s.storer.Create(ctx, &models.Repository{
		Name: request.Name,
		URL:  request.Url,
		VCS:  int(request.Vcs),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create repository")
	}

	url, err := s.generateURL(repository)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate url")
	}
	repository.UniqueURL = url

	return &repov1.Repository{
		Id:        int32(repository.ID),
		Name:      repository.Name,
		Url:       repository.URL,
		Vcs:       int32(repository.VCS),
		UniqueUrl: repository.UniqueURL,
	}, nil
}

// GetRepository gets a repository.
func (s *RepositoryService) GetRepository(ctx context.Context, request *repov1.GetRepositoryRequest) (*repov1.Repository, error) {
	id := request.Id
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert id to number")
	}

	repository, err := s.storer.Get(ctx, n)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get repository")
	}

	url, err := s.generateURL(repository)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate url")
	}
	repository.UniqueURL = url

	response := &repov1.Repository{
		Id:        int32(repository.ID),
		Name:      repository.Name,
		Url:       repository.URL,
		Vcs:       int32(repository.VCS),
		UniqueUrl: repository.UniqueURL,
	}
	return response, nil
}

// generateURL generates the unique URL for the repository.
func (s *RepositoryService) generateURL(repo *models.Repository) (string, error) {
	u, err := url.Parse(s.config.Hostname)
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	u.Path = path.Join(u.Path, strconv.Itoa(repo.ID), strconv.Itoa(repo.VCS), "callback")
	return u.String(), nil
}
