package patient

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/internal/pkg/md"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/patient"
	"github.com/gidyon/umrs/pkg/api/treatment"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
)

const (
	failedToAddMedicalData    = "Failed to add medical data"
	failedToUpdateMedData     = "Failed to update medical data"
	failedToAddTreatmentData  = "Failed to add treatment data"
	failedToGetMedicalData    = "Failed to get medical data"
	failedToGetTreatmentData  = "Failed to get treatment data"
	failedToGetMedicalHistory = "Failed to get medical history"
)

type patientAPIServer struct {
	contractID       string
	ctx              context.Context
	sqlDB            *gorm.DB
	ledgerClient ledger.ledgerClient
}

// Options contains parameters required by NewPatientChaincode
type Options struct {
	ContractID       string
	SQLDB            *gorm.DB
	ledgerClient ledger.ledgerClient
}

// NewPatientChaincode creates an instance of patients API
func NewPatientChaincode(ctx context.Context, opt *Options) (patient.PatientAPIServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case strings.TrimSpace(opt.ContractID) == "":
		return nil, errs.MissingCredential("ContractId")
	case opt.SQLDB == nil:
		return nil, errs.NilObject("SqlDB")
	case opt.ledgerClient == nil:
		return nil, errs.NilObject("ledgerClient")
	}

	pas := &patientAPIServer{
		ctx:              ctx,
		sqlDB:            opt.SQLDB,
		ledgerClient: opt.ledgerClient,
	}

	// Auto migration
	err := pas.sqlDB.AutoMigrate(&patientMedicalData{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate patients medical records table: %w", err)
	}

	// Register the chaincode to ledger server
	ctxReg := auth.AddSuperAdminMD(ctx)
	p, err := auth.AuthenticateSuperAdmin(ctxReg)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate super admin: %v", err)
	}

	regRes, err := pas.ledgerClient.RegisterContract(ctxReg, &ledger.RegisterContractRequest{
		SuperAdminId: p.ID, ContractId: opt.ContractID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register contract with the ledger server: %v", err)
	}

	pas.contractID = regRes.ContractId

	return pas, nil
}

func (patientSrv *patientAPIServer) AddPatientMedData(
	ctx context.Context, addReq *patient.AddPatientMedDataRequest,
) (*patient.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToAddMedicalData)

	// Request must not be nil
	if addReq == nil {
		return nil, WrapError(errs.NilObject("AddPatientMedDataRequest"))
	}

	// Validation
	actor := addReq.GetActor()
	patientID := addReq.GetPatientId()
	medDataPB := addReq.GetMedicalData()
	var err error
	switch {
	case actor == nil:
		err = errs.NilObject("Actor")

	case strings.TrimSpace(actor.ActorId) == "":
		err = errs.MissingCredential("ActorId")

	case strings.TrimSpace(actor.ActorNames) == "":
		err = errs.MissingCredential("ActorNames")

	case actor.Actor == ledger.Actor_UNKNOWN:
		err = errs.ActorUknown()

	case strings.TrimSpace(patientID) == "":
		err = errs.MissingCredential("PatientID")

	case medDataPB == nil:
		err = errs.NilObject("MedicalData")

	case medDataPB.GetDetails() == nil:
		err = errs.NilObject("Details")

	case strings.TrimSpace(medDataPB.GetHospitalId()) == "":
		err = errs.MissingCredential("HospitalID")

	case strings.TrimSpace(medDataPB.GetHospitalName()) == "":
		err = errs.MissingCredential("HospitalName")

	case strings.TrimSpace(medDataPB.GetPatientName()) == "":
		err = errs.MissingCredential("PatientName")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndID(ctx, int32(ledger.Actor_HOSPITAL), actor.ActorId)
	if err != nil {
		return nil, WrapError(err)
	}

	// Add to database
	medDataDB, err := getPatientMedDataDB(medDataPB)
	if err != nil {
		return nil, WrapError(err)
	}
	medDataDB.PatientID = patientID

	tx := patientSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	err = tx.Save(medDataDB).Error
	if err != nil {
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "SAVE"))
	}

	medDataBs, err := proto.Marshal(medDataPB.GetDetails())
	if err != nil {
		return nil, WrapError(err)
	}

	// Add to ledger
	addRes, err := patientSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_ADD_PATIENT_MEDICAL_DATA,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       actor.ActorId,
				ActorNames: actor.ActorNames,
			},
			Patient: &ledger.ActorPayload{
				ActorId:       patientID,
				ActorNames: medDataPB.PatientName,
			},
			Organization: &ledger.ActorPayload{
				ActorId:       medDataPB.HospitalId,
				ActorNames: medDataPB.HospitalName,
			},
			Details: medDataBs,
		},
	}, grpc.WaitForReady(true))
	if err != nil {
		tx.Rollback()
		return nil, WrapError(errs.FailedToAddToledger(err))
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errs.FailedToCommitTx(err)
	}

	return &patient.HashResponse{
		OperationHash: addRes.Hash,
		PatientId:     patientID,
	}, nil
}

func (patientSrv *patientAPIServer) UpdatePatientMedData(
	ctx context.Context, updateReq *patient.AddPatientMedDataRequest,
) (*patient.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToUpdateMedData)

	// Request must not be nil
	if updateReq == nil {
		return nil, WrapError(errs.NilObject("UpdatePatientMedDataRequest"))
	}

	// Validation
	actor := updateReq.GetActor()
	patientID := updateReq.GetPatientId()
	medDataPB := updateReq.GetMedicalData()
	var err error
	switch {
	case actor == nil:
		err = errs.NilObject("Actor")
	case actor.Actor == ledger.Actor_UNKNOWN:
		err = errs.ActorUknown()
	case strings.TrimSpace(actor.ActorId) == "":
		err = errs.MissingCredential("ActorID")
	case strings.TrimSpace(actor.ActorNames) == "":
		err = errs.MissingCredential("ActorNames")
	case strings.TrimSpace(patientID) == "":
		err = errs.MissingCredential("PatientID")
	case medDataPB == nil:
		err = errs.NilObject("MedicalData")
	case strings.TrimSpace(medDataPB.GetHospitalId()) == "":
		err = errs.MissingCredential("HospitalID")
	case strings.TrimSpace(medDataPB.GetHospitalName()) == "":
		err = errs.MissingCredential("HospitalName")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndID(
		ctx, int32(actor.GetActor()), actor.GetActorId(),
	)
	if err != nil {
		return nil, WrapError(err)
	}

	// Add to database
	medDataDB, err := getPatientMedDataDB(medDataPB)
	if err != nil {
		return nil, WrapError(err)
	}
	medDataDB.PatientID = patientID

	tx := patientSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	db := tx.Model(&patientMedicalData{}).Where("patient_id=?", patientID).
		Updates(medDataDB)
	if db.Error != nil {
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "UPDATE"))
	}
	if db.RowsAffected == 0 {
		tx.Rollback()
		errMsg := fmt.Sprintf("no record found for account with id %s", patientID)
		return nil, WrapError(errs.WrapMessage(codes.NotFound, errMsg))
	}

	medDataBs, err := proto.Marshal(medDataPB.GetDetails())
	if err != nil {
		return nil, WrapError(err)
	}

	// Add to ledger
	addRes, err := patientSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_UPDATE_PATIENT_MEDICAL_DATA,
			Creator: &ledger.ActorPayload{
				Actor:         actor.Actor,
				ActorId:       actor.ActorId,
				ActorNames: actor.ActorNames,
			},
			Patient: &ledger.ActorPayload{
				ActorId:       patientID,
				ActorNames: medDataPB.PatientName,
			},
			Organization: &ledger.ActorPayload{
				ActorId:       medDataPB.HospitalId,
				ActorNames: medDataPB.HospitalName,
			},
			Details: medDataBs,
		},
	}, grpc.WaitForReady(true))
	if err != nil {
		tx.Rollback()
		return nil, WrapError(errs.FailedToAddToledger(err))
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, WrapError(errs.FailedToCommitTx(err))
	}

	return &patient.HashResponse{
		PatientId:     patientID,
		OperationHash: addRes.Hash,
	}, nil
}

func (patientSrv *patientAPIServer) GetPatientMedData(
	ctx context.Context, getReq *patient.GetPatientMedDataRequest,
) (*patient.MedicalData, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToGetMedicalData)

	// Request must not be nil
	if getReq == nil {
		return nil, WrapError(errs.NilObject("GetPatientMedDataRequest"))
	}

	// Validation
	var err error
	switch {
	case strings.TrimSpace(getReq.PatientId) == "":
		err = errs.MissingCredential("PatientId")
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
			return nil, WrapError(errs.PermissionDenied("GetPatientMedData"))
		}
	} else {
		claims, err := auth.ParseToken(getReq.AccessToken)
		if err != nil {
			return nil, WrapError(err)
		}
		if claims.ID != getReq.PatientId {
			return nil, WrapError(errs.PermissionDenied("GetPatientMedData"))
		}
	}

	// Get from db
	medDataDB := &patientMedicalData{}
	err = patientSrv.sqlDB.First(medDataDB, "patient_id=?", getReq.PatientId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, WrapError(errs.NoMedRecordFoundForPatient())
	default:
		return nil, WrapError(errs.SQLQueryFailed(err, "SELECT"))
	}

	medDataPB, err := getPatientMedDataPB(medDataDB)
	if err != nil {
		return nil, WrapError(err)
	}

	return medDataPB, nil
}

func (patientSrv *patientAPIServer) GetMedicalHistory(
	ctx context.Context, getReq *patient.GetMedicalHistoryRequest,
) (*patient.MedicalHistory, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToGetMedicalHistory)

	// Request must not be nil
	if getReq == nil {
		return nil, WrapError(errs.NilObject("GetMedicalHistoryRequest"))
	}

	// Validation
	patientID := getReq.GetPatientId()
	switch {
	case strings.TrimSpace(patientID) == "":
		return nil, WrapError(errs.MissingCredential("PatientId"))
	}

	// Authentication
	if getReq.IsOwner {
		payload, err := auth.AuthenticateGroup(ctx, int32(ledger.Actor_PATIENT))
		if err != nil {
			return nil, WrapError(err)
		}
		if payload.ID != getReq.PatientId {
			return nil, WrapError(errs.PermissionDenied("GetPatientMedData"))
		}
	} else {
		claims, err := auth.ParseToken(getReq.AccessToken)
		if err != nil {
			return nil, WrapError(err)
		}
		if claims.ID != getReq.PatientId {
			return nil, WrapError(errs.PermissionDenied("GetPatientMedData"))
		}
	}

	// Filter to query the ledger
	if getReq.GetFilter() == nil {
		getReq.Filter = &ledger.Filter{}
	}
	getReq.GetFilter().Filter = true
	getReq.GetFilter().ByPatient = true
	getReq.GetFilter().PatientIds = []string{patientID}
	getReq.GetFilter().ByOperation = true
	getReq.GetFilter().Operations = []ledger.Operation{
		ledger.Operation_ADD_PATIENT_TREATMENT_RECORD,
		ledger.Operation_ADD_PATIENT_MEDICAL_DATA,
		ledger.Operation_UPDATE_PATIENT_MEDICAL_DATA,
		ledger.Operation_UPDATE_PATIENT_TREATMENT_RECORD,
		ledger.Operation_DELETE_PATIENT_MEDICAL_DATA,
		ledger.Operation_DELETE_PATIENT_MEDICAL_RECORD,
		ledger.Operation_GRANT_PERMISSION,
		ledger.Operation_REPORT_ISSUE,
	}

	// Get history from ledger
	listRes, err := patientSrv.ledgerClient.ListBlocks(md.AddFromCtx(ctx), &ledger.ListBlocksRequest{
		PageNumber: getReq.GetPageNumber(),
		PageSize:   getReq.GetPageSize(),
		Filter:     getReq.GetFilter(),
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, WrapError(err)
	}

	medicalActivities := make([]*patient.MedicalActivity, 0)
	for _, blockPB := range listRes.GetBlocks() {
		medicalActivity, err := getMedicalActivityFromBlock(blockPB)
		if err != nil {
			return nil, WrapError(err)
		}
		medicalActivities = append(medicalActivities, medicalActivity)
	}

	return &patient.MedicalHistory{
		History:        medicalActivities,
		NextPageNumber: listRes.NextPageNumber,
	}, nil
}

func getTreatmentDataFromBlock(blockPB *ledger.Block) (*treatment.TreatmentData, error) {
	treatmentData := &treatment.TreatmentData{}
	details := blockPB.GetPayload().GetDetails()
	if details != nil {
		err := proto.Unmarshal(details, treatmentData)
		if err != nil {
			return nil, err
		}
	}
	return treatmentData, nil
}

func getMedicalDataFromBlock(blockPB *ledger.Block) (*patient.MedicalData, error) {
	data := &patient.Details{}
	blockTransaction := blockPB.GetPayload()
	details := blockTransaction.GetDetails()
	if details != nil {
		err := proto.Unmarshal(details, data)
		if err != nil {
			return nil, err
		}
	}
	medicalData := &patient.MedicalData{
		HospitalId:   blockTransaction.GetOrganization().GetActorId(),
		HospitalName: blockTransaction.GetOrganization().GetActorNames(),
		PatientName:  blockTransaction.GetPatient().GetActorNames(),
		PatientState: patient.State(patient.State_value[data.Details["patient_state"]]),
		Details:      data,
	}
	return medicalData, nil
}

func getMedicalActivityFromBlock(blockPB *ledger.Block) (*patient.MedicalActivity, error) {
	blockTransaction := blockPB.GetPayload()
	medicalActivity := &patient.MedicalActivity{
		BlockHash:    blockPB.GetHash(),
		Operation:    blockTransaction.GetOperation(),
		Patient:      blockTransaction.GetPatient(),
		Organization: blockTransaction.GetOrganization(),
		Creator:      blockTransaction.GetCreator(),
	}

	switch blockTransaction.GetOperation() {
	case ledger.Operation_ADD_PATIENT_MEDICAL_DATA:
		medData, err := getMedicalDataFromBlock(blockPB)
		if err != nil {
			return nil, err
		}
		medicalActivity.Payload = &patient.MedicalActivity_MedicalData{
			MedicalData: medData,
		}
	case ledger.Operation_ADD_PATIENT_TREATMENT_RECORD:
		treatmentPB, err := getTreatmentDataFromBlock(blockPB)
		if err != nil {
			return nil, err
		}
		medicalActivity.Payload = &patient.MedicalActivity_Treatment{
			Treatment: treatmentPB,
		}
	case ledger.Operation_UPDATE_PATIENT_MEDICAL_DATA:
		medData, err := getMedicalDataFromBlock(blockPB)
		if err != nil {
			return nil, err
		}
		medicalActivity.Payload = &patient.MedicalActivity_MedicalData{
			MedicalData: medData,
		}
	case ledger.Operation_UPDATE_PATIENT_TREATMENT_RECORD:
		treatmentPB, err := getTreatmentDataFromBlock(blockPB)
		if err != nil {
			return nil, err
		}
		medicalActivity.Payload = &patient.MedicalActivity_Treatment{
			Treatment: treatmentPB,
		}
	default:
		data := &patient.Details{}
		details := blockPB.GetPayload().GetDetails()
		if details != nil {
			err := proto.Unmarshal(details, data)
			if err != nil {
				return nil, err
			}
		}

		medicalActivity.Payload = &patient.MedicalActivity_OperationPayload{
			OperationPayload: &patient.OperationPayload{
				Details: data,
			},
		}
	}

	return medicalActivity, nil
}
