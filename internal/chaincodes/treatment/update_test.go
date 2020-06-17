package treatment

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/treatment"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Updating treatment data", func() {
	var (
		updateReq *treatment.UpdateTreatmentRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &treatment.UpdateTreatmentRequest{
			TreatmentHash: uuid.New().String(),
			UpdatedTreatment: &treatment.AddTreatmentRequest{
				Patient:       newActorPayload(),
				Uploader:      newActorPayload(),
				Hospital:      newHospital(),
				TreatmentInfo: newTreatmentData(),
			},
		}
		ctx = auth.AddGroupAndIDMD(
			context.Background(), auth.HospitalGroup, updateReq.UpdatedTreatment.Uploader.Id,
		)
	})

	Describe("Adding treatment with malformed request", func() {
		It("should fail when request is nil", func() {
			updateReq = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail if treatment hash is missing", func() {
			updateReq.TreatmentHash = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when patient in request is nil", func() {
			updateReq.UpdatedTreatment.Patient = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when patient id in request is missing", func() {
			updateReq.UpdatedTreatment.Patient.Id = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when patient name in request is missing", func() {
			updateReq.UpdatedTreatment.Patient.FullName = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when uploader in request is nil", func() {
			updateReq.UpdatedTreatment.Uploader = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when uploader id in request is missing", func() {
			updateReq.UpdatedTreatment.Uploader.Id = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when uploader name in request is missing", func() {
			updateReq.UpdatedTreatment.Uploader.FullName = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when uploader id is changed", func() {
			updateReq.UpdatedTreatment.Uploader.Id = uuid.New().String()
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when hospital in request is nil", func() {
			updateReq.UpdatedTreatment.Hospital = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when hospital id in request is nil", func() {
			updateReq.UpdatedTreatment.Hospital.HospitalId = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when hospital name in request is nil", func() {
			updateReq.UpdatedTreatment.Hospital.HospitalName = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when treatment information request is nil", func() {
			updateReq.UpdatedTreatment.TreatmentInfo = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when triage information request is nil", func() {
			updateReq.UpdatedTreatment.TreatmentInfo.TriageDetails = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when symptoms information request is nil", func() {
			updateReq.UpdatedTreatment.TreatmentInfo.Symptoms = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when prescriptions information request is nil", func() {
			updateReq.UpdatedTreatment.TreatmentInfo.Prescriptions = nil
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when observation request is missing", func() {
			updateReq.UpdatedTreatment.TreatmentInfo.Observations = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when observation request is missing", func() {
			updateReq.UpdatedTreatment.TreatmentInfo.Diagnosis = ""
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
	})

	Describe("Adding treatment with correct payload", func() {
		It("should succeed", func() {
			updateRes, err := TreatmentAPI.UpdateTreatment(ctx, updateReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(updateRes).ShouldNot(BeNil())
		})
	})
})
