package messaging

import (
	"context"
	"github.com/gidyon/umrs/pkg/api/messaging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Marking messages as read #read", func() {
	var (
		readReq *messaging.MessageRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		readReq = &messaging.MessageRequest{
			UserId: UserID,
		}
		ctx = context.Background()
	})

	Describe("Marking messages as read with malformed request", func() {
		It("should fail when the request is nil", func() {
			readReq = nil
			getRes, err := MessagingAPI.ReadAll(ctx, readReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when phone number is missing in request", func() {
			readReq.UserId = ""
			getRes, err := MessagingAPI.ReadAll(ctx, readReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Marking messages as read with correct request", func() {
		Context("Lets create a message first", func() {
			It("should succed in creating a message", func() {
				messagePB := fakeMessage()
				messageDB, err := GetMessageDB(messagePB)
				Expect(err).ShouldNot(HaveOccurred())
				err = MessagingServer.sqlDB.Create(messageDB).Error
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("Marking messages as read it should succeed", func() {
			It("should succeed", func() {
				getRes, err := MessagingAPI.ReadAll(ctx, readReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
			})

			Describe("All messages should now be marked as read", func() {
				It("should mark all messages as read", func() {
					listReq := &messaging.ListMessagesRequest{
						UserId:   UserID,
						PageSize: 100,
					}
					getRes, err := MessagingAPI.ListMessages(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes.Messages).ShouldNot(BeNil())

					// Expect all messages to be read
					for _, messagePB := range getRes.Messages {
						Expect(messagePB.Seen).Should(BeTrue())
					}
				})
			})
		})
	})
})
