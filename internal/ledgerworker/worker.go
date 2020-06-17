package ledgerworker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/encryption"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/grpclog"
	"time"
)

const (
	logsList   = "logs:list"
	backupList = "logs:backup:list"
)

type ledgerWorkerAPI struct {
	sqlDB         *gorm.DB
	orderer       *redis.Client
	logger        grpclog.LoggerV2
	encryptionAPI encryption.Interface
	lastLog       *Log
	popTimeout    time.Duration
}

// Options contains parameters for starting worker
type Options struct {
	DB            *gorm.DB
	Orderer       *redis.Client
	Logger        grpclog.LoggerV2
	EncryptionAPI encryption.Interface
	PopTimeout    time.Duration
}

// StartWorker starts ledger API worker
func StartWorker(ctx context.Context, opt *Options) error {
	var err error
	switch {
	case ctx == nil:
		err = errors.New("context must not be nil")
	case opt.DB == nil:
		err = errors.New("sqlDB must not be nil")
	case opt.Orderer == nil:
		err = errors.New("orderer must not be nil")
	case opt.Logger == nil:
		err = errors.New("logger must not be nil")
	case opt.EncryptionAPI == nil:
		err = errors.New("encryptionAPI must not be nil")
	}
	if err != nil {
		return err
	}

	worker := &ledgerWorkerAPI{
		sqlDB:         opt.DB,
		orderer:       opt.Orderer,
		logger:        opt.Logger,
		encryptionAPI: opt.EncryptionAPI,
		popTimeout:    opt.PopTimeout,
		lastLog:       &Log{},
	}

	return worker.start(ctx)
}

func (worker *ledgerWorkerAPI) failOperationAndRollback(data string, errMsg string) {
	worker.rollbackData(data)
}

func (worker *ledgerWorkerAPI) rollbackData(data string) {
	worker.orderer.LPush(backupList, data)
	// Update long running operation
}

func (worker *ledgerWorkerAPI) start(ctx context.Context) error {
	for {
	selectLabel:
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			data, err := worker.orderer.BRPop(worker.popTimeout, LogsList).Result()
			switch {
			case err == nil:
			case errors.Is(err, redis.Nil):
			default:
				worker.logger.Errorf("error during BRPOP operation: %v", err)
				break selectLabel
			}

			if len(data) == 2 {
				worker.addLog(data[1])
			}
		}
	}
}

func (worker *ledgerWorkerAPI) addLog(data string) {
	bs, err := worker.encryptionAPI.Decrypt([]byte(data))
	if err != nil {
		worker.failOperationAndRollback(data, fmt.Sprintf("failed to decrypt Log data: %w", err))
		return
	}

	logPB := &ledger.Log{}
	err = proto.Unmarshal(bs, logPB)
	if err != nil {
		worker.failOperationAndRollback(data, fmt.Sprintf("failed to unmarshal proto message: %w", err))
		return
	}

	logDB, err := GetLogDB(logPB)
	if err != nil {
		worker.failOperationAndRollback(data, err.Error())
		return
	}

	logDB.LogHash = GenerateHash(logDB)

	lastLog := worker.lastLog

	logDB.PrevHash = worker.lastLog.LogHash

	worker.lastLog = logDB

	err = worker.sqlDB.Create(logDB).Error
	if err != nil {
		worker.lastLog = lastLog
		worker.failOperationAndRollback(data, fmt.Sprintf("failed to add to database: %v", err))
		return
	}
}

// GenerateHash generates hash of a block
func GenerateHash(logDB *Log) string {
	record := fmt.Sprintf(
		"%d%s%d%s%s%s",
		logDB.Timestamp,
		logDB.PrevHash,
		logDB.LogIndex,
		logDB.Payload.CreatorID,
		logDB.Payload.PatientID,
		logDB.Payload.OrganizationID,
	)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
