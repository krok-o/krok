// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Vault is an autogenerated mock type for the Vault type
type Vault struct {
	mock.Mock
}

// AddSecret provides a mock function with given fields: key, value
func (_m *Vault) AddSecret(key string, value []byte) {
	_m.Called(key, value)
}

// DeleteSecret provides a mock function with given fields: key
func (_m *Vault) DeleteSecret(key string) {
	_m.Called(key)
}

// GetSecret provides a mock function with given fields: key
func (_m *Vault) GetSecret(key string) ([]byte, error) {
	ret := _m.Called(key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSecrets provides a mock function with given fields:
func (_m *Vault) ListSecrets() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// LoadSecrets provides a mock function with given fields:
func (_m *Vault) LoadSecrets() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveSecrets provides a mock function with given fields:
func (_m *Vault) SaveSecrets() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}