// Code generated by mockery v2.8.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

// CancelRun provides a mock function with given fields: ctx, id
func (_m *Executor) CancelRun(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateRun provides a mock function with given fields: ctx, event, commands
func (_m *Executor) CreateRun(ctx context.Context, event *models.Event, commands []*models.Command) error {
	ret := _m.Called(ctx, event, commands)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Event, []*models.Command) error); ok {
		r0 = rf(ctx, event, commands)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
