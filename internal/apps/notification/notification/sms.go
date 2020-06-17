package notification

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/notification"
)

func sendSMSNotification(
	ctx context.Context,
	notificationAPI *notificationService,
	notificationObj *notification.Notification,
) {
	// Check if context is cancelled before proceeding
	if errs.CtxCancelled(ctx) {
		return
	}

	smsNotification := notificationObj.GetSmsNotification()
	_ = smsNotification
}

func sendBulkSMSNotification(
	ctx context.Context,
	notificationAPI *notificationService,
	notificationObj *notification.Notification,
) {
	// Check if context is cancelled before proceeding
	if errs.CtxCancelled(ctx) {
		return
	}

}
