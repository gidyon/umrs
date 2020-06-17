package account

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	template_util "github.com/gidyon/umrs/internal/pkg/template"
	"github.com/gidyon/umrs/pkg/api/account"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/subscriber"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"strings"
	"time"
)

func (accountAPI *accountAPIServer) validateChangeAccountRequest(
	ctx context.Context, changeReq *account.ChangeAccountRequest,
) (*Account, error) {
	// Request should not be nil
	if changeReq == nil {
		return nil, errs.NilObject("ChangeAccountType")
	}

	// Authenticate super admin
	_, err := auth.AuthenticateSuperAdmin(ctx)
	if err != nil {
		return nil, err
	}

	accountID := changeReq.GetAccountId()
	superAdminID := changeReq.GetSuperAdminId()

	// Validation
	switch {
	case strings.TrimSpace(accountID) == "":
		return nil, errs.MissingCredential("AdminID")
	case strings.TrimSpace(superAdminID) == "":
		return nil, errs.MissingCredential("Super AdminID")
	}

	superAdmin := &Account{}
	// Get super admin
	err = accountAPI.sqlDB.Unscoped().Select("account_state,account_type").
		First(superAdmin, "account_id=?", superAdminID).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.AccountDoesntExist(superAdminID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Super admin account must be active
	if superAdmin.AccountState != account.AccountState_ACTIVE.String() {
		return nil, errs.SuperAdminNotActive()
	}

	// Super Admin must be owner
	if superAdmin.AccountType != account.AccountType_ADMIN_OWNER.String() {
		return nil, errs.OnlyOwnerPermitted()
	}

	accountDB := &Account{}
	// The account must exist
	err = accountAPI.sqlDB.Unscoped().First(accountDB, "account_id=?", accountID).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.AccountDoesntExist(accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	return accountDB, nil
}

func (accountAPI *accountAPIServer) BlockAccount(
	ctx context.Context, blockReq *account.ChangeAccountRequest,
) (*empty.Empty, error) {
	// Validate the request, super admin credentials and account owner
	accountDB, err := accountAPI.validateChangeAccountRequest(ctx, blockReq)
	if err != nil {
		return nil, err
	}

	// The state must be active in order to block it
	if accountDB.AccountState != account.AccountState_ACTIVE.String() {
		return nil, errs.AccountNotActive()
	}

	// Update the model
	err = accountAPI.sqlDB.Model(&Account{}).Where("account_id=?", blockReq.AccountId).
		Update("account_state", account.AccountState_BLOCKED.String()).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "BlockAccount")
	}

	// ================================ NOTIFICATION =============================
	// Template for sending email
	emailContent := template_util.EmailData{
		FirstName:    accountDB.FirstName,
		LastName:     accountDB.LastName,
		AccountID:    accountDB.AccountID,
		AppName:      accountAPI.appName,
		Reason:       blockReq.Reason,
		TemplateName: "block",
	}

	content := bytes.NewBuffer(make([]byte, 0, 64))
	err = accountAPI.tpl.ExecuteTemplate(content, templateName, emailContent)
	if err != nil {
		errs.LogError("failed to execute template: %v", err)
		return &empty.Empty{}, nil
	}

	fullName := accountDB.FirstName + " " + accountDB.LastName
	smsMessage := fmt.Sprintf(
		"Hello %s. We are sad to inform you that your account has been blocked. Contact us for more information", fullName,
	)
	notificationPB, err := accountAPI.createNotification(
		ctx, accountDB, "Your Account Was Blocked", content.String(), smsMessage,
	)
	if err != nil {
		errs.LogWarn("failed to create notification object: %v", err)
		return &empty.Empty{}, nil
	}

	// Send the notification
	_, err = accountAPI.notificationClient.Send(ctx, notificationPB, grpc.WaitForReady(true))
	if err != nil {
		errs.LogError("failed to get Send notification: %v", err)
		return &empty.Empty{}, nil
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) UnBlockAccount(
	ctx context.Context, unblockReq *account.ChangeAccountRequest,
) (*empty.Empty, error) {
	// Validate the request, super admin credentials and account owner
	accountDB, err := accountAPI.validateChangeAccountRequest(ctx, unblockReq)
	if err != nil {
		return nil, err
	}

	// The account state must be blocked
	if accountDB.AccountState != account.AccountState_BLOCKED.String() {
		return nil, errs.AccountNotBlocked()
	}

	// Update the model
	err = accountAPI.sqlDB.Model(&Account{}).Where("account_id=?", unblockReq.AccountId).
		Update("account_state", account.AccountState_ACTIVE.String()).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UnBlockAccount")
	}

	// ================================ SEND NOTIFICATION =============================
	// Template for sending email
	emailContent := template_util.EmailData{
		FirstName:    accountDB.FirstName,
		LastName:     accountDB.LastName,
		AccountID:    accountDB.AccountID,
		AppName:      accountAPI.appName,
		Reason:       unblockReq.Reason,
		TemplateName: "unblock",
	}

	content := bytes.NewBuffer(make([]byte, 0, 64))

	err = accountAPI.tpl.ExecuteTemplate(content, templateName, emailContent)
	if err != nil {
		errs.LogError("failed to execute template: %v", err)
		return &empty.Empty{}, nil
	}

	fullName := accountDB.FirstName + " " + accountDB.LastName
	smsMessage := fmt.Sprintf(
		"Hello %s. Glad to inform you that your account has been unblocked. For support and inquiries visit our website",
		fullName,
	)
	notificationPB, err := accountAPI.createNotification(
		ctx, accountDB, "Your Account Was Unblocked", content.String(), smsMessage,
	)
	if err != nil {
		errs.LogWarn("failed to create notification object: %v", err)
		return &empty.Empty{}, nil
	}

	// Send the notification
	_, err = accountAPI.notificationClient.Send(ctx, notificationPB, grpc.WaitForReady(true))
	if err != nil {
		errs.LogError("failed to get Send notification: %v", err)
		return &empty.Empty{}, nil
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) AdminActivate(
	ctx context.Context, adminActivateReq *account.ChangeAccountRequest,
) (*empty.Empty, error) {
	// Validate the request, super admin credentials and account owner
	accountDB, err := accountAPI.validateChangeAccountRequest(ctx, adminActivateReq)
	if err != nil {
		return nil, err
	}

	// The account state must be inactive
	if accountDB.AccountState != account.AccountState_INACTIVE.String() {
		return nil, errs.AccountNotInactive()
	}

	// Update the model
	err = accountAPI.sqlDB.Model(&Account{}).Where("account_id=?", adminActivateReq.AccountId).
		Update("account_state", account.AccountState_ACTIVE.String()).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UnBlockAccount")
	}

	// ================================ SEND NOTIFICATION =============================
	// Template for sending email
	emailContent := template_util.EmailData{
		FirstName:    accountDB.FirstName,
		LastName:     accountDB.LastName,
		AccountID:    accountDB.AccountID,
		AppName:      accountAPI.appName,
		TemplateName: "activated",
	}

	content := bytes.NewBuffer(make([]byte, 0, 64))
	err = accountAPI.tpl.ExecuteTemplate(content, templateName, emailContent)
	if err != nil {
		errs.LogError("failed to execute template: %v", err)
		return &empty.Empty{}, nil
	}

	fullName := accountDB.FirstName + " " + accountDB.LastName
	smsMessage := fmt.Sprintf(
		"Hello %s. Your account has been activated by administrator", fullName,
	)
	notificationPB, err := accountAPI.createNotification(
		ctx, accountDB, "Account Activated", content.String(), smsMessage,
	)
	if err != nil {
		errs.LogError("failed to create notification object: %v", err)
		return &empty.Empty{}, nil
	}

	// Send the notification
	_, err = accountAPI.notificationClient.Send(ctx, notificationPB, grpc.WaitForReady(true))
	if err != nil {
		errs.LogError("failed to get Send notification: %v", err)
		return &empty.Empty{}, nil
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) ChangeAccountType(
	ctx context.Context, changeLevelReq *account.ChangeAccountTypeRequest,
) (*empty.Empty, error) {
	// Request should not be nil
	if changeLevelReq == nil {
		return nil, errs.NilObject("ChangeAccountTypeRequest")
	}

	accountDB, err := accountAPI.validateChangeAccountRequest(ctx, &account.ChangeAccountRequest{
		AccountId:    changeLevelReq.GetAccountId(),
		SuperAdminId: changeLevelReq.GetSuperAdminId(),
	})
	if err != nil {
		return nil, err
	}

	accountType := changeLevelReq.GetType().String()

	// Update the model
	err = accountAPI.sqlDB.Model(&Account{}).Where("account_id=?", changeLevelReq.AccountId).
		Update("account_type", accountType).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "ChangeAccountType")
	}

	// ================================ SEND NOTIFICATION =============================
	// Template for sending email
	emailContent := template_util.EmailData{
		FirstName:    accountDB.FirstName,
		LastName:     accountDB.LastName,
		AccountID:    accountDB.AccountID,
		AppName:      accountAPI.appName,
		TemplateName: "change",
	}

	content := bytes.NewBuffer(make([]byte, 0, 64))
	err = accountAPI.tpl.ExecuteTemplate(content, templateName, emailContent)
	if err != nil {
		errs.LogError("failed to execute template: %v", err)
		return &empty.Empty{}, nil
	}

	fullName := accountDB.FirstName + " " + accountDB.LastName
	smsMessage := fmt.Sprintf(
		"Hello %s. Your account type has been changed to %s by administrator",
		fullName, changeLevelReq.Type.String(),
	)
	notificationPB, err := accountAPI.createNotification(
		ctx, accountDB, "Account Type changed", content.String(), smsMessage,
	)
	if err != nil {
		errs.LogError("failed to create notification object: %v", err)
		return &empty.Empty{}, nil
	}

	// Send the notification
	_, err = accountAPI.notificationClient.Send(ctx, notificationPB, grpc.WaitForReady(true))
	if err != nil {
		errs.LogError("failed to get Send notification: %v", err)
		return &empty.Empty{}, nil
	}

	return &empty.Empty{}, nil
}

// Restores an account previously deleted
func (accountAPI *accountAPIServer) Undelete(
	ctx context.Context, undelReq *account.ChangeAccountRequest,
) (*empty.Empty, error) {
	accountDB, err := accountAPI.validateChangeAccountRequest(ctx, undelReq)
	if err != nil {
		return nil, err
	}

	// Account must be deleted
	if accountDB.DeletedAt == nil {
		return nil, errs.AccountNotDeleted()
	}

	// Undelete from database
	err = accountAPI.sqlDB.Model(&Account{}).Unscoped().Update("deleted_at", nil).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	// ================================ SEND NOTIFICATION =============================
	// Template for sending email
	emailContent := template_util.EmailData{
		FirstName:    accountDB.FirstName,
		LastName:     accountDB.LastName,
		AccountID:    accountDB.AccountID,
		AppName:      accountAPI.appName,
		TemplateName: "undelete",
		Reason:       undelReq.Reason,
	}

	content := bytes.NewBuffer(make([]byte, 0, 64))
	err = accountAPI.tpl.ExecuteTemplate(content, templateName, emailContent)
	if err != nil {
		errs.LogError("failed to execute template: %v", err)
		return &empty.Empty{}, nil
	}

	fullName := accountDB.FirstName + " " + accountDB.LastName
	smsMessage := fmt.Sprintf(
		"Hello %s, glad to inform you that your account has been restored", fullName,
	)
	notificationPB, err := accountAPI.createNotification(
		ctx, accountDB, "Account Restored", content.String(), smsMessage,
	)
	if err != nil {
		errs.LogError("failed to create notification object: %v", err)
		return &empty.Empty{}, nil
	}

	// Send the notification
	_, err = accountAPI.notificationClient.Send(ctx, notificationPB, grpc.WaitForReady(true))
	if err != nil {
		errs.LogError("failed to Send notification: %v", err)
		return &empty.Empty{}, nil
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) createNotification(
	ctx context.Context, accountDB *Account, subject, emailContent, smsMessage string,
) (*notification.Notification, error) {
	// We get the send method
	sendMethod, err := accountAPI.subscriberClient.GetSendMethod(
		ctx, &subscriber.GetSendMethodRequest{
			AccountId: accountDB.AccountID,
		}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	// Create notification
	notificationPB := &notification.Notification{
		NotificationId: uuid.New().String(),
		OwnerIds:       []string{accountDB.AccountID},
		Priority:       notification.Priority_MEDIUM,
		SendMethod:     sendMethod.SendMethod,
		Content: &notification.NotificationContent{
			Subject: subject,
			Data:    emailContent,
		},
		CreateTimeSec: time.Now().Unix(),
		Save:          true,
	}

	return notificationPB, nil
}
