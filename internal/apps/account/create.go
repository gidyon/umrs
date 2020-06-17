package account

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/internal/pkg/md"
	template_util "github.com/gidyon/umrs/internal/pkg/template"
	"github.com/gidyon/umrs/pkg/api/account"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
	"time"
)

func (accountAPI *accountAPIServer) Create(
	ctx context.Context, createReq *account.CreateRequest,
) (*account.CreateResponse, error) {
	// Request should not be nil
	if createReq == nil {
		return nil, errs.NilObject("CreateRequest")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	accountPB := createReq.GetAccount()
	if accountPB == nil {
		return nil, errs.NilObject("Account")
	}

	// Validation
	switch {
	case strings.TrimSpace(createReq.AccountLabel) == "":
		err = errs.MissingCredential("AccountLabel")
	case strings.TrimSpace(accountPB.FirstName) == "":
		err = errs.MissingCredential("First Name")
	case strings.TrimSpace(accountPB.LastName) == "":
		err = errs.MissingCredential("Last Name")
	case strings.TrimSpace(accountPB.NationalId) == "":
		err = errs.MissingCredential("National/Huduma ID")
		// case strings.TrimSpace(accountPB.Phone) == "" && strings.TrimSpace(accountPB.Email) == "":
		// 	err = errs.MissingCredential("Phone and Email")
	}
	if err != nil {
		return nil, err
	}

	accountPB.AccountLabels = []string{}
	accountPB.AccountLabels = append(accountPB.AccountLabels, createReq.AccountLabel)

	accountDB, err := getAccountDB(accountPB)
	if err != nil {
		return nil, err
	}

	accountState := account.AccountState_INACTIVE.String()

	if createReq.GetByAdmin() {
		// Authenticate the admin
		p, err := auth.AuthenticateAdmin(ctx)
		if err != nil {
			return nil, err
		}
		if p.ID != createReq.AdminId {
			return nil, errs.WrapMessage(codes.Unauthenticated, "Token id and request id do not match")
		}
		accountState = account.AccountState_ACTIVE.String()
	}

	accountDB.AccountState = accountState

	accountPrivate := createReq.GetPrivateAccount()
	if accountPrivate != nil {
		accountDB.SecurityAnswer = accountPrivate.GetSecurityQuestion()
		// Store password as encrypted
		accountDB.SecurityAnswer = accountPrivate.GetSecurityAnswer()
		if strings.TrimSpace(accountPrivate.Password) != "" {
			newPass, err := genHash(accountPrivate.GetPassword())
			if err != nil {
				return nil, errs.FailedToGenHashedPass(err)
			}
			accountDB.Password = newPass
		}
	}

	// If admin, grant them least priviledges at first
	isAdmin := accountPB.AccountType != account.AccountType_USER_OWNER
	if isAdmin {
		accountDB.AccountType = account.AccountType_ADMIN_VIEWER.String()
	}

	// Start a transaction
	tx := accountAPI.sqlDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.IsolationLevel(0),
	})
	defer func() {
		if err := recover(); err != nil {
			errs.LogError("recovering from panic: %v", err)
		}
	}()

	if tx.Error != nil {
		return nil, errs.FailedToBeginTx(err)
	}

	err = tx.Create(accountDB).Error
	switch {
	case err == nil:
	default:
		// Checks whether the error was due to similar account existing
		recordExist := func(err error) bool {
			return strings.Contains(
				strings.ToLower(err.Error()), "duplicate entry",
			)
		}

		// Clarifies what exists between email, phone or user
		emailOrPhone := func(err error) (string, string) {
			if strings.Contains(strings.ToLower(err.Error()), "email") {
				return "email", accountDB.Email
			}
			if strings.Contains(strings.ToLower(err.Error()), "phone") {
				return "phone", accountDB.Phone
			}
			if strings.Contains(strings.ToLower(err.Error()), "national_id") {
				return "national id", accountDB.NationalID
			}
			return "id", accountDB.AccountID
		}

		if recordExist(err) {
			tx.Rollback()
			return nil, errs.AccountDoesExist(emailOrPhone(err))
		}

		tx.Rollback()
		return nil, errs.SQLQueryFailed(err, "CreateUser")
	}

	// Generate user token with expiration of 6 hours
	jwtToken, err := auth.GenToken(ctx, &auth.Payload{
		ID:           accountDB.AccountID,
		FirstName:    accountPB.FirstName,
		LastName:     accountPB.LastName,
		PhoneNumber:  accountPB.Phone,
		EmailAddress: accountPB.Email,
		Group:        ledger.Actor_value[createReq.AccountLabel],
		Label:        createReq.AccountLabel,
	}, 0, time.Now().Add(time.Duration(6*time.Hour)).Unix())
	if err != nil {
		tx.Rollback()
		return nil, errs.FailedToGenToken(err)
	}

	subject := fmt.Sprintf("Account Action Required: Activate Your %s account", accountAPI.appName)

	sendMethod := func() notification.SendMethod {
		if accountPB.Email != "" {
			return notification.SendMethod_EMAIL
		}
		if accountPB.Phone != "" {
			return notification.SendMethod_SMS
		}
		return notification.SendMethod_EMAIL_AND_SMS
	}()

	notificationPB := &notification.Notification{
		NotificationId: uuid.New().String(),
		OwnerIds:       []string{accountDB.AccountID},
		Priority:       notification.Priority_MEDIUM,
		SendMethod:     sendMethod,
		CreateTimeSec:  time.Now().Unix(),
		Content: &notification.NotificationContent{
			Subject: subject,
			Data: fmt.Sprintf(
				"Welcome %s %s. You have created your umrs account successful",
				accountDB.FirstName, accountDB.LastName,
			),
		},
		Save: true,
	}

	switch sendMethod {
	case notification.SendMethod_EMAIL:
		// Template for sending email
		emailContent := template_util.EmailData{
			FirstName: accountDB.FirstName,
			LastName:  accountDB.LastName,
			AccountID: accountDB.AccountID,
			Token:     jwtToken,
			Label:     createReq.AccountLabel,
			Link: fmt.Sprintf(
				"%s?token=%s?&account_id=%s",
				accountAPI.activationURL, jwtToken, accountDB.AccountID,
			),
		}

		content := bytes.NewBuffer(make([]byte, 0, 64))
		err = accountAPI.tpl.ExecuteTemplate(content, templateName, emailContent)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		notificationPB.Payload = &notification.Notification_EmailNotification{
			EmailNotification: &notification.EmailNotification{
				To:              []string{accountDB.Email},
				Subject:         subject,
				BodyContentType: "text/html",
				Body:            content.String(),
			},
		}
	case notification.SendMethod_SMS:
		smsMessage := fmt.Sprintf(
			"Welcome %s %s, you have created your umrs account successful. To activate your account, follow this link %s", accountDB.FirstName, accountDB.LastName,
			fmt.Sprintf(
				"%s?account_id=%s&token=%s",
				accountAPI.activationURL, accountDB.AccountID, jwtToken,
			),
		)
		notificationPB.Payload = &notification.Notification_SmsNotification{
			SmsNotification: &notification.SMSNotification{
				Keyword:          subject,
				Message:          smsMessage,
				DestinationPhone: []string{accountDB.Phone},
			},
		}
	}

	// Create notification account and later send notification
	_, err = accountAPI.notificationClient.CreateNotificationAccount(
		md.AddFromCtx(ctx),
		&notification.CreateNotificationAccountRequest{
			Channels:   []string{"default"},
			AccountId:  accountDB.AccountID,
			Email:      accountDB.Email,
			Phone:      accountDB.Phone,
			SendMethod: notificationPB.SendMethod,
		}, grpc.WaitForReady(true),
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errs.FailedToCommitTx(err)
	}

	ctxSend, cancel := context.WithTimeout(md.AddFromCtx(ctx), 10*time.Second)
	defer cancel()

	// Backoff; retry 5 times in 10 seconds
	for i := 0; i < 5; i++ {
		_, err = accountAPI.notificationClient.Send(
			ctxSend, notificationPB, grpc.WaitForReady(true),
		)
		// We wait for few seconds
		if err != nil {
			<-time.After(time.Duration(i) * time.Second)
			continue
		}
		break
	}

	return &account.CreateResponse{
		AccountId: accountDB.AccountID,
	}, nil
}
