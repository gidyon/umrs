package hospital

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Searching hospitals resource #search", func() {
	var (
		searchReq *hospital.SearchHospitalsRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		searchReq = &hospital.SearchHospitalsRequest{
			PageSize:   15,
			PageNumber: 1,
			Query:      randomdata.Street(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Searching hospitals with malformed request", func() {
		It("should fail when the request is nil", func() {
			searchReq = nil
			searchRes, err := HospitalAPI.SearchHospitals(ctx, searchReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(searchRes).Should(BeNil())
		})
	})

	Describe("Searching hospitals with well formed request", func() {
		Context("Let's create at least one hospital first", func() {
			var query string
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
				query = addReq.Hospital.HospitalName
			})

			Describe("searching for a hospital", func() {

				BeforeEach(func() {
					searchReq.Query = query
				})

				It("should succeed even when page number or page size is weird", func() {
					searchReq.PageNumber = -10
					searchReq.PageSize = -20
					// The server will reset weird values to 1
					searchRes, err := HospitalAPI.SearchHospitals(ctx, searchReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(searchRes).ShouldNot(BeNil())
					Expect(len(searchRes.Hospitals)).ShouldNot(BeZero())
				})
				It("should succeed when everything OK", func() {
					searchRes, err := HospitalAPI.SearchHospitals(ctx, searchReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(searchRes).ShouldNot(BeNil())
					Expect(len(searchRes.Hospitals)).ShouldNot(BeZero())
				})
				It("should succeed when query is empty (will simply call ListHospitals)", func() {
					searchReq.Query = ""
					searchRes, err := HospitalAPI.SearchHospitals(ctx, searchReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(searchRes).ShouldNot(BeNil())
					Expect(len(searchRes.Hospitals)).ShouldNot(BeZero())
				})
			})
		})
	})
})
