// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// RepositoryStorer is an autogenerated mock type for the RepositoryStorer type
type RepositoryStorer struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, c
func (_m *RepositoryStorer) Create(ctx context.Context, c *models.Repository) (*models.Repository, error) {
	ret := _m.Called(ctx, c)

	var r0 *models.Repository
	if rf, ok := ret.Get(0).(func(context.Context, *models.Repository) *models.Repository); ok {
		r0 = rf(ctx, c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Repository) error); ok {
		r1 = rf(ctx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *RepositoryStorer) Delete(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *RepositoryStorer) Get(ctx context.Context, id int) (*models.Repository, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Repository
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Repository); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: ctx, name
func (_m *RepositoryStorer) GetByName(ctx context.Context, name string) (*models.Repository, error) {
	ret := _m.Called(ctx, name)

	var r0 *models.Repository
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Repository); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, opt
func (_m *RepositoryStorer) List(ctx context.Context, opt *models.ListOptions) ([]*models.Repository, error) {
	ret := _m.Called(ctx, opt)

	var r0 []*models.Repository
	if rf, ok := ret.Get(0).(func(context.Context, *models.ListOptions) []*models.Repository); ok {
		r0 = rf(ctx, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.ListOptions) error); ok {
		r1 = rf(ctx, opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, c
func (_m *RepositoryStorer) Update(ctx context.Context, c *models.Repository) (*models.Repository, error) {
	ret := _m.Called(ctx, c)

	var r0 *models.Repository
	if rf, ok := ret.Get(0).(func(context.Context, *models.Repository) *models.Repository); ok {
		r0 = rf(ctx, c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Repository) error); ok {
		r1 = rf(ctx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
