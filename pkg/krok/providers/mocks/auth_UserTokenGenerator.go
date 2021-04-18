// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserTokenGenerator is an autogenerated mock type for the UserTokenGenerator type
type UserTokenGenerator struct {
	mock.Mock
}

// Generate provides a mock function with given fields: length
func (_m *UserTokenGenerator) Generate(length int) (string, error) {
	ret := _m.Called(length)

	var r0 string
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(length)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(length)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
