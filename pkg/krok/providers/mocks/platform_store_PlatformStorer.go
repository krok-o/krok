// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// PlatformStorer is an autogenerated mock type for the PlatformStorer type
type PlatformStorer struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, p
func (_m *PlatformStorer) Create(ctx context.Context, p *models.Platform) (*models.Platform, error) {
	ret := _m.Called(ctx, p)

	var r0 *models.Platform
	if rf, ok := ret.Get(0).(func(context.Context, *models.Platform) *models.Platform); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Platform)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Platform) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *PlatformStorer) Delete(ctx context.Context, id int) error {
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
func (_m *PlatformStorer) Get(ctx context.Context, id int) (*models.Platform, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Platform
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Platform); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Platform)
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
func (_m *PlatformStorer) GetByName(ctx context.Context, name string) (*models.Platform, error) {
	ret := _m.Called(ctx, name)

	var r0 *models.Platform
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Platform); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Platform)
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

// List provides a mock function with given fields: ctx, enabled
func (_m *PlatformStorer) List(ctx context.Context, enabled *bool) ([]*models.Platform, error) {
	ret := _m.Called(ctx, enabled)

	var r0 []*models.Platform
	if rf, ok := ret.Get(0).(func(context.Context, *bool) []*models.Platform); ok {
		r0 = rf(ctx, enabled)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Platform)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *bool) error); ok {
		r1 = rf(ctx, enabled)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, p
func (_m *PlatformStorer) Update(ctx context.Context, p *models.Platform) (*models.Platform, error) {
	ret := _m.Called(ctx, p)

	var r0 *models.Platform
	if rf, ok := ret.Get(0).(func(context.Context, *models.Platform) *models.Platform); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Platform)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Platform) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
