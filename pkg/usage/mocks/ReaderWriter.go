// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import usage "github.com/Optum/dce/pkg/usage"

// ReaderWriter is an autogenerated mock type for the ReaderWriter type
type ReaderWriter struct {
	mock.Mock
}

// Add provides a mock function with given fields: i
func (_m *ReaderWriter) Add(i *usage.Usage) (*usage.Usage, error) {
	ret := _m.Called(i)

	var r0 *usage.Usage
	if rf, ok := ret.Get(0).(func(*usage.Usage) *usage.Usage); ok {
		r0 = rf(i)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usage.Usage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*usage.Usage) error); ok {
		r1 = rf(i)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Write provides a mock function with given fields: i
func (_m *ReaderWriter) Write(i *usage.Usage) (*usage.Usage, error) {
	ret := _m.Called(i)

	var r0 *usage.Usage
	if rf, ok := ret.Get(0).(func(*usage.Usage) *usage.Usage); ok {
		r0 = rf(i)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usage.Usage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*usage.Usage) error); ok {
		r1 = rf(i)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
