package treatment

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/treatment"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
)

func newActor() ledger.Actor {
	index := rand.Intn(len(ledger.Actor_name) - 1)
	if index == 0 {
		index = 1
	}
	return ledger.Actor(index)
}

func newActorPayload() *treatment.Actor {
	return &treatment.Actor{
		Group:    int32(newActor()),
		Id:       uuid.New().String(),
		FullName: randomdata.SillyName(),
	}
}

func newHospital() *treatment.Hospital {
	return &treatment.Hospital{
		HospitalId:   uuid.New().String(),
		HospitalName: randomdata.Street() + " Hospital",
	}
}

func newTreatmentData() *treatment.TreatmentData {
	return &treatment.TreatmentData{
		Date:      randomdata.FullDate(),
		Timestamp: time.Now().Unix(),
		EntryFee:  200.00,
		TriageDetails: map[string]string{
			"height":      "2 M",
			"temperature": "36.9 C",
		},
		Symptoms: []string{
			randomdata.Paragraph(), randomdata.Paragraph(), randomdata.Paragraph(),
		},
		Prescriptions: []*treatment.Prescription{
			{Drug: randomdata.SillyName() + "nine", Cost: float32(randomdata.Decimal(100, 1200))},
			{Drug: randomdata.SillyName() + "cillin", Cost: float32(randomdata.Decimal(100, 1200))},
		},
		Observations:       randomdata.Paragraph(),
		Diagnosis:          randomdata.Paragraph(),
		AdditionalComments: randomdata.Paragraph(),
	}
}

var _ = Describe("Adding treatment data", func() {
	var (
		addReq *treatment.AddTreatmentRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		addReq = &treatment.AddTreatmentRequest{
			Patient:       newActorPayload(),
			Uploader:      newActorPayload(),
			Hospital:      newHospital(),
			TreatmentInfo: newTreatmentData(),
		}
		ctx = auth.AddGroupAndIDMD(context.Background(), auth.HospitalGroup, addReq.Uploader.Id)
	})

	Describe("Adding treatment with malformed request", func() {
		It("should fail when request is nil", func() {
			addReq = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when patient in request is nil", func() {
			addReq.Patient = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when patient id in request is missing", func() {
			addReq.Patient.Id = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when patient name in request is missing", func() {
			addReq.Patient.FullName = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when uploader in request is nil", func() {
			addReq.Uploader = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when uploader id in request is missing", func() {
			addReq.Uploader.Id = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when uploader name in request is missing", func() {
			addReq.Uploader.FullName = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when uploader id is changed", func() {
			addReq.Uploader.Id = uuid.New().String()
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital in request is nil", func() {
			addReq.Hospital = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital id in request is nil", func() {
			addReq.Hospital.HospitalId = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when hospital name in request is nil", func() {
			addReq.Hospital.HospitalName = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when treatment information request is nil", func() {
			addReq.TreatmentInfo = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when triage information request is nil", func() {
			addReq.TreatmentInfo.TriageDetails = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when symptoms information request is nil", func() {
			addReq.TreatmentInfo.Symptoms = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when prescriptions information request is nil", func() {
			addReq.TreatmentInfo.Prescriptions = nil
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when observation request is missing", func() {
			addReq.TreatmentInfo.Observations = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
		It("should fail when observation request is missing", func() {
			addReq.TreatmentInfo.Diagnosis = ""
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(addRes).Should(BeNil())
		})
	})

	Describe("Adding treatment with correct payload", func() {
		It("should succeed", func() {
			addRes, err := TreatmentAPI.AddTreatment(ctx, addReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(addRes).ShouldNot(BeNil())
		})
	})
})
