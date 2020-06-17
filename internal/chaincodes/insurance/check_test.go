package insurance

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Checking insurance suspension status #permission", func() {
	var (
		checkReq *insurance.CheckSuspensionRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		checkReq = &insurance.CheckSuspensionRequest{
			InsuranceId: uuid.New().String(),
		}
		ctx = auth.AddInsuranceMD(context.Background())
	})

	Describe("Getting suspension status with malformed request", func() {
		It("should fail when the request is nil", func() {
			checkReq = nil
			checkRes, err := InsuranceAPI.CheckSuspension(ctx, checkReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(checkRes).Should(BeNil())
		})
		It("should fail when insurance id in request is nil", func() {
			checkReq.InsuranceId = ""
			checkRes, err := InsuranceAPI.CheckSuspension(ctx, checkReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(checkRes).Should(BeNil())
		})
	})

	Describe("Checking permission with well-formed request", func() {
		It("should succeed but with status suspended when insurance id in request is unknown", func() {
			checkReq.InsuranceId = uuid.New().String()
			checkRes, err := InsuranceAPI.CheckSuspension(ctx, checkReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(checkRes).ShouldNot(BeNil())
			Expect(checkRes.Suspended).Should(BeTrue())
		})
		Context("Let's create a insurance first", func() {
			var insuranceID string
			It("should create the insurance record in the database", func() {
				addReq := &insurance.AddInsuranceRequest{
					Insurance: newInsurance(),
				}
				ctx := auth.AddAdminMD(context.Background())
				addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(addRes).ShouldNot(BeNil())
				Expect(addRes.InsuranceId).ShouldNot(BeZero())
				insuranceID = addRes.InsuranceId
			})

			When("Checking for suspension status", func() {
				It("should succeed with status not suspended", func() {
					checkReq.InsuranceId = insuranceID
					checkRes, err := InsuranceAPI.CheckSuspension(ctx, checkReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(checkRes).ShouldNot(BeNil())
					Expect(checkRes.Suspended).Should(BeFalse())
				})
			})
		})
	})
})
