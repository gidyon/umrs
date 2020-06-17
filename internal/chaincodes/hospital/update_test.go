package hospital

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Updating hospital resource to database #update", func() {
	var (
		updateReq *hospital.UpdateHospitalRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &hospital.UpdateHospitalRequest{
			HospitalId: uuid.New().String(),
			Hospital:   newHospital(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Updating hospital with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			addRes, err := HospitalAPI.UpdateHospital(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital in request is nil", func() {
			updateReq.Hospital = nil
			addRes, err := HospitalAPI.UpdateHospital(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital id in request is nil", func() {
			updateReq.HospitalId = ""
			addRes, err := HospitalAPI.UpdateHospital(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital id in request is unknown", func() {
			updateReq.HospitalId = uuid.New().String()
			addRes, err := HospitalAPI.UpdateHospital(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when reason is missing and suspend is true", func() {
			updateReq.Suspend = true
			updateReq.HospitalId = uuid.New().String()
			addRes, err := HospitalAPI.UpdateHospital(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
	})

	Describe("Updating hospital with well formed request", func() {
		Context("Let's create a hospital first", func() {
			var hospitalID string
			var hospitalPB, updateHospital *hospital.Hospital
			It("should create the hospital record in the database", func() {
				addReq := &hospital.AddHospitalRequest{
					Hospital: newHospital(),
				}
				ctx := auth.AddAdminMD(context.Background())
				addRes, err := HospitalAPI.AddHospital(ctx, addReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(addRes).ShouldNot(BeNil())
				Expect(addRes.HospitalId).ShouldNot(BeZero())
				hospitalID = addRes.HospitalId
				hospitalPB = addReq.Hospital
			})

			Describe("Updating the hospital resource", func() {
				It("should update hospital without error", func() {
					updateReq.HospitalId = hospitalID
					updateReq.Hospital.County = ""
					updateReq.Hospital.WebsiteUrl = ""
					addRes, err := HospitalAPI.UpdateHospital(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(addRes).ShouldNot(BeNil())
					Expect(addRes.HospitalId).ShouldNot(BeZero())
					updateHospital = updateReq.Hospital
				})
			})

			Describe("Getting the hospital resource", func() {
				It("should get the hospital resource", func() {
					getReq := &hospital.GetHospitalRequest{
						HospitalId: uuid.New().String(),
					}
					getReq.HospitalId = hospitalID
					getRes, err := HospitalAPI.GetHospital(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())

					// Expectations
					Expect(getRes.County).Should(Equal(hospitalPB.County))
					Expect(getRes.WebsiteUrl).Should(Equal(hospitalPB.WebsiteUrl))
					Expect(getRes.County).Should(Equal(hospitalPB.County))
					Expect(getRes.LogoUrl).Should(Equal(updateHospital.LogoUrl))
					Expect(getRes.LogoUrl).Should(Equal(updateHospital.LogoUrl))
				})
			})

			Describe("Suspending the hospital resource ", func() {
				It("should update hospital without error", func() {
					updateReq.HospitalId = hospitalID
					updateReq.Suspend = true
					updateReq.Reason = randomdata.Paragraph()
					addRes, err := HospitalAPI.UpdateHospital(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(addRes).ShouldNot(BeNil())
					Expect(addRes.HospitalId).ShouldNot(BeZero())
					updateHospital = updateReq.Hospital
				})
			})
		})
	})
})
