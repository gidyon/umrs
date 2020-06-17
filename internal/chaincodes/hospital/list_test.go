package hospital

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Listing hospitals resource #list", func() {
	var (
		listReq *hospital.ListHospitalsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &hospital.ListHospitalsRequest{
			PageSize:   15,
			PageNumber: 1,
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Listing hospitals with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := HospitalAPI.ListHospitals(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("Listing hospitals with well formed request", func() {
		It("should succeed even when page number or page size is weird", func() {
			listReq.PageNumber = -10
			listReq.PageSize = -20
			// The server will reset weird values to 1
			listRes, err := HospitalAPI.ListHospitals(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			Expect(len(listRes.Hospitals)).ShouldNot(BeZero())
		})
		It("should succeed when everything OK", func() {
			listRes, err := HospitalAPI.ListHospitals(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			Expect(len(listRes.Hospitals)).ShouldNot(BeZero())
		})
	})
})
