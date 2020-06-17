package ledger

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/ledgerworker"
	"github.com/gidyon/umrs/internal/pkg/encryption"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"strings"
	"time"

	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ledgerServer struct {
	db            *gorm.DB
	orderer       *redis.Client
	encryptionAPI encryption.Interface
	logger        grpclog.LoggerV2
}

// Options contains parameters for the ledger server
type Options struct {
	SQLDB         *gorm.DB
	RedisClient   *redis.Client
	Logger        grpclog.LoggerV2
	EncryptionAPI encryption.Interface
}

// NewledgerServer creates a new ledger
func NewledgerServer(ctx context.Context, opt *Options) (ledger.LedgerAPIServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errs.NilObject("Context")
	case opt.RedisClient == nil:
		return nil, errs.NilObject("RedisClient")
	case opt.EncryptionAPI == nil:
		return nil, errs.NilObject("EncryptionAPI")
	case opt.SQLDB == nil:
		return nil, errs.NilObject("SqlDB")
	case opt.Logger == nil:
		return nil, errs.NilObject("Logger")
	}

	var err error

	// Create the ledger server
	api := &ledgerServer{
		db:            opt.SQLDB,
		orderer:       opt.RedisClient,
		logger:        opt.Logger,
		encryptionAPI: opt.EncryptionAPI,
	}

	// Automigration of logs table
	err = api.db.AutoMigrate(&ledgerworker.Log{}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate legder model: %w", err)
	}

	return api, nil
}

func validateActor(actor *ledger.ActorPayload) error {
	var err error
	switch {
	case actor == nil:
		err = errs.NilObject("actor")
	case actor.ActorId == "":
		err = errs.MissingCredential("actor id")
	case actor.ActorNames == "":
		err = errs.MissingCredential("actor names")
	case actor.Actor == ledger.Actor_UNKNOWN:
		err = errs.MissingCredential("actor category")
	}
	return err
}

func getActorOpList(actorID string) string {
	return "operations:" + actorID
}

func (api *ledgerServer) AddLog(
	ctx context.Context, addReq *ledger.AddLogRequest,
) (*ledger.AddLogResponse, error) {
	// Request must not be nil
	if addReq == nil {
		return nil, errs.NilObject("AddLogRequest")
	}

	// Validation
	tx := addReq.GetTransaction()
	var err error
	err = validateActor(tx.GetCreator())
	if err != nil {
		return nil, err
	}
	switch {
	case tx == nil:
		err = errs.NilObject("Transaction")
	case tx.GetDetails() == nil:
		err = errs.NilObject("Transaction.Details")
	case tx.GetOperation() == ledger.Operation_UKNOWN:
		err = errs.OperationUknown()
	}
	if err != nil {
		return nil, err
	}

	logPB := &ledger.Log{Payload: tx}

	// Marshal Log data
	logData, err := proto.Marshal(logPB)
	if err != nil {
		return nil, errs.FromProtoMarshal(err, "ledger.Log")
	}

	// Encrypt Log's data
	logCipher, err := api.encryptionAPI.Encrypt(logData)
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to encrypt log data")
	}

	// Create operation
	organizationDetails := fmt.Sprintf("[%s - %s]", tx.GetOrganization().ActorNames, tx.GetOrganization().ActorId)
	patientDetails := fmt.Sprintf("[%s - %s]", tx.GetPatient().ActorNames, tx.GetPatient().ActorId)
	operation := &ledger.AddOperation{
		Id: uuid.New().String(),
		Details: fmt.Sprintf(
			"operation: %s org:%s patient:%s", tx.GetOperation().String(), organizationDetails, patientDetails,
		),
		Status:       ledger.OperationStatus_PENDING,
		TimestampSec: time.Now().Unix(),
	}

	// Marshal operation
	operationData, err := proto.Marshal(operation)
	if err != nil {
		return nil, errs.FromProtoMarshal(err, "ledger.Operation")
	}

	// Encrypt operation
	operationCipher, err := api.encryptionAPI.Encrypt(operationData)
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to encrypt operation")
	}

	// Start transaction
	pipeline := api.orderer.TxPipeline()

	// Add operation to list of user operations
	err = pipeline.LPush(getActorOpList(tx.GetCreator().GetActorId()), operationCipher).Err()
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to add to operations")
	}

	// Add to orderer
	err = pipeline.LPush(ledgerworker.LogsList, logCipher).Err()
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to add to logs")
	}

	// Execute transaction
	_, err = pipeline.ExecContext(ctx)
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to execute transaction")
	}

	return &ledger.AddLogResponse{
		OperationId: operation.Id,
	}, nil
}

func (api *ledgerServer) GetLog(
	ctx context.Context, getReq *ledger.GetLogRequest,
) (*ledger.Log, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetLogRequest")
	}

	// Validation
	if strings.TrimSpace(getReq.Hash) == "" {
		return nil, errs.MissingCredential("Hash")
	}

	// Get log from state database
	logDB := &ledgerworker.Log{}
	err := api.db.First(logDB, "log_hash=?", getReq.Hash).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		return nil, errs.LogNotFound(getReq.Hash)
	default:
		return nil, errs.FailedToGetLog(err)
	}

	// Get Log pb
	logPB, err := ledgerworker.GetLogPB(logDB)
	if err != nil {
		return nil, err
	}

	return logPB, nil
}

const (
	maxOrgsInQuery       = 5
	maxOperationsInQuery = 5
	maxActorsInQuery     = 3
	maxPageSize          = 500
	defaultPageSize      = 10
)

func getCriteriaFromFilter(filter *ledger.Filter, db *gorm.DB) *gorm.DB {
	if filter == nil || !filter.GetFilter() {
		return db
	}

	// By Timestamp filter
	if filter.GetByDate() {
		db = db.Where("created_at BETWEEN ? AND ?", filter.GetStartDate(), filter.GetEndDate())
	}

	// By Creator type
	if filter.GetByCreatorActor() {
		actors := filter.GetCreatorActors()
		actorsInt := make([]int, 0, len(actors))
		for _, actor := range actors {
			actorsInt = append(actorsInt, int(actor))
		}
		db = db.Where("tx_creator IN(?)", actorsInt)
	}

	// By Operation
	if filter.GetByOperation() {
		ops := filter.GetOperations()
		opsInt := make([]int, 0, len(ops))
		for _, op := range ops {
			opsInt = append(opsInt, int(op))
		}
		db = db.Where("tx_operation IN(?)", opsInt)
	}

	// By Creator Id
	if filter.GetByCreatorId() {
		creatorIDs := filter.GetCreatorIds()
		db = db.Where("tx_creator_id IN(?)", creatorIDs)
	}

	// By Organization
	if filter.GetByOrganizationId() {
		orgs := filter.GetOrganizationIds()
		db = db.Where("tx_organization_id IN(?)", orgs)
	}

	// By Patient Id
	if filter.GetByPatientId() {
		patientIDs := filter.GetPatientIds()
		db = db.Where("tx_patient_id IN(?)", patientIDs)
	}

	return db
}

func normalizePageSize(pageSize int32) int32 {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = 500
	}
	return pageSize
}

func (api *ledgerServer) ListLogs(
	ctx context.Context, listReq *ledger.ListLogsRequest,
) (*ledger.Logs, error) {
	// Request must not be nil
	if listReq == nil {
		return nil, errs.NilObject("ListLogsRequest")
	}

	// Parse page size
	pageSize := normalizePageSize(listReq.GetPageSize())

	blocksDB := make([]*ledgerworker.Log, 0, pageSize)

	db := getCriteriaFromFilter(listReq.GetFilter(), api.db)

	err := db.Unscoped().Limit(pageSize).Order("created_at, timestamp").
		Where("log_hash > ?", listReq.PageToken).
		Find(&blocksDB).Error
	switch {
	case err == nil:
	default:
		return nil, errs.FailedToGetLogs(err)
	}

	var pageToken int64

	// Populate response
	blocksPB := make([]*ledger.Log, 0, len(blocksDB))
	for _, blockDB := range blocksDB {
		blockPB, err := ledgerworker.GetLogPB(blockDB)
		if err != nil {
			return nil, err
		}
		blocksPB = append(blocksPB, blockPB)
		pageToken = blockDB.LogIndex
	}

	return &ledger.Logs{
		Logs:          blocksPB,
		NextPageToken: int32(pageToken),
	}, nil
}

func (api *ledgerServer) GetLedgerStat(
	ctx context.Context, getReq *ledger.GetLedgerStatRequest,
) (*ledger.LedgerStats, error) {
	// Request must not be nil
	if getReq == nil {
		return nil, errs.NilObject("GetStatRequest")
	}

	// Get from redis
	statsStr, err := api.orderer.Get(ledgerworker.LedgerStats).Result()
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to get logs")
	}

	// Decrypt content
	statsData, err := api.encryptionAPI.Decrypt([]byte(statsStr))
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to decrypt stats")
	}

	// Proto unmarshal
	statPB := &ledger.LedgerStats{}
	err = proto.Unmarshal(statsData, statPB)
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.Internal, err, "failed to proto unmarshal stats")
	}

	return statPB, nil
}
