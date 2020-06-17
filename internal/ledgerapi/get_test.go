package ledger

import (
	"context"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting Log from the ledger #get", func() {
	var (
		getReq *ledger.GetLogRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &ledger.GetLogRequest{
			Hash: uuid.New().String(),
		}
		ctx = context.Background()
	})

	Describe("adding Log with malformed request", func() {
		It("should fail if the request is nil", func() {
			getReq = nil
			getRes, err := LedgerAPI.GetLog(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail if Log hash is missing in the request", func() {
			getReq.Hash = ""
			getRes, err := LedgerAPI.GetLog(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail if Log hash is incorrect", func() {
			getRes, err := LedgerAPI.GetLog(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting Log with valid request", func() {
		Context("Let's create Log first", func() {
			var operationID string
			It("should succeed if transaction is valid", func() {
				addRes, err := LedgerAPI.AddLog(ctx, &ledger.AddLogRequest{
					Transaction: newTransaction(),
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(addRes).ShouldNot(BeNil())
				Expect(addRes.OperationId).ShouldNot(BeZero())
				operationID = addRes.OperationId
			})

			Describe("Getting the Log", func() {
				It("should fail because the Log is not yet present in the ledger", func() {
					getReq.Hash = operationID
					getRes, err := LedgerAPI.GetLog(ctx, getReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.NotFound))
					Expect(getRes).Should(BeNil())
				})
			})
		})
	})
})
