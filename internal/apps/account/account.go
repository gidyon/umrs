package account

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/umrs/internal/pkg/auth"
	db_util "github.com/gidyon/umrs/internal/pkg/db"
	"github.com/gidyon/umrs/internal/pkg/errs"
	template_util "github.com/gidyon/umrs/internal/pkg/template"
	"github.com/gidyon/umrs/pkg/api/account"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/subscriber"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"html/template"
	"math/rand"
	"os"
	"strings"
	"time"
)

const templateName = "base"

type accountAPIServer struct {
	ctx                context.Context
	activationURL      string
	appName            string
	sqlDB              *gorm.DB
	redisDB            *redis.Client
	notificationClient notification.NotificationServiceClient
	subscriberClient   subscriber.SubscriberAPIClient
	tpl                *template.Template
}

// Options contain parameters for NewAccountAPI
type Options struct {
	AppName            string
	SQLDB              *gorm.DB
	RedisDB            *redis.Client
	NotificationClient notification.NotificationServiceClient
	SubscriberClient   subscriber.SubscriberAPIClient
}

// NewAccountAPI creates a singleton of an AccountAPIServer
func NewAccountAPI(
	ctx context.Context, activationURL string, opt *Options,
) (account.AccountAPIServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case strings.TrimSpace(opt.AppName) == "":
		return nil, errs.MissingCredential("AppName")
	case opt.SQLDB == nil:
		return nil, errs.NilObject("SqlDB")
	case opt.RedisDB == nil:
		return nil, errs.NilObject("RedisDB")
	case opt.NotificationClient == nil:
		return nil, errs.NilObject("NotificationClient")
	case opt.SubscriberClient == nil:
		return nil, errs.NilObject("SubscriberClient")
	}

	accountAPI := &accountAPIServer{
		ctx:                ctx,
		activationURL:      activationURL,
		sqlDB:              opt.SQLDB,
		redisDB:            opt.RedisDB,
		notificationClient: opt.NotificationClient,
		subscriberClient:   opt.SubscriberClient,
	}

	// Perform auto migration
	err := accountAPI.sqlDB.AutoMigrate(&Account{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate accounts table: %w", err)
	}

	// Create a full text search index
	err = db_util.CreateFullTextIndex(accountAPI.sqlDB, tableName, "national_id", "email", "phone")
	if err != nil {
		return nil, fmt.Errorf("failed to create full text index: %v", err)
	}

	// Read template files from directory
	tFiles, err := template_util.ReadFiles(os.Getenv("TEMPLATES_DIR"))
	if err != nil {
		return nil, fmt.Errorf("failed to read template files in directory: %w", err)
	}

	// Parse template
	accountAPI.tpl, err = template_util.ParseTemplate(tFiles...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return accountAPI, nil
}

func (accountAPI *accountAPIServer) Activate(
	ctx context.Context, activateReq *account.ActivateRequest,
) (*account.ActivateResponse, error) {
	// Request should not be nil
	if activateReq == nil {
		return nil, errs.NilObject("ActivateRequest")
	}

	var err error
	// Validation 1
	switch {
	case activateReq.Token == "":
		err = errs.MissingCredential("Token")
	case activateReq.AccountId == "":
		err = errs.MissingCredential("AccountID")
	}
	if err != nil {
		return nil, err
	}

	// Retrieve token claims
	claims, err := auth.ParseToken(activateReq.Token)
	if err != nil {
		return nil, errs.FailedToParseToken(err)
	}

	payload := claims.Payload
	if payload == nil {
		return nil, errs.NilObject("TokenPayload")
	}

	// Validation 2
	switch {
	case payload.ID == "":
		err = errs.MissingCredential("TokenID")
	}
	if err != nil {
		return nil, err
	}

	// Compare that account id match
	if payload.ID != activateReq.AccountId {
		return nil, errs.TokenCredentialNotMatching("AccountID")
	}

	// Check that account Exists
	if notFound := accountAPI.sqlDB.Select("account_state,account_type").
		First(&Account{}, "account_id=?", payload.ID).RecordNotFound(); notFound {
		return nil, errs.AccountDoesntExist(payload.ID)
	}

	accountDB := &Account{
		AccountState: account.AccountState_ACTIVE.String(),
	}

	// Update the model of the user to activate their account
	err = accountAPI.sqlDB.Model(accountDB).Where("account_id=?", payload.ID).
		Update("account_state", account.AccountState_ACTIVE.String()).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &account.ActivateResponse{}, nil
}

func (accountAPI *accountAPIServer) Update(
	ctx context.Context, updateReq *account.UpdateRequest,
) (*empty.Empty, error) {
	// Request should not be nil
	if updateReq == nil {
		return nil, errs.NilObject("UpdateRequest")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	accountID := updateReq.GetAccountId()
	accountPB := updateReq.GetAccount()

	// Validation
	switch {
	case strings.TrimSpace(accountID) == "":
		err = errs.MissingCredential("AccountID")
	case accountPB == nil:
		err = errs.NilObject("Account")
	}
	if err != nil {
		return nil, err
	}

	// Get the account details from database
	accountDB := &Account{}
	err = accountAPI.sqlDB.Select("account_state,account_type").
		First(accountDB, "account_id=?", accountID).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.AccountDoesntExist(accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Check that account is not blocked
	if accountDB.AccountState == account.AccountState_BLOCKED.String() {
		return nil, errs.AccountBlocked()
	}

	accountDBX, err := getAccountDB(accountPB)
	if err != nil {
		return nil, err
	}

	// Update the model; omit "account_id", "account_type", "account_state"
	err = accountAPI.sqlDB.Model(accountDBX).
		Omit("account_id", "account_type", "account_state", "password", "security_answer", "security_question").
		Where("account_id=?", accountID).
		Updates(accountDBX).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UPDATE")
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) UpdatePrivate(
	ctx context.Context, updatePrivateReq *account.UpdatePrivateRequest,
) (*empty.Empty, error) {
	// Request should not be nil
	if updatePrivateReq == nil {
		return nil, errs.NilObject("UpdatePrivateRequest")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	accountID := updatePrivateReq.GetAccountId()
	accountPrivate := updatePrivateReq.GetPrivateAccount()

	// Validation
	switch {
	case strings.TrimSpace(accountID) == "":
		err = errs.MissingCredential("AccountID")
	case accountPrivate == nil:
		err = errs.NilObject("PrivateAccount")
	}
	if err != nil {
		return nil, err
	}

	// Get the account details from database
	accountDB := &Account{}
	err = accountAPI.sqlDB.Select("account_state,account_type").
		First(accountDB, "account_id=?", accountID).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.AccountDoesntExist(accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Check that account is not blocked
	if accountDB.AccountState == account.AccountState_BLOCKED.String() {
		return nil, errs.AccountBlocked()
	}

	// Hash the password if not empty
	if strings.TrimSpace(accountPrivate.Password) != "" {
		token, err := accountAPI.redisDB.Get(accountID).Result()
		switch {
		case err == nil:
		case err == redis.Nil:
			return nil, errs.WrapMessage(codes.PermissionDenied, "No update token found")
		default:
			return nil, errs.RedisCmdFailed(err, "Get")
		}

		if token != updatePrivateReq.ChangeToken {
			return nil, errs.WrapMessage(codes.PermissionDenied, "Token is incorrect")
		}

		// Passwords must be similar
		if accountPrivate.ConfirmPassword != accountPrivate.Password {
			return nil, errs.PasswordNoMatch()
		}

		accountPrivate.Password, err = genHash(accountPrivate.Password)
		if err != nil {
			return nil, errs.FailedToGenHashedPass(err)
		}
	}

	// Create database model of the new account
	privateDB := &Account{
		SecurityQuestion: accountPrivate.SecurityQuestion,
		SecurityAnswer:   accountPrivate.SecurityAnswer,
		Password:         accountPrivate.Password,
	}

	// Update the model
	err = accountAPI.sqlDB.Model(privateDB).Unscoped().
		Select("security_question,security_answer,password").
		Where("account_id=?", accountID).Updates(privateDB).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "UpdatePrivateProfile")
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) Delete(
	ctx context.Context, delReq *account.DeleteRequest,
) (*empty.Empty, error) {
	// Request should not be nil
	if delReq == nil {
		return nil, errs.NilObject("DeleteRequest")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Check AccountID is provided
	accountID := delReq.GetAccountId()
	if strings.TrimSpace(accountID) == "" {
		return nil, errs.MissingCredential("AccountID")
	}

	// Get the account details from database
	accountDB := &Account{}
	err = accountAPI.sqlDB.Select("account_state,account_type").
		First(accountDB, "account_id=?", accountID).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.AccountDoesntExist(accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Check that account is not blocked
	if accountDB.AccountState == account.AccountState_BLOCKED.String() {
		return nil, errs.AccountBlocked()
	}

	// Soft delete their account
	err = accountAPI.sqlDB.Delete(accountDB, "account_id=?", accountID).Error
	if err != nil {
		return nil, errs.SQLQueryFailed(err, "DELETE")
	}

	return &empty.Empty{}, nil
}

func (accountAPI *accountAPIServer) Get(
	ctx context.Context, getReq *account.GetRequest,
) (*account.Account, error) {
	// Request should not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetRequest")
	}

	var err error
	// Authenticate the request
	if getReq.Privileged {
		// Authenticate groups
		err = auth.AuthenticateGroups(
			ctx, int32(ledger.Actor_ADMIN), int32(ledger.Actor_HOSPITAL),
		)
		if err != nil {
			return nil, err
		}
	} else {
		err = auth.AuthenticateRequest(ctx)
	}
	if err != nil {
		return nil, err
	}

	// Check AccountID is provided
	accountID := getReq.GetAccountId()
	if strings.TrimSpace(accountID) == "" {
		return nil, errs.MissingCredential("AccountID")
	}

	// Get account from database
	accountDB := &Account{}
	if getReq.WithNationalId {
		if getReq.Privileged {
			err = accountAPI.sqlDB.Unscoped().First(accountDB, "national_id=?", accountID).Error
		} else {
			err = accountAPI.sqlDB.First(accountDB, "national_id=?", accountID).Error
		}
	} else {
		if getReq.Privileged {
			err = accountAPI.sqlDB.Unscoped().First(accountDB, "account_id=?", accountID).Error
		} else {
			err = accountAPI.sqlDB.First(accountDB, "account_id=?", accountID).Error
		}
	}

	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.AccountDoesntExist(accountID)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	if !getReq.Privileged {
		// Account should not be blocked
		if accountDB.AccountState == account.AccountState_BLOCKED.String() {
			return nil, errs.AccountBlocked()
		}
	}

	accountPB, err := getAccountPB(accountDB)
	if err != nil {
		return nil, err
	}

	return getAccountPBView(accountPB, getReq.GetView()), nil
}

func (accountAPI *accountAPIServer) Exist(
	ctx context.Context, existReq *account.ExistRequest,
) (*account.ExistResponse, error) {
	// Request should not be nil
	if existReq == nil {
		return nil, errs.NilObject("ExistRequest")
	}

	// Authenticate the request
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	var (
		email      = existReq.GetEmail()
		phone      = existReq.GetPhone()
		nationalID = existReq.GetNationalId()
	)

	// Validation
	if email == "" && phone == "" && nationalID == "" {
		return nil, errs.MissingCredential("Email, Phone Or NationalID")
	}

	accountDB := &Account{}

	// Query for account with email or phone
	err = accountAPI.sqlDB.Select("email,phone,national_id").
		First(accountDB, "national_id=? OR email=? OR phone=?", nationalID, email, phone).Error
	switch {
	case err == nil:
		// They exist
		return &account.ExistResponse{
			Exists: true,
		}, nil
	case gorm.IsRecordNotFoundError(err):
		// They don't exist
		return &account.ExistResponse{
			Exists: false,
		}, nil
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}
}

func (accountAPI *accountAPIServer) RequestChangePrivateAccount(
	ctx context.Context, req *account.RequestChangePrivateAccountRequest,
) (*account.RequestChangePrivateAccountResponse, error) {
	// Request must not be nil
	if req == nil {
		return nil, errs.NilObject("ListRequest")
	}

	// Authentication
	err := auth.AuthenticateRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case strings.TrimSpace(req.Payload) == "":
		err = errs.MissingCredential("Payload")
	case strings.TrimSpace(req.FallbackUrl) == "":
		err = errs.MissingCredential("Fallback URL")
	}
	if err != nil {
		return nil, err
	}

	// Get the user from database
	accountDB := &Account{}
	err = accountAPI.sqlDB.Find(accountDB, "email=? OR phone=?", req.Payload, req.Payload).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.AccountDoesntExist(req.Payload)
	default:
		return nil, errs.SQLQueryFailed(err, "FIND")
	}

	jwtToken, err := auth.GenToken(ctx, nil, auth.PatientGroup, int64(6*time.Hour))
	if err != nil {
		return nil, errs.FailedToGenToken(err)
	}

	uniqueNumber := rand.Intn(6)

	// Set token with expiration of 6 hours
	err = accountAPI.redisDB.Set(accountDB.AccountID, uniqueNumber, time.Duration(time.Hour*6)).Err()
	if err != nil {
		return nil, errs.RedisCmdFailed(err, "SET")
	}

	subject := "Reset Account Password"
	notificationPB := &notification.Notification{
		NotificationId: uuid.New().String(),
		OwnerIds:       []string{req.Payload},
		Priority:       notification.Priority_HIGH,
		Content: &notification.NotificationContent{
			Subject: subject,
			Data:    fmt.Sprintf("Reset token: %d", uniqueNumber),
		},
		CreateTimeSec: time.Now().Unix(),
	}
	switch {
	case strings.Contains(req.Payload, "@"):
		content := bytes.NewBuffer(make([]byte, 0, 64))
		accountAPI.tpl.ExecuteTemplate(content, templateName, &template_util.EmailData{
			FirstName: accountDB.FirstName,
			LastName:  accountDB.LastName,
			AccountID: accountDB.AccountID,
			Token:     jwtToken,
			Link: fmt.Sprintf(
				"%s?token=%s&account_id=%s&passphrase=%d",
				req.FallbackUrl, jwtToken, accountDB.AccountID, uniqueNumber,
			),
			AppName:      accountAPI.appName,
			TemplateName: "reset",
		})
		// Email notification
		notificationPB.Payload = &notification.Notification_EmailNotification{
			EmailNotification: &notification.EmailNotification{
				To:              []string{req.Payload},
				Subject:         subject,
				BodyContentType: "text/html",
				Body:            content.String(),
			},
		}
	default:
		// Send sms notification
		notificationPB.Payload = &notification.Notification_SmsNotification{
			SmsNotification: &notification.SMSNotification{
				DestinationPhone: []string{req.Payload},
				Keyword:          "Reset Password",
				Message:          notificationPB.Content.Data,
			},
		}
	}
	// Send notification
	_, err = accountAPI.notificationClient.Send(ctx, notificationPB, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	method := func() string {
		if strings.Contains(req.Payload, "@") {
			return "email"
		}
		return "phone"
	}()

	return &account.RequestChangePrivateAccountResponse{
		ResponseMessage: fmt.Sprintf("Check your %s for reset token", method),
	}, nil
}

const defaultPageSize = 15

func normalizePage(pageToken, pageSize int32) (int, int) {
	if pageSize <= 0 {
		pageSize = 1
	}
	if pageToken <= 0 {
		pageToken = 1
	}
	return int(pageToken), int(pageSize)
}

func (accountAPI *accountAPIServer) ListAccounts(
	ctx context.Context, listReq *account.ListAccountsRequest,
) (*account.Accounts, error) {
	// Request should not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListRequest")
	}

	// Authenticate admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	db := generateWhereCondition(accountAPI.sqlDB, listReq.GetListCriteria())

	pageNumber, pageSize := normalizePage(listReq.GetPageToken(), listReq.GetPageSize())
	offset := (pageNumber * pageSize) - pageSize

	accountsDB := make([]*Account, 0, pageSize)

	err = db.Unscoped().Offset(offset).Limit(pageSize).
		Find(&accountsDB).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.SQLQueryNoRows(err)
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Populate response
	accountsPB := make([]*account.Account, 0, len(accountsDB))
	for _, accountDB := range accountsDB {
		accountPB, err := getAccountPB(accountDB)
		if err != nil {
			return nil, err
		}
		accountPB = getAccountPBView(accountPB, listReq.GetView())
		accountsPB = append(accountsPB, accountPB)
	}

	return &account.Accounts{
		NextPageToken: int32(pageNumber + 1),
		Accounts:      accountsPB,
	}, nil
}

// Searches for accounts
func (accountAPI *accountAPIServer) SearchAccounts(
	ctx context.Context, searchReq *account.SearchAccountsRequest,
) (*account.Accounts, error) {
	// Request should not be nil
	if searchReq == nil {
		return nil, errs.NilObject("SearchRequest")
	}

	// Authenticate admin
	_, err := auth.AuthenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// For empty queries
	if searchReq.Query == "" {
		return &account.Accounts{
			Accounts: []*account.Account{},
		}, nil
	}

	// If filter is present, double the page size
	if searchReq.GetSearchCriteria().GetFilter() {
		searchReq.PageSize *= 2
	}

	pageNumber, pageSize := normalizePage(searchReq.GetPageToken(), searchReq.GetPageSize())
	offset := (pageNumber * pageSize) - pageSize

	parsedQuery := db_util.ParseQuery(searchReq.Query)

	accountsDB := make([]*Account, 0, pageSize)

	logrus.Infoln("QUERY: ", parsedQuery)
	err = accountAPI.sqlDB.Unscoped().Offset(offset).Limit(pageSize).
		Find(&accountsDB, "MATCH(national_id, email, phone) AGAINST(?)", parsedQuery).
		Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "SELECT")
	}

	// Populate response
	accountsPB := make([]*account.Account, 0, len(accountsDB))
	for _, accountDB := range accountsDB {
		if !fullfillsCriteria(searchReq.SearchCriteria, accountDB) {
			continue
		}
		accountPB, err := getAccountPB(accountDB)
		if err != nil {
			return nil, err
		}
		accountPB = getAccountPBView(accountPB, searchReq.GetView())
		accountsPB = append(accountsPB, accountPB)
	}

	return &account.Accounts{
		NextPageToken: int32(pageNumber + 1),
		Accounts:      accountsPB,
	}, nil
}

func fullfillsCriteria(criteria *account.Criteria, accountDB *Account) bool {
	if criteria == nil || !criteria.Filter {
		return true
	}
	// Filter by account type
	switch {
	case criteria.ShowAdmins && criteria.ShowUsers:
	case !criteria.ShowAdmins && !criteria.ShowUsers:
	case criteria.ShowAdmins:
		if accountDB.AccountType == account.AccountType_USER_OWNER.String() {
			return false
		}
	case criteria.ShowUsers:
		if accountDB.AccountType != account.AccountType_USER_OWNER.String() {
			return false
		}
	}

	// Filter by gender
	switch {
	case criteria.ShowFemales && criteria.ShowMales:
	case !criteria.ShowFemales && !criteria.ShowMales:
	case criteria.ShowFemales:
		if accountDB.Gender != "female" {
			return false
		}
	case criteria.ShowMales:
		if accountDB.Gender != "male" {
			return false
		}
	}

	// Filter by account state
	switch {
	case criteria.ShowBlockedAccounts && criteria.ShowActiveAccounts && criteria.ShowInactiveAccounts:
	case !criteria.ShowBlockedAccounts && !criteria.ShowActiveAccounts && !criteria.ShowInactiveAccounts:
	case criteria.ShowBlockedAccounts:
		if accountDB.AccountState != account.AccountState_BLOCKED.String() {
			return false
		}
	case criteria.ShowActiveAccounts:
		if accountDB.AccountState != account.AccountState_ACTIVE.String() {
			return false
		}
	case criteria.ShowInactiveAccounts:
		if accountDB.AccountState != account.AccountState_INACTIVE.String() {
			return false
		}
	}

	// Filter by date
	if criteria.FilterCreationDate {
		createdAt := accountDB.CreatedAt.Unix()
		switch {
		case criteria.CreatedFrom > 0 && criteria.CreatedUntil > 0:
			if createdAt <= criteria.CreatedFrom || createdAt >= criteria.CreatedUntil {
				return false
			}
		case criteria.CreatedUntil > 0:
			if createdAt >= criteria.CreatedUntil {
				return false
			}
		case criteria.CreatedFrom > 0:
			if createdAt <= criteria.CreatedFrom {
				return false
			}
		}
	}

	// Filter by account_labels
	if criteria.FilterAccountLabels {
		var include bool
		for _, label := range criteria.GetLabels() {
			if strings.Contains(strings.ToUpper(accountDB.AccountLabels), strings.ToUpper(label)) {
				include = true
			}
		}
		if !include {
			return false
		}
	}

	return true
}

func generateWhereCondition(db *gorm.DB, criteria *account.Criteria) *gorm.DB {
	if criteria == nil || !criteria.Filter {
		return db
	}

	// Filter by account type
	switch {
	case criteria.ShowAdmins && criteria.ShowUsers:
	case !criteria.ShowAdmins && !criteria.ShowUsers:
	case criteria.ShowAdmins:
		db = db.Where("account_type != ?", account.AccountType_USER_OWNER.String())
	case criteria.ShowUsers:
		db = db.Where("account_type = ?", account.AccountType_USER_OWNER.String())
	}

	// Filter by gender
	switch {
	case criteria.ShowFemales && criteria.ShowMales:
	case !criteria.ShowFemales && !criteria.ShowMales:
	case criteria.ShowFemales:
		db = db.Where("gender = ?", "female")
	case criteria.ShowMales:
		db = db.Where("gender = ?", "male")
	}

	// Filter by account state
	accountStates := make([]string, 0, 3)
	if criteria.ShowBlockedAccounts {
		accountStates = append(accountStates, account.AccountState_BLOCKED.String())
	}
	if criteria.ShowActiveAccounts {
		accountStates = append(accountStates, account.AccountState_ACTIVE.String())
	}
	if criteria.ShowInactiveAccounts {
		accountStates = append(accountStates, account.AccountState_INACTIVE.String())
	}
	if len(accountStates) != 0 {
		db = db.Where("account_state IN (?)", accountStates)
	}

	// Filter by date
	if criteria.FilterCreationDate {
		switch {
		case criteria.CreatedFrom > 0 && criteria.CreatedUntil > 0:
			db = db.Where(
				"UNIX_TIMESTAMP(created_at) BETWEEN ? AND ?",
				criteria.CreatedFrom, criteria.CreatedUntil,
			)
		case criteria.CreatedUntil > 0:
			db = db.Where(
				"UNIX_TIMESTAMP(created_at) < ?", criteria.CreatedUntil,
			)
		case criteria.CreatedFrom > 0:
			db = db.Where(
				"UNIX_TIMESTAMP(created_at) > ?", criteria.CreatedFrom,
			)
		}
	}

	// Filter by account_labels
	if criteria.FilterAccountLabels {
		db = db.Where("account_labels IN (?)", criteria.GetLabels())
	}

	return db
}
