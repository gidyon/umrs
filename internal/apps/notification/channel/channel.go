package channel

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/db"

	channel_model "github.com/gidyon/umrs/internal/apps/notification/channel/model"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/channel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"strings"
)

const channelsPageSize = 20

type channelAPIServer struct {
	ctx context.Context
	DB  *gorm.DB
}

// NewChannelAPIServerServer factory creates a singleton channel.ChannelAPIServer
func NewChannelAPIServerServer(
	ctx context.Context,
	DB *gorm.DB,
) (channel.ChannelAPIServer, error) {
	// Validation'
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case DB == nil:
		return nil, errs.NilObject("SqlDB")
	}

	channelAPI := &channelAPIServer{
		ctx: ctx,
		DB:  DB,
	}

	// Automigration
	err := channelAPI.DB.AutoMigrate(&channel_model.Channel{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate channel table: %w", err)
	}

	return channelAPI, nil
}

func (channelAPI *channelAPIServer) CreateChannel(
	ctx context.Context, createReq *channel.CreateChannelRequest,
) (*channel.CreateChannelResponse, error) {
	// Request must not be nil
	if createReq == nil {
		return nil, errs.NilObject("CreateChannelRequest")
	}

	// Authenticate the admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// Validate the channel payload
	channelPB := createReq.GetChannel()
	switch {
	case channelPB == nil:
		err = errs.NilObject("Channel")
	case strings.TrimSpace(channelPB.OwnerId) == "":
		err = errs.MissingCredential("Owner Id")
	case strings.TrimSpace(channelPB.Title) == "":
		err = errs.MissingCredential("Channel Title")
	case strings.TrimSpace(channelPB.Description) == "":
		err = errs.MissingCredential("Channel Description")
	}
	if err != nil {
		return nil, err
	}

	channelDB := getChannelDB(channelPB)

	err = channelAPI.DB.Create(&channelDB).Error
	if err != nil {
		switch {
		case db.IsDuplicate(err):
			err = errs.ChannelExists(channelDB.ChannelName)
		default:
			err = errs.SQLQueryFailed(err, "CreateChannel")
		}
		return nil, err
	}

	return &channel.CreateChannelResponse{
		ChannelId: channelDB.ChannelID,
	}, nil
}

func (channelAPI *channelAPIServer) UpdateChannel(
	ctx context.Context, updateReq *channel.UpdateChannelRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateChannelRequest")
	}

	// Authenticate the admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	channelID := updateReq.GetChannelId()

	if strings.TrimSpace(channelID) == "" {
		return nil, errs.MissingCredential("channel id")
	}

	channelPB := updateReq.GetChannel()
	if channelPB == nil {
		return nil, errs.NilObject("Channel")
	}

	channelDB := getChannelDB(channelPB)

	err = channelAPI.DB.Model(channelDB).Where("channel_id=?", channelID).
		Omit("channel_id").Updates(channelDB).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "UpdateChannel")
	}

	return &empty.Empty{}, nil
}

func (channelAPI *channelAPIServer) DeleteChannel(
	ctx context.Context, delReq *channel.DeleteChannelRequest,
) (*empty.Empty, error) {
	// Request must not be nil
	if delReq == nil {
		return nil, errs.NilObject("DeleteChannelRequest")
	}

	// Authenticate the admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	channelID := delReq.GetChannelId()

	if strings.TrimSpace(channelID) == "" {
		return nil, errs.MissingCredential("channel id")
	}

	err = channelAPI.DB.Unscoped().Delete(&channel_model.Channel{}, "channel_id=?", channelID).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "DeleteChannel")
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

func (channelAPI *channelAPIServer) ListChannels(
	ctx context.Context, listReq *channel.ListChannelsRequest,
) (*channel.ListChannelsResponse, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListChannelsRequest")
	}

	// Authenticate the admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	pageNumber, pageSize := normalizePage(listReq.GetPageToken(), listReq.GetPageSize())
	offset := pageNumber*pageSize - pageSize

	channelsDB := make([]*channel_model.Channel, 0, pageSize)

	err = channelAPI.DB.Offset(offset).Limit(pageSize).Find(&channelsDB).Error
	switch {
	case err == nil:
	default:
		if err != nil {
			return nil, errs.SQLQueryFailed(err, "LIST")
		}
	}

	channelsPB := make([]*channel.Channel, 0, len(channelsDB))

	for _, channelDB := range channelsDB {
		channelPB := getChannelPB(channelDB)
		channelsPB = append(channelsPB, channelPB)
	}

	return &channel.ListChannelsResponse{
		NextPageToken: int32(pageNumber + 1),
		Channels:      channelsPB,
	}, nil
}

func (channelAPI *channelAPIServer) GetChannel(
	ctx context.Context, getReq *channel.GetChannelRequest,
) (*channel.Channel, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetChannelRequest")
	}

	// Authenticate the admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	channelID := getReq.GetChannelId()

	if strings.TrimSpace(channelID) == "" {
		return nil, errs.MissingCredential("channel id")
	}

	channelDB := &channel_model.Channel{}

	err = channelAPI.DB.Find(channelDB, "channel_id=?", getReq.ChannelId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.ChannelDoesntExist()
	default:
		return nil, errs.SQLQueryFailed(err, "FIND")
	}

	channelPB := getChannelPB(channelDB)

	return channelPB, nil
}
