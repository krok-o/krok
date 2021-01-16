package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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

		repository, err := svc.GetRepository(context.Background(), &repov1.GetRepositoryRequest{Id: "1234"})
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

		repository, err := svc.GetRepository(context.Background(), &repov1.GetRepositoryRequest{Id: "1234"})
		storer.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Internal desc = failed to get repository")
		assert.Nil(t, repository)
	})
}
