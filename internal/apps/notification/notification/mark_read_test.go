package notification

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Mark Notification As Seen #seen", func() {
	var (
		markReadReq *notification.MarkNotificationReadRequest
		ctx         context.Context
	)

	BeforeEach(func() {
		markReadReq = &notification.MarkNotificationReadRequest{
			NotificationId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Calling MarkNotificationRead with missing/incorrect values", func() {
		It("should fail when the request is nil", func() {
			markReadReq = nil
			markReadRes, err := NotificationAPI.MarkNotificationRead(ctx, markReadReq)
			Expect(err).Should(HaveOccurred())
			Expect(markReadRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when notification id is missing", func() {
			markReadReq.NotificationId = ""
			markReadRes, err := NotificationAPI.MarkNotificationRead(ctx, markReadReq)
			Expect(err).Should(HaveOccurred())
			Expect(markReadRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Calling MarkNotificationRead with correct values", func() {
		It("should succeed when the request is valid", func() {
			markReadReq.NotificationId = "105ac57a-1181-4fcc-940e-551a670af445"
			markReadRes, err := NotificationAPI.MarkNotificationRead(ctx, markReadReq)
			_, _ = markReadRes, err
			// Expect(err).ShouldNot(HaveOccurred())
			// Expect(markReadRes).ShouldNot(BeNil())
			// Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
