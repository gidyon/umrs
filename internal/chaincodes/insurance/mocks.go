package insurance

import (
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/notification"
)

// ledgerAPIMock is mock interface for ledger API
type ledgerAPIMock interface {
	ledger.ledgerClient
}

// NotificationAPIMock embeds notification API
type NotificationAPIMock interface {
	notification.NotificationServiceClient
}
