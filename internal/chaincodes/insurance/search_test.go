package insurance

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Searching insurances resource #search", func() {
	var (
		searchReq *insurance.SearchInsurancesRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		searchReq = &insurance.SearchInsurancesRequest{
			PageSize:   15,
			PageNumber: 1,
			Query:      randomdata.Street(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Searching insurances with malformed request", func() {
		It("should fail when the request is nil", func() {
			searchReq = nil
			searchRes, err := InsuranceAPI.SearchInsurances(ctx, searchReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(searchRes).Should(BeNil())
		})
	})

	Describe("Searching insurances with well formed request", func() {
		Context("Let's create at least one insurance first", func() {
			var query string
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
				query = addReq.Insurance.InsuranceName
			})

			Describe("searching for a insurance", func() {

				BeforeEach(func() {
					searchReq.Query = query
				})

				It("should succeed even when page number or page size is weird", func() {
					searchReq.PageNumber = -10
					searchReq.PageSize = -20
					// The server will reset weird values to 1
					searchRes, err := InsuranceAPI.SearchInsurances(ctx, searchReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(searchRes).ShouldNot(BeNil())
					Expect(len(searchRes.Insurances)).ShouldNot(BeZero())
				})
				It("should succeed when everything OK", func() {
					searchRes, err := InsuranceAPI.SearchInsurances(ctx, searchReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(searchRes).ShouldNot(BeNil())
					Expect(len(searchRes.Insurances)).ShouldNot(BeZero())
				})
				It("should succeed when query is empty (will simply call ListInsurances)", func() {
					searchReq.Query = ""
					searchRes, err := InsuranceAPI.SearchInsurances(ctx, searchReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(searchRes).ShouldNot(BeNil())
					Expect(len(searchRes.Insurances)).ShouldNot(BeZero())
				})
			})
		})
	})
})
