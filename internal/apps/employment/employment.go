package employment

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/employment"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"strings"
)

const (
	failedToAdd       = "Failed to add employment"
	failedToUpdate    = "Failed to update employment"
	failedToDel       = "Failed to delete employment"
	failedToGet       = "Failed to get employment"
	failedToGetRecent = "Failed to get recent employments"
	failedToList      = "Failed to get employments"
	failedToCheck     = "Failed to check if employed"
)

type employmentAPIServer struct {
	sqlDB *gorm.DB
}

// NewEmploymentAPI creates a new employment API singleton
func NewEmploymentAPI(ctx context.Context, sqlDB *gorm.DB) (employment.EmploymentAPIServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case sqlDB == nil:
		return nil, errs.NilObject("SqlDB")
	}

	employmentServer := &employmentAPIServer{
		sqlDB: sqlDB,
	}

	// Automigration
	err := employmentServer.sqlDB.AutoMigrate(&employmentModel{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate employment table: %w", err)
	}

	return employmentServer, nil
}

func (employmentServer *employmentAPIServer) AddEmployment(
	ctx context.Context, addReq *employment.AddEmploymentRequest,
) (*employment.AddEmploymentResponse, error) {
	FailedToAddWrapper := errs.WrapErrorWithMsgFunc(failedToAdd)

	// Request must not be nil
	if addReq == nil {
		return nil, FailedToAddWrapper(errs.NilObject("AddEmploymentRequest"))
	}

	// Authentication
	_, err := auth.AuthenticateGroup(ctx, addReq.GetActor().GetActor())
	if err != nil {
		return nil, FailedToAddWrapper(err)
	}

	// Validation
	employmentPB := addReq.GetEmployment()
	switch {
	case employmentPB == nil:
		err = errs.NilObject("Employment")
	case strings.TrimSpace(employmentPB.AccountId) == "":
		err = errs.MissingCredential("AccountID")
	case strings.TrimSpace(employmentPB.FullName) == "":
		err = errs.MissingCredential("FullName")
	case employmentPB.EmploymentType == employment.EmploymentType_UNKNOWN:
		err = errs.MissingCredential("EmploymentType")
	case strings.TrimSpace(employmentPB.JoinedDate) == "":
		err = errs.MissingCredential("JoinedDate")
	case strings.TrimSpace(employmentPB.OrganizationType) == "":
		err = errs.MissingCredential("OrganizationType")
	case strings.TrimSpace(employmentPB.OrganizationName) == "":
		err = errs.MissingCredential("OrganizationName")
	case strings.TrimSpace(employmentPB.OrganizationId) == "":
		err = errs.MissingCredential("OrganizationId")
	case strings.TrimSpace(employmentPB.RoleAtOrganization) == "":
		err = errs.MissingCredential("RoleAtOrganization")
	case strings.TrimSpace(employmentPB.WorkEmail) == "":
		err = errs.MissingCredential("WorkEmail")
	}
	if err != nil {
		return nil, FailedToAddWrapper(err)
	}

	// Save to database
	employmentDB, err := getEmploymentDB(employmentPB)
	if err != nil {
		return nil, err
	}

	employmentDB.IsRecent = true
	employmentDB.StillEmployed = true
	employmentDB.EmploymentVerified = false

	err = employmentServer.sqlDB.Create(employmentDB).Error
	if err != nil {
		return nil, FailedToAddWrapper(err)
	}

	return &employment.AddEmploymentResponse{
		EmploymentId: employmentDB.EmploymentID,
	}, nil
}

func (employmentServer *employmentAPIServer) UpdateEmployment(
	ctx context.Context, updateReq *employment.UpdateEmploymentRequest,
) (*empty.Empty, error) {
	FailedToUpdateWrapper := errs.WrapErrorWithMsgFunc(failedToUpdate)

	// Request must not be nil
	if updateReq == nil {
		return nil, FailedToUpdateWrapper(errs.NilObject("UpdateEmploymentRequest"))
	}

	// Authentication
	_, err := auth.AuthenticateGroup(ctx, updateReq.GetActor().GetActor())
	if err != nil {
		return nil, FailedToUpdateWrapper(err)
	}

	// Validation
	employmentPB := updateReq.GetEmployment()
	switch {
	case employmentPB == nil:
		err = errs.NilObject("Employment")
	case strings.TrimSpace(employmentPB.EmploymentId) == "":
		err = errs.MissingCredential("EmploymentId")
	case strings.TrimSpace(employmentPB.AccountId) == "":
		err = errs.MissingCredential("AccountId")
	}
	if err != nil {
		return nil, FailedToUpdateWrapper(err)
	}

	// Update model in database
	employmentDB, err := getEmploymentDB(employmentPB)
	if err != nil {
		return nil, FailedToUpdateWrapper(err)
	}

	db := employmentServer.sqlDB.Model(employmentDB).
		Omit("still_employed", "employment_verified", "is_recent").
		Where("employment_id=?", employmentPB.EmploymentId).
		Updates(employmentDB)
	if db.RowsAffected == 0 {
		return nil, FailedToUpdateWrapper(errs.EmploymentNotFound(employmentPB.EmploymentId))
	}
	if db.Error != nil {
		return nil, FailedToUpdateWrapper(errs.SQLQueryFailed(err, "UPDATE"))
	}

	return &empty.Empty{}, nil
}

func (employmentServer *employmentAPIServer) GetEmployment(
	ctx context.Context, getReq *employment.GetEmploymentRequest,
) (*employment.Employment, error) {
	FailedToGetWrapper := errs.WrapErrorWithMsgFunc(failedToGet)

	// Request must not be nil
	if getReq == nil {
		return nil, FailedToGetWrapper(errs.NilObject("GetEmploymentRequest"))
	}

	// Authentication
	_, err := auth.AuthenticateGroup(ctx, getReq.GetActor().GetActor())
	if err != nil {
		return nil, FailedToGetWrapper(err)
	}

	// Validation
	switch {
	case strings.TrimSpace(getReq.EmploymentId) == "":
		return nil, FailedToGetWrapper(errs.MissingCredential("EmploymentID"))
	}

	// Get employment from model
	employmentDB := &employmentModel{}

	err = employmentServer.sqlDB.First(employmentDB, "employment_id=?", getReq.EmploymentId).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, FailedToGetWrapper(errs.EmploymentNotFound(getReq.EmploymentId))
	default:
		return nil, FailedToGetWrapper(errs.SQLQueryFailed(err, "SELECT"))
	}

	employmentPB, err := getEmploymentPB(employmentDB)
	if err != nil {
		return nil, FailedToGetWrapper(err)
	}

	return employmentPB, nil
}

func (employmentServer *employmentAPIServer) DeleteEmployment(
	ctx context.Context, delReq *employment.DeleteEmploymentRequest,
) (*empty.Empty, error) {
	FailedToDelWrapper := errs.WrapErrorWithMsgFunc(failedToDel)

	// Request must not be nil
	if delReq == nil {
		return nil, FailedToDelWrapper(errs.NilObject("DeleteEmploymentRequest"))
	}

	// Authentication
	_, err := auth.AuthenticateGroup(ctx, delReq.GetActor().GetActor())
	if err != nil {
		return nil, FailedToDelWrapper(err)
	}

	// Validation
	switch {
	case strings.TrimSpace(delReq.EmploymentId) == "":
		return nil, FailedToDelWrapper(errs.MissingCredential("EmploymentID"))
	}

	// Soft Delete employment from model
	db := employmentServer.sqlDB.Delete(&employmentModel{}, "employment_id=?", delReq.EmploymentId)
	if err != nil {
		return nil, FailedToDelWrapper(errs.SQLQueryFailed(err, "DELETE"))
	}
	if db.RowsAffected == 0 {
		return nil, FailedToDelWrapper(errs.EmploymentNotFound(delReq.EmploymentId))
	}
	if db.Error != nil {
		return nil, FailedToDelWrapper(errs.SQLQueryFailed(err, "DELETE"))
	}

	return &empty.Empty{}, nil
}

func (employmentServer *employmentAPIServer) GetRecentEmployment(
	ctx context.Context, getReq *employment.GetRecentEmploymentRequest,
) (*employment.Employment, error) {
	FailedToGetWrapper := errs.WrapErrorWithMsgFunc(failedToGetRecent)

	// Request must not be nil
	if getReq == nil {
		return nil, FailedToGetWrapper(errs.NilObject("GetRecentEmploymentRequest"))
	}

	// Authentication
	_, err := auth.AuthenticateGroup(ctx, getReq.GetActor().GetActor())
	if err != nil {
		return nil, FailedToGetWrapper(err)
	}

	// Validation
	switch {
	case strings.TrimSpace(getReq.AccountId) == "":
		return nil, FailedToGetWrapper(errs.MissingCredential("AccountId"))
	}

	// Get employment from model
	employmentDBs := make([]*employmentModel, 0, 1)
	err = employmentServer.sqlDB.Order("created_at DESC").Limit(1).
		Find(&employmentDBs, "account_id=? AND is_recent=1 AND employment_verified=1", getReq.AccountId).
		Error
	switch {
	case err == nil:
	default:
		return nil, FailedToGetWrapper(errs.SQLQueryFailed(err, "SELECT"))
	}

	if len(employmentDBs) == 0 {
		return nil, FailedToGetWrapper(errs.UserEmploymentsNotFound(getReq.AccountId))
	}

	employmentPB, err := getEmploymentPB(employmentDBs[0])
	if err != nil {
		return nil, FailedToGetWrapper(err)
	}

	return employmentPB, nil
}

const defaultPageSize = 10

func normalizePage(pageToken, pageSize int32) (int, int) {
	if pageToken <= 0 {
		pageToken = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return int(pageToken), int(pageSize)
}

func (employmentServer *employmentAPIServer) GetEmployments(
	ctx context.Context, getReq *employment.GetEmploymentsRequest,
) (*employment.Employments, error) {
	FailedToGetWrapper := errs.WrapErrorWithMsgFunc(failedToList)

	// Request must not be nil
	if getReq == nil {
		return nil, FailedToGetWrapper(errs.NilObject("GetEmploymentsRequest"))
	}

	// Authentication
	_, err := auth.AuthenticateGroup(ctx, getReq.GetActor().GetActor())
	if err != nil {
		return nil, FailedToGetWrapper(err)
	}

	// Validation
	switch {
	case strings.TrimSpace(getReq.AccountId) == "":
		return nil, FailedToGetWrapper(errs.MissingCredential("AccountId"))
	}

	pageNumber, pageSize := normalizePage(getReq.GetPageNumber(), getReq.GetPageSize())
	offset := pageNumber*pageSize - pageSize

	// Get employments from model
	employmentDBs := make([]*employmentModel, 0)
	err = employmentServer.sqlDB.Order("created_at DESC").Offset(offset).Limit(pageSize).
		Find(&employmentDBs, "account_id=?", getReq.AccountId).Error
	switch {
	case err == nil:
	default:
		return nil, FailedToGetWrapper(errs.SQLQueryFailed(err, "SELECT"))
	}

	if len(employmentDBs) == 0 {
		return nil, FailedToGetWrapper(errs.UserEmploymentsNotFound(getReq.AccountId))
	}

	employmentPBs := make([]*employment.Employment, 0, len(employmentDBs))
	for _, employmentDB := range employmentDBs {
		employmentPB, err := getEmploymentPB(employmentDB)
		if err != nil {
			return nil, FailedToGetWrapper(err)
		}
		employmentPBs = append(employmentPBs, employmentPB)
	}

	return &employment.Employments{
		Employments: employmentPBs,
	}, nil
}

func (employmentServer *employmentAPIServer) CheckEmploymentStatus(
	ctx context.Context, checkReq *employment.CheckEmploymentStatusRequest,
) (*employment.CheckEmploymentStatusResponse, error) {
	FailedToCheckWrapper := errs.WrapErrorWithMsgFunc(failedToCheck)

	// Request must not be nil
	if checkReq == nil {
		return nil, FailedToCheckWrapper(errs.NilObject("CheckEmploymentStatusRequest"))
	}

	// Authentication
	_, err := auth.AuthenticateGroup(ctx, checkReq.GetActor().GetActor())
	if err != nil {
		return nil, FailedToCheckWrapper(err)
	}

	// Validation
	switch {
	case strings.TrimSpace(checkReq.AccountId) == "":
		return nil, FailedToCheckWrapper(errs.MissingCredential("AccountId"))
	}

	// Query db to check if employed
	employmentDBs := make([]*employmentModel, 0, 3)
	err = employmentServer.sqlDB.Limit(2).Order("still_employed, employment_verified, is_recent DESC").
		Find(&employmentDBs, "account_id=?", checkReq.AccountId).Error
	switch {
	case err == nil:
	default:
		return nil, FailedToCheckWrapper(errs.SQLQueryFailed(err, "SELECT"))
	}

	if len(employmentDBs) > 0 {
		for _, employmentDB := range employmentDBs {
			if employmentDB.EmploymentVerified {
				return &employment.CheckEmploymentStatusResponse{
					IsEmployed: true,
					IsVerified: true,
				}, nil
			}
		}
		return &employment.CheckEmploymentStatusResponse{
			IsEmployed: true,
			IsVerified: false,
		}, nil
	}

	return &employment.CheckEmploymentStatusResponse{
		IsEmployed: false,
		IsVerified: false,
	}, nil
}
