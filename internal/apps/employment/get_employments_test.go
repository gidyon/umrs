package employment

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/employment"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting employments data #list", func() {
	var (
		getReq *employment.GetEmploymentsRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &employment.GetEmploymentsRequest{
			AccountId: uuid.New().String(),
			Actor:     newActor(),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	Describe("Getting employments with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := EmploymentAPI.GetEmployments(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when the actor is changed", func() {
			getReq.Actor.Actor = int32(ledger.Actor_INSURANCE)
			getRes, err := EmploymentAPI.GetEmployments(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when account id is missing", func() {
			getReq.AccountId = ""
			getRes, err := EmploymentAPI.GetEmployments(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when account id is unknown", func() {
			getRes, err := EmploymentAPI.GetEmployments(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting employment with well formed request", func() {
		var accountID string
		Context("Lets create employment first", func() {
			It("should succeed when request is valid", func() {
				addReq := &employment.AddEmploymentRequest{
					Employment: newEmployment(),
					Actor:      newActor(),
				}
				ctx := auth.AddHospitalMD(context.Background())
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(addRes).ShouldNot(BeNil())
				accountID = addReq.Employment.AccountId
			})
		})
		It("should succeed when request is valid", func() {
			getReq.AccountId = accountID
			getRes, err := EmploymentAPI.GetEmployments(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
			Expect(len(getRes.Employments)).ShouldNot(BeZero())
		})
	})
})
