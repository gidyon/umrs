package hospital

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting hospital resource #get", func() {
	var (
		getReq *hospital.GetHospitalRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &hospital.GetHospitalRequest{
			HospitalId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(ctx)
	})

	Describe("Getting hospital resource with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := HospitalAPI.GetHospital(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when hospital id in the request is nil", func() {
			getReq.HospitalId = ""
			getRes, err := HospitalAPI.GetHospital(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when hospital id in the request is unknown", func() {
			getReq.HospitalId = uuid.New().String()
			getRes, err := HospitalAPI.GetHospital(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting hospital resource with well formed request", func() {
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

			Describe("Getting the hospital resource", func() {
				It("should get the hospital resource", func() {
					getReq.HospitalId = hospitalID
					getRes, err := HospitalAPI.GetHospital(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())
				})
			})
		})
	})
})
