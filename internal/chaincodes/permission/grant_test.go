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
)

func grantPaylod() *permission.GrantPermissionTokenPayload {
	return &permission.GrantPermissionTokenPayload{
		Reason:           randomdata.Paragraph(),
		RequesterProfile: newRequesterProfile(),
	}
}

var _ = Describe("Granting permission #grant", func() {
	var (
		grantReq *permission.GrantPermissionTokenRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		grantReq = &permission.GrantPermissionTokenRequest{
			Requester:    newActorPayload(),
			Organization: newActorPayload(),
			Patient:      newActorPayload(),
			Payload:      grantPaylod(),
		}
		grantReq.Patient.Group = int32(ledger.Actor_PATIENT)
		patientID := grantReq.Patient.Id
		grantReq.AuthorizationToken = auth.GenPatientAccessToken(context.Background(), patientID)
		ctx = context.Background()
	})

	Describe("Granting permission with malformed request", func() {
		It("should fail when the request is nil", func() {
			grantReq = nil
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when patient is nil", func() {
			grantReq.Patient = nil
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when patient id is missing", func() {
			grantReq.Patient.Id = ""
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when requester is nil", func() {
			grantReq.Requester = nil
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when requester id is missing", func() {
			grantReq.Requester.Id = ""
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when authorization token is missing", func() {
			grantReq.AuthorizationToken = ""
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when organization is nil", func() {
			grantReq.Organization = nil
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when organization id is missing", func() {
			grantReq.Organization.Id = ""
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when actor is unknown", func() {
			grantReq.Requester.Group = int32(ledger.Actor_UNKNOWN)
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(grantRes).Should(BeNil())
		})
		It("should fail when patient id is different than one in token", func() {
			grantReq.Patient.Id = uuid.New().String()
			grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(grantRes).Should(BeNil())
		})
	})

	Describe("Granting permission with correct request", func() {
		var patientID, requesterID string
		Describe("Lets request for permission access", func() {
			It("should succeed", func() {
				requestReq := &permission.RequestPermissionTokenRequest{
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
				requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(requestRes).ShouldNot(BeNil())
				requesterID = tokenID
				patientID = requestReq.GetPatient().GetId()
			})
			Describe("Lets grant permission now", func() {
				It("should succeed", func() {
					grantReq.Patient.Id = patientID
					grantReq.Requester.Id = requesterID
					grantReq.AuthorizationToken = auth.GenPatientAccessToken(context.Background(), patientID)
					grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(grantRes).ShouldNot(BeNil())
				})
				It("should succeed getting the permission token", func() {
					getReq := &permission.GetPermissionTokenRequest{
						PatientId: patientID,
						Actor:     newActorPayload(),
					}
					getReq.Actor.Group = int32(ledger.Actor_INSURANCE)
					getReq.Actor.Id = requesterID
					ctx := auth.AddGroupAndIDMD(
						context.Background(), getReq.Actor.Group, requesterID,
					)
					getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())
					Expect(getRes.Allowed).Should(BeTrue())
					Expect(getRes.AccessToken).ShouldNot(BeZero())
				})
			})
		})

	})
})
