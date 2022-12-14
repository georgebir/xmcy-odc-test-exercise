// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	handlers "test-exercise/api/messagebroker/handlers"

	mock "github.com/stretchr/testify/mock"
)

// MessageBroker is an autogenerated mock type for the MessageBroker type
type MessageBroker struct {
	mock.Mock
}

// ListenAndServe provides a mock function with given fields: topic, group, handler
func (_m *MessageBroker) ListenAndServe(topic string, group string, handler handlers.Handler) {
	_m.Called(topic, group, handler)
}

// Produce provides a mock function with given fields: topic, value
func (_m *MessageBroker) Produce(topic string, value interface{}) error {
	ret := _m.Called(topic, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(topic, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMessageBroker interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageBroker creates a new instance of MessageBroker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageBroker(t mockConstructorTestingTNewMessageBroker) *MessageBroker {
	mock := &MessageBroker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
