// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// EnvironmentConverter is an autogenerated mock type for the EnvironmentConverter type
type EnvironmentConverter struct {
	mock.Mock
}

// LoadValueFromFile provides a mock function with given fields: f
func (_m *EnvironmentConverter) LoadValueFromFile(f string) (string, error) {
	ret := _m.Called(f)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(f)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
