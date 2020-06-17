package patient

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/patient"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting patient medical data #getMedicalData", func() {
	var (
		getReq *patient.GetPatientMedDataRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &patient.GetPatientMedDataRequest{
			PatientId:   uuid.New().String(),
			AccessToken: uuid.New().String(),
			IsOwner:     true,
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	Describe("Getting patient medical data with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is missing", func() {
			ctx := auth.AddPatientMD(context.Background(), "")
			getReq.PatientId = ""
			getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is different from id in context", func() {
			ctx := auth.AddPatientMD(context.Background(), uuid.New().String())
			getReq.PatientId = uuid.New().String()
			getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id doesn't exist", func() {
			ctx := auth.AddPatientMD(context.Background(), getReq.PatientId)
			getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(getRes).Should(BeNil())
		})
	})

	When("Patient is accessing their own medical data", func() {
		BeforeEach(func() {
			getReq.IsOwner = true
		})
		It("should fail when group is not patient", func() {
			getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		Context("Lets create patient medical data", func() {
			var patientID string
			It("should add patient medical data without error", func() {
				addReq := &patient.AddPatientMedDataRequest{
					PatientId:   uuid.New().String(),
					MedicalData: newMedicalData(),
					Actor:       newActorPayload(),
				}
				addReq.Actor.Actor = ledger.Actor_HOSPITAL
				ctx := auth.AddGroupAndIDMD(
					context.Background(), int32(addReq.Actor.Actor), addReq.Actor.ActorId,
				)
				addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(addRes).ShouldNot(BeNil())
				patientID = addReq.GetPatientId()
			})
			It("should succeed when ID in context and request are similar", func() {
				getReq.PatientId = patientID
				ctx := auth.AddPatientMD(context.Background(), patientID)
				getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
			})
		})
	})

	When("Different actor is accessing a patient medical history", func() {
		BeforeEach(func() {
			getReq.IsOwner = false
		})
		Context("Lets create patient medical data", func() {
			var patientID string
			It("should add patient medical data without error", func() {
				addReq := &patient.AddPatientMedDataRequest{
					PatientId:   uuid.New().String(),
					MedicalData: newMedicalData(),
					Actor:       newActorPayload(),
				}
				addReq.Actor.Actor = ledger.Actor_HOSPITAL
				ctx := auth.AddGroupAndIDMD(
					context.Background(), int32(addReq.Actor.Actor), addReq.Actor.ActorId,
				)
				addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(addRes).ShouldNot(BeNil())
				patientID = addReq.GetPatientId()
			})

			Describe("Getting the medical data for the patient", func() {
				It("should fail when access token is incorrect", func() {
					getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.Unauthenticated))
					Expect(getRes).Should(BeNil())
				})
				It("should fail when ID in access token is different than request ID", func() {
					getReq.PatientId = patientID
					getReq.AccessToken = auth.GenPatientAccessToken(ctx, uuid.New().String())
					getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
					Expect(getRes).Should(BeNil())
				})
				It("should succeed when ID in access token is similar with request ID", func() {
					getReq.PatientId = patientID
					getReq.AccessToken = auth.GenPatientAccessToken(ctx, patientID)
					getRes, err := PatientAPI.GetPatientMedData(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())
				})
			})
		})
	})
})
