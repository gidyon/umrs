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

var _ = Describe("Deleting employment data #delete", func() {
	var (
		delReq *employment.DeleteEmploymentRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &employment.DeleteEmploymentRequest{
			EmploymentId: uuid.New().String(),
			Actor:        newActor(),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	Describe("Deleting employment data with malformed request", func() {
		It("should fail when the request is nil", func() {
			delReq = nil
			delRes, err := EmploymentAPI.DeleteEmployment(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when the actor is changed", func() {
			delReq.Actor.Actor = int32(ledger.Actor_INSURANCE)
			delRes, err := EmploymentAPI.DeleteEmployment(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when employment id is missing", func() {
			delReq.EmploymentId = ""
			delRes, err := EmploymentAPI.DeleteEmployment(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when employment id is unknown", func() {
			delRes, err := EmploymentAPI.DeleteEmployment(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(delRes).Should(BeNil())
		})
	})

	Describe("Deleting employment with well formed request", func() {
		var employmentID string
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
				employmentID = addRes.EmploymentId
			})
		})
		It("should succeed when request is valid", func() {
			delReq.EmploymentId = employmentID
			delRes, err := EmploymentAPI.DeleteEmployment(ctx, delReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(delRes).ShouldNot(BeNil())
		})
	})
})
