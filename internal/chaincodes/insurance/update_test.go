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

var _ = Describe("Updating insurance resource to database #update", func() {
	var (
		updateReq *insurance.UpdateInsuranceRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &insurance.UpdateInsuranceRequest{
			InsuranceId: uuid.New().String(),
			Insurance:   newInsurance(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Updating insurance with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			addRes, err := InsuranceAPI.UpdateInsurance(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance in request is nil", func() {
			updateReq.Insurance = nil
			addRes, err := InsuranceAPI.UpdateInsurance(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance id in request is nil", func() {
			updateReq.InsuranceId = ""
			addRes, err := InsuranceAPI.UpdateInsurance(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance id in request is unknown", func() {
			updateReq.InsuranceId = uuid.New().String()
			addRes, err := InsuranceAPI.UpdateInsurance(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(addRes).Should(BeNil())
		})
	})

	Describe("Updating insurance with well formed request", func() {
		Context("Let's create a insurance first", func() {
			var insuranceID string
			var insurancePB, updateInsurance *insurance.Insurance
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
				insurancePB = addReq.Insurance
			})

			Describe("Updating the insurance resource", func() {
				It("should update insurance without error", func() {
					updateReq.InsuranceId = insuranceID
					updateReq.Insurance.About = ""
					updateReq.Insurance.WebsiteUrl = ""
					addRes, err := InsuranceAPI.UpdateInsurance(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(addRes).ShouldNot(BeNil())
					Expect(addRes.InsuranceId).ShouldNot(BeZero())
					updateInsurance = updateReq.Insurance
				})
			})

			Describe("Getting the insurance resource", func() {
				It("should get the insurance resource", func() {
					getReq := &insurance.GetInsuranceRequest{
						InsuranceId: uuid.New().String(),
					}
					getReq.InsuranceId = insuranceID
					getRes, err := InsuranceAPI.GetInsurance(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())

					// Expectations
					Expect(getRes.About).Should(Equal(insurancePB.About))
					Expect(getRes.WebsiteUrl).Should(Equal(insurancePB.WebsiteUrl))
					Expect(getRes.About).Should(Equal(insurancePB.About))
					Expect(getRes.LogoUrl).Should(Equal(updateInsurance.LogoUrl))
					Expect(getRes.LogoUrl).Should(Equal(updateInsurance.LogoUrl))
				})
			})

			Describe("Suspending the insurance resource ", func() {
				It("should update insurance without error", func() {
					updateReq.InsuranceId = insuranceID
					updateReq.Suspend = true
					updateReq.Reason = randomdata.Paragraph()
					addRes, err := InsuranceAPI.UpdateInsurance(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(addRes).ShouldNot(BeNil())
					Expect(addRes.InsuranceId).ShouldNot(BeZero())
					updateInsurance = updateReq.Insurance
				})
			})
		})
	})
})
