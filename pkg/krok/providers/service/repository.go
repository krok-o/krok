package service

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
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

// RepositoryServiceDependencies represents the RepositoryService dependencies.
type RepositoryServiceDependencies struct {
	Logger zerolog.Logger
	Storer providers.RepositoryStorer
}

// RepositoryService is the gRPC server implementation for repository interactions.
type RepositoryService struct {
	RepositoryServiceConfig
	RepositoryServiceDependencies

	repov1.UnimplementedRepositoryServiceServer
}

// NewRepositoryService creates a new RepositoryService.
func NewRepositoryService(cfg RepositoryServiceConfig, deps RepositoryServiceDependencies) *RepositoryService {
	return &RepositoryService{RepositoryServiceConfig: cfg, RepositoryServiceDependencies: deps}
}

// CreateRepository creates a repository.
func (s *RepositoryService) CreateRepository(ctx context.Context, request *repov1.CreateRepositoryRequest) (*repov1.Repository, error) {
	log := s.Logger.Debug()

	repository, err := s.Storer.Create(ctx, &models.Repository{
		Name: request.Name,
		URL:  request.Url,
		VCS:  int(request.Vcs),
	})
	if err != nil {
		log.Err(err).Msg("error creating repo in store")
		return nil, status.Error(codes.Internal, "failed to create repository")
	}

	uurl, err := s.generateURL(repository)
	if err != nil {
		log.Err(err).Msg("error generating url")
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
	id := request.GetId()
	if request.Id == nil {
		s.Logger.Debug().Msg("missing id")
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	log := s.Logger.Debug().Int32("id", id.GetValue())

	repository, err := s.Storer.Update(ctx, &models.Repository{
		ID:   int(id.GetValue()),
		Name: request.Name,
	})
	if err != nil {
		log.Err(err).Msg("error updating repo in store")
		return nil, status.Error(codes.Internal, "failed to update repository")
	}

	uurl, err := s.generateURL(repository)
	if err != nil {
		log.Err(err).Msg("error generating url")
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
	id := request.GetId()
	if request.Id == nil {
		s.Logger.Debug().Msg("missing id")
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	log := s.Logger.Debug().Int32("id", id.GetValue())

	repository, err := s.Storer.Get(ctx, int(id.GetValue()))
	if err != nil {
		log.Err(err).Msg("error getting repo from store")
		return nil, status.Error(codes.Internal, "failed to get repository")
	}

	uurl, err := s.generateURL(repository)
	if err != nil {
		log.Err(err).Msg("error generating url")
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
	log := s.Logger.Debug()

	repositories, err := s.Storer.List(ctx, &models.ListOptions{
		Name: request.Name,
		VCS:  int(request.Vcs),
	})
	if err != nil {
		log.Err(err).Msg("error listing repos from store")
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
	id := request.GetId()
	if request.Id == nil {
		s.Logger.Debug().Msg("missing id")
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	log := s.Logger.Debug().Int32("id", id.GetValue())

	if err := s.Storer.Delete(ctx, int(request.Id.Value)); err != nil {
		log.Err(err).Msg("error deleting repo from store")
		return nil, status.Error(codes.Internal, "failed to delete repository")
	}

	return &empty.Empty{}, nil
}

// generateURL generates the unique URL for the repository.
func (s *RepositoryService) generateURL(repo *models.Repository) (string, error) {
	u, err := url.Parse(s.Hostname)
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	u.Path = path.Join(u.Path, "hooks", strconv.Itoa(repo.ID), strconv.Itoa(repo.VCS), "callback")
	return u.String(), nil
}
