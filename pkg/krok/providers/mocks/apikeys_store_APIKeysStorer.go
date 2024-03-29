// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// APIKeysStorer is an autogenerated mock type for the APIKeysStorer type
type APIKeysStorer struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, key
func (_m *APIKeysStorer) Create(ctx context.Context, key *models.APIKey) (*models.APIKey, error) {
	ret := _m.Called(ctx, key)

	var r0 *models.APIKey
	if rf, ok := ret.Get(0).(func(context.Context, *models.APIKey) *models.APIKey); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.APIKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.APIKey) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id, userID
func (_m *APIKeysStorer) Delete(ctx context.Context, id int, userID int) error {
	ret := _m.Called(ctx, id, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, id, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id, userID
func (_m *APIKeysStorer) Get(ctx context.Context, id int, userID int) (*models.APIKey, error) {
	ret := _m.Called(ctx, id, userID)

	var r0 *models.APIKey
	if rf, ok := ret.Get(0).(func(context.Context, int, int) *models.APIKey); ok {
		r0 = rf(ctx, id, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.APIKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, id, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByAPIKeyID provides a mock function with given fields: ctx, id
func (_m *APIKeysStorer) GetByAPIKeyID(ctx context.Context, id string) (*models.APIKey, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.APIKey
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.APIKey); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.APIKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, userID
func (_m *APIKeysStorer) List(ctx context.Context, userID int) ([]*models.APIKey, error) {
	ret := _m.Called(ctx, userID)

	var r0 []*models.APIKey
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.APIKey); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.APIKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
