// Code generated by mockery v2.45.0. DO NOT EDIT.

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

type ApplicationInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *ApplicationInterface) EXPECT() *ApplicationInterface_Expecter {
	return &ApplicationInterface_Expecter{mock: &_m.Mock}
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

// ApplicationInterface_GetMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMessage'
type ApplicationInterface_GetMessage_Call struct {
	*mock.Call
}

// GetMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *ApplicationInterface_Expecter) GetMessage(ctx interface{}, id interface{}) *ApplicationInterface_GetMessage_Call {
	return &ApplicationInterface_GetMessage_Call{Call: _e.mock.On("GetMessage", ctx, id)}
}

func (_c *ApplicationInterface_GetMessage_Call) Run(run func(ctx context.Context, id string)) *ApplicationInterface_GetMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ApplicationInterface_GetMessage_Call) Return(msg string, createdAt time.Time, err error) *ApplicationInterface_GetMessage_Call {
	_c.Call.Return(msg, createdAt, err)
	return _c
}

func (_c *ApplicationInterface_GetMessage_Call) RunAndReturn(run func(context.Context, string) (string, time.Time, error)) *ApplicationInterface_GetMessage_Call {
	_c.Call.Return(run)
	return _c
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

// ApplicationInterface_ProcessRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessRequest'
type ApplicationInterface_ProcessRequest_Call struct {
	*mock.Call
}

// ProcessRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - url string
//   - service string
//   - msg string
func (_e *ApplicationInterface_Expecter) ProcessRequest(ctx interface{}, url interface{}, service interface{}, msg interface{}) *ApplicationInterface_ProcessRequest_Call {
	return &ApplicationInterface_ProcessRequest_Call{Call: _e.mock.On("ProcessRequest", ctx, url, service, msg)}
}

func (_c *ApplicationInterface_ProcessRequest_Call) Run(run func(ctx context.Context, url string, service string, msg string)) *ApplicationInterface_ProcessRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *ApplicationInterface_ProcessRequest_Call) Return(err error) *ApplicationInterface_ProcessRequest_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *ApplicationInterface_ProcessRequest_Call) RunAndReturn(run func(context.Context, string, string, string) error) *ApplicationInterface_ProcessRequest_Call {
	_c.Call.Return(run)
	return _c
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
