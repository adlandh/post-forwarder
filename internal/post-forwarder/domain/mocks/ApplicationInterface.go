// Code generated by mockery v2.44.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ApplicationInterface is an autogenerated mock type for the ApplicationInterface type
type ApplicationInterface struct {
	mock.Mock
}

// GetMessage provides a mock function with given fields: ctx, id
func (_m *ApplicationInterface) GetMessage(ctx context.Context, id string) (string, time.Time, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetMessage")
	}

	var r0 string
	var r1 time.Time
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, time.Time, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) time.Time); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(time.Time)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ProcessRequest provides a mock function with given fields: ctx, url, service, msg
func (_m *ApplicationInterface) ProcessRequest(ctx context.Context, url string, service string, msg string) error {
	ret := _m.Called(ctx, url, service, msg)

	if len(ret) == 0 {
		panic("no return value specified for ProcessRequest")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, url, service, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewApplicationInterface creates a new instance of ApplicationInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplicationInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ApplicationInterface {
	mock := &ApplicationInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
