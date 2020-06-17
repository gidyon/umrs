package insurance

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Deleting a insurance resource #del", func() {
	var (
		delReq *insurance.DeleteInsuranceRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &insurance.DeleteInsuranceRequest{
			InsuranceId: uuid.New().String(),
			Reason:      randomdata.Paragraph(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Deleting a insurance with malformed request", func() {
		It("should fail when the request is nil", func() {
			delReq = nil
			delRes, err := InsuranceAPI.DeleteInsurance(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when the insurance id is missing", func() {
			delReq.InsuranceId = ""
			delRes, err := InsuranceAPI.DeleteInsurance(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when the insurance id is incorrect", func() {
			delRes, err := InsuranceAPI.DeleteInsurance(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when not an admin wants to delete", func() {
			ctx := auth.AddHospitalMD(context.Background())
			delRes, err := InsuranceAPI.DeleteInsurance(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(delRes).Should(BeNil())
		})
	})

	Describe("Deleting insurance with well frmed request", func() {
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

			Describe("Deleting the insurance resource", func() {
				It("should delete the insurance resource without error", func() {
					delReq.InsuranceId = insuranceID
					delRes, err := InsuranceAPI.DeleteInsurance(ctx, delReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(delRes).ShouldNot(BeNil())
				})
			})

			Describe("Getting the insurance should fail", func() {
				It("should delete the insurance resource without error", func() {
					getReq := &insurance.GetInsuranceRequest{
						InsuranceId: insuranceID,
					}
					getRes, err := InsuranceAPI.GetInsurance(ctx, getReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.NotFound))
					Expect(getRes).Should(BeNil())
					exist, ok := InsuranceAPIServer.allowedInsurances[insuranceID]
					Expect(exist).Should(BeFalse())
					Expect(ok).Should(BeFalse())
				})
			})
		})
	})
})
