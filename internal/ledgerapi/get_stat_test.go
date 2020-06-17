package ledger

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/ledgerworker"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var _ = Describe("Getting Log from the ledger #getstat", func() {
	var (
		getReq *ledger.GetLedgerStatRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &ledger.GetLedgerStatRequest{}
		ctx = context.Background()
	})

	Describe("adding Log with malformed request", func() {
		It("should fail if the request is nil", func() {
			getReq = nil
			getRes, err := LedgerAPI.GetLedgerStat(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Getting Log with valid request", func() {
		It("should succeed", func() {
			// Lets add log to redis
			Expect(addStat()).ShouldNot(HaveOccurred())

			getRes, err := LedgerAPI.GetLedgerStat(ctx, getReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(getRes).ShouldNot(BeNil())
		})
	})
})

func addStat() error {
	// Create stat
	statPB := &ledger.LedgerStats{
		TotalTxLogs:              12,
		LastUpdatedTimestampSec:  time.Now().Unix(),
		LastVerifiedTimestampSec: time.Now().Unix(),
		LastInsertHash:           randomdata.RandStringRunes(32),
		Valid:                    false,
	}

	// Marshal
	statBuf, err := proto.Marshal(statPB)
	if err != nil {
		return err
	}

	// Encrypt
	statCipher, err := LedgerAPIServer.encryptionAPI.Encrypt(statBuf)
	if err != nil {
		return err
	}

	// Save to redis
	return RedisDB.Set(ledgerworker.LedgerStats, statCipher, 0).Err()
}
