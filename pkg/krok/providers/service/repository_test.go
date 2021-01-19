package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	"github.com/krok-o/krok/pkg/models"
	repov1 "github.com/krok-o/krok/proto/repository/v1"
)

func TestRepositoryService_CreateRepository(t *testing.T) {
	t.Run("successful creation of repository", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Create", context.Background(), &models.Repository{
			Name: "test",
			URL:  "test-url",
			VCS:  models.BITBUCKET,
		}).Return(&models.Repository{
			ID:   1,
			Name: "test",
			URL:  "test-url",
			VCS:  models.BITBUCKET,
		}, nil).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		created, err := svc.CreateRepository(context.Background(), &repov1.CreateRepositoryRequest{
			Name: "test",
			Url:  "test-url",
			Vcs:  models.BITBUCKET,
		})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), created.Id)
		assert.Equal(t, "test", created.Name)
		assert.Equal(t, "test-url", created.Url)
		assert.Equal(t, "hostname/1/4/callback", created.UniqueUrl)
	})

	t.Run("store create repository error", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		created, err := svc.CreateRepository(context.Background(), &repov1.CreateRepositoryRequest{
			Name: "test",
			Url:  "test-url",
			Vcs:  models.BITBUCKET,
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to create repository")
		assert.Nil(t, created)
	})
}

func TestRepositoryService_UpdateRepository(t *testing.T) {
	t.Run("successful update of repository", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Update", context.Background(), &models.Repository{
			ID:   1,
			Name: "new-name",
		}).Return(&models.Repository{
			ID:   1,
			Name: "test",
			URL:  "test-url",
			VCS:  models.BITBUCKET,
		}, nil).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		created, err := svc.UpdateRepository(context.Background(), &repov1.UpdateRepositoryRequest{
			Id:   wrapperspb.Int32(1),
			Name: "new-name",
		})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), created.Id)
		assert.Equal(t, "test", created.Name)
		assert.Equal(t, "test-url", created.Url)
		assert.Equal(t, "hostname/1/4/callback", created.UniqueUrl)
	})

	t.Run("store update repository error", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		created, err := svc.UpdateRepository(context.Background(), &repov1.UpdateRepositoryRequest{
			Id:   wrapperspb.Int32(1),
			Name: "new-name",
		})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to update repository")
		assert.Nil(t, created)
	})
}

func TestRepositoryService_GetRepository(t *testing.T) {
	t.Run("successful creation of repository", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Get", context.Background(), 1234).Return(&models.Repository{
			ID:   1234,
			Name: "test",
			URL:  "test-url",
			VCS:  models.BITBUCKET,
		}, nil).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		repository, err := svc.GetRepository(context.Background(), &repov1.GetRepositoryRequest{Id: wrapperspb.Int32(1234)})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, int32(1234), repository.Id)
		assert.Equal(t, "test", repository.Name)
		assert.Equal(t, "test-url", repository.Url)
		assert.Equal(t, "hostname/1234/4/callback", repository.UniqueUrl)
	})

	t.Run("store get repository error", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Get", context.Background(), 1234).Return(nil, errors.New("err")).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		repository, err := svc.GetRepository(context.Background(), &repov1.GetRepositoryRequest{Id: wrapperspb.Int32(1234)})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to get repository")
		assert.Nil(t, repository)
	})
}

func TestRepositoryService_ListRepositories(t *testing.T) {
	t.Run("successful listing of repositories", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("List", context.Background(), &models.ListOptions{VCS: models.BITBUCKET}).Return([]*models.Repository{
			{
				ID:   1,
				Name: "test-1",
				URL:  "test-url",
				VCS:  models.BITBUCKET,
			},
			{
				ID:   2,
				Name: "test-2",
				URL:  "test-url",
				VCS:  models.BITBUCKET,
			},
		}, nil).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		repos, err := svc.ListRepositories(context.Background(), &repov1.ListRepositoryRequest{Vcs: models.BITBUCKET})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Len(t, repos.Items, 2)
		assert.Equal(t, int32(1), repos.Items[0].Id)
		assert.Equal(t, "test-1", repos.Items[0].Name)
		assert.Equal(t, "test-url", repos.Items[0].Url)
		assert.Equal(t, int32(4), repos.Items[0].Vcs)
		assert.Equal(t, int32(2), repos.Items[1].Id)
		assert.Equal(t, "test-2", repos.Items[1].Name)
		assert.Equal(t, "test-url", repos.Items[1].Url)
		assert.Equal(t, int32(4), repos.Items[1].Vcs)
	})

	t.Run("store list repositories error", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("List", context.Background(), mock.Anything).Return(nil, errors.New("err")).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		repository, err := svc.ListRepositories(context.Background(), &repov1.ListRepositoryRequest{Name: "test"})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to list repositories")
		assert.Nil(t, repository)
	})
}

func TestRepositoryService_DeleteRepository(t *testing.T) {
	t.Run("successful deletion of repository", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Delete", context.Background(), 1234).Return(nil).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		_, err := svc.DeleteRepository(context.Background(), &repov1.DeleteRepositoryRequest{Id: wrapperspb.Int32(1234)})
		storer.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("store delete repository error", func(t *testing.T) {
		storer := &mocks.RepositoryStorer{}
		storer.On("Delete", context.Background(), 1234).Return(errors.New("err")).Once()

		svc := NewRepositoryService(RepositoryServiceConfig{Hostname: "hostname"}, storer)

		_, err := svc.DeleteRepository(context.Background(), &repov1.DeleteRepositoryRequest{Id: wrapperspb.Int32(1234)})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to delete repository")
	})
}
