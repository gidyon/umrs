package subscriber

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/subscriber"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Subsribing A User To A Channel #subscribe", func() {
	var (
		subscribeReq *subscriber.SubscriberRequest
		ctx          context.Context
	)

	BeforeEach(func() {
		subscribeReq = &subscriber.SubscriberRequest{
			AccountId: uuid.New().String(),
			Channel:   randomdata.Month(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Subscribing to a channel with incorrect/missing values", func() {
		It("should fail when the request to subscribe is nil", func() {
			subscribeReq = nil
			subscribeRes, err := SubsriberAPI.Subscribe(ctx, subscribeReq)
			Expect(err).Should(HaveOccurred())
			Expect(subscribeRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account id is is missing", func() {
			subscribeReq.AccountId = ""
			subscribeRes, err := SubsriberAPI.Subscribe(ctx, subscribeReq)
			Expect(err).Should(HaveOccurred())
			Expect(subscribeRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel is is missing", func() {
			subscribeReq.Channel = ""
			subscribeRes, err := SubsriberAPI.Subscribe(ctx, subscribeReq)
			Expect(err).Should(HaveOccurred())
			Expect(subscribeRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Subscribing to a channel with correct/valid request", func() {
		It("should subscribe user when the request is valid", func() {
			subscribeRes, err := SubsriberAPI.Subscribe(ctx, subscribeReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(subscribeRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
