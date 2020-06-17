package insurance

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Listing insurances resource #list", func() {
	var (
		listReq *insurance.ListInsurancesRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &insurance.ListInsurancesRequest{
			PageSize:   15,
			PageNumber: 1,
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Listing insurances with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := InsuranceAPI.ListInsurances(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("Listing insurances with well formed request", func() {
		It("should succeed even when page number or page size is weird", func() {
			listReq.PageNumber = -10
			listReq.PageSize = -20
			// The server will reset weird values to 1
			listRes, err := InsuranceAPI.ListInsurances(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			Expect(len(listRes.Insurances)).ShouldNot(BeZero())
		})
		It("should succeed when everything OK", func() {
			listRes, err := InsuranceAPI.ListInsurances(ctx, listReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(listRes).ShouldNot(BeNil())
			Expect(len(listRes.Insurances)).ShouldNot(BeZero())
		})
	})
})
