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

var _ = Describe("Getting permission token #getToken", func() {
	var (
		getReq *permission.GetPermissionTokenRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &permission.GetPermissionTokenRequest{
			PatientId: uuid.New().String(),
			Actor:     newActorPayload(),
		}
		getReq.Actor.Group = int32(ledger.Actor_INSURANCE)
		ctx = auth.AddGroupAndIDMD(
			context.Background(), getReq.Actor.Group, getReq.Actor.Id,
		)
	})

	Describe("Getting permission token with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when actor is nil", func() {
			getReq.Actor = nil
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is missing", func() {
			getReq.PatientId = ""
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when requester id is missing", func() {
			getReq.Actor.Id = ""
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when requester id is different than one in token", func() {
			getReq.Actor.Id = uuid.New().String()
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when requester group and id in token is changed", func() {
			ctx := auth.AddGroupAndIDMD(ctx, auth.AdminGroup, getReq.Actor.Id)
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when requester group in request is changed", func() {
			getReq.Actor.Group = int32(ledger.Actor_GOVERNMENT)
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting permission token with correct request", func() {
		It("should succeed even when requester has not been granted permission", func() {
			getRes, err := PermissionAPI.GetPermissionToken(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
			Expect(getRes.Allowed).Should(BeFalse())
			Expect(getRes.AccessToken).Should(BeZero())
		})

		Describe("Getting permission token for allowed actor", func() {
			var (
				patientID   string
				requesterID string
				group       int32
			)
			Describe("Requesting for permission", func() {
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
					group = requestReq.Requester.Group
				})
			})

			Context("Lets grant access to medical data first", func() {
				It("should succeed", func() {
					grantReq := &permission.GrantPermissionTokenRequest{
						Requester:    newActorPayload(),
						Organization: newActorPayload(),
						Patient:      newActorPayload(),
					}
					grantReq.Patient.Group = int32(ledger.Actor_PATIENT)
					grantReq.Patient.Id = patientID
					grantReq.Requester.Id = requesterID
					grantReq.Requester.Group = group
					grantReq.AuthorizationToken = auth.GenPatientAccessToken(
						context.Background(), grantReq.Patient.Id,
					)
					ctx := context.Background()
					grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(grantRes).ShouldNot(BeNil())
				})
			})

			It("should succeed getting the permission token", func() {
				getReq.PatientId = patientID
				getReq.Actor.Group = group
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
