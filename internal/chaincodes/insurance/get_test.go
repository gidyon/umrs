package insurance

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting insurance resource #get", func() {
	var (
		getReq *insurance.GetInsuranceRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &insurance.GetInsuranceRequest{
			InsuranceId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(ctx)
	})

	Describe("Getting insurance resource with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := InsuranceAPI.GetInsurance(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when insurance id in the request is nil", func() {
			getReq.InsuranceId = ""
			getRes, err := InsuranceAPI.GetInsurance(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting insurance resource with well formed request", func() {
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

			Describe("Getting the insurance resource", func() {
				It("should get the insurance resource", func() {
					getReq.InsuranceId = insuranceID
					getRes, err := InsuranceAPI.GetInsurance(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())
				})
			})
		})
	})
})
