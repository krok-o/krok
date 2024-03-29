// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// CommandStorer is an autogenerated mock type for the CommandStorer type
type CommandStorer struct {
	mock.Mock
}

// AddCommandRelForPlatform provides a mock function with given fields: ctx, commandID, platformID
func (_m *CommandStorer) AddCommandRelForPlatform(ctx context.Context, commandID int, platformID int) error {
	ret := _m.Called(ctx, commandID, platformID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, commandID, platformID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddCommandRelForRepository provides a mock function with given fields: ctx, commandID, repositoryID
func (_m *CommandStorer) AddCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error {
	ret := _m.Called(ctx, commandID, repositoryID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, commandID, repositoryID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, c
func (_m *CommandStorer) Create(ctx context.Context, c *models.Command) (*models.Command, error) {
	ret := _m.Called(ctx, c)

	var r0 *models.Command
	if rf, ok := ret.Get(0).(func(context.Context, *models.Command) *models.Command); ok {
		r0 = rf(ctx, c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Command) error); ok {
		r1 = rf(ctx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSetting provides a mock function with given fields: ctx, settings
func (_m *CommandStorer) CreateSetting(ctx context.Context, settings *models.CommandSetting) (*models.CommandSetting, error) {
	ret := _m.Called(ctx, settings)

	var r0 *models.CommandSetting
	if rf, ok := ret.Get(0).(func(context.Context, *models.CommandSetting) *models.CommandSetting); ok {
		r0 = rf(ctx, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CommandSetting)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.CommandSetting) error); ok {
		r1 = rf(ctx, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *CommandStorer) Delete(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSetting provides a mock function with given fields: ctx, id
func (_m *CommandStorer) DeleteSetting(ctx context.Context, id int) error {
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
func (_m *CommandStorer) Get(ctx context.Context, id int) (*models.Command, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Command
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Command); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Command)
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
func (_m *CommandStorer) GetByName(ctx context.Context, name string) (*models.Command, error) {
	ret := _m.Called(ctx, name)

	var r0 *models.Command
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Command); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Command)
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

// GetSetting provides a mock function with given fields: ctx, id
func (_m *CommandStorer) GetSetting(ctx context.Context, id int) (*models.CommandSetting, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.CommandSetting
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.CommandSetting); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CommandSetting)
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

// IsPlatformSupported provides a mock function with given fields: ctx, commandID, platformID
func (_m *CommandStorer) IsPlatformSupported(ctx context.Context, commandID int, platformID int) (bool, error) {
	ret := _m.Called(ctx, commandID, platformID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int, int) bool); ok {
		r0 = rf(ctx, commandID, platformID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, commandID, platformID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, opts
func (_m *CommandStorer) List(ctx context.Context, opts *models.ListOptions) ([]*models.Command, error) {
	ret := _m.Called(ctx, opts)

	var r0 []*models.Command
	if rf, ok := ret.Get(0).(func(context.Context, *models.ListOptions) []*models.Command); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSettings provides a mock function with given fields: ctx, commandID
func (_m *CommandStorer) ListSettings(ctx context.Context, commandID int) ([]*models.CommandSetting, error) {
	ret := _m.Called(ctx, commandID)

	var r0 []*models.CommandSetting
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.CommandSetting); ok {
		r0 = rf(ctx, commandID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.CommandSetting)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, commandID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveCommandRelForPlatform provides a mock function with given fields: ctx, commandID, platformID
func (_m *CommandStorer) RemoveCommandRelForPlatform(ctx context.Context, commandID int, platformID int) error {
	ret := _m.Called(ctx, commandID, platformID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, commandID, platformID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveCommandRelForRepository provides a mock function with given fields: ctx, commandID, repositoryID
func (_m *CommandStorer) RemoveCommandRelForRepository(ctx context.Context, commandID int, repositoryID int) error {
	ret := _m.Called(ctx, commandID, repositoryID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, commandID, repositoryID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, c
func (_m *CommandStorer) Update(ctx context.Context, c *models.Command) (*models.Command, error) {
	ret := _m.Called(ctx, c)

	var r0 *models.Command
	if rf, ok := ret.Get(0).(func(context.Context, *models.Command) *models.Command); ok {
		r0 = rf(ctx, c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Command)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Command) error); ok {
		r1 = rf(ctx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSetting provides a mock function with given fields: ctx, setting
func (_m *CommandStorer) UpdateSetting(ctx context.Context, setting *models.CommandSetting) error {
	ret := _m.Called(ctx, setting)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.CommandSetting) error); ok {
		r0 = rf(ctx, setting)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
