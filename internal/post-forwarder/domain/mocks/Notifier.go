// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Notifier is an autogenerated mock type for the Notifier type
type Notifier struct {
	mock.Mock
}

// Send provides a mock function with given fields: ctx, service, msg
func (_m *Notifier) Send(ctx context.Context, service string, msg string) error {
	ret := _m.Called(ctx, service, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, service, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewNotifier creates a new instance of Notifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *Notifier {
	mock := &Notifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
