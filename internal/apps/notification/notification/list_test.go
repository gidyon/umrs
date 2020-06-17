package notification

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ListNotifications For A User #list", func() {
	var (
		listReq *notification.ListNotificationsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &notification.ListNotificationsRequest{
			AccountId: uuid.New().String(),
			PageToken: 0,
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Calling ListNotifications with missing/incorrect values", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := NotificationAPI.ListNotifications(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(listRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account id is missing", func() {
			listReq.AccountId = ""
			listRes, err := NotificationAPI.ListNotifications(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(listRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Calling ListNotifications with valid values", func() {
		It("should succeed when the request is valid", func() {
			listRes, err := NotificationAPI.ListNotifications(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
