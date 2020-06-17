package messaging

import (
	"context"
	"github.com/gidyon/umrs/pkg/api/messaging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting how many new messages #count", func() {
	var (
		getReq *messaging.MessageRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &messaging.MessageRequest{
			UserId: UserID,
		}
		ctx = context.Background()
	})

	Describe("Getting how many new messages with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := MessagingAPI.GetNewMessagesCount(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when user id is missing in request", func() {
			getReq.UserId = ""
			getRes, err := MessagingAPI.GetNewMessagesCount(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting how many new messages with correct request", func() {
		Context("Lets create a message first", func() {
			It("should succed in creating a message", func() {
				messagePB := fakeMessage()
				messageDB, err := GetMessageDB(messagePB)
				Expect(err).ShouldNot(HaveOccurred())
				err = MessagingServer.sqlDB.Create(messageDB).Error
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Describe("Getting how many new messages should succeed", func() {
			It("should succeed", func() {
				getReq.UserId = UserID
				getRes, err := MessagingAPI.GetNewMessagesCount(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
				Expect(getRes.Count).ShouldNot(BeZero())
			})
		})
	})
})
