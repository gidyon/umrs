package channel

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/channel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ListChannels For A User #list", func() {
	var (
		listReq *channel.ListChannelsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &channel.ListChannelsRequest{
			PageToken: 0,
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Calling ListChannels with missing/incorrect values", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := ChannelAPI.ListChannels(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("Calling ListChannels with valid values", func() {
		Context("Lets create one channel first", func() {
			It("should succeed", func() {
				ctx = auth.AddAdminMD(context.Background())
				createReq := &channel.CreateChannelRequest{
					Channel: mockChannel(),
				}
				createRes, err := ChannelAPI.CreateChannel(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
			})
		})

		It("should succeed when the request is valid", func() {
			listRes, err := ChannelAPI.ListChannels(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
		})

		It("should succeed even when the page token is too big", func() {
			listReq.PageToken = 20000
			listRes, err := ChannelAPI.ListChannels(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
		})
	})
})
