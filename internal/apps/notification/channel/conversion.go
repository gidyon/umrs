package channel

import (
	"github.com/gidyon/umrs/internal/apps/notification/channel/model"
	"github.com/gidyon/umrs/pkg/api/channel"
)

func getChannelPB(channelDB *model.Channel) *channel.Channel {
	return &channel.Channel{
		ChannelId:   channelDB.ChannelID,
		Title:       channelDB.ChannelName,
		Description: channelDB.Description,
		OwnerId:     channelDB.OwnerID,
		CreatedTime: channelDB.CreatedAt.UTC().Local().String(),
	}
}

func getChannelDB(channelPB *channel.Channel) *model.Channel {
	return &model.Channel{
		ChannelID:   channelPB.ChannelId,
		ChannelName: channelPB.Title,
		Description: channelPB.Description,
		OwnerID:     channelPB.OwnerId,
	}
}
