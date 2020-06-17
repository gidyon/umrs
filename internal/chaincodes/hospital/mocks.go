package hospital

import (
	"github.com/gidyon/umrs/internal/chaincodes/hospital/mocks"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// ledgerAPIMock is mock interface for ledger API
type ledgerAPIMock interface {
	ledger.ledgerClient
}

// NotificationAPIMock embeds notification API
type NotificationAPIMock interface {
	notification.NotificationServiceClient
}

func fakeledger() ledger.ledgerClient {
	// Create ledger mock
	ledgerAPI := &mocks.ledgerAPIMock{}
	ledgerAPI.On("AddBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.AddBlockResponse{Hash: uuid.New().String()}, nil)
	ledgerAPI.On("GetBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Block{Hash: uuid.New().String()}, nil)
	ledgerAPI.On("ListBlocks", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Blocks{NextPageNumber: 3, Blocks: []*ledger.Block{}}, nil)
	ledgerAPI.On("RegisterContract", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.RegisterContractResponse{ContractId: uuid.New().String()}, nil)
	return ledgerAPI
}

func fakeNotificationAPI() notification.NotificationServiceClient {
	notificationAPI := new(mocks.NotificationAPIMock)
	notificationAPI.On("CreateNotificationAccount", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	notificationAPI.On("Send", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	notificationAPI.On("ChannelSend", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	return notificationAPI
}
