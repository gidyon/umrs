package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gidyon/account/pkg/api/admin"
	"github.com/gidyon/account/pkg/api/user"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"google.golang.org/grpc/codes"
	"time"

	notification_model "github.com/gidyon/umrs/internal/apps/notification/notification/model"
	subscriber_model "github.com/gidyon/umrs/internal/apps/notification/subscriber/model"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
	"strings"
)

const (
	highPriorityQueue    = "notifications:queue:1"
	mediumPriorityQueue  = "notifications:queue:2"
	lowPriorityQueue     = "notifications:queue:3"
	processedQueue       = "notifications:processed"
	listSubscribersSize  = 50
	listNotificationSize = 15
)

type redisErrHandlerDS struct {
	statusCMD redis.Cmder
	list      string
	operation string
}

type sqlErrHandlerDS struct {
	operation string
	payload   interface{}
	err       error
}

type emailNotificationWorker struct {
	emailNotification *notification.EmailNotification
	err               error
}

type smsNotificationWorker struct {
	smsNotification *notification.SMSNotification
	err             error
}

type dialer interface {
	DialAndSend(...*gomail.Message) error
}

type smsClient struct{}

type ussdClient struct{}

type callClient struct{}

type notificationService struct {
	ctx                   context.Context
	sqlDB                 *gorm.DB
	redisDB               *redis.Client
	notificationChan      chan *notification.Notification
	redisErrHandlerChan   chan *redisErrHandlerDS
	sqlErrHandlerChan     chan *sqlErrHandlerDS
	smtpDialer            dialer
	emailNotificationChan chan *emailNotificationWorker
	smsClient             *smsClient
	smsNotificationChan   chan *smsNotificationWorker
	ussdClient            *ussdClient
	callClient            *callClient
	adminAccountClient    admin.AdminAPIClient
	userAccountClient     user.UserAPIClient
}

// Options contains parameters for NewNotificationServiceServer function
type Options struct {
	SQLDB       *gorm.DB
	RedisClient *redis.Client
	SMTPDialer  interface {
		DialAndSend(...*gomail.Message) error
	}
}

// NewNotificationServiceServer creates an API singleton for notification service
func NewNotificationServiceServer(
	ctx context.Context, opt *Options,
) (notification.NotificationServiceServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case opt.SQLDB == nil:
		return nil, errs.NilObject("SqlDB")
	case opt.RedisClient == nil:
		return nil, errs.NilObject("RedisDB")
	case opt.SMTPDialer == nil:
		return nil, errs.NilObject("SMTPDialer")
	}

	notificationAPI := &notificationService{
		ctx:                   ctx,
		sqlDB:                 opt.SQLDB,
		redisDB:               opt.RedisClient,
		notificationChan:      make(chan *notification.Notification, 0),
		redisErrHandlerChan:   make(chan *redisErrHandlerDS, 0),
		sqlErrHandlerChan:     make(chan *sqlErrHandlerDS, 0),
		smtpDialer:            opt.SMTPDialer,
		smsClient:             &smsClient{},
		emailNotificationChan: make(chan *emailNotificationWorker, 0),
		smsNotificationChan:   make(chan *smsNotificationWorker, 0),
	}

	// Automigration of notifications tables
	err := notificationAPI.sqlDB.AutoMigrate(
		&subscriber_model.Subscriber{}, &notification_model.Notification{},
	).Error
	if err != nil {
		return nil, fmt.Errorf("failed to perform automigration: %w", err)
	}

	// Start err workers
	go notificationAPI.handleRedisErr()
	go notificationAPI.handleSQLErr()
	go notificationAPI.emailNotificationErrWorker()
	go notificationAPI.smsNotificationErrWorker()

	// Consumer process
	go notificationAPI.worker()

	return notificationAPI, nil
}

func (notificationAPI *notificationService) CreateNotificationAccount(
	ctx context.Context, createReq *notification.CreateNotificationAccountRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if createReq == nil {
		return nil, errs.NilObject("CreateNotificationAccountRequest")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validatation
	channels := createReq.GetChannels()
	switch {
	case strings.TrimSpace(createReq.AccountId) == "":
		return nil, errs.MissingCredential("Account ID")
	case strings.TrimSpace(createReq.Email) == "" && strings.TrimSpace(createReq.Phone) == "":
		return nil, errs.MissingCredential("Email or Phone")
	case len(channels) == 0:
		channels = []string{"default"}
	}

	// Marshal subscriber channels to json
	channelsJSON, err := json.Marshal(channels)
	if err != nil {
		return nil, errs.FromJSONMarshal(err, "Channels")
	}

	// Subscriber model
	subscriberDB := &subscriber_model.Subscriber{
		AccountID:  createReq.AccountId,
		Email:      createReq.Email,
		Phone:      createReq.Phone,
		SendMethod: createReq.SendMethod.String(),
		Channels:   channelsJSON,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = notificationAPI.sqlDB.Create(subscriberDB).Error
	switch {
	case err == nil:
	case recordExist(err):
	default:
		return nil, errs.SQLQueryFailed(err, "Create")
	}

	return &empty.Empty{}, nil
}

func handleRedisIntCMD(intCMD *redis.IntCmd) error {
	if intCMD.Err() != nil {
		return errs.RedisCmdFailed(intCMD.Err(), "LPUSH")
	}
	return nil
}

func (notificationAPI *notificationService) Send(
	ctx context.Context, sendReq *notification.Notification,
) (*empty.Empty, error) {
	// Request must not be nil
	if sendReq == nil {
		return nil, errs.NilObject("Notification")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	sendReq.Bulk = false

	// Validatation
	err = validateNotification(sendReq)
	if err != nil {
		return nil, err
	}

	notificationBytes, err := proto.Marshal(sendReq)
	if err != nil {
		return nil, errs.FromProtoMarshal(err, "Notification")
	}

	// Send the notification to head of redis list
	switch sendReq.GetPriority() {
	case notification.Priority_HIGH:
		// Send to first queue
		err = handleRedisIntCMD(
			notificationAPI.redisDB.LPush(highPriorityQueue, notificationBytes),
		)

	case notification.Priority_MEDIUM:
		// Send to middle queue
		err = handleRedisIntCMD(
			notificationAPI.redisDB.LPush(mediumPriorityQueue, notificationBytes),
		)

	case notification.Priority_LOW:
		// Send to last queue
		err = handleRedisIntCMD(
			notificationAPI.redisDB.LPush(lowPriorityQueue, notificationBytes),
		)
	}

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (notificationAPI *notificationService) ChannelSend(
	ctx context.Context, ChannelSendReq *notification.Notification,
) (*empty.Empty, error) {
	// Request smust not be nil
	if ChannelSendReq == nil {
		return nil, errs.NilObject("ChannelSendRequest")
	}

	// Authentication: Only admin user can broadcast notifications
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	ChannelSendReq.Bulk = true

	// Validates that appropriate content is provided
	err = validateNotification(ChannelSendReq)
	if err != nil {
		return nil, err
	}

	notificationBytes, err := proto.Marshal(ChannelSendReq)
	if err != nil {
		return nil, errs.FromProtoMarshal(err, "Notification")
	}

	// Send the notification to head of appropriate redis list
	switch ChannelSendReq.GetPriority() {
	case notification.Priority_HIGH:
		// Send to first queue
		err = handleRedisIntCMD(
			notificationAPI.redisDB.LPush(highPriorityQueue, notificationBytes),
		)

	case notification.Priority_MEDIUM:
		// Send to middle queue
		err = handleRedisIntCMD(
			notificationAPI.redisDB.LPush(mediumPriorityQueue, notificationBytes),
		)

	case notification.Priority_LOW:
		// Send to last queue
		err = handleRedisIntCMD(
			notificationAPI.redisDB.LPush(lowPriorityQueue, notificationBytes),
		)
	}
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

const defaultPageSize = 100

func normalizePage(pageToken, pageSize int32) (int, int) {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageToken <= 0 {
		pageToken = 1
	}
	return int(pageToken), int(pageSize)
}

func (notificationAPI *notificationService) ListNotifications(
	ctx context.Context, listReq *notification.ListNotificationsRequest,
) (*notification.ListNotificationsResponse, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListNotificationsRequest")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validate that account ID is provided
	accountID := listReq.GetAccountId()
	if strings.TrimSpace(accountID) == "" {
		return nil, errs.MissingCredential("account id")
	}

	// Page token
	pageNumber, pageSize := normalizePage(listReq.GetPageToken(), listReq.GetPageNumber())
	offset := pageNumber*pageSize - pageSize

	notificationsDB := make([]*notification_model.Notification, 0, pageSize)

	// Get notifications
	err = notificationAPI.sqlDB.Offset(offset).Limit(pageSize).
		Find(&notificationsDB, "account_id=?", accountID).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "ListNotifications")
	}

	notificationsPB := make([]*notification.Notification, 0, len(notificationsDB))

	for _, notificationDB := range notificationsDB {
		notificationPB := &notification.Notification{
			NotificationId: notificationDB.NotificationID,
			OwnerIds:       []string{notificationDB.AccountID},
			Seen:           notificationDB.Seen,
		}

		// json.Unmarshal is not nil safe
		if len(notificationDB.Notification) > 0 {
			err = json.Unmarshal(notificationDB.Notification, &notificationPB.Content)
			if err != nil {
				return nil, errs.FromJSONMarshal(err, "Notification.Content")
			}
		}

		notificationsPB = append(notificationsPB, notificationPB)
	}

	return &notification.ListNotificationsResponse{
		Notifications: notificationsPB,
		NextPageToken: int32(pageNumber + 1),
	}, nil
}

func (notificationAPI *notificationService) MarkNotificationRead(
	ctx context.Context, readReq *notification.MarkNotificationReadRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if readReq == nil {
		return nil, errs.NilObject("MarkNotificationReadRequest")
	}

	// Authenticate the request
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	notificationID := readReq.GetNotificationId()

	// Check that account id and notification id are not empty
	switch {
	case strings.TrimSpace(notificationID) == "":
		return nil, errs.MissingCredential("NotificationId")
	}

	// Update the model
	db := notificationAPI.sqlDB.Model(&notification_model.Notification{}).
		Where("notification_id=?", notificationID).
		Update("seen", true)
	if db.RowsAffected == 0 {
		return nil, errs.WrapMessage(
			codes.NotFound, fmt.Sprintf("notification with id %s not found", notificationID),
		)
	}
	if db.Error != nil {
		return nil, errs.SQLQueryFailed(err, "MarkNotificationRead")
	}

	return &empty.Empty{}, nil
}
