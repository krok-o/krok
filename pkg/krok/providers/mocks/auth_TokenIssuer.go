// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/krok-o/krok/pkg/models"
	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// TokenIssuer is an autogenerated mock type for the TokenIssuer type
type TokenIssuer struct {
	mock.Mock
}

// Create provides a mock function with given fields: token
func (_m *TokenIssuer) Create(token *models.User) (*oauth2.Token, error) {
	ret := _m.Called(token)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(*models.User) *oauth2.Token); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Refresh provides a mock function with given fields: ctx, refreshToken
func (_m *TokenIssuer) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
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
