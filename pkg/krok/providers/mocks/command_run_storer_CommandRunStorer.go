// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// CommandRunStorer is an autogenerated mock type for the CommandRunStorer type
type CommandRunStorer struct {
	mock.Mock
}

// CreateRun provides a mock function with given fields: ctx, run
func (_m *CommandRunStorer) CreateRun(ctx context.Context, run *models.CommandRun) (*models.CommandRun, error) {
	ret := _m.Called(ctx, run)

	var r0 *models.CommandRun
	if rf, ok := ret.Get(0).(func(context.Context, *models.CommandRun) *models.CommandRun); ok {
		r0 = rf(ctx, run)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CommandRun)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.CommandRun) error); ok {
		r1 = rf(ctx, run)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRunStatus provides a mock function with given fields: ctx, id, status, outcome
func (_m *CommandRunStorer) UpdateRunStatus(ctx context.Context, id int, status string, outcome string) error {
	ret := _m.Called(ctx, id, status, outcome)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string, string) error); ok {
		r0 = rf(ctx, id, status, outcome)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
