package employment

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/employment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Updating employment data #update", func() {
	var (
		updateReq *employment.UpdateEmploymentRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &employment.UpdateEmploymentRequest{
			Employment: newEmployment(),
			Actor:      newActor(),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	Describe("Updating employment data with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := EmploymentAPI.UpdateEmployment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the actor is changed", func() {
			updateReq.Actor.Actor = int32(ledger.Actor_INSURANCE)
			updateRes, err := EmploymentAPI.UpdateEmployment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when employment is nil", func() {
			updateReq.Employment = nil
			updateRes, err := EmploymentAPI.UpdateEmployment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		Describe("Updating employment with missing data", func() {
			It("should fail when employment id is missing", func() {
				updateReq.Employment.EmploymentId = ""
				updateRes, err := EmploymentAPI.UpdateEmployment(ctx, updateReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(updateRes).Should(BeNil())
			})
			It("should fail when account id is missing", func() {
				updateReq.Employment.AccountId = ""
				updateRes, err := EmploymentAPI.UpdateEmployment(ctx, updateReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(updateRes).Should(BeNil())
			})
		})
	})

	Describe("Updating employment with well formed request", func() {
		var accountID, employmentID string
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
				employmentID = addRes.EmploymentId
			})
		})
		It("should succeed when request is valid", func() {
			updateReq.Employment.AccountId = accountID
			updateReq.Employment.EmploymentId = employmentID
			updateRes, err := EmploymentAPI.UpdateEmployment(ctx, updateReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(updateRes).ShouldNot(BeNil())
		})
	})
})
