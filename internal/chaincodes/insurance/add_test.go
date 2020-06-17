package insurance

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newInsurance() *insurance.Insurance {
	about := randomdata.Paragraph()
	if len(about) > 256 {
		about = about[:256]
	}
	return &insurance.Insurance{
		InsuranceName:    randomdata.SillyName() + " Insurance",
		About:            about,
		WebsiteUrl:       randomdata.Address(),
		LogoUrl:          randomdata.Address(),
		SupportTelNumber: randomdata.PhoneNumber(),
		SupportEmail:     randomdata.Email(),
		AdminEmails: []string{
			randomdata.Email(), randomdata.Email(), randomdata.Email(),
		},
		Permission: true,
	}
}

var _ = Describe("Adding insurance resource to database #add", func() {
	var (
		addReq *insurance.AddInsuranceRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		addReq = &insurance.AddInsuranceRequest{
			Insurance: newInsurance(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Adding insurance with malformed request", func() {
		It("should fail when the request is nil", func() {
			addReq = nil
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance in request is nil", func() {
			addReq.Insurance = nil
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance name in request is missing", func() {
			addReq.Insurance.InsuranceName = ""
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when website url in request is missing", func() {
			addReq.Insurance.WebsiteUrl = ""
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance about in request is missing", func() {
			addReq.Insurance.About = ""
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance about in request is too large", func() {
			addReq.Insurance.About += addReq.Insurance.About
			addReq.Insurance.About += addReq.Insurance.About
			addReq.Insurance.About += addReq.Insurance.About
			addReq.Insurance.About += addReq.Insurance.About
			addReq.Insurance.About += addReq.Insurance.About
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance support email in request is missing", func() {
			addReq.Insurance.SupportEmail = ""
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when insurance support telephone in request is missing", func() {
			addReq.Insurance.SupportTelNumber = ""
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when admin emails in request is missing", func() {
			addReq.Insurance.AdminEmails = nil
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
	})

	Describe("Adding insurance with well formed request", func() {
		It("should add insurance without error", func() {
			addRes, err := InsuranceAPI.AddInsurance(ctx, addReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(addRes).ShouldNot(BeNil())
			Expect(addRes.InsuranceId).ShouldNot(BeZero())
		})
	})
})
