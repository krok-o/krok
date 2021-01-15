package service

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/krok-o/krok/api"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

type RepositoryService struct {
	// TODO: Use dependency object?
	RepositoryStorer providers.RepositoryStorer
	URLGenerator     providers.URLGenerator

	api.UnimplementedRepositoryServiceServer
}

func NewRepositoryService(repositoryStorer providers.RepositoryStorer, urlGenerator providers.URLGenerator) *RepositoryService {
	return &RepositoryService{
		RepositoryStorer: repositoryStorer,
		URLGenerator:     urlGenerator,
	}
}

func (r *RepositoryService) CreateRepository(ctx context.Context, repository *api.Repository) (*api.Repository, error) {
	// TODO: Get user from token.

	vcs, ok := models.VCS[repository.Vcs.String()]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Invalid VSC provided.")
	}

	repo := &models.Repository{
		Name: repository.Name,
		URL:  repository.Url,
		VCS:  vcs,
	}

	created, err := r.RepositoryStorer.Create(ctx, repo)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create repository.")
	}

	url, err := r.URLGenerator.Generate(created)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate url.")
	}
	created.UniqueURL = url

	return &api.Repository{
		Id:        int32(created.ID),
		Vcs:       api.VCS(vcs + 1), // TODO: remove... ew...
		Name:      created.Name,
		Url:       created.URL,
		UniqueUrl: created.UniqueURL,
	}, nil
}

func (r *RepositoryService) GetRepository(ctx context.Context, request *api.GetRepositoryRequest) (*api.Repository, error) {
	// TODO: Get user from token.

	id := request.Id
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert id to number")
	}

	repo, err := r.RepositoryStorer.Get(ctx, n)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get repository")
	}

	// TODO: Populate full response.
	response := &api.Repository{
		Id:   int32(repo.ID),
		Name: repo.Name,
	}
	return response, nil
}
