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

var _ = Describe("Getting active permissions token #getActive", func() {
	var (
		getReq *permission.GetActivePermissionsRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &permission.GetActivePermissionsRequest{
			PatientId: uuid.New().String(),
		}
		ctx = auth.AddPatientMD(ctx, getReq.PatientId)
	})

	Describe("Getting permission token with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := PermissionAPI.GetActivePermissions(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is missing", func() {
			getReq.PatientId = ""
			getRes, err := PermissionAPI.GetActivePermissions(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when actor in patient is changed", func() {
			ctx := auth.AddInsuranceMD(context.Background())
			getRes, err := PermissionAPI.GetActivePermissions(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is different than one in token", func() {
			getReq.PatientId = uuid.New().String()
			getRes, err := PermissionAPI.GetActivePermissions(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when group is not patient group", func() {
			ctx := auth.AddGroupAndIDMD(ctx, auth.AdminGroup, getReq.PatientId)
			getRes, err := PermissionAPI.GetActivePermissions(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting active permission with correct request", func() {
		It("should succeed even when requester has not been granted permission", func() {
			getRes, err := PermissionAPI.GetActivePermissions(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
			Expect(len(getRes.ActiveProfiles)).Should(BeZero())
		})

		Describe("Getting active permission for patient who granted access to their data", func() {
			var (
				patientID   string
				requesterID string
				group       int32
			)
			Context("The requester should request for permission token", func() {
				It("should succeed if the request is valid", func() {
					requestReq := &permission.RequestPermissionTokenRequest{
						RequesterProfile: newRequesterProfile(),
						Patient:          newActorPayload(),
						Reason:           "Medical checkup",
						PermissionMethod: allowedPermissionMethod(),
						Requester:        newActorPayload(),
					}
					requestReq.Requester.Group = int32(ledger.Actor_INSURANCE)
					requestReq.Requester.Id = uuid.New().String()
					requestReq.Patient.Id = uuid.New().String()
					ctx := auth.AddGroupAndIDMD(
						context.Background(), requestReq.Requester.Group, requestReq.Requester.Id,
					)
					requestRes, err := PermissionAPI.RequestPermissionToken(ctx, requestReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(requestRes).ShouldNot(BeNil())
					patientID = requestReq.Patient.Id
					requesterID = requestReq.Requester.Id
					group = requestReq.Requester.Group
				})
			})

			Context("Lets grant access to medical data", func() {
				It("should succeed", func() {
					grantReq := &permission.GrantPermissionTokenRequest{
						Requester:    newActorPayload(),
						Organization: newActorPayload(),
						Patient:      newActorPayload(),
					}
					grantReq.Patient.Group = int32(ledger.Actor_PATIENT)
					grantReq.Patient.Id = patientID
					grantReq.Requester.Group = group
					grantReq.Requester.Id = requesterID
					grantReq.AuthorizationToken = auth.GenPatientAccessToken(
						context.Background(), patientID,
					)
					ctx := context.Background()
					grantRes, err := PermissionAPI.GrantPermissionToken(ctx, grantReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(grantRes).ShouldNot(BeNil())
				})
			})
			It("should have atleast one active requester", func() {
				getReq.PatientId = patientID
				ctx := auth.AddGroupAndIDMD(
					context.Background(), int32(ledger.Actor_PATIENT), patientID,
				)
				getRes, err := PermissionAPI.GetActivePermissions(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
				Expect(len(getRes.ActiveProfiles)).ShouldNot(BeZero())
			})
		})
	})
})
