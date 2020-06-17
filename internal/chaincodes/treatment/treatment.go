package treatment

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/internal/pkg/md"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/treatment"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"strings"
)

type treatmentAPIServer struct {
	ledgerClient ledger.ledgerClient
}

// Options contains paramaters for NewTreatmentChaincode
type Options struct {
	ContractID       string
	ledgerClient ledger.ledgerClient
}

// NewTreatmentChaincode creates a new treatment API singleton
func NewTreatmentChaincode(
	ctx context.Context, opt *Options,
) (treatment.TreatmentAPIServer, error) {
	// Validation
	var err error
	switch {
	case ctx == nil:
		err = errs.NilObject("Context")
	case opt == nil:
		err = errs.NilObject("Options")
	case opt.ledgerClient == nil:
		err = errs.NilObject("ledgerAPI")
	case strings.TrimSpace(opt.ContractID) == "":
		err = errs.MissingCredential("ContractId")
	}
	if err != nil {
		return nil, err
	}

	treatmentAPI := &treatmentAPIServer{
		ledgerClient: opt.ledgerClient,
	}

	// Register the chaincode to ledger server
	ctxReg := auth.AddSuperAdminMD(ctx)
	p, err := auth.AuthenticateSuperAdmin(ctxReg)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate super admin: %v", err)
	}

	_, err = treatmentAPI.ledgerClient.RegisterContract(ctxReg, &ledger.RegisterContractRequest{
		SuperAdminId: p.ID, ContractId: opt.ContractID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register contract with the ledger server: %v", err)
	}

	return treatmentAPI, nil
}

func (treatmentAPI *treatmentAPIServer) AddTreatment(
	ctx context.Context, addReq *treatment.AddTreatmentRequest,
) (*treatment.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc("Failed to add treatment")

	// Request must not be nil
	if addReq == nil {
		return nil, WrapError(errs.NilObject("AddTreatmentRequest"))
	}

	// Validation
	uploader := addReq.GetUploader()
	patientInfo := addReq.GetPatient()
	hospitalInfo := addReq.GetHospital()
	treatmentInfo := addReq.GetTreatmentInfo()

	err := validateRequest(uploader, patientInfo, hospitalInfo, treatmentInfo)
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndID(ctx, auth.HospitalGroup, uploader.Id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Marshal treatmentInfo to bytes
	treatmentInfoBs, err := proto.Marshal(treatmentInfo)
	if err != nil {
		return nil, WrapError(err)
	}

	// Add treatment to ledger
	addRes, err := treatmentAPI.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_ADD_PATIENT_TREATMENT_RECORD,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       uploader.Id,
				ActorNames: uploader.FullName,
			},
			Patient: &ledger.ActorPayload{
				Actor:         ledger.Actor_PATIENT,
				ActorId:       patientInfo.Id,
				ActorNames: patientInfo.FullName,
			},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       hospitalInfo.HospitalId,
				ActorNames: hospitalInfo.HospitalName,
			},
			Details: treatmentInfoBs,
		},
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, WrapError(err)
	}

	return &treatment.HashResponse{
		OperationHash: addRes.GetHash(),
		PatientId:     patientInfo.Id,
	}, nil
}

func (treatmentAPI *treatmentAPIServer) UpdateTreatment(
	ctx context.Context, updateReq *treatment.UpdateTreatmentRequest,
) (*treatment.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc("Update treatment failed")

	// Request must not be nil
	if updateReq == nil {
		return nil, WrapError(errs.NilObject("UpdateTreatmentRequest"))
	}

	// Validation
	uploader := updateReq.GetUpdatedTreatment().GetUploader()
	patientInfo := updateReq.GetUpdatedTreatment().GetPatient()
	hospitalInfo := updateReq.GetUpdatedTreatment().GetHospital()
	treatmentInfo := updateReq.GetUpdatedTreatment().GetTreatmentInfo()

	if strings.TrimSpace(updateReq.TreatmentHash) == "" {
		return nil, errs.MissingCredential("TreatmentHash")
	}

	err := validateRequest(uploader, patientInfo, hospitalInfo, treatmentInfo)
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndID(ctx, auth.HospitalGroup, uploader.Id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Marshal treatmentInfo to bytes
	treatmentInfoBs, err := proto.Marshal(treatmentInfo)
	if err != nil {
		return nil, WrapError(err)
	}

	// Add treatment to ledger
	addRes, err := treatmentAPI.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_UPDATE_PATIENT_TREATMENT_RECORD,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       uploader.Id,
				ActorNames: uploader.FullName,
			},
			Patient: &ledger.ActorPayload{
				Actor:         ledger.Actor_PATIENT,
				ActorId:       patientInfo.Id,
				ActorNames: patientInfo.FullName,
			},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       hospitalInfo.HospitalId,
				ActorNames: hospitalInfo.HospitalName,
			},
			Details: treatmentInfoBs,
		},
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, WrapError(err)
	}

	return &treatment.HashResponse{
		OperationHash: addRes.GetHash(),
		PatientId:     patientInfo.Id,
	}, nil
}

func (treatmentAPI *treatmentAPIServer) GetTreatment(
	ctx context.Context, getReq *treatment.GetTreatmentRequest,
) (*treatment.GetTreatmentResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc("Failed to get treatment")

	// Request must not be nil
	if getReq == nil {
		return nil, WrapError(errs.NilObject("GetTreatmentRequest"))
	}

	// Validation
	var err error
	switch {
	case strings.TrimSpace(getReq.TreatmentHash) == "":
		err = errs.MissingCredential("TreatmentHash")
	case strings.TrimSpace(getReq.PatientId) == "":
		err = errs.MissingCredential("PatientId")
	case strings.TrimSpace(getReq.AccessToken) == "" && !getReq.IsOwner:
		err = errs.MissingCredential("AccessToken")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	if getReq.IsOwner {
		payload, err := auth.AuthenticateGroup(ctx, int32(ledger.Actor_PATIENT))
		if err != nil {
			return nil, WrapError(err)
		}
		if payload.ID != getReq.PatientId {
			return nil, WrapError(errs.PermissionDenied("GetTreatment"))
		}
	} else {
		claims, err := auth.ParseToken(getReq.AccessToken)
		if err != nil {
			return nil, WrapError(err)
		}
		if claims.ID != getReq.PatientId {
			return nil, WrapError(errs.PermissionDenied("GetTreatment"))
		}
	}

	// Get from ledger
	getRes, err := treatmentAPI.ledgerClient.GetBlock(md.AddFromCtx(ctx), &ledger.GetBlockRequest{
		Hash: getReq.TreatmentHash,
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, WrapError(err)
	}
	res, err := getTreatmentResponseFromBlock(getRes)
	if err != nil {
		return nil, WrapError(err)
	}

	return res, nil
}

func validateRequest(
	uploader, patientInfo *treatment.Actor,
	hospitalInfo *treatment.Hospital,
	treatmentInfo *treatment.TreatmentData,
) error {
	var err error

	switch {
	case patientInfo == nil:
		err = errs.NilObject("Patient")
	case strings.TrimSpace(patientInfo.Id) == "":
		err = errs.MissingCredential("PatientId")
	case strings.TrimSpace(patientInfo.FullName) == "":
		err = errs.MissingCredential("PatientFullName")
	case uploader == nil:
		err = errs.NilObject("Uploader")
	case strings.TrimSpace(uploader.Id) == "":
		err = errs.MissingCredential("UploaderId")
	case strings.TrimSpace(uploader.FullName) == "":
		err = errs.MissingCredential("UploaderFullName")
	case hospitalInfo == nil:
		err = errs.NilObject("Hospital")
	case strings.TrimSpace(hospitalInfo.HospitalId) == "":
		err = errs.MissingCredential("HospitalId")
	case strings.TrimSpace(hospitalInfo.HospitalName) == "":
		err = errs.MissingCredential("HospitalName")
	case treatmentInfo == nil:
		err = errs.NilObject("TreatmentInfo")
	case len(treatmentInfo.TriageDetails) == 0:
		err = errs.NilObject("TriageDetails")
	case len(treatmentInfo.Symptoms) == 0:
		err = errs.NilObject("Symptoms")
	case strings.TrimSpace(treatmentInfo.Observations) == "":
		err = errs.MissingCredential("Observations")
	case strings.TrimSpace(treatmentInfo.Diagnosis) == "":
		err = errs.MissingCredential("Diagnosis")
	case len(treatmentInfo.Prescriptions) == 0:
		err = errs.NilObject("Prescription")
	}

	return err
}

func getTreatmentResponseFromBlock(blockPB *ledger.Block) (*treatment.GetTreatmentResponse, error) {
	if blockPB.GetPayload() == nil {
		return nil, errs.NilObject("BlockTransaction")
	}

	treatmentRes := &treatment.GetTreatmentResponse{
		TreatmentHash: blockPB.GetHash(),
		Timestamp:     blockPB.GetTimestamp(),
		Patient: &treatment.Actor{
			Group:    int32(blockPB.GetPayload().GetPatient().GetActor()),
			Id:       blockPB.GetPayload().GetPatient().GetActorId(),
			FullName: blockPB.GetPayload().GetPatient().GetActorNames(),
		},
		Uploader: &treatment.Actor{
			Group:    int32(blockPB.GetPayload().GetCreator().GetActor()),
			Id:       blockPB.GetPayload().GetCreator().GetActorId(),
			FullName: blockPB.GetPayload().GetCreator().GetActorNames(),
		},
		Hospital: &treatment.Hospital{
			HospitalId:   blockPB.GetPayload().GetOrganization().GetActorId(),
			HospitalName: blockPB.GetPayload().GetOrganization().GetActorNames(),
		},
		TreatmentInfo: &treatment.TreatmentData{},
	}

	details := blockPB.GetPayload().GetDetails()
	if details != nil {
		err := proto.Unmarshal(blockPB.GetPayload().GetDetails(), treatmentRes.TreatmentInfo)
		if err != nil {
			return nil, err
		}
	}

	return treatmentRes, nil
}
