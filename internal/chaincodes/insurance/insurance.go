package insurance

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/db"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/internal/pkg/md"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"html/template"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	failedToAdd               = "Failed to add insurance"
	failedToGet               = "Failed to get insurance"
	failedToDelete            = "Failed to delete insurance"
	failedToList              = "Failed to list insurances"
	failedToSearch            = "Failed to search insurances"
	failedToUpdate            = "Failed to update insurance"
	failedToCheck             = "Failed to check suspension status"
	updatePermissionsDuration = time.Duration(12 * time.Hour)
)

type insuranceAPIServer struct {
	ctx                context.Context
	sqlDB              *gorm.DB
	ledgerClient   ledger.ledgerClient
	notificationClient notification.NotificationServiceClient
	tplAddInsurance    *template.Template
	tplUpdateInsurance *template.Template
	tplDeleteInsurance *template.Template
	mu                 *sync.RWMutex // guards allowedInsurances
	allowedInsurances  map[string]bool
}

// Options contains parameters for NewInsuranceAPIServer function
type Options struct {
	ContractID         string
	SQLDB              *gorm.DB
	ledgerClient   ledger.ledgerClient
	NotificationClient notification.NotificationServiceClient
}

// NewInsuranceAPIServer creates a single insurance singleton
func NewInsuranceAPIServer(ctx context.Context, opt *Options) (insurance.InsuranceAPIServer, error) {
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
	case opt.NotificationClient == nil:
		return nil, errs.NilObject("NotificationClient")
	}

	insuranceSrv := &insuranceAPIServer{
		ctx:                ctx,
		sqlDB:              opt.SQLDB,
		ledgerClient:   opt.ledgerClient,
		notificationClient: opt.NotificationClient,
		mu:                 &sync.RWMutex{},
		allowedInsurances:  make(map[string]bool, 0),
	}

	var err error
	// AddInsurance template file
	insuranceSrv.tplAddInsurance, err = template.ParseFiles(
		os.Getenv("ADD_INSURANCE_TEMPLATE_FILE"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	// UpdateInsurance template file
	insuranceSrv.tplUpdateInsurance, err = template.ParseFiles(
		os.Getenv("UPDATE_INSURANCE_TEMPLATE_FILE"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	// DeleteInsurance template file
	insuranceSrv.tplDeleteInsurance, err = template.ParseFiles(
		os.Getenv("DELETE_INSURANCE_TEMPLATE_FILE"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Auto-migration
	err = insuranceSrv.sqlDB.AutoMigrate(&insuranceModel{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate table: %w", err)
	}

	// Create full-text index
	err = db.CreateFullTextIndex(insuranceSrv.sqlDB, tableName, "insurance_name")
	if err != nil {
		return nil, fmt.Errorf("failed to create full-text index: %v", err)
	}

	// Register the chaincode to ledger server
	ctxReg := auth.AddSuperAdminMD(ctx)
	p, err := auth.AuthenticateSuperAdmin(ctxReg)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate super admin: %v", err)
	}

	_, err = insuranceSrv.ledgerClient.RegisterContract(ctxReg, &ledger.RegisterContractRequest{
		SuperAdminId: p.ID, ContractId: opt.ContractID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register contract with the ledger server: %v", err)
	}

	// Run permission worker
	go func() {
		insuranceSrv.updatePermissions()
		for {
			select {
			case <-time.After(updatePermissionsDuration):
				insuranceSrv.updatePermissions()
			}
		}
	}()

	return insuranceSrv, nil
}

type emailData struct {
	InsuranceID   string
	InsuranceName string
	Content       string
}

func (insuranceSrv *insuranceAPIServer) AddInsurance(
	ctx context.Context, addReq *insurance.AddInsuranceRequest,
) (*insurance.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToAdd)

	// Request must not be nil
	if addReq == nil {
		return nil, WrapError(errs.NilObject("AddInsuranceRequest"))
	}

	// Authentication
	adminPayload, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	insurancePB := addReq.GetInsurance()
	switch {
	case insurancePB == nil:
		err = errs.NilObject("Insurance")
	case strings.TrimSpace(insurancePB.InsuranceName) == "":
		err = errs.MissingCredential("Insurance Name")
	case strings.TrimSpace(insurancePB.WebsiteUrl) == "":
		err = errs.MissingCredential("Insurance Website")
	case strings.TrimSpace(insurancePB.About) == "":
		err = errs.MissingCredential("Insurance About")
	case len(insurancePB.About) > 256:
		err = errs.IncorrectVal("Insurance About (>256")
	case strings.TrimSpace(insurancePB.SupportEmail) == "":
		err = errs.MissingCredential("Insurance Support Email")
	case strings.TrimSpace(insurancePB.SupportTelNumber) == "":
		err = errs.MissingCredential("Insurance Support Number")
	case len(insurancePB.AdminEmails) == 0:
		err = errs.MissingCredential("InsuranceAdmins")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	insuranceDB, err := getInsuranceDB(insurancePB)
	if err != nil {
		return nil, WrapError(err)
	}

	// Start a transaction
	tx := insuranceSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	insuranceDB.Permission = true

	// Save to database
	err = tx.Save(insuranceDB).Error
	switch {
	case err == nil:
	default:
		tx.Rollback()
		// Checks whether the error was due to similar insurance existing
		if db.IsDuplicate(err) {
			return nil, WrapError(
				errs.WrapMessage(
					codes.ResourceExhausted,
					fmt.Sprintf("Insurance with name %s exists", insurancePB.InsuranceName),
				),
			)
		}
		return nil, WrapError(err)
	}

	// Marshal insurance content
	insuranceReqBs, err := proto.Marshal(insurancePB)
	if err != nil {
		return nil, WrapError(err)
	}

	addRes, err := insuranceSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_ADD_INSURANCE,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_ADMIN,
				ActorId:       adminPayload.ID,
				ActorNames: fmt.Sprintf("%s %s", adminPayload.FirstName, adminPayload.LastName),
			},
			Patient: &ledger.ActorPayload{},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_INSURANCE,
				ActorId:       insuranceDB.InsuranceID,
				ActorNames: insuranceDB.InsuranceName,
			},
			Details: insuranceReqBs,
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

	// Update map
	insuranceSrv.mu.Lock()
	insuranceSrv.allowedInsurances[insuranceDB.InsuranceID] = true
	insuranceSrv.mu.Unlock()

	// ================================= SENDING NOTIFICATION ===================================
	subject := fmt.Sprintf("You've been made the insurance admin for %s", insuranceDB.InsuranceName)
	notificationContent := &notification.NotificationContent{
		Subject: subject,
		Data: fmt.Sprintf(
			"You have been made the insurance admin for %s in order to manage their activities in the network",
			insuranceDB.InsuranceName,
		),
	}

	// Template for sending email
	content := bytes.NewBuffer(make([]byte, 0, 64))
	err = insuranceSrv.tplAddInsurance.Execute(content, emailData{
		InsuranceID:   insuranceDB.InsuranceID,
		InsuranceName: insuranceDB.InsuranceName,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	emailNotification := &notification.EmailNotification{
		To:              insurancePB.AdminEmails,
		Subject:         subject,
		BodyContentType: "text/html",
		Body:            content.String(),
	}

	// Send notification with Backoff; retry 5 times in 10 seconds
	insuranceSrv.sendNotification(ctx, &notification.Notification{
		OwnerIds:      insurancePB.AdminEmails,
		Priority:      notification.Priority_MEDIUM,
		SendMethod:    notification.SendMethod_EMAIL,
		Content:       notificationContent,
		CreateTimeSec: time.Now().Unix(),
		Payload: &notification.Notification_EmailNotification{
			EmailNotification: emailNotification,
		},
		Save: true,
	})

	return &insurance.HashResponse{
		InsuranceId:   insuranceDB.InsuranceID,
		OperationHash: addRes.Hash,
	}, nil
}

func (insuranceSrv *insuranceAPIServer) UpdateInsurance(
	ctx context.Context, updateReq *insurance.UpdateInsuranceRequest,
) (*insurance.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToUpdate)

	// Request must not be nil
	if updateReq == nil {
		return nil, WrapError(errs.NilObject("UpdateInsuranceRequest"))
	}

	// Authentication
	adminPayload, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	insurancePB := updateReq.GetInsurance()
	switch {
	case strings.TrimSpace(updateReq.InsuranceId) == "":
		err = errs.MissingCredential("InsuranceID")
	case insurancePB == nil && !updateReq.Suspend:
		err = errs.NilObject("Insurance")
	case updateReq.Suspend && strings.TrimSpace(updateReq.Reason) == "":
		err = errs.MissingCredential("Reason")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	var insuranceDB *insuranceModel
	if !updateReq.Suspend {
		insuranceDB, err = getInsuranceDB(insurancePB)
		if err != nil {
			return nil, WrapError(err)
		}
	}

	// Start a transaction
	tx := insuranceSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	// Get insurance from db
	insuranceDBX := &insuranceModel{}
	err = tx.Find(insuranceDBX, "insurance_id=?", updateReq.InsuranceId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		tx.Rollback()
		return nil, WrapError(errs.InsuranceNotFound(updateReq.InsuranceId))
	default:
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "FIND"))
	}

	insurancePBX, err := getInsurancePB(insuranceDBX)
	if err != nil {
		tx.Rollback()
		return nil, WrapError(err)
	}

	// Update model
	if updateReq.Suspend {
		err = tx.Table(tableName).Unscoped().Where("insurance_id=?", updateReq.InsuranceId).Update("permission", false).Error
	} else {
		err = tx.Table(tableName).Unscoped().Where("insurance_id=?", updateReq.InsuranceId).Updates(insuranceDB).Error
	}
	switch {
	case err == nil:
	default:
		tx.Rollback()
		// Checks whether the error was due to similar insurance existing
		if db.IsDuplicate(err) {
			return nil, WrapError(
				errs.WrapMessage(
					codes.ResourceExhausted,
					fmt.Sprintf("Insurance with name %s exists", insurancePBX.InsuranceName),
				),
			)
		}
		return nil, WrapError(errs.SQLQueryFailed(err, "SAVE"))
	}

	// Marshal insurance
	insuranceReqBs, err := proto.Marshal(updateReq)
	if err != nil {
		return nil, WrapError(err)
	}

	// Add event log to ledger
	addRes, err := insuranceSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_UPDATE_INSURANCE,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_ADMIN,
				ActorId:       adminPayload.ID,
				ActorNames: fmt.Sprintf("%s %s", adminPayload.FirstName, adminPayload.LastName),
			},
			Patient: &ledger.ActorPayload{},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_INSURANCE,
				ActorId:       insuranceDBX.InsuranceID,
				ActorNames: insuranceDBX.InsuranceName,
			},
			Details: insuranceReqBs,
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

	// ================================= SENDING NOTIFICATION ===================================
	subject := func() string {
		if updateReq.Suspend {
			return fmt.Sprintf("%s details was suspended", insurancePBX.InsuranceName)
		}
		return fmt.Sprintf("%s details was updated", insurancePBX.InsuranceName)
	}()
	notificationContent := &notification.NotificationContent{
		Subject: subject,
		Data: func() string {
			if updateReq.Suspend {
				return fmt.Sprintf(
					"%s details was suspended by %s %s for %s",
					insurancePBX.InsuranceName,
					adminPayload.FirstName,
					adminPayload.LastName,
					updateReq.Reason,
				)
			}
			return fmt.Sprintf(
				"%s details was updated by %s %s",
				insurancePBX.InsuranceName, adminPayload.FirstName, adminPayload.LastName,
			)
		}(),
	}

	// Template for sending email
	content := bytes.NewBuffer(make([]byte, 0, 64))
	err = insuranceSrv.tplUpdateInsurance.Execute(content, &emailData{
		InsuranceID:   insurancePBX.InsuranceId,
		InsuranceName: insurancePBX.InsuranceName,
		Content: func() string {
			if updateReq.GetSuspend() {
				return updateReq.GetReason()
			}
			return "Insurance Updated"
		}(),
	})

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	emails := func() []string {
		if updateReq.GetSuspend() {
			return insurancePBX.AdminEmails
		}
		if len(updateReq.GetInsurance().GetAdminEmails()) > 0 {
			return updateReq.Insurance.AdminEmails
		}
		return insurancePBX.AdminEmails
	}()

	emailNotification := &notification.EmailNotification{
		To:              emails,
		Subject:         subject,
		BodyContentType: "text/html",
		Body:            content.String(),
	}

	// Send notification with Backoff; retry 5 times in 10 seconds
	insuranceSrv.sendNotification(ctx, &notification.Notification{
		OwnerIds:      emails,
		Priority:      notification.Priority_MEDIUM,
		SendMethod:    notification.SendMethod_EMAIL,
		Content:       notificationContent,
		CreateTimeSec: time.Now().Unix(),
		Payload: &notification.Notification_EmailNotification{
			EmailNotification: emailNotification,
		},
		Save: true,
	})

	// Update map
	insuranceSrv.mu.Lock()
	if updateReq.Suspend {
		insuranceSrv.allowedInsurances[updateReq.InsuranceId] = false
	} else {
		insuranceSrv.allowedInsurances[updateReq.InsuranceId] = insuranceDB.Permission
	}
	insuranceSrv.mu.Unlock()

	return &insurance.HashResponse{
		InsuranceId:   updateReq.InsuranceId,
		OperationHash: addRes.Hash,
	}, nil
}

func (insuranceSrv *insuranceAPIServer) GetInsurance(
	ctx context.Context, getReq *insurance.GetInsuranceRequest,
) (*insurance.Insurance, error) {
	FailedToGetWrapper := errs.WrapErrorWithMsgFunc(failedToGet)

	// Request must not be nil
	if getReq == nil {
		return nil, FailedToGetWrapper(errs.NilObject("GetInsuranceRequest"))
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, FailedToGetWrapper(err)
	}

	// Validation
	if strings.TrimSpace(getReq.InsuranceId) == "" {
		return nil, FailedToGetWrapper(errs.MissingCredential("InsuranceID"))
	}

	// Get from database
	insuranceDB := &insuranceModel{}
	err = insuranceSrv.sqlDB.First(insuranceDB, "insurance_id=?", getReq.InsuranceId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, FailedToGetWrapper(errs.InsuranceNotFound(getReq.InsuranceId))
	default:
		return nil, FailedToGetWrapper(errs.WrapError(err))
	}

	insurancePB, err := getInsurancePB(insuranceDB)
	if err != nil {
		return nil, FailedToGetWrapper(err)
	}

	return insurancePB, nil
}

func (insuranceSrv *insuranceAPIServer) DeleteInsurance(
	ctx context.Context, delReq *insurance.DeleteInsuranceRequest,
) (*empty.Empty, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToDelete)

	// Request must not be nil
	if delReq == nil {
		return nil, WrapError(errs.NilObject("DeleteInsuranceRequest"))
	}

	// Authentication
	adminPayload, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	switch {
	case strings.TrimSpace(delReq.InsuranceId) == "":
		err = errs.MissingCredential("Insurance ID")
	case strings.TrimSpace(delReq.Reason) == "":
		err = errs.MissingCredential("Reason")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Start a transaction
	tx := insuranceSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	// Get insurance from db
	insuranceDB := &insuranceModel{}
	err = tx.Find(insuranceDB, "insurance_id=?", delReq.InsuranceId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		tx.Rollback()
		return nil, WrapError(errs.InsuranceNotFound(delReq.InsuranceId))
	default:
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "FIND"))
	}

	insurancePB, err := getInsurancePB(insuranceDB)
	if err != nil {
		return nil, WrapError(err)
	}

	// Soft delete from database
	err = tx.Table(tableName).Delete(&insuranceModel{}, "insurance_id=?", delReq.InsuranceId).Error
	if err != nil {
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "DELETE"))
	}

	// Marshal request
	insuranceReqBs, err := proto.Marshal(delReq)
	if err != nil {
		return nil, WrapError(err)
	}

	_, err = insuranceSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_DELETE_INSURANCE,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_ADMIN,
				ActorId:       adminPayload.ID,
				ActorNames: fmt.Sprintf("%s %s", adminPayload.FirstName, adminPayload.LastName),
			},
			Patient: &ledger.ActorPayload{},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_INSURANCE,
				ActorId:       delReq.InsuranceId,
				ActorNames: insuranceDB.InsuranceName,
			},
			Details: insuranceReqBs,
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

	// ================================= SENDING NOTIFICATION ===================================
	subject := fmt.Sprintf("%s Has Been Removed From Network", insurancePB.InsuranceName)
	notificationContent := &notification.NotificationContent{
		Subject: subject,
		Data: fmt.Sprintf(
			"%s was removed from the consortium of insurances network for %s",
			insurancePB.InsuranceName, delReq.Reason,
		),
	}

	// Template for sending email
	content := bytes.NewBuffer(make([]byte, 0, 64))
	err = insuranceSrv.tplDeleteInsurance.Execute(content, emailData{
		InsuranceID:   delReq.InsuranceId,
		InsuranceName: insurancePB.InsuranceName,
		Content:       delReq.Reason,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	emailNotification := &notification.EmailNotification{
		To:              insurancePB.AdminEmails,
		Subject:         subject,
		BodyContentType: "text/html",
		Body:            content.String(),
	}

	// Send notification with Backoff; retry 5 times in 10 seconds
	insuranceSrv.sendNotification(ctx, &notification.Notification{
		OwnerIds:      insurancePB.AdminEmails,
		Priority:      notification.Priority_MEDIUM,
		SendMethod:    notification.SendMethod_EMAIL,
		Content:       notificationContent,
		CreateTimeSec: time.Now().Unix(),
		Payload: &notification.Notification_EmailNotification{
			EmailNotification: emailNotification,
		},
		Save: true,
	})

	// Delete from map
	insuranceSrv.mu.Lock()
	delete(insuranceSrv.allowedInsurances, delReq.InsuranceId)
	insuranceSrv.mu.Unlock()

	return &empty.Empty{}, nil
}

const defaultCollectionSize = 50

func normalizePage(pageToken, pageSize int32) (int, int) {
	if pageSize <= 0 {
		pageSize = 1
	}
	if pageToken <= 0 {
		pageToken = 1
	}
	if pageSize > 50 {
		pageSize = defaultCollectionSize
	}
	return int(pageToken), int(pageSize)
}

func (insuranceSrv *insuranceAPIServer) ListInsurances(
	ctx context.Context, listReq *insurance.ListInsurancesRequest,
) (*insurance.Insurances, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToList)

	// Request must not be nil
	if listReq == nil {
		return nil, WrapError(errs.NilObject("ListInsurancesRequest"))
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	pageNumber, pageSize := normalizePage(listReq.GetPageNumber(), listReq.GetPageSize())
	offset := pageNumber*pageSize - pageSize

	insuranceDBs := make([]*insuranceModel, 0, pageSize)
	err = insuranceSrv.sqlDB.Offset(offset).Limit(pageSize).Find(&insuranceDBs).Error
	switch {
	case err == nil:
	default:
		return nil, WrapError(errs.WrapErrorWithCode(codes.Internal, err))
	}

	insurancePBs := make([]*insurance.Insurance, 0, len(insuranceDBs))
	for _, insuranceDB := range insuranceDBs {
		insurancePB, err := getInsurancePB(insuranceDB)
		if err != nil {
			return nil, WrapError(err)
		}
		insurancePBs = append(insurancePBs, insurancePB)
	}

	return &insurance.Insurances{
		Insurances: insurancePBs,
	}, nil
}

func (insuranceSrv *insuranceAPIServer) SearchInsurances(
	ctx context.Context, searchReq *insurance.SearchInsurancesRequest,
) (*insurance.Insurances, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToSearch)

	// requets must not be nil
	if searchReq == nil {
		return nil, WrapError(errs.NilObject("SearchInsurancesRequest"))
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// List insurances in case search term is empty
	if searchReq.Query == "" {
		return insuranceSrv.ListInsurances(ctx, &insurance.ListInsurancesRequest{
			PageNumber: searchReq.GetPageNumber(),
			PageSize:   searchReq.GetPageSize(),
		})
	}

	pageNumber, pageSize := normalizePage(searchReq.GetPageNumber(), searchReq.GetPageSize())
	offset := pageNumber*pageSize - pageSize

	parsedQuery := db.ParseQuery(searchReq.Query, "insurance", "insurances")

	insuranceDBs := make([]*insuranceModel, 0, pageSize)
	err = insuranceSrv.sqlDB.Offset(offset).Limit(pageSize).
		Find(&insuranceDBs, "MATCH(insurance_name) AGAINST(? IN BOOLEAN MODE )", parsedQuery).
		Error
	switch {
	case err == nil:
	default:
		return nil, WrapError(errs.SQLQueryFailed(err, "LIST"))
	}

	insurancePBs := make([]*insurance.Insurance, 0, len(insuranceDBs))
	for _, insuranceDB := range insuranceDBs {
		insurancePB, err := getInsurancePB(insuranceDB)
		if err != nil {
			return nil, WrapError(err)
		}
		insurancePBs = append(insurancePBs, insurancePB)
	}

	return &insurance.Insurances{
		Insurances: insurancePBs,
	}, nil
}
func (insuranceSrv *insuranceAPIServer) CheckSuspension(
	ctx context.Context, checkReq *insurance.CheckSuspensionRequest,
) (*insurance.CheckSuspensionResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToUpdate)

	// Request must not be nil
	if checkReq == nil {
		return nil, WrapError(errs.NilObject("CheckSuspensionRequest"))
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	if strings.TrimSpace(checkReq.InsuranceId) == "" {
		return nil, WrapError(errs.MissingCredential("InsuranceId"))
	}

	// Read from map
	insuranceSrv.mu.RLock()
	allowed := insuranceSrv.allowedInsurances[checkReq.InsuranceId]
	insuranceSrv.mu.RUnlock()

	if allowed {
		return &insurance.CheckSuspensionResponse{
			Suspended: false,
		}, nil
	}

	return &insurance.CheckSuspensionResponse{
		Suspended: true,
	}, nil
}

func (insuranceSrv *insuranceAPIServer) updatePermissions() {
	limit := 1000
	insurances := make([]*insuranceModel, 0, limit)
	db := insuranceSrv.sqlDB.Table(tableName)
	fetchMore := true

	for fetchMore {
		db = db.Limit(limit).Find(&insurances)
		if db.Error != nil {
			errs.LogError("failed to fetch insurances: %v", db.Error)
			fetchMore = false
			break
		}
		if len(insurances) < limit || db.RowsAffected < int64(limit) {
			fetchMore = false
		}

		// Lock mutex
		insuranceSrv.mu.Lock()
		for _, insuranceDB := range insurances {
			if insuranceDB.Permission {
				insuranceSrv.allowedInsurances[insuranceDB.InsuranceID] = true
			}
		}
		// Release to other routines
		insuranceSrv.mu.Unlock()
	}

	errs.LogInfo("finished updating insurances permission")
}

func (insuranceSrv *insuranceAPIServer) sendNotification(
	ctx context.Context, notificationPB *notification.Notification,
) {
	ctxSend, cancel := context.WithTimeout(md.AddFromCtx(ctx), 10*time.Second)
	defer cancel()

	var err error
	for i := 0; i < 5; i++ {
		_, err = insuranceSrv.notificationClient.Send(
			ctxSend, notificationPB, grpc.WaitForReady(true),
		)
		// We wait for few seconds
		if err != nil {
			errs.LogError("failed to Send notification: %v", err)
			<-time.After(time.Duration(i) * time.Second)
			continue
		}
		break
	}
}
