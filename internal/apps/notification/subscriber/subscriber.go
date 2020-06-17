package subscriber

import (
	"context"
	"encoding/json"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/notification"

	"github.com/gidyon/umrs/internal/apps/notification/subscriber/model"
	"github.com/gidyon/umrs/pkg/api/subscriber"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"strings"
)

const (
	limitRange   = 100
	listPageSize = 20
)

type subscriberAPIServer struct {
	ctx   context.Context
	sqlDB *gorm.DB
}

// NewSubscriberAPIServer factory creates a singleton subscriber.SubscriberAPIServer
func NewSubscriberAPIServer(
	ctx context.Context,
	sqlDB *gorm.DB,
) (subscriber.SubscriberAPIServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case sqlDB == nil:
		return nil, errs.NilObject("SqlDB")
	}

	subscriberAPI := &subscriberAPIServer{
		ctx:   ctx,
		sqlDB: sqlDB,
	}

	return subscriberAPI, nil
}

func (subscriberAPI *subscriberAPIServer) Subscribe(
	ctx context.Context, subReq *subscriber.SubscriberRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if subReq == nil {
		return nil, errs.NilObject("SubscriberRequest")
	}

	// Authenticate the actor
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	accountID := subReq.GetAccountId()
	channel := subReq.GetChannel()
	switch {
	case strings.TrimSpace(accountID) == "":
		err = errs.MissingCredential("Account Id")
	case strings.TrimSpace(channel) == "":
		err = errs.MissingCredential("Channel")
	}
	if err != nil {
		return nil, err
	}

	// Subscribe user to the channel
	err = subscriberAPI.sqlDB.Exec(
		`UPDATE subscribers SET channels=JSON_ARRAY_APPEND(channels, '$', ?) WHERE account_id=?`,
		channel, accountID,
	).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "Subscribe")
	}

	return &empty.Empty{}, nil
}

func (subscriberAPI *subscriberAPIServer) Unsubscribe(
	ctx context.Context, unSubReq *subscriber.SubscriberRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if unSubReq == nil {
		return nil, errs.NilObject("UnSubscriberRequest")
	}

	// Authenticate the actor
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	accountID := unSubReq.GetAccountId()
	channel := unSubReq.GetChannel()

	// Check that account id and channel is provided
	switch {
	case strings.TrimSpace(accountID) == "":
		err = errs.MissingCredential("Account Id")
	case strings.TrimSpace(channel) == "":
		err = errs.MissingCredential("Channel")
	}
	if err != nil {
		return nil, err
	}

	sub := &model.Subscriber{}
	// Get user channel
	err = subscriberAPI.sqlDB.Find(sub, "account_id=?", accountID).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.SubscriberDoesntExist(accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "GetChannels")
	}

	channels := []string{}
	// safe json unmarshal
	if len(sub.Channels) > 0 {
		err = json.Unmarshal(sub.Channels, &channels)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Channels")
		}
	}

	for pos, ch := range channels {
		if channel == ch {
			// Remove with append
			channels = append(channels[:pos-1], channels[pos:]...)
		}
	}

	if len(sub.Channels) > 0 {
		sub.Channels, err = json.Marshal(channels)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Channels")
		}
	}

	// Unsubscribe user to the channel
	err = subscriberAPI.sqlDB.Model(sub).
		Where("account_id=?", accountID).
		Select("channels").
		Updates(sub).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UnSubscribe")
	}

	return &empty.Empty{}, nil
}

func normalizePage(pageToken, pageSize int32) (int, int) {
	if pageSize <= 0 {
		pageSize = 1
	}
	if pageToken <= 0 {
		pageToken = 1
	}
	return int(pageToken), int(pageSize)
}

func (subscriberAPI *subscriberAPIServer) ListSubscribers(
	ctx context.Context, listReq *subscriber.ListSubscribersRequest,
) (*subscriber.ListSubscribersResponse, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListSubscribersRequest")
	}

	// Authenticate the admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// Validate that channnel is provided
	channel := listReq.GetChannel()
	if strings.TrimSpace(channel) == "" {
		return nil, errs.MissingCredential("Channel Name")
	}

	// Page token
	pageNumber, pageSize := normalizePage(listReq.GetPageToken(), listReq.GetPageSize())
	offset := pageNumber*pageSize - pageSize

	subscribersDB := make([]*model.Subscriber, 0, pageSize)

	// Get subscribers
	err = subscriberAPI.sqlDB.Offset(offset).Limit(pageSize).Order("created_at").
		Find(&subscribersDB).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "ListSubscribers")
	}

	subscribersPB := make([]*subscriber.Subscriber, 0, len(subscribersDB))

	for _, subscriberDB := range subscribersDB {

		subscriberPB, err := getSubscriberPB(subscriberDB)
		if err != nil {
			return nil, err
		}

		// Append to list only if channel exists
		for _, ch := range subscriberPB.Channels {
			if ch == channel {
				subscribersPB = append(subscribersPB, subscriberPB)
				break
			}
		}
	}

	return &subscriber.ListSubscribersResponse{
		NextPageToken: int32(pageNumber + 1),
		Subscribers:   subscribersPB,
	}, nil
}

func (subscriberAPI *subscriberAPIServer) GetSubscriber(
	ctx context.Context, getReq *subscriber.GetSubscriberRequest,
) (*subscriber.Subscriber, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetSubscriberRequest")
	}

	// Authenticate request
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	if strings.TrimSpace(getReq.AccountId) == "" {
		return nil, errs.MissingCredential("Account ID")
	}

	// Get subsriber
	subscriberDB := &model.Subscriber{}
	err = subscriberAPI.sqlDB.First(subscriberDB, "account_id=?", getReq.AccountId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.SubscriberDoesntExist(getReq.AccountId)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	return getSubscriberPB(subscriberDB)
}

func (subscriberAPI *subscriberAPIServer) GetSendMethod(
	ctx context.Context, getReq *subscriber.GetSendMethodRequest,
) (*subscriber.GetSendMethodResponse, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetSendMethodRequest")
	}

	// Authenticate request
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	accountID := getReq.GetAccountId()
	if strings.TrimSpace(accountID) == "" {
		return nil, errs.MissingCredential("Account ID")
	}

	// Get subsriber
	subscriberDB := &model.Subscriber{}
	err = subscriberAPI.sqlDB.Select("send_method").
		First(subscriberDB, "account_id=?", accountID).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.SubscriberDoesntExist(accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	return &subscriber.GetSendMethodResponse{
		SendMethod: notification.SendMethod(notification.SendMethod_value[subscriberDB.SendMethod]),
	}, nil
}

func getSubscriberPB(subscriberDB *model.Subscriber) (*subscriber.Subscriber, error) {
	subscriberPB := &subscriber.Subscriber{
		AccountId:  subscriberDB.AccountID,
		Email:      subscriberDB.Email,
		Phone:      subscriberDB.Phone,
		SendMethod: notification.SendMethod(notification.SendMethod_value[subscriberDB.SendMethod]),
		Channels:   []string{},
	}

	// safe json
	if len(subscriberDB.Channels) > 0 {
		err := json.Unmarshal(subscriberDB.Channels, &subscriberPB.Channels)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Subscriber Channels")
		}
	}

	return subscriberPB, nil
}
