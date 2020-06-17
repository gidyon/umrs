package hospital

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/db"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/internal/pkg/md"
	templateutil "github.com/gidyon/umrs/internal/pkg/templateutil"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/hospital"
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
	failedToAdd               = "Failed to add hospital"
	failedToGet               = "Failed to get hospital"
	failedToDelete            = "Failed to delete hospital"
	failedToUpdate            = "Failed to update hospital"
	failedToList              = "Failed to list hospitals"
	failedToSearch            = "Failed to search hospitals"
	failedToCheck             = "Failed to check suspension status"
	updatePermissionsDuration = time.Duration(12 * time.Hour)
)

type hospitalAPIServer struct {
	ctx                context.Context
	sqlDB              *gorm.DB
	ledgerClient   ledger.ledgerClient
	notificationClient notification.NotificationServiceClient
	tpl                *template.Template
	mu                 *sync.RWMutex // guards allowedHospitals
	allowedHospitals   map[string]bool
}

// Options contains parameters for NewHospitalAPIServer function
type Options struct {
	ContractID         string
	SQLDB              *gorm.DB
	ledgerClient   ledger.ledgerClient
	NotificationClient notification.NotificationServiceClient
}

// NewHospitalAPIServer creates a singleton for hospitals API
func NewHospitalAPIServer(ctx context.Context, opt *Options) (hospital.HospitalAPIServer, error) {
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

	hospitalSrv := &hospitalAPIServer{
		ctx:                ctx,
		sqlDB:              opt.SQLDB,
		ledgerClient:   opt.ledgerClient,
		notificationClient: opt.NotificationClient,
		mu:                 &sync.RWMutex{},
		allowedHospitals:   make(map[string]bool, 0),
	}

	// Testing only
	hospitalSrv.ledgerClient = fakeledger()
	// hospitalSrv.notificationClient = fakeNotificationAPI()

	// Parse templates
	tfiles, err := templateutil.ReadFiles(os.Getenv(templateutil.TemplateDirsEnv))
	if err != nil {
		return nil, fmt.Errorf("failed to read files in dir: %w", err)
	}

	hospitalSrv.tpl, err = templateutil.ParseTemplate(tfiles...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Auto migration
	err = hospitalSrv.sqlDB.AutoMigrate(&hospitalModel{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate table: %w", err)
	}

	// Create full-text index
	err = db.CreateFullTextIndex(hospitalSrv.sqlDB, tableName, "hospital_name")
	if err != nil {
		return nil, fmt.Errorf("failed to create full-text index: %v", err)
	}

	// Register the chaincode to ledger server
	ctxReg := auth.AddSuperAdminMD(ctx)
	p, err := auth.AuthenticateSuperAdmin(ctxReg)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate super admin: %v", err)
	}

	_, err = hospitalSrv.ledgerClient.RegisterContract(ctxReg, &ledger.RegisterContractRequest{
		SuperAdminId: p.ID, ContractId: opt.ContractID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register contract with the ledger server: %v", err)
	}

	// Run permission worker
	go func() {
		hospitalSrv.updatePermissions()
		for {
			select {
			case <-time.After(updatePermissionsDuration):
				hospitalSrv.updatePermissions()
			}
		}
	}()

	return hospitalSrv, nil
}

type emailData struct {
	HospitalID   string
	HospitalName string
	Content      string
	*templateutil.EmailData
}

func addEmailDataDefaults(v *emailData) *emailData {
	v.EmailData = templateutil.DefaultEmailData()
	return v
}

const (
	addTemplate     = "add"
	updateTemplate  = "update"
	deleteTemplate  = "delete"
	suspendTemplate = "suspend"
)

func (hospitalSrv *hospitalAPIServer) AddHospital(
	ctx context.Context, addReq *hospital.AddHospitalRequest,
) (*hospital.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToAdd)

	// Request must not be nil
	if addReq == nil {
		return nil, WrapError(errs.NilObject("AddHospitalRequest"))
	}

	// Authentication
	adminPayload, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	hospitalPB := addReq.GetHospital()
	switch {
	case hospitalPB == nil:
		err = errs.NilObject("Hospital")
	case strings.TrimSpace(hospitalPB.HospitalName) == "":
		err = errs.MissingCredential("HospitalName")
	case strings.TrimSpace(hospitalPB.County) == "":
		err = errs.MissingCredential("HospitalCounty")
	case strings.TrimSpace(hospitalPB.SubCounty) == "":
		err = errs.MissingCredential("HospitalSubCounty")
	case len(hospitalPB.AdminEmails) == 0:
		err = errs.MissingCredential("HospitalAdmins")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	hospitalDB, err := getHospitalDB(hospitalPB)
	if err != nil {
		return nil, WrapError(err)
	}

	// Start a transaction
	tx := hospitalSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	hospitalDB.Permission = true

	// Save to database
	err = tx.Save(hospitalDB).Error
	switch {
	case err == nil:
	default:
		tx.Rollback()
		// Checks whether the error was due to similar hospital existing
		if db.IsDuplicate(err) {
			return nil, WrapError(
				errs.WrapMessage(
					codes.ResourceExhausted,
					fmt.Sprintf("Hospital with name %s exists", hospitalPB.HospitalName),
				),
			)
		}
		return nil, WrapError(errs.SQLQueryFailed(err, "SAVE"))
	}

	// Marshal hospital
	hospitalBs, err := proto.Marshal(hospitalPB)
	if err != nil {
		tx.Rollback()
		return nil, WrapError(err)
	}

	// Add event log to ledger
	addRes, err := hospitalSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_ADD_HOSPITAL,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_ADMIN,
				ActorId:       adminPayload.ID,
				ActorNames: fmt.Sprintf("%s %s", adminPayload.FirstName, adminPayload.LastName),
			},
			Patient: &ledger.ActorPayload{},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       hospitalDB.HospitalID,
				ActorNames: hospitalDB.HospitalName,
			},
			Details: hospitalBs,
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
	hospitalSrv.mu.Lock()
	hospitalSrv.allowedHospitals[hospitalDB.HospitalID] = true
	hospitalSrv.mu.Unlock()

	// ================================= SENDING NOTIFICATION ===================================
	subject := fmt.Sprintf("You've been made the hospital admin for %s", hospitalDB.HospitalName)
	notificationContent := &notification.NotificationContent{
		Subject: subject,
		Data: fmt.Sprintf(
			"You have been made a hospital admin for %s in order to manage their activities in the network",
			hospitalDB.HospitalName,
		),
	}

	// Template for sending email
	notificationBuffer := bytes.NewBuffer(make([]byte, 0, 64))
	err = hospitalSrv.tpl.ExecuteTemplate(notificationBuffer, addTemplate, addEmailDataDefaults(&emailData{
		HospitalID:   hospitalDB.HospitalID,
		HospitalName: hospitalDB.HospitalName,
	}))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	emailNotification := &notification.EmailNotification{
		To:              hospitalPB.AdminEmails,
		Subject:         subject,
		BodyContentType: "text/html",
		Body:            notificationBuffer.String(),
	}

	// Send notification with Backoff; retry 5 times in 10 seconds
	hospitalSrv.sendNotification(ctx, &notification.Notification{
		OwnerIds:      hospitalPB.AdminEmails,
		Priority:      notification.Priority_MEDIUM,
		SendMethod:    notification.SendMethod_EMAIL,
		Content:       notificationContent,
		CreateTimeSec: time.Now().Unix(),
		Payload: &notification.Notification_EmailNotification{
			EmailNotification: emailNotification,
		},
		Save: true,
	})

	return &hospital.HashResponse{
		HospitalId:    hospitalDB.HospitalID,
		OperationHash: addRes.Hash,
	}, nil
}

func (hospitalSrv *hospitalAPIServer) GetHospital(
	ctx context.Context, getReq *hospital.GetHospitalRequest,
) (*hospital.Hospital, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToGet)

	// Request must not be nil
	if getReq == nil {
		return nil, WrapError(errs.NilObject("GetHospitalRequest"))
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	if strings.TrimSpace(getReq.HospitalId) == "" {
		return nil, WrapError(errs.MissingCredential("HospitalID"))
	}

	// Get from database
	hospitalDB := &hospitalModel{}
	err = hospitalSrv.sqlDB.First(hospitalDB, "hospital_id=?", getReq.HospitalId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, WrapError(errs.HospitalNotFound(getReq.HospitalId))
	default:
		return nil, WrapError(errs.SQLQueryFailed(err, "SELECT"))
	}

	hospitalPB, err := getHospitalPB(hospitalDB)
	if err != nil {
		return nil, WrapError(err)
	}

	return hospitalPB, nil
}

func (hospitalSrv *hospitalAPIServer) DeleteHospital(
	ctx context.Context, delReq *hospital.DeleteHospitalRequest,
) (*empty.Empty, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToDelete)

	// Request must not be nil
	if delReq == nil {
		return nil, WrapError(errs.NilObject("DeleteHospitalRequest"))
	}

	// Authentication
	adminPayload, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	switch {
	case strings.TrimSpace(delReq.HospitalId) == "":
		err = errs.MissingCredential("Hospital ID")
	case strings.TrimSpace(delReq.Reason) == "":
		err = errs.MissingCredential("Reason")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Start a transaction
	tx := hospitalSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	// Get hospital from db
	hospitalDB := &hospitalModel{}
	err = tx.Find(hospitalDB, "hospital_id=?", delReq.HospitalId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		tx.Rollback()
		return nil, WrapError(errs.HospitalNotFound(delReq.HospitalId))
	default:
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "FIND"))
	}

	hospitalPB, err := getHospitalPB(hospitalDB)
	if err != nil {
		return nil, WrapError(err)
	}

	// Soft delete from database
	err = tx.Table(tableName).Delete(&hospitalModel{}, "hospital_id=?", delReq.HospitalId).Error
	if err != nil {
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "DELETE"))
	}

	// Marshal request
	delBs, err := proto.Marshal(delReq)
	if err != nil {
		return nil, WrapError(err)
	}

	_, err = hospitalSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_DELETE_HOSPITAL,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_ADMIN,
				ActorId:       adminPayload.ID,
				ActorNames: fmt.Sprintf("%s %s", adminPayload.FirstName, adminPayload.LastName),
			},
			Patient: &ledger.ActorPayload{},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       delReq.HospitalId,
				ActorNames: hospitalPB.HospitalName,
			},
			Details: delBs,
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
	subject := fmt.Sprintf("%s Has Been Removed From Network", hospitalPB.HospitalName)
	notificationContent := &notification.NotificationContent{
		Subject: subject,
		Data: fmt.Sprintf(
			"%s was removed from the consortium of hospitals network for %s",
			hospitalPB.HospitalName, delReq.Reason,
		),
	}

	// Template for sending email
	notificationBuffer := bytes.NewBuffer(make([]byte, 0, 64))
	err = hospitalSrv.tpl.ExecuteTemplate(notificationBuffer, deleteTemplate, addEmailDataDefaults(&emailData{
		HospitalID:   delReq.HospitalId,
		HospitalName: hospitalPB.HospitalName,
		Content:      delReq.Reason,
	}))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	emailNotification := &notification.EmailNotification{
		To:              hospitalPB.AdminEmails,
		Subject:         subject,
		BodyContentType: "text/html",
		Body:            notificationBuffer.String(),
	}

	// Send notification with Backoff; retry 5 times in 10 seconds
	hospitalSrv.sendNotification(ctx, &notification.Notification{
		OwnerIds:      hospitalPB.AdminEmails,
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
	hospitalSrv.mu.Lock()
	delete(hospitalSrv.allowedHospitals, delReq.HospitalId)
	hospitalSrv.mu.Unlock()

	return &empty.Empty{}, nil
}

func (hospitalSrv *hospitalAPIServer) UpdateHospital(
	ctx context.Context, updateReq *hospital.UpdateHospitalRequest,
) (*hospital.HashResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToUpdate)

	// Request must not be nil
	if updateReq == nil {
		return nil, WrapError(errs.NilObject("UpdateHospitalRequest"))
	}

	// Authentication
	adminPayload, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// Validation
	hospitalPB := updateReq.GetHospital()
	switch {
	case strings.TrimSpace(updateReq.HospitalId) == "":
		err = errs.MissingCredential("HospitalID")
	case hospitalPB == nil && !updateReq.Suspend:
		err = errs.NilObject("Hospital")
	case updateReq.Suspend && strings.TrimSpace(updateReq.Reason) == "":
		err = errs.MissingCredential("Reason")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	var hospitalDB *hospitalModel
	if !updateReq.Suspend {
		hospitalDB, err = getHospitalDB(hospitalPB)
		if err != nil {
			return nil, WrapError(err)
		}
	}

	// Start a transaction
	tx := hospitalSrv.sqlDB.Begin()
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		tx.Rollback()
		return nil, WrapError(errs.FailedToBeginTx(err))
	}

	// Get hospital from db
	hospitalDBX := &hospitalModel{}
	err = tx.Find(hospitalDBX, "hospital_id=?", updateReq.HospitalId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		tx.Rollback()
		return nil, WrapError(errs.HospitalNotFound(updateReq.HospitalId))
	default:
		tx.Rollback()
		return nil, WrapError(errs.SQLQueryFailed(err, "FIND"))
	}

	hospitalPBX, err := getHospitalPB(hospitalDBX)
	if err != nil {
		tx.Rollback()
		return nil, WrapError(err)
	}

	// Update model
	if updateReq.Suspend {
		err = tx.Table(tableName).Unscoped().Where("hospital_id=?", updateReq.HospitalId).Update("permission", false).Error
	} else {
		err = tx.Table(tableName).Unscoped().Where("hospital_id=?", updateReq.HospitalId).Updates(hospitalDB).Error
	}
	switch {
	case err == nil:
	default:
		tx.Rollback()
		// Checks whether the error was due to similar hospital existing
		if db.IsDuplicate(err) {
			return nil, WrapError(
				errs.WrapMessage(
					codes.ResourceExhausted,
					fmt.Sprintf("Hospital with name %s exists", hospitalPB.HospitalName),
				),
			)
		}
		return nil, WrapError(errs.SQLQueryFailed(err, "SAVE"))
	}

	// Marshal hospital
	updateReqBs, err := proto.Marshal(updateReq)
	if err != nil {
		tx.Rollback()
		return nil, WrapError(err)
	}

	// Add event log to ledger
	addRes, err := hospitalSrv.ledgerClient.AddBlock(md.AddFromCtx(ctx), &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_UPDATE_HOSPITAL,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor_ADMIN,
				ActorId:       adminPayload.ID,
				ActorNames: fmt.Sprintf("%s %s", adminPayload.FirstName, adminPayload.LastName),
			},
			Patient: &ledger.ActorPayload{},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor_HOSPITAL,
				ActorId:       hospitalPBX.HospitalId,
				ActorNames: hospitalPBX.HospitalName,
			},
			Details: updateReqBs,
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
			return fmt.Sprintf("%s was suspended", hospitalPBX.HospitalName)
		}
		return fmt.Sprintf("%s details was updated", hospitalPBX.HospitalName)
	}()
	notificationContent := &notification.NotificationContent{
		Subject: subject,
		Data: func() string {
			if updateReq.Suspend {
				return fmt.Sprintf(
					"%s was suspended by %s %s for %s",
					hospitalPBX.HospitalName,
					adminPayload.FirstName,
					adminPayload.LastName,
					updateReq.Reason,
				)
			}
			return fmt.Sprintf(
				"%s details was updated by %s %s",
				hospitalPBX.HospitalName, adminPayload.FirstName, adminPayload.LastName,
			)
		}(),
	}

	// Template for sending email
	notificationBuffer := bytes.NewBuffer(make([]byte, 0, 64))
	content := func() string {
		if updateReq.GetSuspend() {
			return updateReq.GetReason()
		}
		return "Hospital Updated"
	}()
	if updateReq.Suspend {
		err = hospitalSrv.tpl.ExecuteTemplate(notificationBuffer, suspendTemplate, addEmailDataDefaults(&emailData{
			HospitalID:   hospitalPBX.HospitalId,
			HospitalName: hospitalPBX.HospitalName,
			Content:      content,
		}))
	} else {
		err = hospitalSrv.tpl.ExecuteTemplate(notificationBuffer, updateTemplate, addEmailDataDefaults(&emailData{
			HospitalID:   hospitalPBX.HospitalId,
			HospitalName: hospitalPBX.HospitalName,
			Content:      content,
		}))
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	emails := func() []string {
		if updateReq.GetSuspend() {
			return hospitalPBX.AdminEmails
		}
		if len(updateReq.GetHospital().GetAdminEmails()) > 0 {
			return updateReq.Hospital.AdminEmails
		}
		return hospitalPBX.AdminEmails
	}()

	emailNotification := &notification.EmailNotification{
		To:              emails,
		Subject:         subject,
		BodyContentType: "text/html",
		Body:            notificationBuffer.String(),
	}

	// Send notification with Backoff; retry 5 times in 10 seconds
	hospitalSrv.sendNotification(ctx, &notification.Notification{
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
	hospitalSrv.mu.Lock()
	if updateReq.Suspend {
		hospitalSrv.allowedHospitals[updateReq.HospitalId] = false
	} else {
		hospitalSrv.allowedHospitals[updateReq.HospitalId] = hospitalDB.Permission
	}
	hospitalSrv.mu.Unlock()

	return &hospital.HashResponse{
		HospitalId:    updateReq.HospitalId,
		OperationHash: addRes.Hash,
	}, nil
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

func (hospitalSrv *hospitalAPIServer) ListHospitals(
	ctx context.Context, listReq *hospital.ListHospitalsRequest,
) (*hospital.Hospitals, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToList)

	// Request must not be nil
	if listReq == nil {
		return nil, WrapError(errs.NilObject("ListHospitalsRequest"))
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	pageNumber, pageSize := normalizePage(listReq.GetPageNumber(), listReq.GetPageSize())
	offset := pageNumber*pageSize - pageSize

	hospitalDBs := make([]*hospitalModel, 0, pageSize)
	err = hospitalSrv.sqlDB.Offset(offset).Limit(pageSize).Find(&hospitalDBs).Error
	switch {
	case err == nil:
	default:
		return nil, WrapError(errs.SQLQueryFailed(err, "LIST"))
	}

	hospitalPBs := make([]*hospital.Hospital, 0, len(hospitalDBs))
	for _, hospitalDB := range hospitalDBs {
		hospitalPB, err := getHospitalPB(hospitalDB)
		if err != nil {
			return nil, WrapError(err)
		}
		hospitalPBs = append(hospitalPBs, hospitalPB)
	}

	return &hospital.Hospitals{
		Hospitals: hospitalPBs,
	}, nil
}

func (hospitalSrv *hospitalAPIServer) SearchHospitals(
	ctx context.Context, searchReq *hospital.SearchHospitalsRequest,
) (*hospital.Hospitals, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToSearch)

	// requets must not be nil
	if searchReq == nil {
		return nil, WrapError(errs.NilObject("SearchHospitalsRequest"))
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	// List hospitals in case search term is empty
	if searchReq.Query == "" {
		return hospitalSrv.ListHospitals(ctx, &hospital.ListHospitalsRequest{
			PageNumber: searchReq.GetPageNumber(),
			PageSize:   searchReq.GetPageSize(),
		})
	}

	pageNumber, pageSize := normalizePage(searchReq.GetPageNumber(), searchReq.GetPageSize())
	offset := pageNumber*pageSize - pageSize

	parsedQuery := db.ParseQuery(searchReq.Query, "hospital", "hospitals")

	hospitalDBs := make([]*hospitalModel, 0, pageSize)
	err = hospitalSrv.sqlDB.Offset(offset).Limit(pageSize).
		Find(&hospitalDBs, "MATCH(hospital_name) AGAINST(? IN BOOLEAN MODE )", parsedQuery).
		Error
	switch {
	case err == nil:
	default:
		return nil, WrapError(errs.SQLQueryFailed(err, "LIST"))
	}

	hospitalPBs := make([]*hospital.Hospital, 0, len(hospitalDBs))
	for _, hospitalDB := range hospitalDBs {
		hospitalPB, err := getHospitalPB(hospitalDB)
		if err != nil {
			return nil, WrapError(err)
		}
		hospitalPBs = append(hospitalPBs, hospitalPB)
	}

	return &hospital.Hospitals{
		Hospitals: hospitalPBs,
	}, nil
}

func (hospitalSrv *hospitalAPIServer) CheckSuspension(
	ctx context.Context, checkReq *hospital.CheckSuspensionRequest,
) (*hospital.CheckSuspensionResponse, error) {
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
	if strings.TrimSpace(checkReq.HospitalId) == "" {
		return nil, WrapError(errs.MissingCredential("HospitalId"))
	}

	// Read from map
	hospitalSrv.mu.RLock()
	allowed := hospitalSrv.allowedHospitals[checkReq.HospitalId]
	hospitalSrv.mu.RUnlock()

	if allowed {
		return &hospital.CheckSuspensionResponse{
			Suspended: false,
		}, nil
	}

	return &hospital.CheckSuspensionResponse{
		Suspended: true,
	}, nil
}

func (hospitalSrv *hospitalAPIServer) updatePermissions() {
	limit := 1000
	hospitals := make([]*hospitalModel, 0, limit)
	db := hospitalSrv.sqlDB.Table(tableName)
	fetchMore := true

	for fetchMore {
		db = db.Limit(limit).Find(&hospitals)
		if db.Error != nil {
			errs.LogError("failed to fetch hospitals: %v", db.Error)
			fetchMore = false
			break
		}
		if len(hospitals) < limit || db.RowsAffected < int64(limit) {
			fetchMore = false
		}

		// Lock mutex
		hospitalSrv.mu.Lock()
		for _, hospitalDB := range hospitals {
			if hospitalDB.Permission {
				hospitalSrv.allowedHospitals[hospitalDB.HospitalID] = true
			}
		}
		// Release to other routines
		hospitalSrv.mu.Unlock()
	}

	errs.LogInfo("finished updating hospitals permission")
}

func (hospitalSrv *hospitalAPIServer) sendNotification(
	ctx context.Context, notificationPB *notification.Notification,
) {
	ctxSend, cancel := context.WithTimeout(md.AddFromCtx(ctx), 10*time.Second)
	defer cancel()

	var err error
	for i := 0; i < 5; i++ {
		_, err = hospitalSrv.notificationClient.Send(
			ctxSend, notificationPB, grpc.WaitForReady(true),
		)
		if err != nil {
			// We wait for few seconds
			errs.LogError("failed to Send notification: %v", err)
			<-time.After(time.Duration(i) * time.Second)
			continue
		}
		break
	}
}
