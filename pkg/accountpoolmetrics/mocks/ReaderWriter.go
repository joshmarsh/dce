// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import accountpoolmetrics "github.com/Optum/dce/pkg/accountpoolmetrics"
import mock "github.com/stretchr/testify/mock"

// ReaderWriter is an autogenerated mock type for the ReaderWriter type
type ReaderWriter struct {
	mock.Mock
}

// GetSingleton provides a mock function with given fields:
func (_m *ReaderWriter) GetSingleton() (*accountpoolmetrics.AccountPoolMetrics, error) {
	ret := _m.Called()

	var r0 *accountpoolmetrics.AccountPoolMetrics
	if rf, ok := ret.Get(0).(func() *accountpoolmetrics.AccountPoolMetrics); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accountpoolmetrics.AccountPoolMetrics)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Write provides a mock function with given fields: i, lastModifiedOn
func (_m *ReaderWriter) Write(i *accountpoolmetrics.AccountPoolMetrics, lastModifiedOn *int64) error {
	ret := _m.Called(i, lastModifiedOn)

	var r0 error
	if rf, ok := ret.Get(0).(func(*accountpoolmetrics.AccountPoolMetrics, *int64) error); ok {
		r0 = rf(i, lastModifiedOn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
