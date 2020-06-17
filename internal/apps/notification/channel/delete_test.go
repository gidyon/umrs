package channel

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/channel"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Deleting A Channel #delete", func() {
	var (
		deleteReq *channel.DeleteChannelRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		deleteReq = &channel.DeleteChannelRequest{
			ChannelId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Deleting a channel with incorrect/missing values", func() {
		It("should fail when the request is nil", func() {
			deleteReq = nil
			deleteRes, err := ChannelAPI.DeleteChannel(ctx, deleteReq)
			Expect(err).Should(HaveOccurred())
			Expect(deleteRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel id is missing", func() {
			deleteReq.ChannelId = ""
			deleteRes, err := ChannelAPI.DeleteChannel(ctx, deleteReq)
			Expect(err).Should(HaveOccurred())
			Expect(deleteRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Deleting a channel with correct/valid request", func() {
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
			Describe("Deleting the channel", func() {
				It("should succeed when the request is valid", func() {
					deleteReq.ChannelId = channelID
					deleteRes, err := ChannelAPI.DeleteChannel(ctx, deleteReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(deleteRes).ShouldNot(BeNil())
					Expect(status.Code(err)).Should(Equal(codes.OK))
				})
			})

			Describe("Getting the deleted channel", func() {
				It("should not exist in database", func() {
					getRes, err := ChannelAPI.GetChannel(ctx, &channel.GetChannelRequest{
						ChannelId: channelID,
					})
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.NotFound))
					Expect(getRes).Should(BeNil())
				})
			})
		})

	})
})
