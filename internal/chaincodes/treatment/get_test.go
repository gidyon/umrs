package treatment

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/treatment"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Getting treatment treatment data #getTreatment", func() {
	var (
		getReq *treatment.GetTreatmentRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &treatment.GetTreatmentRequest{
			TreatmentHash: uuid.New().String(),
			PatientId:     uuid.New().String(),
			IsOwner:       true,
			AccessToken:   uuid.New().String(),
		}
		ctx = auth.AddPatientMD(context.Background(), uuid.New().String())
	})

	Describe("Getting treatment treatment data with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is missing", func() {
			ctx := auth.AddPatientMD(context.Background(), uuid.New().String())
			getReq.PatientId = ""
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when treatment hash is missing", func() {
			ctx := auth.AddPatientMD(context.Background(), uuid.New().String())
			getReq.TreatmentHash = ""
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is different from id in context", func() {
			ctx := auth.AddPatientMD(context.Background(), uuid.New().String())
			getReq.PatientId = uuid.New().String()
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
	})

	When("Patient is viewing their own treatment", func() {
		var patientID string
		BeforeEach(func() {
			getReq.IsOwner = true
			patientID = uuid.New().String()
		})
		It("should fail when group is not patient", func() {
			ctx := auth.AddHospitalMD(context.Background())
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should succeed when ID in context and request are similar", func() {
			getReq.PatientId = patientID
			ctx := auth.AddPatientMD(context.Background(), patientID)
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
		})
	})

	When("Different actor is accessing treatment medical history", func() {
		var patientID string
		BeforeEach(func() {
			getReq.IsOwner = false
			patientID = uuid.New().String()
		})
		It("should fail when access token is incorrect", func() {
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.Unauthenticated))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when ID in access token is different than patient ID", func() {
			getReq.PatientId = patientID
			getReq.AccessToken = auth.GenPatientAccessToken(ctx, uuid.New().String())
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should succeed when ID in access token is similar with request ID", func() {
			getReq.PatientId = patientID
			getReq.AccessToken = auth.GenPatientAccessToken(ctx, patientID)
			getRes, err := TreatmentAPI.GetTreatment(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
		})
	})
})
