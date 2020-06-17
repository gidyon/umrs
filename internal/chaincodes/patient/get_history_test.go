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

func newFilter() *ledger.Filter {
	return &ledger.Filter{
		Filter: true,
	}
}

var _ = Describe("Getting patient treatment data #getHistory", func() {
	var (
		getReq *patient.GetMedicalHistoryRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &patient.GetMedicalHistoryRequest{
			PatientId:  uuid.New().String(),
			PageNumber: 2,
			PageSize:   40,
			Filter:     newFilter(),
		}
		ctx = auth.AddInsuranceMD(context.Background())
	})

	Describe("Getting patient treatment history with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is missing", func() {
			ctx := auth.AddPatientMD(context.Background(), "")
			getReq.PatientId = ""
			getReq.IsOwner = true
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when patient id is different from id in context", func() {
			ctx := auth.AddPatientMD(context.Background(), uuid.New().String())
			getReq.PatientId = uuid.New().String()
			getReq.IsOwner = true
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
	})

	When("Patient is viewing their own medical history", func() {
		var patientID string
		BeforeEach(func() {
			getReq.IsOwner = true
			patientID = uuid.New().String()
		})
		It("should fail when group is not patient", func() {
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should succeed when ID in context and request are similar", func() {
			getReq.PatientId = patientID
			ctx := auth.AddPatientMD(context.Background(), patientID)
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
		})
	})

	When("Different actor is accessing patient medical history", func() {
		var patientID string
		BeforeEach(func() {
			getReq.IsOwner = false
			patientID = uuid.New().String()
		})
		It("should fail when access token is incorrect", func() {
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.Unauthenticated))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when ID in access token is different than request ID", func() {
			getReq.PatientId = patientID
			getReq.AccessToken = auth.GenPatientAccessToken(ctx, uuid.New().String())
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(getRes).Should(BeNil())
		})
		It("should succeed when ID in access token is similar with request ID", func() {
			getReq.PatientId = patientID
			getReq.AccessToken = auth.GenPatientAccessToken(ctx, patientID)
			getRes, err := PatientAPI.GetMedicalHistory(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
		})
	})
})
