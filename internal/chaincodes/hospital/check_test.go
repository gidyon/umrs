package hospital

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Checking hospital suspension status #permission", func() {
	var (
		checkReq *hospital.CheckSuspensionRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		checkReq = &hospital.CheckSuspensionRequest{
			HospitalId: uuid.New().String(),
		}
		ctx = auth.AddInsuranceMD(context.Background())
	})

	Describe("Getting suspension status with malformed request", func() {
		It("should fail when the request is nil", func() {
			checkReq = nil
			checkRes, err := HospitalAPI.CheckSuspension(ctx, checkReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(checkRes).Should(BeNil())
		})
		It("should fail when hospital id in request is nil", func() {
			checkReq.HospitalId = ""
			checkRes, err := HospitalAPI.CheckSuspension(ctx, checkReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(checkRes).Should(BeNil())
		})
	})

	Describe("Checking permission with well-formed request", func() {
		It("should succeed but with status suspended when hospital id in request is unknown", func() {
			checkReq.HospitalId = uuid.New().String()
			checkRes, err := HospitalAPI.CheckSuspension(ctx, checkReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(checkRes).ShouldNot(BeNil())
			Expect(checkRes.Suspended).Should(BeTrue())
		})
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

			When("Checking for suspension status", func() {
				It("should succeed with status not suspended", func() {
					checkReq.HospitalId = hospitalID
					checkRes, err := HospitalAPI.CheckSuspension(ctx, checkReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(checkRes).ShouldNot(BeNil())
					Expect(checkRes.Suspended).Should(BeFalse())
				})
			})
		})
	})
})
