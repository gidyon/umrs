package notification

import (
	"context"
	"encoding/json"
	"github.com/gidyon/umrs/internal/apps/notification/notification/model"
	"github.com/gidyon/umrs/internal/pkg/db"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

const (
	dbSaveOperation     = "SAVE_NOTIFICATION"
	dbDeleteOperation   = "DELETE_NOTIFICATION"
	dbMarkReadOperation = "MARK_READ_NOTIFICATION"
)

func recordExist(err error) bool {
	return strings.Contains(
		strings.ToLower(err.Error()), strings.ToLower("Duplicate entry"),
	)
}

func saveNotificationInDB(
	ctx context.Context,
	DB *gorm.DB,
	notificationPB *notification.Notification,
) error {

	// Check that the context is not cancelled
	if errs.CtxCancelled(ctx) {
		return ctx.Err()
	}

	if notificationPB == nil {
		return errs.NilObject("Notification")
	}

	var (
		notificationJSON []byte
		err              error
	)

	if notificationPB.Content != nil {
		notificationJSON, err = json.Marshal(notificationPB.Content)
		if err != nil {
			return errs.FromJSONMarshal(err, "notification.Content")
		}
	}

	for _, ownerID := range notificationPB.OwnerIds {
		notificationDB := &model.Notification{
			NotificationID: notificationPB.NotificationId,
			AccountID:      ownerID,
			Notification:   notificationJSON,
			Seen:           false,
			CreatedAt:      time.Now(),
		}

		// Save to db
		err = DB.Create(notificationDB).Error
		if err != nil {
			switch {
			case db.IsDuplicate(err):
			default:
				errs.LogError("failed to save notification: %v", err)
			}
		}
	}

	return nil
}

func saveNotificationInDBAndHandleErr(
	ctx context.Context,
	sqlErrHandlerChan chan<- *sqlErrHandlerDS,
	db *gorm.DB,
	notificationPB *notification.Notification,
) {
	err := saveNotificationInDB(ctx, db, notificationPB)
	if err != nil {
		// Send to worker and prevent goroutine leak
		go sendSQLErrToHandler(ctx, sqlErrHandlerChan, &sqlErrHandlerDS{
			operation: dbSaveOperation,
			payload:   notificationPB,
			err:       err,
		})
	}
}

func sendSQLErrToHandler(
	ctx context.Context, sqlErrHandlerChan chan<- *sqlErrHandlerDS, sqlErrPayload *sqlErrHandlerDS,
) {
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
	case sqlErrHandlerChan <- sqlErrPayload:
	}
}
