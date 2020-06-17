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

var _ = Describe("Deleting a hospital resource #del", func() {
	var (
		delReq *hospital.DeleteHospitalRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &hospital.DeleteHospitalRequest{
			HospitalId: uuid.New().String(),
			Reason:     randomdata.Paragraph(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Deleting a hospital with malformed request", func() {
		It("should fail when the request is nil", func() {
			delReq = nil
			delRes, err := HospitalAPI.DeleteHospital(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when the hospital id is missing", func() {
			delReq.HospitalId = ""
			delRes, err := HospitalAPI.DeleteHospital(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when the hospital id is incorrect", func() {
			delRes, err := HospitalAPI.DeleteHospital(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(delRes).Should(BeNil())
		})
		It("should fail when not an admin wants to delete", func() {
			ctx := auth.AddHospitalMD(context.Background())
			delRes, err := HospitalAPI.DeleteHospital(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(delRes).Should(BeNil())
		})
	})

	Describe("Deleting hospital with well formed request", func() {
		Context("Let's create a hospital first", func() {
			var hospitalID string
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
			})

			Describe("Deleting the hospital resource", func() {
				It("should delete the hospital resource without error", func() {
					delReq.HospitalId = hospitalID
					delRes, err := HospitalAPI.DeleteHospital(ctx, delReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(delRes).ShouldNot(BeNil())
				})
			})

			Describe("Getting the hospital should fail", func() {
				It("should delete the hospital resource without error", func() {
					getReq := &hospital.GetHospitalRequest{
						HospitalId: hospitalID,
					}
					getRes, err := HospitalAPI.GetHospital(ctx, getReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.NotFound))
					Expect(getRes).Should(BeNil())
					exist, ok := HospitalAPIServer.allowedHospitals[hospitalID]
					Expect(exist).Should(BeFalse())
					Expect(ok).Should(BeFalse())
				})
			})
		})
	})
})
