package channel

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/apps/notification/channel/model"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/channel"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mockChannel() *channel.Channel {
	return &channel.Channel{
		ChannelId:   uuid.New().String(),
		Title:       randomdata.SillyName(),
		Description: randomdata.Paragraph(),
		OwnerId:     uuid.New().String(),
		CreatedTime: randomdata.FullDate(),
		Subscribers: 0,
	}
}

var _ = Describe("Creating A Channel #create", func() {
	var (
		createReq *channel.CreateChannelRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &channel.CreateChannelRequest{
			Channel: mockChannel(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Creating a channel with incorrect/missing values", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel is nil", func() {
			createReq.Channel = nil
			createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel title is missing", func() {
			createReq.Channel.Title = ""
			createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when owner id is missing", func() {
			createReq.Channel.OwnerId = ""
			createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel description is missing", func() {
			createReq.Channel.Description = ""
			createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Creating a channel with correct/valid request", func() {
		var channelID string
		It("should succeed when the request is valid", func() {
			createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			channelID = createRes.ChannelId
		})

		Describe("Created channel", func() {
			It("should exist in database", func() {
				channelDB := &model.Channel{}
				err := DB.First(channelDB, "channel_id=?", channelID).Error
				Expect(err).ShouldNot(HaveOccurred())
				Expect(channelDB).ShouldNot(BeNil())
			})
		})
	})
})