package hospital

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newHospital() *hospital.Hospital {
	return &hospital.Hospital{
		HospitalName: randomdata.SillyName() + " Hospital",
		WebsiteUrl:   randomdata.Address(),
		LogoUrl:      randomdata.Address(),
		County:       randomdata.State(randomdata.Large),
		SubCounty:    randomdata.Street(),
		AdminEmails: []string{
			randomdata.Email(), randomdata.Email(), randomdata.Email(),
		},
		Permission: true,
	}
}

var _ = Describe("Adding hospital resource to database #add", func() {
	var (
		addReq *hospital.AddHospitalRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		addReq = &hospital.AddHospitalRequest{
			Hospital: newHospital(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Adding hospital with malformed request", func() {
		It("should fail when the request is nil", func() {
			addReq = nil
			addRes, err := HospitalAPI.AddHospital(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital in request is nil", func() {
			addReq.Hospital = nil
			addRes, err := HospitalAPI.AddHospital(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital name in request is missing", func() {
			addReq.Hospital.HospitalName = ""
			addRes, err := HospitalAPI.AddHospital(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when county in request is missing", func() {
			addReq.Hospital.County = ""
			addRes, err := HospitalAPI.AddHospital(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when sub-county in request is missing", func() {
			addReq.Hospital.SubCounty = ""
			addRes, err := HospitalAPI.AddHospital(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when admin emails in request is missing", func() {
			addReq.Hospital.AdminEmails = nil
			addRes, err := HospitalAPI.AddHospital(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
	})

	Describe("Adding hospital with well formed request", func() {
		It("should add hospital without error", func() {
			addRes, err := HospitalAPI.AddHospital(ctx, addReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(addRes).ShouldNot(BeNil())
			Expect(addRes.HospitalId).ShouldNot(BeZero())
		})
	})
})
