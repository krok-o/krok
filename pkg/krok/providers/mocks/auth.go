// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	jwt "github.com/dgrijalva/jwt-go"
	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// OAuthProvider is an autogenerated mock type for the OAuthProvider type
type OAuthProvider struct {
	mock.Mock
}

// Exchange provides a mock function with given fields: ctx, code
func (_m *OAuthProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	ret := _m.Called(ctx, code)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(context.Context, string) *oauth2.Token); ok {
		r0 = rf(ctx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateState provides a mock function with given fields:
func (_m *OAuthProvider) GenerateState() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAuthCodeURL provides a mock function with given fields: state
func (_m *OAuthProvider) GetAuthCodeURL(state string) string {
	ret := _m.Called(state)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Verify provides a mock function with given fields: rawToken
func (_m *OAuthProvider) Verify(rawToken string) (jwt.StandardClaims, error) {
	ret := _m.Called(rawToken)

	var r0 jwt.StandardClaims
	if rf, ok := ret.Get(0).(func(string) jwt.StandardClaims); ok {
		r0 = rf(rawToken)
	} else {
		r0 = ret.Get(0).(jwt.StandardClaims)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(rawToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyState provides a mock function with given fields: rawToken
func (_m *OAuthProvider) VerifyState(rawToken string) error {
	ret := _m.Called(rawToken)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(rawToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}