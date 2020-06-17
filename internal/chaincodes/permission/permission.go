package permission

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/internal/pkg/md"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/permission"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"html/template"
	"os"
	"strings"
	"time"
)

const (
	expTimeForPermissionToken  = time.Duration(time.Hour * 6)
	expTimeForRequesterProfile = time.Duration(time.Hour * 12)
	expTimeForActiveSet        = time.Duration(time.Hour * 24)
)

const (
	failedToRequestPermission    = "Failed to request permission"
	failedToGrantPermission      = "Failed to grant permission"
	failedToRevokePermission     = "Failed to revoke permission"
	failedToGetPermission        = "Failed to get permission token"
	failedToGetActivePermissions = "Failed to get active permission profiles"
)

type permissionAPIServer struct {
	baseURL              string
	redisClient          *redis.Client
	ledgerClient     ledger.ledgerClient
	notificationClient   notification.NotificationServiceClient
	tplRequestPermission *template.Template
}

// Options contains parameters for NewPermissionAPI
type Options struct {
	ContractID         string
	RedisClient        *redis.Client
	ledgerClient   ledger.ledgerClient
	NotificationClient notification.NotificationServiceClient
}

// NewPermissionAPI creates a permission API singleton service
func NewPermissionAPI(
	ctx context.Context, opt *Options,
) (permission.PatientPermissionAPIServer, error) {
	var err error
	// Validation
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case strings.TrimSpace(opt.ContractID) == "":
		return nil, errs.MissingCredential("ContractId")
	case opt.RedisClient == nil:
		return nil, errs.NilObject("RedisClient")
	case opt.ledgerClient == nil:
		return nil, errs.NilObject("ledgerClient")
	case opt.NotificationClient == nil:
		return nil, errs.NilObject("NotificationClient")
	}

	permisionAPI := &permissionAPIServer{
		redisClient:        opt.RedisClient,
		ledgerClient:   opt.ledgerClient,
		notificationClient: opt.NotificationClient,
	}

	// parse base URL from environment variables
	baseURL := os.Getenv("PERMISSION_BASE_URL")
	if strings.TrimSpace(baseURL) == "" {
		return nil, errs.MissingCredential("BaseUrl")
	}
	permisionAPI.baseURL = baseURL

	// parse <request email notification> email template
	tplFile := os.Getenv("REQUEST_ACCESS_TEMPLATE_FILE")
	if strings.TrimSpace(tplFile) == "" {
		return nil, errs.MissingCredential("template file")
	}

	permisionAPI.tplRequestPermission, err = template.ParseFiles(tplFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'request access' template file: %w", err)
	}

	// Register the chaincode to ledger server
	ctxReg := auth.AddSuperAdminMD(ctx)
	p, err := auth.AuthenticateSuperAdmin(ctxReg)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate super admin: %v", err)
	}

	_, err = permisionAPI.ledgerClient.RegisterContract(ctxReg, &ledger.RegisterContractRequest{
		SuperAdminId: p.ID, ContractId: opt.ContractID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register contract with the ledger server: %v", err)
	}

	return permisionAPI, nil
}

func getRequesterKey(requesterID string) string {
	return "requester:" + requesterID
}

func getPatientActiveRequesterSet(patientID string) string {
	return fmt.Sprintf("allowed_access:%s:%d", patientID, time.Now().Day())
}

func getTokenKey(patientID, requesterID string) string {
	return fmt.Sprintf("token:%s:%s", patientID, requesterID)
}

var patientGroup = int32(ledger.Actor_PATIENT)

func (permissionServer *permissionAPIServer) RequestPermissionToken(
	ctx context.Context, req *permission.RequestPermissionTokenRequest,
) (*permission.RequestPermissionTokenResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToRequestPermission)

	// Request must not be nil
	if req == nil {
		return nil, WrapError(errs.NilObject("RequestPermissionTokenRequest"))
	}

	// Validation
	var err error
	requesterProfile := req.GetRequesterProfile()
	patientActor := req.GetPatient()
	patientID := patientActor.GetId()
	requesterActor := req.GetRequester()
	requesterID := requesterActor.GetId()
	switch {
	case requesterProfile == nil:
		err = errs.NilObject("RequesterProfile")
	case strings.TrimSpace(requesterProfile.AccountId) == "":
		err = errs.MissingCredential("RequesterAccountId")
	case strings.TrimSpace(requesterProfile.FullName) == "":
		err = errs.MissingCredential("RequesterFullName")
	case strings.TrimSpace(requesterProfile.OrganizationName) == "":
		err = errs.MissingCredential("RequesterOrganizationName")
	case strings.TrimSpace(requesterProfile.OrganizationId) == "":
		err = errs.MissingCredential("RequesterOrganizationId")
	case strings.TrimSpace(requesterProfile.RoleAtOrganization) == "":
		err = errs.MissingCredential("RequesterRoleAtOrganization")
	case patientActor == nil:
		err = errs.NilObject("PatientActor")
	case strings.TrimSpace(patientID) == "":
		err = errs.MissingCredential("PatientId")
	case requesterActor == nil:
		err = errs.NilObject("RequesterActor")
	case strings.TrimSpace(requesterID) == "":
		err = errs.MissingCredential("RequesterID")
	case requesterActor.Group == int32(ledger.Actor_UNKNOWN):
		err = errs.ActorUknown()
	case req.GetPermissionMethod() == nil:
		err = errs.NilObject("PermissionMethod")
	case req.GetPermissionMethod().Method == permission.RequestPermissionMethod_UNKNOWN:
		err = errs.MissingCredential("RequestPermissionMethod")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndID(ctx, requesterActor.GetGroup(), requesterID)
	if err != nil {
		return nil, WrapError(err)
	}

	// Save requester profile to cache
	requesterBs, err := proto.Marshal(requesterProfile)
	if err != nil {
		return nil, WrapError(errs.FromProtoMarshal(err, "RequesterProfile"))
	}
	requesterKey := getRequesterKey(requesterID)
	err = permissionServer.redisClient.Set(
		requesterKey, requesterBs, expTimeForRequesterProfile,
	).Err()
	if err != nil {
		return nil, WrapError(errs.RedisCmdFailed(err, "SET"))
	}

	// Handle requesting for permission
	subject := "Permission To Access Your Data"
	notificationContent := &notification.NotificationContent{
		Subject: subject,
		Data: fmt.Sprintf(
			"%s from %s want to access your medical record",
			requesterProfile.FullName, requesterProfile.OrganizationName,
		),
	}
	notificationPB := &notification.Notification{
		NotificationId: uuid.New().String(),
		OwnerIds:       []string{patientID},
		Priority:       notification.Priority_HIGH,
		Content:        notificationContent,
		CreateTimeSec:  time.Now().Unix(),
		Save:           true,
	}
	permMethod := req.GetPermissionMethod()
	switch permMethod.GetMethod() {
	case permission.RequestPermissionMethod_EMAIL:
		acessToken := auth.GenPatientAccessToken(ctx, patientActor.GetId())
		queryParams := fmt.Sprintf(
			"requester.id=%s&requester.group=%d&requester.full_name=%s&organization.full_name=%s&patient.id=%s&patient.full_name=%s&reason=%s&authorization_token=%s",
			requesterID,
			requesterActor.GetGroup(),
			requesterActor.GetFullName(),
			requesterProfile.GetOrganizationName(),
			patientActor.GetId(),
			patientActor.GetFullName(),
			req.GetReason(),
			acessToken,
		)
		urlPath := fmt.Sprintf(
			"/api/umrs/permissions/patient/%s/grants/%s?%s",
			patientID, requesterProfile.AccountId, queryParams,
		)
		grantURL := fmt.Sprintf(
			"%s%s", strings.TrimSuffix(permissionServer.baseURL, "/"), urlPath,
		)
		// Template for sending email
		emailContent := emailData{
			RequesterProfile: requesterProfile,
			Reason:           req.Reason,
			GrantAccessURL:   grantURL,
		}
		content := bytes.NewBuffer(make([]byte, 0, 64))
		err = permissionServer.tplRequestPermission.Execute(content, emailContent)
		if err != nil {
			return nil, errs.FailedToExecuteTemplate(err)
		}
		notificationPB.SendMethod = notification.SendMethod_EMAIL
		notificationPB.Payload = &notification.Notification_EmailNotification{
			EmailNotification: &notification.EmailNotification{
				To:              []string{req.GetPermissionMethod().GetPayload()},
				Subject:         subject,
				BodyContentType: "text/html",
				Body:            content.String(),
			},
		}
		_, err = permissionServer.notificationClient.Send(
			md.AddFromCtx(ctx), notificationPB, grpc.WaitForReady(true),
		)
		if err != nil {
			return nil, WrapError(err)
		}
	case permission.RequestPermissionMethod_SMS:
		notificationPB.SendMethod = notification.SendMethod_SMS
		notificationPB.Payload = &notification.Notification_SmsNotification{
			SmsNotification: &notification.SMSNotification{
				Keyword:          "Permission Required",
				DestinationPhone: []string{req.GetPermissionMethod().GetPayload()},
				Message:          notificationContent.Data,
			},
		}
		_, err = permissionServer.notificationClient.Send(
			md.AddFromCtx(ctx), notificationPB, grpc.WaitForReady(true),
		)
		if err != nil {
			return nil, WrapError(err)
		}
	case permission.RequestPermissionMethod_USSD:
		notificationPB.SendMethod = notification.SendMethod_USSD
		notificationPB.Payload = &notification.Notification_UssdNotification{
			UssdNotification: &notification.USSDNotification{},
		}
		_, err = permissionServer.notificationClient.Send(
			md.AddFromCtx(ctx), notificationPB, grpc.WaitForReady(true),
		)
		if err != nil {
			return nil, WrapError(err)
		}
	case permission.RequestPermissionMethod_FINGERPRINT:
	case permission.RequestPermissionMethod_VOICE:
	case permission.RequestPermissionMethod_FACIAL:
	}

	return &permission.RequestPermissionTokenResponse{}, nil
}

func (permissionServer *permissionAPIServer) GrantPermissionToken(
	ctx context.Context, grantReq *permission.GrantPermissionTokenRequest,
) (*permission.GrantPermissionTokenResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToGrantPermission)

	// Request must not be nil
	if grantReq == nil {
		return nil, WrapError(errs.NilObject("GrantPermissionTokenRequest"))
	}

	// Validation
	patient := grantReq.GetPatient()
	patientID := patient.GetId()
	requester := grantReq.GetRequester()
	requesterID := requester.GetId()
	organization := grantReq.GetOrganization()
	organizationID := organization.GetId()
	var err error
	switch {
	case patient == nil:
		err = errs.NilObject("Patient")
	case requester == nil:
		err = errs.NilObject("Requester")
	case organization == nil:
		err = errs.NilObject("Organization")
	case strings.TrimSpace(requesterID) == "":
		err = errs.MissingCredential("RequesterId")
	case strings.TrimSpace(patientID) == "":
		err = errs.MissingCredential("PatientId")
	case strings.TrimSpace(organizationID) == "":
		err = errs.MissingCredential("OrganizationId")
	case strings.TrimSpace(grantReq.GetAuthorizationToken()) == "":
		err = errs.MissingCredential("AuthToken")
	case grantReq.GetRequester().GetGroup() == int32(ledger.Actor_UNKNOWN):
		err = errs.ActorUknown()
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndIDFromToken(
		grantReq.GetAuthorizationToken(), patientGroup, patientID,
	)
	if err != nil {
		return nil, WrapError(err)
	}

	tokenKey := getTokenKey(patientID, requesterID)
	setKey := getPatientActiveRequesterSet(patientID)

	// Start pipeline transaction to grant permision
	tx := permissionServer.redisClient.TxPipeline()
	err = tx.Set(tokenKey, grantReq.AuthorizationToken, expTimeForPermissionToken).Err()
	if err != nil {
		return nil, WrapError(err)
	}
	err = tx.SAdd(setKey, getRequesterKey(requesterID)).Err()
	if err != nil {
		return nil, WrapError(err)
	}
	data, err := tx.Get(getRequesterKey(requesterID)).Result()
	if err != nil {
		return nil, WrapError(err)
	}
	err = tx.Expire(setKey, expTimeForActiveSet).Err()
	if err != nil {
		return nil, WrapError(err)
	}

	_, err = tx.ExecContext(ctx)
	if err != nil {
		return nil, WrapError(err)
	}

	requesterProfile := &permission.BasicProfile{}
	err = proto.Unmarshal([]byte(data), requesterProfile)
	if err != nil {
		return nil, WrapError(err)
	}

	grantPayload := &permission.GrantPermissionTokenPayload{
		RequesterProfile: requesterProfile,
		Reason:           grantReq.GetPayload().GetReason(),
	}

	grantPayloadBs, err := proto.Marshal(grantPayload)
	if err != nil {
		return nil, WrapError(err)
	}

	ctx = auth.AddGroupAndIDMD(ctx, grantReq.GetRequester().GetGroup(), requesterID)

	// Add event log to ledger
	addRes, err := permissionServer.ledgerClient.AddBlock(ctx, &ledger.AddBlockRequest{
		Transaction: &ledger.Transaction{
			Operation: ledger.Operation_GRANT_PERMISSION,
			Creator: &ledger.ActorPayload{
				Actor:         ledger.Actor(grantReq.GetRequester().GetGroup()),
				ActorId:       requesterID,
				ActorNames: grantReq.GetRequester().GetFullName(),
			},
			Patient: &ledger.ActorPayload{
				Actor:         ledger.Actor(grantReq.GetPatient().GetGroup()),
				ActorId:       patientID,
				ActorNames: grantReq.GetPatient().GetFullName(),
			},
			Organization: &ledger.ActorPayload{
				Actor:         ledger.Actor(grantReq.GetPatient().GetGroup()),
				ActorId:       grantReq.GetOrganization().GetId(),
				ActorNames: grantReq.GetOrganization().GetFullName(),
			},
			Details: grantPayloadBs,
		},
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, WrapError(err)
	}

	return &permission.GrantPermissionTokenResponse{
		AllowedMessage: fmt.Sprintf(
			"Successfully allowed %s to access your data", grantReq.GetRequester().GetFullName(),
		),
		OperationHash: addRes.Hash,
	}, nil
}

func (permissionServer *permissionAPIServer) RevokePermissionToken(
	ctx context.Context, revokeReq *permission.RevokePermissionTokenRequest,
) (*permission.RevokePermissionTokenResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToRevokePermission)

	// Request must not be nil
	if revokeReq == nil {
		return nil, WrapError(errs.NilObject("RevokePermissionTokenRequest"))
	}

	// Validation
	var err error
	switch {
	case strings.TrimSpace(revokeReq.PatientId) == "":
		err = errs.MissingCredential("PatientId")
	case strings.TrimSpace(revokeReq.RequesterId) == "":
		err = errs.MissingCredential("RequesterId")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndID(ctx, patientGroup, revokeReq.GetPatientId())
	if err != nil {
		return nil, WrapError(err)
	}

	tokenKey := getTokenKey(revokeReq.PatientId, revokeReq.RequesterId)
	setKey := getPatientActiveRequesterSet(revokeReq.PatientId)

	// Revoke the requester by removing him from set and deleting token
	tx := permissionServer.redisClient.TxPipeline()
	err = tx.SRem(setKey, getRequesterKey(revokeReq.RequesterId)).Err()
	if err != nil {
		return nil, WrapError(err)
	}
	err = tx.Del(tokenKey).Err()
	if err != nil {
		return nil, WrapError(err)
	}

	_, err = tx.Exec()
	if err != nil {
		return nil, WrapError(err)
	}

	return &permission.RevokePermissionTokenResponse{}, nil
}

func (permissionServer *permissionAPIServer) GetPermissionToken(
	ctx context.Context, getReq *permission.GetPermissionTokenRequest,
) (*permission.GetPermissionTokenResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToGetPermission)

	// Request must not be nil
	if getReq == nil {
		return nil, WrapError(errs.NilObject("GetPermissionTokenRequest"))
	}

	// Validation
	actor := getReq.GetActor()
	var err error
	switch {
	case actor == nil:
		err = errs.NilObject("Actor")
	case strings.TrimSpace(actor.GetId()) == "":
		err = errs.MissingCredential("RequesterId")
	case strings.TrimSpace(getReq.PatientId) == "":
		err = errs.MissingCredential("PatientId")
	}
	if err != nil {
		return nil, WrapError(err)
	}

	// Authentication
	err = auth.AuthenticateGroupAndID(ctx, actor.GetGroup(), actor.GetId())
	if err != nil {
		return nil, WrapError(err)
	}

	// Get token from cache
	tokenKey := getTokenKey(getReq.PatientId, actor.GetId())
	token, err := permissionServer.redisClient.Get(tokenKey).Result()
	switch {
	case err == nil:
	case errors.Is(err, redis.Nil):
		return &permission.GetPermissionTokenResponse{
			Allowed: false,
		}, nil
	default:
		return nil, WrapError(err)
	}

	return &permission.GetPermissionTokenResponse{
		AccessToken: token,
		Allowed:     true,
	}, nil
}

func (permissionServer *permissionAPIServer) GetActivePermissions(
	ctx context.Context, getReq *permission.GetActivePermissionsRequest,
) (*permission.GetActivePermissionsResponse, error) {
	WrapError := errs.WrapErrorWithMsgFunc(failedToGetActivePermissions)

	// Request must not be nil
	if getReq == nil {
		return nil, WrapError(errs.NilObject("GetActivePermissionsRequest"))
	}

	// Validation
	switch {
	case strings.TrimSpace(getReq.PatientId) == "":
		return nil, WrapError(errs.MissingCredential("PatientId"))
	}

	// Authentication
	err := auth.AuthenticateGroupAndID(ctx, patientGroup, getReq.PatientId)
	if err != nil {
		return nil, WrapError(err)
	}

	// Get from set
	setKey := getPatientActiveRequesterSet(getReq.PatientId)

	allowedRequesters, err := permissionServer.redisClient.SMembers(setKey).Result()
	if err != nil {
		return nil, WrapError(err)
	}

	activeProfiles := make([]*permission.BasicProfile, 0, len(allowedRequesters))
	for _, allowedRequesterID := range allowedRequesters {
		requesterStr, err := permissionServer.redisClient.Get(allowedRequesterID).Result()
		if err != nil {
			continue
		}
		if requesterStr == "" {
			continue
		}

		requesterProfile := &permission.BasicProfile{}
		err = proto.Unmarshal([]byte(requesterStr), requesterProfile)
		if err != nil {
			continue
		}
		activeProfiles = append(activeProfiles, requesterProfile)
	}

	return &permission.GetActivePermissionsResponse{
		ActiveProfiles: activeProfiles,
	}, nil
}
