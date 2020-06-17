package messaging

import (
	"context"
	"github.com/gidyon/umrs/pkg/api/messaging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting messages #list", func() {
	var (
		listReq *messaging.ListMessagesRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &messaging.ListMessagesRequest{
			UserId: UserID,
		}
		ctx = context.Background()
	})

	Describe("Getting messages with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			getRes, err := MessagingAPI.ListMessages(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when user id is missing in request", func() {
			listReq.UserId = ""
			getRes, err := MessagingAPI.ListMessages(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting messages with correct request", func() {
		Context("Lets create a message first", func() {
			It("should succed in creating a message", func() {
				messagePB := fakeMessage()
				messageDB, err := GetMessageDB(messagePB)
				Expect(err).ShouldNot(HaveOccurred())
				err = MessagingServer.sqlDB.Create(messageDB).Error
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("Getting messages it should succeed", func() {
			It("should succeed", func() {
				getRes, err := MessagingAPI.ListMessages(ctx, listReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes.Messages).ShouldNot(BeNil())
				Expect(len(getRes.Messages)).ShouldNot(BeZero())
			})
		})

		When("Getting messages with type filters should succeed", func() {
			It("should succeed", func() {
				listReq.TypeFilters = []messaging.MessageType{messaging.MessageType_ALL}
				getRes, err := MessagingAPI.ListMessages(ctx, listReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes.Messages).ShouldNot(BeNil())
				Expect(len(getRes.Messages)).ShouldNot(BeZero())
			})
		})
	})
})
