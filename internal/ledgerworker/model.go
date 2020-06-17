package ledgerworker

import (
	"github.com/gidyon/umrs/internal/pkg/errs"
	"time"

	"github.com/gidyon/umrs/pkg/api/ledger"
)

const tableName = "logs"

const (
	// LogsTable is table for logs
	LogsTable = "logs"
	// LogsList is logs list
	LogsList = "logs:store"
	// LogsBackup is logs backup
	LogsBackup = "logs:backup"
	// LedgerStats is key containing ledger stats
	LedgerStats = "stats:ledger"
)

// Log is model transaction log
type Log struct {
	LogIndex  int64    `gorm:"auto_increment;"`
	LogHash   string   `gorm:"primary_key;type:varchar(256);not null"`
	PrevHash  string   `gorm:"type:varchar(256);not null"`
	Timestamp int64    `gorm:"type:int(32);not null"`
	Payload   *payload `gorm:"embedded;embedded_prefix:tx_"`
	CreatedAt time.Time
	DeletedAt *time.Time
}

// TableName refers to the name of the table
func (*Log) TableName() string {
	return LogsTable
}

type payload struct {
	Operation        int32  `gorm:"index:bq_index;type:tinyint(1)"`
	Creator          int32  `gorm:"index:bq_index;type:tinyint(1)"`
	CreatorID        string `gorm:"index:bq_index;type:varchar(50);not null"`
	CreatorName      string `gorm:"type:varchar(50);not null"`
	PatientID        string `gorm:"index:bq_index;type:varchar(50);not null"`
	PatientName      string `gorm:"type:varchar(50);not null"`
	OrganizationID   string `gorm:"index:bq_index;type:varchar(50);not null"`
	OrganizationName string `gorm:"type:varchar(50);not null"`
	Details          []byte `gorm:"type:blob;not null"`
}

// GetLogDB converts block protobuf message into a db model
func GetLogDB(blockPB *ledger.Log) (*Log, error) {
	if blockPB == nil {
		return nil, errs.NilObject("Log")
	}

	payloadPB := blockPB.GetPayload()
	if payloadPB == nil {
		return nil, errs.NilObject("Payload")
	}

	blockDB := &Log{
		LogHash:   blockPB.Hash,
		PrevHash:  blockPB.PrevHash,
		Timestamp: blockPB.Timestamp,
		Payload: &payload{
			Operation:        int32(payloadPB.GetOperation()),
			Creator:          int32(payloadPB.GetCreator().GetActor()),
			CreatorName:      payloadPB.GetCreator().GetActorNames(),
			CreatorID:        payloadPB.GetCreator().GetActorId(),
			PatientID:        payloadPB.GetPatient().GetActorId(),
			PatientName:      payloadPB.GetPatient().GetActorNames(),
			OrganizationID:   payloadPB.GetOrganization().GetActorId(),
			OrganizationName: payloadPB.GetOrganization().GetActorNames(),
			Details:          payloadPB.GetDetails(),
		},
	}

	return blockDB, nil
}

// GetLogPB returns the block as a protobuf message
func GetLogPB(blockDB *Log) (*ledger.Log, error) {
	if blockDB == nil {
		return nil, errs.NilObject("Log")
	}
	if blockDB.Payload == nil {
		return nil, errs.NilObject("Payload")
	}

	blockPB := &ledger.Log{
		Timestamp: blockDB.Timestamp,
		Hash:      blockDB.LogHash,
		PrevHash:  blockDB.PrevHash,
		Payload: &ledger.Transaction{
			Operation: ledger.Operation(blockDB.Payload.Operation),
			Creator: &ledger.ActorPayload{
				Actor:      ledger.Actor(blockDB.Payload.Creator),
				ActorId:    blockDB.Payload.CreatorID,
				ActorNames: blockDB.Payload.CreatorName,
			},
			Patient: &ledger.ActorPayload{
				ActorId:    blockDB.Payload.PatientID,
				ActorNames: blockDB.Payload.PatientName,
			},
			Organization: &ledger.ActorPayload{
				ActorId:    blockDB.Payload.OrganizationID,
				ActorNames: blockDB.Payload.OrganizationName,
			},
			Details: blockDB.Payload.Details,
		},
	}

	return blockPB, nil
}
