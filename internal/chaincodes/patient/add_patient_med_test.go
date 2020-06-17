package patient

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/patient"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

var bloodGroups = []string{"O+", "O-", "B+", "B-", "AB"}
var skinColor = []string{"Pale", "White", "Brown", "Red", "Black", "Chocolate"}

func newMedicalData() *patient.MedicalData {
	return &patient.MedicalData{
		HospitalId:   uuid.New().String(),
		HospitalName: randomdata.State(randomdata.Large) + " hospital",
		PatientName:  randomdata.SillyName(),
		PatientState: patient.State_ALIVE,
		Details: &patient.Details{
			Details: map[string]string{
				"blood_group": bloodGroups[rand.Intn(len(bloodGroups)-1)],
				"skin_color":  skinColor[rand.Intn(len(skinColor)-1)],
			},
		},
	}
}

func newActor() ledger.Actor {
	index := rand.Intn(len(ledger.Actor_name) - 1)
	if index == 0 {
		index = 1
	}
	return ledger.Actor(index)
}

func newActorPayload() *ledger.ActorPayload {
	return &ledger.ActorPayload{
		Actor:         newActor(),
		ActorId:       uuid.New().String(),
		ActorNames: randomdata.SillyName(),
	}
}

var patientOps = []ledger.Operation{
	ledger.Operation_ADD_PATIENT_TREATMENT_RECORD,
	ledger.Operation_ADD_PATIENT_MEDICAL_DATA,
	ledger.Operation_UPDATE_PATIENT_MEDICAL_DATA,
	ledger.Operation_UPDATE_PATIENT_TREATMENT_RECORD,
	ledger.Operation_DELETE_PATIENT_MEDICAL_DATA,
}

func newOperation() ledger.Operation {
	return patientOps[rand.Intn(len(patientOps)-1)]
}

var _ = Describe("Adding patient medical data #add", func() {
	var (
		addReq *patient.AddPatientMedDataRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		addReq = &patient.AddPatientMedDataRequest{
			PatientId:   uuid.New().String(),
			MedicalData: newMedicalData(),
			Actor:       newActorPayload(),
		}
		addReq.Actor.Actor = ledger.Actor_HOSPITAL
		ctx = auth.AddGroupAndIDMD(
			context.Background(), int32(ledger.Actor_HOSPITAL), addReq.Actor.ActorId,
		)
	})

	Describe("Adding patient with malformed request", func() {
		It("should fail when the request is nil", func() {
			addReq = nil
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the patient medical data is nil", func() {
			addReq.MedicalData = nil
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the patient details is nil", func() {
			addReq.MedicalData.Details = nil
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the patient id is missing", func() {
			addReq.PatientId = ""
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the hospital id is missing", func() {
			addReq.GetMedicalData().HospitalId = ""
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the hospital name is missing", func() {
			addReq.GetMedicalData().HospitalName = ""
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the patient name is missing", func() {
			addReq.GetMedicalData().PatientName = ""
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the actor is nil", func() {
			addReq.Actor = nil
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the actor id is missing", func() {
			addReq.Actor.ActorId = ""
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the actor name is missing", func() {
			addReq.Actor.ActorNames = ""
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when the actor is unknown", func() {
			addReq.Actor.Actor = ledger.Actor_UNKNOWN
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when actor id is changed", func() {
			addReq.Actor.ActorId = uuid.New().String()
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(addRes).Should(BeNil())
		})
	})

	Describe("Adding patient with valid request", func() {
		It("should add patient medical data without error", func() {
			addReq.Actor.Actor = ledger.Actor_HOSPITAL
			addRes, err := PatientAPI.AddPatientMedData(ctx, addReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(addRes).ShouldNot(BeNil())
		})
	})
})
