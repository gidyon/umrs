package ledger

import (
	"context"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newFilter() *ledger.Filter {
	return &ledger.Filter{
		Filter: true,
	}
}

var _ = Describe("Listing blocks from the ledger #list", func() {
	var (
		listReq *ledger.ListLogsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &ledger.ListLogsRequest{
			PageToken: 0,
			PageSize:  10,
			Filter:    newFilter(),
		}
		ctx = context.Background()
	})

	Describe("Listing blocks with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
		It("should succeed when page number or page size is lower than 0", func() {
			listReq.PageToken = -10
			listReq.PageSize = -10
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
		})
		It("should succeed when filter is nil", func() {
			listReq.Filter = nil
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
		})
	})

	Describe("Listing blocks with correct request", func() {
		var (
			creatorID string
			patientID string
			actor     ledger.Actor
		)
		Describe("Adding a new Log", func() {
			It("should succeed", func() {
				addReq := &ledger.AddLogRequest{
					Transaction: newTransaction(),
				}
				addRes, err := LedgerAPI.AddLog(ctx, addReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(addRes).ShouldNot(BeNil())
				Expect(addRes.OperationId).ShouldNot(BeZero())
				creatorID = addReq.Transaction.GetCreator().GetActorId()
				patientID = addReq.Transaction.GetPatient().GetActorId()
				actor = addReq.Transaction.GetCreator().GetActor()
			})
		})

		It("should succeed when the request is valid", func() {
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
		})
		It("should succeed when filter is false", func() {
			listReq.Filter.Filter = false
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
		})
		It("should succeed when filter by date is true", func() {
			listReq.Filter.ByDate = true
			listReq.Filter.StartDate = "01/12/2019"
			listReq.Filter.EndDate = "01/03/2020"
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
		})
		It("should succeed when filter by organization is true", func() {
			listReq.Filter.ByOrganizationId = true
			listReq.Filter.OrganizationIds = []string{randomOrganization(), randomOrganization()}
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			for _, blockPB := range listRes.Logs {
				orgID := blockPB.GetPayload().GetOrganization().GetActorId()
				Expect(orgID).Should(BeMemberOfStringSlice(listReq.Filter.OrganizationIds))
			}
		})
		It("should succeed when filter by operation is true", func() {
			listReq.Filter.ByOperation = true
			listReq.Filter.Operations = []ledger.Operation{
				randomOperation(), randomOperation(), randomOperation(),
			}
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			for _, blockPB := range listRes.Logs {
				opID := blockPB.GetPayload().GetOperation()
				Expect(opID).Should(BeElementOf(listReq.Filter.Operations))
			}
		})
		It("should succeed when filter by creator is true", func() {
			listReq.Filter.ByCreatorId = true
			listReq.Filter.CreatorIds = []string{
				creatorID, uuid.New().String(), uuid.New().String(),
			}
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			for _, blockPB := range listRes.Logs {
				opID := blockPB.GetPayload().GetCreator().GetActorId()
				Expect(opID).Should(BeElementOf(listReq.Filter.CreatorIds))
			}
		})
		It("should succeed when filter by creator actor is true", func() {
			listReq.Filter.ByCreatorActor = true
			listReq.Filter.CreatorActors = []ledger.Actor{
				randomActor(), randomActor(), actor,
			}
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			for _, blockPB := range listRes.Logs {
				opID := blockPB.GetPayload().GetCreator().GetActor()
				Expect(opID).Should(BeElementOf(listReq.Filter.CreatorActors))
			}
		})
		It("should succeed when filter by patient id is true", func() {
			listReq.Filter.ByPatientId = true
			listReq.Filter.PatientIds = []string{patientID}
			listRes, err := LedgerAPI.ListLogs(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			for _, blockPB := range listRes.Logs {
				opID := blockPB.GetPayload().GetPatient().GetActorId()
				Expect(opID).Should(BeElementOf(listReq.Filter.PatientIds))
			}
		})
	})
})
