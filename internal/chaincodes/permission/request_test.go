package permission

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/permission"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

func newRequesterProfile() *permission.BasicProfile {
	return &permission.BasicProfile{
		AccountId:          uuid.New().String(),
		FullName:           randomdata.SillyName(),
		ProfileThumbUrl:    "https://encrypted-tbn0.gstatic.com/images?q=tbn%3AANd9GcSvbmQ-tQLMbfDHL5MgHX2a7ar-p1mpjN_V-3GI78jE1MDJSU9P",
		RoleAtOrganization: "Medical Officer",
		OrganizationName:   randomdata.State(randomdata.Large) + " Hospital",
		OrganizationId:     uuid.New().String(),
		WorkEmail:          randomdata.Email(),
	}
}

func newActor() ledger.Actor {
	index := rand.Intn(len(ledger.Actor_name) - 1)
	if index == 0 {
		index = 1
	}
	return ledger.Actor(index)
}

func newActorPayload() *permission.Actor {
	return &permission.Actor{
		Group:    int32(newActor()),
		Id:       uuid.New().String(),
		FullName: randomdata.SillyName(),
	}
}

func allowedPermissionMethod() *permission.PermissionMethod {
	index := rand.Intn(len(permission.RequestPermissionMethod_name))
	if index == 0 {
		index = 1
	}
	return &permission.PermissionMethod{
		Method:   permission.RequestPermissionMethod(index),
		Payload:  "uknown yet",
		Metadata: map[string]string{},
	}
}

var _ = Describe("Requesting access to patient medical data #request", func() {
	var (
		requestReq *permission.RequestPermissionTokenRequest
		ctx        context.Context
	)

	BeforeEach(func() {
		requestReq = &permission.RequestPermissionTokenRequest{
			RequesterProfile: newRequesterProfile(),
			Patient:          newActorPayload(),
			Reason:           "Medical checkup",
			PermissionMethod: allowedPermissionMethod(),
			Requester:        newActorPayload(),
		}
		requestReq.Patient.Group = int32(ledger.Actor_PATIENT)
		requestReq.Requester.Group = int32(ledger.Actor_INSURANCE)
		tokenID := requestReq.Requester.Id
		ctx = auth.AddGroupAndIDMD(
			context.Background(), int32(ledger.Actor_INSURANCE), tokenID,
		)
	})

	Describe("Requesting access to patient medical data with malformed request", func() {
		It("should fail when the request is nil", func() {
			requestReq = nil
			requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(requestRes).Should(BeNil())
		})
		It("should fail when the actor is changed", func() {
			requestReq.Requester.Group = int32(ledger.Actor_GOVERNMENT)
			requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(requestRes).Should(BeNil())
		})
		It("should fail when the requester ID changed", func() {
			requestReq.Requester.Id = uuid.New().String()
			requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(requestRes).Should(BeNil())
		})
		Describe("Requesting for permission with missing values", func() {
			It("should fail when requester profile is nil", func() {
				requestReq.RequesterProfile = nil
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester account id is missing", func() {
				requestReq.RequesterProfile.AccountId = ""
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester account fullname is missing", func() {
				requestReq.RequesterProfile.FullName = ""
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester organization name is missing", func() {
				requestReq.RequesterProfile.OrganizationName = ""
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester organization id is missing", func() {
				requestReq.RequesterProfile.OrganizationId = ""
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester role ar organization is missing", func() {
				requestReq.RequesterProfile.RoleAtOrganization = ""
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when patient actor is nil", func() {
				requestReq.Patient = nil
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when patient id is missing", func() {
				requestReq.Patient.Id = ""
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester actor is nil", func() {
				requestReq.Requester = nil
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester id is missing", func() {
				requestReq.Requester.Id = ""
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when requester actor is unknown", func() {
				requestReq.Requester.Group = int32(ledger.Actor_UNKNOWN)
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when permission method is nil", func() {
				requestReq.PermissionMethod = nil
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
			It("should fail when permission method is unknown", func() {
				requestReq.PermissionMethod.Method = permission.RequestPermissionMethod_UNKNOWN
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(requestRes).Should(BeNil())
			})
		})
	})

	Describe("Requesting for permission with different requesting method", func() {
		It("should succeed if method is email", func() {
			requestReq.PermissionMethod.Method = permission.RequestPermissionMethod_EMAIL
			requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(requestRes).ShouldNot(BeNil())
		})
		It("should succeed if method is sms", func() {
			requestReq.PermissionMethod.Method = permission.RequestPermissionMethod_SMS
			requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(requestRes).ShouldNot(BeNil())
		})
		It("should succeed if method is ussd", func() {
			requestReq.PermissionMethod.Method = permission.RequestPermissionMethod_USSD
			requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(requestRes).ShouldNot(BeNil())
		})
	})

	Describe("Requesting for permission with correct request payload", func() {
		It("should succeed if the request is valid", func() {
			requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(requestRes).ShouldNot(BeNil())
		})
	})
})
