package notification

import (
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/notification"
	"strings"
)

// Checks that a notification content is provided and is valid
func validateNotification(notificationObj *notification.Notification) error {
	// Notification should not be nil; may cause panics
	if notificationObj == nil {
		return errs.NilObject("Notification")
	}
	// Owner of the notification must be provided
	if len(notificationObj.OwnerIds) == 0 {
		return errs.MissingCredential("Owner Ids")
	}
	// If bulk is true, bulk channel must be provided
	if notificationObj.Bulk {
		if strings.TrimSpace(notificationObj.BulkChannel) == "" {
			return errs.MissingCredential("bulk channel")
		}
	}

	var err error
	switch notificationObj.GetSendMethod() {
	case notification.SendMethod_EMAIL:
		err = validateEmailNotification(notificationObj)
	case notification.SendMethod_SMS:
		err = validateSMSNotification(notificationObj)
	case notification.SendMethod_EMAIL_AND_SMS:
		err = validateEmailNotification(notificationObj)
		if err == nil {
			err = validateSMSNotification(notificationObj)
		}
	}

	return err
}

func validateEmailNotification(notificationObj *notification.Notification) error {
	emailNotification := notificationObj.GetEmailNotification()
	// email notification should not be nil
	if emailNotification == nil {
		return errs.NilObject("EmailNotification")
	}
	switch {
	case strings.TrimSpace(emailNotification.Body) == "":
		return errs.MissingCredential("Email Body")
	case strings.TrimSpace(emailNotification.BodyContentType) == "":
		return errs.MissingCredential("Email Content Type")
	case strings.TrimSpace(emailNotification.Subject) == "":
		return errs.MissingCredential("Email Subject")
	}
	return nil
}

// Validate SMS notification content
func validateSMSNotification(notificationObj *notification.Notification) error {
	smsNotification := notificationObj.GetSmsNotification()
	// sms should not be nil
	if smsNotification == nil {
		return errs.NilObject("SmsNotification")
	}
	switch {
	case strings.TrimSpace(smsNotification.Keyword) == "":
		return errs.MissingCredential("SMS Keyword")
	case strings.TrimSpace(smsNotification.Message) == "":
		return errs.MissingCredential("SMS Message")
	}
	return nil
}
