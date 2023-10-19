// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MessageStorage is an autogenerated mock type for the MessageStorage type
type MessageStorage struct {
	mock.Mock
}

// Read provides a mock function with given fields: ctx, id
func (_m *MessageStorage) Read(ctx context.Context, id string) (string, time.Time, error) {
	ret := _m.Called(ctx, id)

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

// Store provides a mock function with given fields: ctx, msg
func (_m *MessageStorage) Store(ctx context.Context, msg string) (string, error) {
	ret := _m.Called(ctx, msg)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, msg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMessageStorage creates a new instance of MessageStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMessageStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MessageStorage {
	mock := &MessageStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}