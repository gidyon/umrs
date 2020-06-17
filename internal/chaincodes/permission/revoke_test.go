package permission

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/permission"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Revoking permission #revoke", func() {
	var (
		revokeReq *permission.RevokePermissionTokenRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		revokeReq = &permission.RevokePermissionTokenRequest{
			RequesterId: uuid.New().String(),
			PatientId:   uuid.New().String(),
		}
		ctx = auth.AddPatientMD(ctx, revokeReq.PatientId)
	})

	Describe("Revoking permission with malformed request #grant", func() {
		It("should fail when the request is nil", func() {
			revokeReq = nil
			revokeRes, err := PermissionAPI.RevokePermissionToken(ctx, revokeReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(revokeRes).Should(BeNil())
		})
		It("should fail when patient id is missing", func() {
			revokeReq.PatientId = ""
			revokeRes, err := PermissionAPI.RevokePermissionToken(ctx, revokeReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(revokeRes).Should(BeNil())
		})
		It("should fail when requester id is missing", func() {
			revokeReq.RequesterId = ""
			revokeRes, err := PermissionAPI.RevokePermissionToken(ctx, revokeReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(revokeRes).Should(BeNil())
		})
		It("should fail when patient id is different than one in token", func() {
			revokeReq.PatientId = uuid.New().String()
			revokeRes, err := PermissionAPI.RevokePermissionToken(ctx, revokeReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(revokeRes).Should(BeNil())
		})
		It("should fail when group in context is not patient group", func() {
			ctx := auth.AddInsuranceMD(context.Background())
			revokeRes, err := PermissionAPI.RevokePermissionToken(ctx, revokeReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(revokeRes).Should(BeNil())
		})
	})

	Describe("Revoking permission with correct request", func() {
		// Grant
		// Try To Get Token - Succeed
		// Revoke
		// Try To Get Token - Fails
		var (
			requesterID string
			patientID   string
			group       int32
		)
		Describe("Lets request for permission", func() {
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
				group = requestReq.GetRequester().Group
			})
		})

		Describe("Lets grant requester permision first", func() {
			It("should succeed granting access", func() {
				grantReq := &permission.GrantPermissionTokenRequest{
					Requester:    newActorPayload(),
					Organization: newActorPayload(),
					Patient:      newActorPayload(),
				}
				grantReq.Patient.Id = patientID
				grantReq.Requester.Group = group
				grantReq.Requester.Id = requesterID
				grantReq.AuthorizationToken = auth.GenPatientAccessToken(context.Background(), patientID)
				ctx := context.Background()
				grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(grantRes).ShouldNot(BeNil())
			})

			Describe("Lets try to get permission token", func() {
				It("should succeed getting the permission token", func() {
					getReq := &permission.GetPermissionTokenRequest{
						PatientId: patientID,
						Actor:     newActorPayload(),
					}
					getReq.Actor.Group = group
					getReq.Actor.Id = requesterID
					ctx = auth.AddGroupAndIDMD(
						context.Background(), getReq.Actor.Group, getReq.Actor.Id,
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

		Describe("Lets revoke their permission", func() {
			It("should succeed in revoking their permission", func() {
				revokeReq.PatientId = patientID
				revokeReq.RequesterId = requesterID
				ctx := auth.AddPatientMD(context.Background(), revokeReq.PatientId)
				revokeRes, err := PermissionAPI.RevokePermissionToken(ctx, revokeReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(revokeRes).ShouldNot(BeNil())
			})
			Describe("Lets try to get permission token", func() {
				It("should succeed getting the permission token but condition not permitted", func() {
					getReq := &permission.GetPermissionTokenRequest{
						PatientId: patientID,
						Actor:     newActorPayload(),
					}
					getReq.Actor.Group = int32(ledger.Actor_INSURANCE)
					ctx = auth.AddGroupAndIDMD(
						context.Background(), getReq.Actor.Group, getReq.Actor.Id,
					)

					getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())
					Expect(getRes.Allowed).Should(BeFalse())
					Expect(getRes.AccessToken).Should(BeZero())
				})
			})
		})
	})
})
