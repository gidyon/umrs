package permission

import (
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/notification"
)

// ledgerClientMock is mock API of Blochchain for testin
type ledgerClientMock interface {
	ledger.ledgerClient
}

// NotificationClient is mocked API of notification API for testing
type NotificationClient interface {
	notification.NotificationServiceClient
}
