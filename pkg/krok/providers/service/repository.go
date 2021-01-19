package service

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
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

	uurl, err := s.generateURL(repository)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate url")
	}
	repository.UniqueURL = uurl

	return &repov1.Repository{
		Id:        int32(repository.ID),
		Name:      repository.Name,
		Url:       repository.URL,
		Vcs:       int32(repository.VCS),
		UniqueUrl: repository.UniqueURL,
	}, nil
}

// UpdateRepository updates a repository.
func (s *RepositoryService) UpdateRepository(ctx context.Context, request *repov1.UpdateRepositoryRequest) (*repov1.Repository, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	repository, err := s.storer.Update(ctx, &models.Repository{
		ID:   int(request.Id.Value),
		Name: request.Name,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update repository")
	}

	uurl, err := s.generateURL(repository)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate url")
	}
	repository.UniqueURL = uurl

	response := &repov1.Repository{
		Id:        int32(repository.ID),
		Name:      repository.Name,
		Url:       repository.URL,
		Vcs:       int32(repository.VCS),
		UniqueUrl: repository.UniqueURL,
	}
	return response, nil
}

// GetRepository gets a repository.
func (s *RepositoryService) GetRepository(ctx context.Context, request *repov1.GetRepositoryRequest) (*repov1.Repository, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	repository, err := s.storer.Get(ctx, int(request.Id.Value))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get repository")
	}

	uurl, err := s.generateURL(repository)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate url")
	}
	repository.UniqueURL = uurl

	response := &repov1.Repository{
		Id:        int32(repository.ID),
		Name:      repository.Name,
		Url:       repository.URL,
		Vcs:       int32(repository.VCS),
		UniqueUrl: repository.UniqueURL,
	}
	return response, nil
}

// ListRepositories lists repositories.
func (s *RepositoryService) ListRepositories(ctx context.Context, request *repov1.ListRepositoryRequest) (*repov1.Repositories, error) {
	repositories, err := s.storer.List(ctx, &models.ListOptions{
		Name: request.Name,
		VCS:  int(request.Vcs),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list repositories")
	}

	items := make([]*repov1.Repository, len(repositories))
	for i := range items {
		items[i] = &repov1.Repository{
			Id:   int32(repositories[i].ID),
			Name: repositories[i].Name,
			Vcs:  int32(repositories[i].VCS),
			Url:  repositories[i].URL,
		}
	}
	return &repov1.Repositories{Items: items}, nil
}

// DeleteRepository deletes a repository.
func (s *RepositoryService) DeleteRepository(ctx context.Context, request *repov1.DeleteRepositoryRequest) (*empty.Empty, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	if err := s.storer.Delete(ctx, int(request.Id.Value)); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete repository")
	}

	return &empty.Empty{}, nil
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
