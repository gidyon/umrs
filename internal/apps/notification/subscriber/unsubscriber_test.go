package subscriber

import (
	"context"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/apps/notification/subscriber/model"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/subscriber"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Unsubsribing A User From A Channel #unsubscribe", func() {
	var (
		unSubscribeReq *subscriber.SubscriberRequest
		ctx            context.Context
	)

	BeforeEach(func() {
		unSubscribeReq = &subscriber.SubscriberRequest{
			AccountId: uuid.New().String(),
			Channel:   randomdata.Month(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Unsubscribing from a channel with incorrect/missing values", func() {
		It("should fail when the request to unsubscribe is nil", func() {
			unSubscribeReq = nil
			unSubscribeRes, err := SubsriberAPI.Unsubscribe(ctx, unSubscribeReq)
			Expect(err).Should(HaveOccurred())
			Expect(unSubscribeRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account id is missing", func() {
			unSubscribeReq.AccountId = ""
			unSubscribeRes, err := SubsriberAPI.Unsubscribe(ctx, unSubscribeReq)
			Expect(err).Should(HaveOccurred())
			Expect(unSubscribeRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel is missing", func() {
			unSubscribeReq.Channel = ""
			unSubscribeRes, err := SubsriberAPI.Unsubscribe(ctx, unSubscribeReq)
			Expect(err).Should(HaveOccurred())
			Expect(unSubscribeRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("subscribing to a channel with correct/valid request", func() {
		var channel, accountID string
		Context("Lets subscribe the account to a chanel first", func() {
			Context("Lets create subscriber account first", func() {
				It("should create subscriber account", func() {
					bs, err := json.Marshal([]string{"defualt"})
					Expect(err).ShouldNot(HaveOccurred())
					subscriberDB := &model.Subscriber{
						AccountID:  uuid.New().String(),
						Email:      randomdata.Email(),
						Phone:      randomdata.PhoneNumber(),
						SendMethod: notification.SendMethod_EMAIL_AND_SMS.String(),
						Channels:   bs,
					}
					err = DB.Create(subscriberDB).Error
					Expect(err).ShouldNot(HaveOccurred())
					accountID = subscriberDB.AccountID
				})
			})

			Context("Subscribing the account to a channel", func() {
				It("should subscribe user to channel", func() {
					subscribeReq := &subscriber.SubscriberRequest{
						AccountId: accountID,
						Channel:   randomdata.Month(),
					}
					ctx := auth.AddAdminMD(context.Background())
					subscribeRes, err := SubsriberAPI.Subscribe(ctx, subscribeReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(subscribeRes).ShouldNot(BeNil())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					channel = subscribeReq.Channel
					accountID = subscribeReq.AccountId
				})
			})

			Describe("Unsubscribing the account", func() {
				It("should subscribe admin when the request is valid", func() {
					unSubscribeReq.AccountId = accountID
					unSubscribeReq.Channel = channel
					unSubscribeRes, err := SubsriberAPI.Unsubscribe(ctx, unSubscribeReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(unSubscribeRes).ShouldNot(BeNil())
					Expect(status.Code(err)).Should(Equal(codes.OK))
				})
			})
		})
	})
})
