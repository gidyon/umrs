package mocks

import (
	"github.com/gidyon/umrs/internal/mocks/mocks"
	"github.com/gidyon/umrs/pkg/api/messaging/channel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
)

// ChannelAPIClientMock is a mock for channel API client
type ChannelAPIClientMock interface {
	channel.ChannelAPIClient
}

// ChannelAPI is a fake channel API
var ChannelAPI = &mocks.ChannelAPIClientMock{}

func init() {
	ChannelAPI.On("CreateChannel", mock.Anything, mock.Anything, mock.Anything).
		Return(mock.Anything, nil)
	ChannelAPI.On("UpdateChannel", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	ChannelAPI.On("DeleteChannel", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	ChannelAPI.On("ListChannels", mock.Anything, mock.Anything, mock.Anything).
		Return(mock.Anything, nil)
	ChannelAPI.On("GetChannel", mock.Anything, mock.Anything, mock.Anything).
		Return(mock.Anything, nil)
}