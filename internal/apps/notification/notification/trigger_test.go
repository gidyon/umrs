package notification

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/notification"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Send Notification #Send", func() {
	var (
		SendReq *notification.Notification
		ctx        context.Context
	)

	BeforeEach(func() {
		SendReq = mockNotification()
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Sending notification with malformed request", func() {
		It("should fail when the notification is nil", func() {
			SendReq = nil
			SendRes, err := NotificationAPI.Send(ctx, SendReq)
			Expect(err).Should(HaveOccurred())
			Expect(SendRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when the notification owner id is missing", func() {
			SendReq.OwnerIds = nil
			SendRes, err := NotificationAPI.Send(ctx, SendReq)
			Expect(err).Should(HaveOccurred())
			Expect(SendRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		Context("Email Notification", func() {
			BeforeEach(func() {
				SendReq.SendMethod = notification.SendMethod_EMAIL
				SendReq.Payload = &notification.Notification_EmailNotification{
					EmailNotification: mockEmailNotification(),
				}
			})
			It("should fail when the notification send method is email and email notification is nil", func() {
				SendReq.Payload = nil
				SendRes, err := NotificationAPI.Send(ctx, SendReq)
				Expect(err).Should(HaveOccurred())
				Expect(SendRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			})
			It("should fail when the notification send method is email and email body is missing", func() {
				SendReq.GetEmailNotification().Body = ""
				SendRes, err := NotificationAPI.Send(ctx, SendReq)
				Expect(err).Should(HaveOccurred())
				Expect(SendRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			})
			It("should fail when the notification send method is email and email body ctype is missing", func() {
				SendReq.GetEmailNotification().BodyContentType = ""
				SendRes, err := NotificationAPI.Send(ctx, SendReq)
				Expect(err).Should(HaveOccurred())
				Expect(SendRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			})
			It("should fail when the notification send method is email and email subject is missing", func() {
				SendReq.GetEmailNotification().Subject = ""
				SendRes, err := NotificationAPI.Send(ctx, SendReq)
				Expect(err).Should(HaveOccurred())
				Expect(SendRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			})
		})
		Context("SMS Notification", func() {
			BeforeEach(func() {
				SendReq.SendMethod = notification.SendMethod_SMS
				SendReq.Payload = &notification.Notification_SmsNotification{
					SmsNotification: mockSMSNotification(),
				}
			})
			It("should fail when the notification send method is sms and sms notification is nil", func() {
				SendReq.Payload = nil
				SendRes, err := NotificationAPI.Send(ctx, SendReq)
				Expect(err).Should(HaveOccurred())
				Expect(SendRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			})
			It("should fail when the notification send method is sms and sms keyword is missing", func() {
				SendReq.GetSmsNotification().Keyword = ""
				SendRes, err := NotificationAPI.Send(ctx, SendReq)
				Expect(err).Should(HaveOccurred())
				Expect(SendRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			})
			It("should fail when the notification send method is sms and sms message is missing", func() {
				SendReq.GetSmsNotification().Message = ""
				SendRes, err := NotificationAPI.Send(ctx, SendReq)
				Expect(err).Should(HaveOccurred())
				Expect(SendRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			})
		})
	})

	Describe("Sending Notification with valid content", func() {
		It("should succeed when the notification is bulk and bulk channel is missing", func() {
			SendReq.Bulk = true
			SendReq.BulkChannel = ""
			SendRes, err := NotificationAPI.Send(ctx, SendReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(SendRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
		It("should succeed when the notification meets constraints", func() {
			SendRes, err := NotificationAPI.Send(ctx, SendReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(SendRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
		It("should succeed and save notification in database when the notification save is true", func() {
			SendReq.Save = true
			SendRes, err := NotificationAPI.Send(ctx, SendReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(SendRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
		It("should send notification to a single destination even when bulk is true", func() {
			SendReq.Bulk = true
			SendReq.BulkChannel = "devs-chat"
			SendRes, err := NotificationAPI.Send(ctx, SendReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(SendRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})