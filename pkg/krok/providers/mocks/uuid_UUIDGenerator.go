// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UUIDGenerator is an autogenerated mock type for the UUIDGenerator type
type UUIDGenerator struct {
	mock.Mock
}

// Generate provides a mock function with given fields:
func (_m *UUIDGenerator) Generate() (string, error) {
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