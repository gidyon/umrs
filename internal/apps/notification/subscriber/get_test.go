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

var _ = Describe("Getting subscriber #get", func() {
	var (
		getReq *subscriber.GetSubscriberRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &subscriber.GetSubscriberRequest{
			AccountId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Getting subscriber with misformed request", func() {
		It("should fail when the request is absent/nil", func() {
			getReq = nil
			getRes, err := SubsriberAPI.GetSubscriber(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(getRes).To(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when the account_id is missing in request", func() {
			getReq.AccountId = ""
			getRes, err := SubsriberAPI.GetSubscriber(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(getRes).To(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when the account_id does not exist", func() {
			getRes, err := SubsriberAPI.GetSubscriber(ctx, getReq)
			Expect(err).To(HaveOccurred())
			Expect(getRes).To(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
		})
	})

	Describe("Getting subscriber with valid request", func() {
		var accountID string
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

		Describe("Getting the subscriber", func() {
			It("should succeed when account id exists and is valid", func() {
				getReq.AccountId = accountID
				getRes, err := SubsriberAPI.GetSubscriber(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
			})
		})
	})
})
