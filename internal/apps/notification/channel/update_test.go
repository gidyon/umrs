package channel

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/channel"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Updating A Channel #update", func() {
	var (
		updateReq *channel.UpdateChannelRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &channel.UpdateChannelRequest{
			ChannelId: uuid.New().String(),
			Channel: &channel.Channel{
				Title:       randomdata.Month(),
				Description: randomdata.Paragraph(),
			},
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Updating a channel with incorrect/missing values", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := ChannelAPI.UpdateChannel(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(updateRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel is nil", func() {
			updateReq.Channel = nil
			updateRes, err := ChannelAPI.UpdateChannel(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(updateRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel id is missing", func() {
			updateReq.ChannelId = ""
			updateRes, err := ChannelAPI.UpdateChannel(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(updateRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Updating a channel with correct/valid request", func() {
		var channelID string
		Describe("Lets create the channel first", func() {
			It("should succeed", func() {
				ctx = auth.AddAdminMD(context.Background())
				createReq := &channel.CreateChannelRequest{
					Channel: mockChannel(),
				}
				createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(createRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				channelID = createRes.ChannelId
			})

			Describe("updating the channel", func() {
				var title string
				BeforeEach(func() {
					updateReq.ChannelId = channelID
					updateReq.Channel.Title = randomdata.SillyName()
				})

				It("should succeed when the request is valid", func() {
					updateRes, err := ChannelAPI.UpdateChannel(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(updateRes).ShouldNot(BeNil())
				})
				It("should succeed when owner id is missing", func() {
					updateReq.Channel.OwnerId = ""
					updateRes, err := ChannelAPI.UpdateChannel(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(updateRes).ShouldNot(BeNil())
				})
				It("should succeed when channel description is missing", func() {
					updateReq.Channel.Description = ""
					updateRes, err := ChannelAPI.UpdateChannel(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(updateRes).ShouldNot(BeNil())
					title = updateReq.Channel.Title
				})

				Describe("Getting an updated channel", func() {
					It("should reflect updated fields", func() {
						getRes, err := ChannelAPI.GetChannel(ctx, &channel.GetChannelRequest{
							ChannelId: channelID,
						})
						Expect(err).ShouldNot(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.OK))
						Expect(getRes).ShouldNot(BeNil())
						Expect(getRes.Title).Should(Equal(title))
					})
				})
			})
		})
	})
})
