package notification

import (
	// "bytes"
	"context"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"time"
)

// Counting semaphore for Rate Limiting
var sem chan struct{}

// consumes published content from redis list
func (notificationAPI *notificationService) worker() {
	sem = make(chan struct{}, 1000)
	defer close(notificationAPI.notificationChan)

	// Runs forever like good
	for {
		// Get the notification object from redis list using BRPOP
		strSliceCMD := notificationAPI.redisDB.BRPop(
			time.Duration(1*time.Hour),
			highPriorityQueue,
			mediumPriorityQueue,
			lowPriorityQueue,
		)

		// Acquire token
		sem <- struct{}{}

		// We start as goroutine so that we don't block worker
		go process(notificationAPI, strSliceCMD)
	}
}

const (
	operationRedisCache        = "redis cache" // default
	operationProtoMarshaling   = "protobufs marshaling"
	operationProtoUnMarshaling = "protbufs unmarshaling"
)

func process(notificationAPI *notificationService, strSliceCMD *redis.StringSliceCmd) {
	// Release token when done
	defer func() {
		<-sem
	}()

	// First item in array is list name, second item is data
	res, err := strSliceCMD.Result()
	if strSliceCMD.Err() != nil {
		errs.LogError("redis non-nil error: %v", strSliceCMD.Err())
		return
	}

	// Save the element that has been popped to processed queue for eventual consistency
	intCMD := notificationAPI.redisDB.LPush(processedQueue, res[1])
	if intCMD.Err() != nil {
		// Send the failed command to error handler
		sendRedisErrToHandler(
			notificationAPI.redisErrHandlerChan, &redisErrHandlerDS{
				statusCMD: intCMD,
				list:      processedQueue,
				operation: "LPUSH",
			},
		)
	}

	notificationPB := &notification.Notification{}

	// Unmarshal from proto to notification object
	err = proto.Unmarshal([]byte(res[1]), notificationPB)
	if err != nil {
		// Send the failed command to error hander for proper handling
		sendRedisErrToHandler(
			notificationAPI.redisErrHandlerChan, &redisErrHandlerDS{
				statusCMD: strSliceCMD,
				list:      res[0],
				operation: operationProtoUnMarshaling,
			},
		)
		return
	}

	// Save notification in db if it's required
	if notificationPB.GetSave() {
		// Save the notification in DB and handle error
		saveNotificationInDBAndHandleErr(
			notificationAPI.ctx,
			notificationAPI.sqlErrHandlerChan,
			notificationAPI.sqlDB,
			notificationPB,
		)
	}

	// Pick the appropriate send method
	switch notificationPB.GetSendMethod() {
	case notification.SendMethod_EMAIL:
		// Send Bulk Notification
		if notificationPB.GetBulk() {
			sendBulkEmailNotification(
				notificationAPI.ctx,
				notificationAPI,
				notificationPB,
			)
		} else {
			// Single destination notification
			sendEmailNotification(
				notificationAPI.ctx,
				notificationAPI,
				notificationPB,
			)
		}
	case notification.SendMethod_SMS:
		if notificationPB.GetBulk() {
			// Bulk Notification
			sendBulkSMSNotification(
				notificationAPI.ctx,
				notificationAPI,
				notificationPB,
			)
		} else {
			// Single destination notification
			sendSMSNotification(
				notificationAPI.ctx,
				notificationAPI,
				notificationPB,
			)
		}
	case notification.SendMethod_EMAIL_AND_SMS:
		if notificationPB.GetBulk() {
			// Bulk Notification
			sendBulkEmailAndSMSNotification(
				notificationAPI.ctx,
				notificationAPI,
				notificationPB,
			)
			break
		}
		// Single destination notification
		sendEmailAndSMSNotification(
			notificationAPI.ctx,
			notificationAPI,
			notificationPB,
		)
	}
}

func sendRedisErrToHandler(
	redisErrHandlerDS chan<- *redisErrHandlerDS,
	redisErrHandlerDSPayload *redisErrHandlerDS,
) {
	select {
	case <-time.After(time.Duration(3 * time.Second)):
	case redisErrHandlerDS <- redisErrHandlerDSPayload:
	}
}

func sendEmailAndSMSNotification(
	ctx context.Context,
	notificationAPI *notificationService,
	notificationPB *notification.Notification,
) {
	go sendEmailNotification(ctx, notificationAPI, notificationPB)
	go sendSMSNotification(ctx, notificationAPI, notificationPB)
}

func sendBulkEmailAndSMSNotification(
	ctx context.Context,
	notificationAPI *notificationService,
	notificationPB *notification.Notification,
) {
	go sendBulkEmailNotification(ctx, notificationAPI, notificationPB)
	go sendBulkSMSNotification(ctx, notificationAPI, notificationPB)
}
