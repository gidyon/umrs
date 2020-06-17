// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	subscriber "github.com/gidyon/umrs/pkg/api/subscriber"
)

// SubscriberAPIMock is an autogenerated mock type for the SubscriberAPIMock type
type SubscriberAPIMock struct {
	mock.Mock
}

// GetSendMethod provides a mock function with given fields: ctx, in, opts
func (_m *SubscriberAPIMock) GetSendMethod(ctx context.Context, in *subscriber.GetSendMethodRequest, opts ...grpc.CallOption) (*subscriber.GetSendMethodResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *subscriber.GetSendMethodResponse
	if rf, ok := ret.Get(0).(func(context.Context, *subscriber.GetSendMethodRequest, ...grpc.CallOption) *subscriber.GetSendMethodResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*subscriber.GetSendMethodResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *subscriber.GetSendMethodRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriber provides a mock function with given fields: ctx, in, opts
func (_m *SubscriberAPIMock) GetSubscriber(ctx context.Context, in *subscriber.GetSubscriberRequest, opts ...grpc.CallOption) (*subscriber.Subscriber, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *subscriber.Subscriber
	if rf, ok := ret.Get(0).(func(context.Context, *subscriber.GetSubscriberRequest, ...grpc.CallOption) *subscriber.Subscriber); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*subscriber.Subscriber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *subscriber.GetSubscriberRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSubscribers provides a mock function with given fields: ctx, in, opts
func (_m *SubscriberAPIMock) ListSubscribers(ctx context.Context, in *subscriber.ListSubscribersRequest, opts ...grpc.CallOption) (*subscriber.ListSubscribersResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *subscriber.ListSubscribersResponse
	if rf, ok := ret.Get(0).(func(context.Context, *subscriber.ListSubscribersRequest, ...grpc.CallOption) *subscriber.ListSubscribersResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*subscriber.ListSubscribersResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *subscriber.ListSubscribersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subscribe provides a mock function with given fields: ctx, in, opts
func (_m *SubscriberAPIMock) Subscribe(ctx context.Context, in *subscriber.SubscriberRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *empty.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *subscriber.SubscriberRequest, ...grpc.CallOption) *empty.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*empty.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *subscriber.SubscriberRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unsubscribe provides a mock function with given fields: ctx, in, opts
func (_m *SubscriberAPIMock) Unsubscribe(ctx context.Context, in *subscriber.SubscriberRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *empty.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *subscriber.SubscriberRequest, ...grpc.CallOption) *empty.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*empty.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *subscriber.SubscriberRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
