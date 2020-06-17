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

var _ = Describe("Listing Subscribers From A Channel #list", func() {
	var (
		listReq *subscriber.ListSubscribersRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &subscriber.ListSubscribersRequest{
			Channel:   "tree",
			PageToken: 0,
			PageSize:  100,
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Listing subscribers from a channel with incorrect/missing values", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := SubsriberAPI.ListSubscribers(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(listRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when channel is is missing", func() {
			listReq.Channel = ""
			listRes, err := SubsriberAPI.ListSubscribers(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(listRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Listing subscribers from a channel with correct/valid request", func() {
		channel := randomdata.Adjective()
		Context("Lets create a subscriber account first", func() {
			It("should create subscriber account", func() {
				bs, err := json.Marshal([]string{channel})
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
			})
		})
		It("should list subscribers for a channel", func() {
			listReq.Channel = channel
			listRes, err := SubsriberAPI.ListSubscribers(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			Expect(len(listRes.Subscribers)).ShouldNot(BeZero())
		})
	})
})
