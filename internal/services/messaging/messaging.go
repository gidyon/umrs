package messaging

import (
	"context"
	"errors"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/messaging/call"
	"github.com/gidyon/umrs/pkg/api/messaging/emailing"
	"github.com/gidyon/umrs/pkg/api/messaging/push"
	"github.com/gidyon/umrs/pkg/api/messaging/sms"
	"github.com/gidyon/umrs/pkg/api/messaging/subscriber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/grpclog"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/jinzhu/gorm"

	"github.com/gidyon/umrs/pkg/api/messaging"
)

var emptyMsg = &empty.Empty{}

type messagingServer struct {
	sqlDB            *gorm.DB
	logger           grpclog.LoggerV2
	emailClient      emailing.EmailingClient
	emailSender      string
	pushClient       push.PushMessagingClient
	smsClient        sms.SMSAPIClient
	callClient       call.CallAPIClient
	subscriberClient subscriber.SubscriberAPIClient
	authAPI          auth.Interface
}

// Options contains options passed while calling NewMessagingServer
type Options struct {
	SQLDB            *gorm.DB
	Logger           grpclog.LoggerV2
	EmailClient      emailing.EmailingClient
	EmailSender      string
	JWTSigningKey    string
	PushClient       push.PushMessagingClient
	SMSClient        sms.SMSAPIClient
	CallClient       call.CallAPIClient
	SubscriberClient subscriber.SubscriberAPIClient
}

// NewMessagingServer creates a new MessagingServer API server
func NewMessagingServer(
	ctx context.Context, opt *Options,
) (messaging.MessagingServer, error) {
	// Validation
	var err error
	switch {
	case ctx == nil:
		err = errors.New("context is required")
	case opt.SQLDB == nil:
		err = errors.New("sqlDB is required")
	case opt.Logger == nil:
		err = errors.New("logger is required")
	case opt.EmailClient == nil:
		err = errors.New("email client is required")
	case opt.EmailSender == "":
		err = errors.New("email sender is required")
	case opt.JWTSigningKey == "":
		err = errors.New("jwt signing key is required")
	case opt.PushClient == nil:
		err = errors.New("push client is required")
	case opt.SMSClient == nil:
		err = errors.New("sms client is required")
	case opt.CallClient == nil:
		err = errors.New("call client is required")
	case opt.SubscriberClient == nil:
		err = errors.New("subscriber client is required")
	}
	if err != nil {
		return nil, err
	}

	authAPI, err := auth.NewAPI(opt.JWTSigningKey)
	if err != nil {
		return nil, err
	}

	api := &messagingServer{
		sqlDB:            opt.SQLDB,
		logger:           opt.Logger,
		emailClient:      opt.EmailClient,
		emailSender:      opt.EmailSender,
		pushClient:       opt.PushClient,
		smsClient:        opt.SMSClient,
		callClient:       opt.CallClient,
		subscriberClient: opt.SubscriberClient,
		authAPI:          authAPI,
	}

	// Auto migration
	err = api.sqlDB.AutoMigrate(&Message{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate: %v", err)
	}

	return api, nil
}

func validateMessage(msg *messaging.Message) error {
	// Validation
	var err error
	switch {
	case msg == nil:
		err = errs.NilObject("Message")
	case msg.UserId == "":
		err = errs.MissingField("user id")
	case msg.Title == "":
		err = errs.MissingField("title")
	case msg.Data == "":
		err = errs.MissingField("data")
	case len(msg.Details) == 0:
		err = errs.MissingField("payload")
	case len(msg.SendMethods) == 0:
		err = errs.MissingField("send methods")
	default:
		unknown := true
		for _, sendMethod := range msg.SendMethods {
			if sendMethod != messaging.SendMethod_UNKNOWN {
				unknown = false
				break
			}
		}
		if unknown {
			err = errs.MissingField("send methods")
		}
	}
	return err
}

func (api *messagingServer) BroadCastMessage(
	ctx context.Context, req *messaging.BroadCastMessageRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if req == nil {
		return nil, errs.NilObject("BroadCastMessageRequest")
	}

	// Authorize request
	err := api.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case len(req.Channels) == 0:
		err = errs.MissingField("topics")
	default:
		err = validateMessage(req.GetMessage())
	}
	if err != nil {
		return nil, err
	}

	err = api.sendBroadCastMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return emptyMsg, nil
}

func (api *messagingServer) sendBroadCastMessage(
	ctx context.Context, req *messaging.BroadCastMessageRequest,
) error {
	msg := req.GetMessage()
	msgDB, err := GetMessageDB(msg)
	if err != nil {
		return errs.WrapErrorWithMsg(err, "failed to get message model")
	}

	md, ok := metadata.FromIncomingContext(ctx)

	var pageSize int32 = 1000
	var pageToken int32 = -1

	fetch := true

	for fetch {
		// Get subscribers
		subscribersRes, err := api.subscriberClient.ListSubscribers(ctx, &subscriber.ListSubscribersRequest{
			Channels:  req.GetChannels(),
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			return errs.WrapErrorWithMsg(err, "failed to fetch subscribers")
		}

		// Update page token
		pageToken = subscribersRes.GetNextPageToken()

		if len(subscribersRes.GetSubscribers()) < int(pageSize) {
			fetch = false
		}

		// Send using anonymous goroutine
		go func(subscribers []*subscriber.Subscriber) {

			ctx2, cancel := context.WithCancel(context.Background())
			defer cancel()

			if ok {
				ctx2 = metadata.NewOutgoingContext(ctx2, md)
			}

			phones := make([]string, 0, len(subscribers))
			deviceTokens := make([]string, 0, len(subscribers))
			emails := make([]string, 0, len(subscribers))

			for _, subscriberPB := range subscribers {
				emails = append(emails, subscriberPB.GetEmail())
				deviceTokens = append(deviceTokens, subscriberPB.GetDeviceToken())
				phones = append(phones, subscriberPB.GetPhone())

				// Save message
				if msg.GetSave() {
					err = api.sqlDB.Create(msgDB).Error
					if err != nil {
						api.logger.Errorf("failed to save message model: %v", err)
						return
					}
				}
			}

			for _, sendMethod := range msg.GetSendMethods() {
				switch sendMethod {
				case messaging.SendMethod_UNKNOWN:
				case messaging.SendMethod_EMAIL:
					_, err = api.emailClient.SendEmail(ctx2, &emailing.Email{
						Destinations:    emails,
						From:            api.emailSender,
						Subject:         msg.Title,
						Body:            msg.Data,
						BodyContentType: "text/html",
					})
					if err != nil {
						api.logger.Errorf("failed to send email message: %v", err)
					}
				case messaging.SendMethod_SMS:
					_, err = api.smsClient.SendSMS(ctx2, &sms.SMS{
						DestinationPhones: phones,
						Keyword:           msg.Title,
						Message:           msg.Data,
					})
					if err != nil {
						api.logger.Errorf("failed to send sms message: %v", err)
					}
				case messaging.SendMethod_CALL:
					_, err = api.callClient.Call(ctx2, &call.CallPayload{
						DestinationPhones: phones,
						Keyword:           msg.Title,
						Message:           msg.Data,
					})
					if err != nil {
						api.logger.Errorf("failed to call recipients message: %v", err)
					}
				case messaging.SendMethod_PUSH:
					_, err = api.pushClient.SendPushMessage(ctx2, &push.PushMessage{
						DeviceTokens: deviceTokens,
						Title:        msg.Title,
						Message:      msg.Data,
						Details:      msg.Details,
					})
					if err != nil {
						api.logger.Errorf("failed to call recipients message: %v", err)
					}
				}
			}
		}(subscribersRes.GetSubscribers())
	}

	return nil
}

func (api *messagingServer) SendMessage(
	ctx context.Context, msg *messaging.Message,
) (*messaging.SendMessageResponse, error) {
	// Request must not be nil
	if msg == nil {
		return nil, errs.NilObject("Message")
	}

	// Authorize request
	err := api.authAPI.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	err = validateMessage(msg)
	if err != nil {
		return nil, err
	}

	// Get subscriber
	subscriberPB, err := api.subscriberClient.GetSubscriber(ctx, &subscriber.GetSubscriberRequest{
		AccountId: msg.UserId,
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errs.WrapErrorWithMsg(err, "failed to get subscriber")
	}

	// Send message
	for _, sendMethod := range msg.GetSendMethods() {
		switch sendMethod {
		case messaging.SendMethod_UNKNOWN:
		case messaging.SendMethod_EMAIL:
			_, err = api.emailClient.SendEmail(ctx, &emailing.Email{
				Destinations:    []string{subscriberPB.GetEmail()},
				From:            api.emailSender,
				Subject:         msg.Title,
				Body:            msg.Data,
				BodyContentType: "text/html",
			}, grpc.WaitForReady(true))
			if err != nil {
				api.logger.Errorf("failed to send email message: %v", err)
			}
		case messaging.SendMethod_SMS:
			_, err = api.smsClient.SendSMS(ctx, &sms.SMS{
				DestinationPhones: []string{subscriberPB.GetPhone()},
				Keyword:           msg.Title,
				Message:           msg.Data,
			}, grpc.WaitForReady(true))
			if err != nil {
				api.logger.Errorf("failed to send sms message: %v", err)
			}
		case messaging.SendMethod_CALL:
			_, err = api.callClient.Call(ctx, &call.CallPayload{
				DestinationPhones: []string{subscriberPB.GetPhone()},
				Keyword:           msg.Title,
				Message:           msg.Data,
			}, grpc.WaitForReady(true))
			if err != nil {
				api.logger.Errorf("failed to call recipients message: %v", err)
			}
		case messaging.SendMethod_PUSH:
			_, err = api.pushClient.SendPushMessage(ctx, &push.PushMessage{
				DeviceTokens: []string{subscriberPB.GetDeviceToken()},
				Title:        msg.Title,
				Message:      msg.Data,
				Details:      msg.Details,
			}, grpc.WaitForReady(true))
			if err != nil {
				api.logger.Errorf("failed to send push message: %v", err)
			}
		}
	}

	var msgID uint

	// Save message
	if msg.GetSave() {
		msgDB, err := GetMessageDB(msg)
		if err != nil {
			return nil, err
		}
		err = api.sqlDB.Create(msgDB).Error
		if err != nil {
			return nil, errs.WrapErrorWithMsg(err, "failed to save message")
		}
		msgID = msgDB.ID
	}

	return &messaging.SendMessageResponse{
		MessageId: fmt.Sprint(msgID),
	}, nil
}

const (
	defaultPageSize  = 10
	defaultPageToken = 1000000000
)

func normalizePage(pageToken, pageSize int32) (int, int) {
	if pageToken <= 0 {
		pageToken = defaultPageToken
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > 20 {
		pageSize = 20
	}

	return int(pageToken), int(pageSize)
}

func (api *messagingServer) ListMessages(
	ctx context.Context, listReq *messaging.ListMessagesRequest,
) (*messaging.Messages, error) {
	// Requst must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListMessagesRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, listReq.UserId)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case listReq.UserId == "":
		err = errs.MissingField("user id")
	}
	if err != nil {
		return nil, err
	}

	// Normalize page
	pageToken, pageSize := normalizePage(listReq.GetPageToken(), listReq.GetPageSize())

	messagesDB := make([]*Message, 0, pageSize)

	db := api.sqlDB.Order("id DESC").Limit(pageSize).Where("id<?", pageToken)

	if len(listReq.TypeFilters) > 0 {
		types := make([]int8, 0)
		filter := true
		for _, msgType := range listReq.GetTypeFilters() {
			types = append(types, int8(msgType))
			if msgType == messaging.MessageType_ALL {
				filter = false
				break
			}
		}
		if filter {
			db = db.Where("type IN(?)", types)
		}
	}

	err = db.Find(&messagesDB, "user_id=?", listReq.UserId).Error
	if err != nil {
		return nil, errs.WrapErrorWithMsg(err, "failed to find messages")
	}

	messagesPB := make([]*messaging.Message, 0, len(messagesDB))

	for _, messageDB := range messagesDB {
		messagePB, err := GetMessagePB(messageDB)
		if err != nil {
			return nil, err
		}

		messagesPB = append(messagesPB, messagePB)
		pageToken = int(messageDB.ID)
	}

	return &messaging.Messages{
		Messages:      messagesPB,
		NextPageToken: int32(pageToken),
	}, nil
}

func (api *messagingServer) ReadAll(
	ctx context.Context, readReq *messaging.MessageRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if readReq == nil {
		return nil, errs.NilObject("MessageRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, readReq.UserId)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case readReq.UserId == "":
		err = errs.MissingField("user id")
	}
	if err != nil {
		return nil, err
	}

	// Update messages
	err = api.sqlDB.Model(Message{}).Where("user_id=? AND seen=?", readReq.UserId, false).
		Update("seen", true).Error
	if err != nil {
		return nil, errs.WrapErrorWithMsg(err, "failed to mark messages as read")
	}

	return emptyMsg, nil
}

func (api *messagingServer) GetNewMessagesCount(
	ctx context.Context, getReq *messaging.MessageRequest,
) (*messaging.NewMessagesCount, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("MessageRequest")
	}

	// Authorize request
	_, err := api.authAPI.AuthorizeActor(ctx, getReq.UserId)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case getReq.UserId == "":
		err = errs.MissingField("user id")
	}
	if err != nil {
		return nil, err
	}

	var count int32
	err = api.sqlDB.Model(Message{}).Where("user_id=? AND seen=?", getReq.UserId, false).
		Count(&count).Error
	if err != nil {
		return nil, errs.WrapErrorWithMsg(err, "failed to get new messages count")
	}

	return &messaging.NewMessagesCount{
		Count: count,
	}, nil
}
