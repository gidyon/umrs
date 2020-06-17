package employment

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/employment"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newActor() *employment.Actor {
	return &employment.Actor{
		Actor:   int32(ledger.Actor_HOSPITAL),
		ActorId: uuid.New().String(),
	}
}

func newEmployment() *employment.Employment {
	return &employment.Employment{
		AccountId:          uuid.New().String(),
		EmploymentId:       uuid.New().String(),
		FullName:           randomdata.SillyName(),
		ProfileThumbUrl:    "https://encrypted-tbn0.gstatic.com/images?q=tbn%3AANd9GcSvbmQ-tQLMbfDHL5MgHX2a7ar-p1mpjN_V-3GI78jE1MDJSU9P",
		EmploymentType:     employment.EmploymentType_PERMANENT,
		JoinedDate:         "03-05-2017",
		OrganizationType:   "HOSPITAL",
		OrganizationName:   randomdata.State(randomdata.Large) + " Hospital",
		OrganizationId:     uuid.New().String(),
		RoleAtOrganization: "Medical Officer",
		WorkEmail:          randomdata.Email(),
		EmploymentVerified: true,
		StillEmployed:      true,
		IsRecent:           true,
	}
}

var _ = Describe("Adding employment data #add", func() {
	var (
		addReq *employment.AddEmploymentRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		addReq = &employment.AddEmploymentRequest{
			Employment: newEmployment(),
			Actor:      newActor(),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	Describe("Adding employment data with malformed request", func() {
		It("should fail when the request is nil", func() {
			addReq = nil
			addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the actor is changed", func() {
			addReq.Actor.Actor = int32(ledger.Actor_INSURANCE)
			addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the request is nil", func() {
			addReq.Employment = nil
			addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		Describe("Adding employment with missing data", func() {
			It("should fail when account id is missing", func() {
				addReq.Employment.AccountId = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when full name is missing", func() {
				addReq.Employment.FullName = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when employment type is unknown", func() {
				addReq.Employment.EmploymentType = employment.EmploymentType_UNKNOWN
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when join date is missing", func() {
				addReq.Employment.JoinedDate = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when organization type is missing", func() {
				addReq.Employment.OrganizationType = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when organization name is missing", func() {
				addReq.Employment.OrganizationName = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when organization id is missing", func() {
				addReq.Employment.OrganizationId = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when role at organization is missing", func() {
				addReq.Employment.RoleAtOrganization = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
			It("should fail when work email is missing", func() {
				addReq.Employment.WorkEmail = ""
				addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(addRes).Should(BeNil())
			})
		})
	})

	Describe("Adding employment with well formed request", func() {
		It("should succeed when request is valid", func() {
			addRes, err := EmploymentAPI.AddEmployment(ctx, addReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(addRes).ShouldNot(BeNil())
		})
	})
})
