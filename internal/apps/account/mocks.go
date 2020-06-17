package account

import (
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/subscriber"
)

// NotificationAPIMock embeds notification API
type NotificationAPIMock interface {
	notification.NotificationServiceClient
}

// SubscriberAPIMock embeds subscriber API
type SubscriberAPIMock interface {
	subscriber.SubscriberAPIClient
}
