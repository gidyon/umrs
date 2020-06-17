package notification

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	subscriber_model "github.com/gidyon/umrs/internal/apps/notification/subscriber/model"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Create Notification #create", func() {
	var (
		createReq *notification.CreateNotificationAccountRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &notification.CreateNotificationAccountRequest{
			Channels:   []string{randomdata.Adjective(), randomdata.Adjective(), randomdata.Adjective()},
			AccountId:  uuid.New().String(),
			Email:      randomdata.Email(),
			Phone:      randomdata.PhoneNumber(),
			SendMethod: notification.SendMethod_EMAIL_AND_SMS,
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Creating Notification Account malformed request", func() {
		It("should fail when the create request is nil", func() {
			createReq = nil
			createRes, err := NotificationAPI.CreateNotificationAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account id is missing", func() {
			createReq.AccountId = ""
			createRes, err := NotificationAPI.CreateNotificationAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when the email and phone are missing", func() {
			createReq.Email = ""
			createReq.Phone = ""
			createRes, err := NotificationAPI.CreateNotificationAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Creating Notification with correct values", func() {
		var accountID string
		It("should succeed when the create request is valid", func() {
			createRes, err := NotificationAPI.CreateNotificationAccount(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			accountID = createReq.AccountId
		})

		Describe("Getting the account", func() {
			It("should exist in database", func() {
				accountDB := &subscriber_model.Subscriber{}
				err := DB.First(accountDB, "account_id=?", accountID).Error
				Expect(err).ShouldNot(HaveOccurred())
				Expect(accountDB).ShouldNot(BeNil())
			})
		})

		It("should succeed when the email is missing but phone is not", func() {
			createReq.Email = ""
			createRes, err := NotificationAPI.CreateNotificationAccount(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
		It("should succeed when the phone is missing but email is not", func() {
			createReq.Phone = ""
			createRes, err := NotificationAPI.CreateNotificationAccount(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
