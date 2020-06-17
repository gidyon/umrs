package notification

import (
	"context"
	"github.com/gidyon/umrs/internal/apps/notification/subscriber/model"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/notification"
	"gopkg.in/gomail.v2"
)

func sendEmailNotification(
	ctx context.Context,
	notificationAPI *notificationService,
	notificationObj *notification.Notification,
) {
	// Check if context is cancelled before proceeding
	if errs.CtxCancelled(ctx) {
		return
	}

	// get the dialer
	dialer, ok := notificationAPI.smtpDialer.(*gomail.Dialer)
	if !ok {
		return
	}

	emailNotification := notificationObj.GetEmailNotification()

	m := gomail.NewMessage()
	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", emailNotification.GetTo()...)
	m.SetHeader("Subject", emailNotification.GetSubject())
	m.SetBody(emailNotification.GetBodyContentType(), emailNotification.GetBody())

	if err := dialer.DialAndSend(m); err != nil {
		select {
		case <-ctx.Done():
		case notificationAPI.emailNotificationChan <- &emailNotificationWorker{
			err:               err,
			emailNotification: emailNotification,
		}:
		}
	}
}

const limit = 100

func sendBulkEmailNotification(
	ctx context.Context,
	notificationAPI *notificationService,
	notificationObj *notification.Notification,
) {
	var (
		err                  error
		offset               = 0
		subscribers          = make([]*model.Subscriber, 0, limit)
		subscribersAvailable = true
	)

	// Get list account ids from bulk channel; send 100 per time;
	for subscribersAvailable {
		err = notificationAPI.sqlDB.Table("subscribers").Offset(offset).Limit(limit).
			Where("? MEMBER OF(channels)", notificationObj.BulkChannel).
			Find(&subscribers).
			Error
		switch {
		case err == nil:
		default:
			errs.LogError("failed to get list of subscribers: %v", err)
			return
		}

		if len(subscribers) < limit {
			subscribersAvailable = false
		}

		subscriberSlice := make([]string, 0, limit)

		for _, subscriber := range subscribers {
			subscriberSlice = append(subscriberSlice, subscriber.Email)
		}

		notificationObjLocal := *notificationObj

		notificationObjLocal.GetEmailNotification().To = subscriberSlice

		sendEmailNotification(ctx, notificationAPI, &notificationObjLocal)
		offset += limit
	}
}
