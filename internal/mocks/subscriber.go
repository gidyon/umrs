package mocks

import (
	"github.com/gidyon/umrs/internal/mocks/mocks"
	"github.com/gidyon/umrs/pkg/api/messaging/subscriber"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
)

// SubscriberAPIClientMock is a mock for subscriber API client
type SubscriberAPIClientMock interface {
	subscriber.SubscriberAPIClient
}

// SubscriberAPI is a fake subscriber API
var SubscriberAPI = &mocks.SubscriberAPIClientMock{}

func init() {
	SubscriberAPI.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).
		Return(empty.Empty{}, nil)
	SubscriberAPI.On("Unsubscribe", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	SubscriberAPI.On("ListSubscribers", mock.Anything, mock.Anything, mock.Anything).
		Return(subscriber.ListSubscribersResponse{}, nil)
	SubscriberAPI.On("GetSubscriber", mock.Anything, mock.Anything, mock.Anything).
		Return(&subscriber.Subscriber{}, nil)
}
