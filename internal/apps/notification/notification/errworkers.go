package notification

import (
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/umrs/internal/apps/notification/notification/model"
	"math/rand"
	"time"
)

func (notificationAPI *notificationService) handleSQLErr() {
	for {
		select {
		case <-notificationAPI.ctx.Done():
			return
		case sqlErr := <-notificationAPI.sqlErrHandlerChan:
			logrus.Errorf("Error: %v Operation: %s\n", sqlErr.err, sqlErr.operation)
			go func(sqlErr *sqlErrHandlerDS) {
				var err error
				for i := 0; i < 5; i++ {
					sleepTime := time.Millisecond * time.Duration(rand.Float32()*1000) * time.Duration(i)
					time.Sleep(sleepTime)
					switch sqlErr.operation {
					case dbSaveOperation:
						err = notificationAPI.sqlDB.Create(sqlErr.payload.(*model.Notification)).Error
						if err != nil {
							continue
						}
						return
					case dbDeleteOperation:
					case dbMarkReadOperation:
					}
				}
			}(sqlErr)
		}
	}
}

func (notificationAPI *notificationService) handleRedisErr() {
	for {
		select {
		case <-notificationAPI.ctx.Done():
			return
		case redisErr := <-notificationAPI.redisErrHandlerChan:
			logrus.Errorf("Error: %v Operation: %s\n", redisErr.statusCMD.Err(), redisErr.operation)
		}
	}
}

func (notificationAPI *notificationService) emailNotificationErrWorker() {
	for {
		select {
		case <-notificationAPI.ctx.Done():
			return
		case emailErr := <-notificationAPI.emailNotificationChan:
			logrus.Errorf(
				"Error: %v; From: %s; To: %s; Subject: %s;\n",
				emailErr.err,
				emailErr.emailNotification.GetFrom(),
				emailErr.emailNotification.GetTo(),
				emailErr.emailNotification.GetSubject(),
			)
		}
	}
}

func (notificationAPI *notificationService) smsNotificationErrWorker() {
	for {
		select {
		case <-notificationAPI.ctx.Done():
			return
		case smsErr := <-notificationAPI.smsNotificationChan:
			logrus.Errorf(
				"Error: %v; Destination: %s; Keyword: %s; Message: %s;\n",
				smsErr.err,
				smsErr.smsNotification.GetDestinationPhone(),
				smsErr.smsNotification.GetKeyword(),
				smsErr.smsNotification.GetMessage(),
			)
		}
	}
}
