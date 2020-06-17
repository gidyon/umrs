package channel

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/channel"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("GetChannel For A User #get", func() {
	var (
		getReq *channel.GetChannelRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &channel.GetChannelRequest{
			ChannelId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Calling GetChannel with missing/incorrect values", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := ChannelAPI.GetChannel(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel id is missing", func() {
			getReq.ChannelId = ""
			getRes, err := ChannelAPI.GetChannel(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel does not exist", func() {
			getReq.ChannelId = "i-don=ty-eewbu8i-egszist"
			getRes, err := ChannelAPI.GetChannel(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
		})
	})

	Describe("Calling GetChannel with correct request", func() {
		var channelID string
		Context("Let's create a channel first", func() {
			It("should succeed in creating a channel", func() {
				createReq := &channel.CreateChannelRequest{
					Channel: mockChannel(),
				}
				createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(createRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				channelID = createRes.ChannelId
			})
		})

		When("Getting the channel", func() {
			It("should succeed", func() {
				getReq.ChannelId = channelID
				getRes, err := ChannelAPI.GetChannel(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
			})
		})
	})
})
