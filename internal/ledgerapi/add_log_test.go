package ledger

import (
	"context"
	"github.com/golang/protobuf/proto"
	"math/rand"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/google/uuid"
)

func randomActor() ledger.Actor {
	index := rand.Intn(len(ledger.Actor_name) - 1)
	if index == 0 {
		index = 1
	}
	return ledger.Actor(index)
}

func randomOperation() ledger.Operation {
	index := rand.Intn(len(ledger.Actor_name) - 1)
	if index == 0 {
		index = 1
	}
	return ledger.Operation(index)
}

var organizations = []string{
	"ORG1", "ORG2", "ORG3", "ORG4", "ORG5",
}

func randomOrganization() string {
	return organizations[rand.Intn(len(organizations)-1)]
}

func newTransaction() *ledger.Transaction {
	tx := &ledger.Transaction{
		Operation: randomOperation(),
		Creator: &ledger.ActorPayload{
			Actor:      randomActor(),
			ActorId:    uuid.New().String(),
			ActorNames: randomdata.SillyName(),
		},
		Patient: &ledger.ActorPayload{
			ActorId:    uuid.New().String(),
			ActorNames: randomdata.SillyName(),
		},
		Organization: &ledger.ActorPayload{
			ActorId:    uuid.New().String(),
			ActorNames: randomOrganization(),
		},
	}

	bs, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	tx.Details = bs

	return tx
}

var _ = Describe("Adding Log to ledger #add", func() {
	var (
		addReq *ledger.AddLogRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		addReq = &ledger.AddLogRequest{
			Transaction: newTransaction(),
		}
		ctx = context.Background()
	})

	Describe("Adding Log with malformed request", func() {
		It("should fail when the request is nil", func() {
			addReq = nil
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(addRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when transaction is nil", func() {
			addReq.Transaction = nil
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(addRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when actor is unknown", func() {
			addReq.Transaction.GetCreator().Actor = ledger.Actor_UNKNOWN
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(addRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when actor id is missiing", func() {
			addReq.Transaction.GetCreator().ActorId = ""
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(addRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when actor names is missiing", func() {
			addReq.Transaction.GetCreator().ActorNames = ""
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(addRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when operation is uknown", func() {
			addReq.Transaction.Operation = ledger.Operation_UKNOWN
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(addRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when transaction details is nil", func() {
			addReq.Transaction.Details = nil
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(addRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Adding a Log with well formed request", func() {
		It("should succeed if transaction is valid", func() {
			addRes, err := LedgerAPI.AddLog(ctx, addReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(addRes).ShouldNot(BeNil())
			Expect(addRes.OperationId).ShouldNot(BeZero())
		})
	})
})
