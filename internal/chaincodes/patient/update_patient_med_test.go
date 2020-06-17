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

var _ = Describe("Updating patient medical data #update", func() {
	var (
		updateReq *patient.AddPatientMedDataRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &patient.AddPatientMedDataRequest{
			PatientId:   uuid.New().String(),
			MedicalData: newMedicalData(),
			Actor:       newActorPayload(),
		}
		updateReq.Actor.Actor = ledger.Actor_HOSPITAL
		ctx = auth.AddGroupAndIDMD(
			context.Background(), int32(ledger.Actor_HOSPITAL), updateReq.Actor.ActorId,
		)
	})

	Describe("Updating patient with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the patient medical data nil", func() {
			updateReq.MedicalData = nil
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the patient id is missing", func() {
			updateReq.PatientId = ""
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the hospital id is missing", func() {
			updateReq.GetMedicalData().HospitalId = ""
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the hospital name is missing", func() {
			updateReq.GetMedicalData().HospitalName = ""
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the actor id is missing", func() {
			updateReq.Actor.ActorId = ""
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the actor id nil", func() {
			updateReq.Actor = nil
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the actor name is missing", func() {
			updateReq.Actor.ActorNames = ""
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the actor is unknown", func() {
			updateReq.Actor.Actor = ledger.Actor_UNKNOWN
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when actor id is changed", func() {
			updateReq.Actor.ActorId = uuid.New().String()
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when actor is changed", func() {
			updateReq.Actor.Actor = ledger.Actor_GOVERNMENT
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when patient is not in database", func() {
			updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(updateRes).Should(BeNil())
		})
	})

	Describe("Updating patient with valid request", func() {
		Context("Lets create patient medical data", func() {
			var patientID string
			var actor ledger.Actor
			It("should add patient medical data without error", func() {
				updateReq := &patient.AddPatientMedDataRequest{
					PatientId:   uuid.New().String(),
					MedicalData: newMedicalData(),
					Actor:       newActorPayload(),
				}
				updateReq.Actor.Actor = ledger.Actor_HOSPITAL
				ctx = auth.AddGroupAndIDMD(
					context.Background(), int32(updateReq.Actor.Actor), updateReq.Actor.ActorId,
				)
				updateRes, err := PatientAPI.AddPatientMedData(ctx, updateReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(updateRes).ShouldNot(BeNil())
				patientID = updateReq.GetPatientId()
				actor = updateReq.GetActor().GetActor()
			})
			Describe("Lets update the medical data", func() {
				It("should update patient medical data without error", func() {
					updateReq.Actor.Actor = actor
					updateReq.PatientId = patientID
					updateRes, err := PatientAPI.UpdatePatientMedData(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(updateRes).ShouldNot(BeNil())
				})
			})
		})
	})
})
