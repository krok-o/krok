// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// UserTokenIssuer is an autogenerated mock type for the UserTokenIssuer type
type UserTokenIssuer struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, token
func (_m *UserTokenIssuer) Create(ctx context.Context, token *models.UserAuthDetails) (*oauth2.Token, error) {
	ret := _m.Called(ctx, token)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(context.Context, *models.UserAuthDetails) *oauth2.Token); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.UserAuthDetails) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Refresh provides a mock function with given fields: ctx, refreshToken
func (_m *UserTokenIssuer) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	ret := _m.Called(ctx, refreshToken)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(context.Context, string) *oauth2.Token); ok {
		r0 = rf(ctx, refreshToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, refreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
