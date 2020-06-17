package ledger

import (
	"context"
	"fmt"
	"github.com/gidyon/micros"
	"github.com/gidyon/umrs/internal/ledgerworker"
	"github.com/gidyon/umrs/internal/pkg/encryption"
	"math/rand"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"

	"github.com/gidyon/umrs/pkg/api/ledger"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestLedger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ledger Suite")
}

const (
	dbAddress  = "localhost"
	dbFilePath = "./db"
)

var (
	DB              *gorm.DB
	RedisDB         *redis.Client
	LedgerAPI       ledger.LedgerAPIServer
	LedgerAPIServer *ledgerServer
	err             error
	errChan         = make(chan error, 1)
)

func startSqlite() (*gorm.DB, error) {
	return gorm.Open("sqlite3", dbFilePath)
}

func staryMySQL() (*gorm.DB, error) {
	param := "charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("root:hakty11@tcp(%s:3306)/umrs?%s", dbAddress, param)
	return gorm.Open("mysql", dsn)
}

var _ = BeforeSuite(func() {
	// Start real database
	DB, err = staryMySQL()
	Expect(err).ShouldNot(HaveOccurred())

	DB.LogMode(true)

	// Connect to redis
	RedisDB = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:6379", dbAddress),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create encryption API
	encryptionAPI, err := encryption.NewInterface([]byte(randomdata.RandStringRunes(32)))
	Expect(err).ShouldNot(HaveOccurred())

	opt := &Options{
		SQLDB:         DB,
		RedisClient:   RedisDB,
		Logger:        micros.NewLogger("ledger_api"),
		EncryptionAPI: encryptionAPI,
	}

	LedgerAPI, err = NewledgerServer(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	LedgerAPIServer, ok = LedgerAPI.(*ledgerServer)
	Expect(ok).Should(BeTrue())

	// When creating a ledger server with missing credentials
	_, err = NewledgerServer(nil, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = nil
	_, err = NewledgerServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = DB
	opt.RedisClient = nil
	_, err = NewledgerServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.RedisClient = RedisDB
	opt.EncryptionAPI = nil
	_, err = NewledgerServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.EncryptionAPI = encryptionAPI
	opt.Logger = nil
	_, err = NewledgerServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.Logger = micros.NewLogger("legder_api")

	rand.Seed(time.Now().UnixNano())

	// Start ledger worker to adding blocks in atomic fashion
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		errChan <- ledgerworker.StartWorker(ctx, &ledgerworker.Options{
			DB:            opt.SQLDB,
			Orderer:       opt.RedisClient,
			Logger:        opt.Logger,
			EncryptionAPI: encryptionAPI,
			PopTimeout:    3 * time.Second,
		})
	}()
})

var _ = AfterSuite(func() {
	LedgerAPIServer.logger.Infoln("waiting for worker to complete")
	err := <-errChan
	Expect(err).Should(HaveOccurred())
	Expect(RedisDB.FlushAll().Err()).ShouldNot(HaveOccurred())
	Expect(DB.Close()).ShouldNot(HaveOccurred())
})

// Declarations for Ginkgo DSL
type Done ginkgo.Done
type Benchmarker ginkgo.Benchmarker

var GinkgoWriter = ginkgo.GinkgoWriter
var GinkgoRandomSeed = ginkgo.GinkgoRandomSeed
var GinkgoParallelNode = ginkgo.GinkgoParallelNode
var GinkgoT = ginkgo.GinkgoT
var CurrentGinkgoTestDescription = ginkgo.CurrentGinkgoTestDescription
var RunSpecs = ginkgo.RunSpecs
var RunSpecsWithDefaultAndCustomReporters = ginkgo.RunSpecsWithDefaultAndCustomReporters
var RunSpecsWithCustomReporters = ginkgo.RunSpecsWithCustomReporters
var Skip = ginkgo.Skip
var Fail = ginkgo.Fail
var GinkgoRecover = ginkgo.GinkgoRecover
var Describe = ginkgo.Describe
var FDescribe = ginkgo.FDescribe
var PDescribe = ginkgo.PDescribe
var XDescribe = ginkgo.XDescribe
var Context = ginkgo.Context
var FContext = ginkgo.FContext
var PContext = ginkgo.PContext
var XContext = ginkgo.XContext
var When = ginkgo.When
var FWhen = ginkgo.FWhen
var PWhen = ginkgo.PWhen
var XWhen = ginkgo.XWhen
var It = ginkgo.It
var FIt = ginkgo.FIt
var PIt = ginkgo.PIt
var XIt = ginkgo.XIt
var Specify = ginkgo.Specify
var FSpecify = ginkgo.FSpecify
var PSpecify = ginkgo.PSpecify
var XSpecify = ginkgo.XSpecify
var By = ginkgo.By
var Measure = ginkgo.Measure
var FMeasure = ginkgo.FMeasure
var PMeasure = ginkgo.PMeasure
var XMeasure = ginkgo.XMeasure
var BeforeSuite = ginkgo.BeforeSuite
var AfterSuite = ginkgo.AfterSuite
var SynchronizedBeforeSuite = ginkgo.SynchronizedBeforeSuite
var SynchronizedAfterSuite = ginkgo.SynchronizedAfterSuite
var BeforeEach = ginkgo.BeforeEach
var JustBeforeEach = ginkgo.JustBeforeEach
var JustAfterEach = ginkgo.JustAfterEach
var AfterEach = ginkgo.AfterEach

// Declarations for Gomega DSL
var RegisterFailHandler = gomega.RegisterFailHandler
var RegisterFailHandlerWithT = gomega.RegisterFailHandlerWithT
var RegisterTestingT = gomega.RegisterTestingT
var InterceptGomegaFailures = gomega.InterceptGomegaFailures
var Ω = gomega.Ω
var Expect = gomega.Expect
var ExpectWithOffset = gomega.ExpectWithOffset
var Eventually = gomega.Eventually
var EventuallyWithOffset = gomega.EventuallyWithOffset
var Consistently = gomega.Consistently
var ConsistentlyWithOffset = gomega.ConsistentlyWithOffset
var SetDefaultEventuallyTimeout = gomega.SetDefaultEventuallyTimeout
var SetDefaultEventuallyPollingInterval = gomega.SetDefaultEventuallyPollingInterval
var SetDefaultConsistentlyDuration = gomega.SetDefaultConsistentlyDuration
var SetDefaultConsistentlyPollingInterval = gomega.SetDefaultConsistentlyPollingInterval
var NewWithT = gomega.NewWithT
var NewGomegaWithT = gomega.NewGomegaWithT

// Declarations for Gomega Matchers
var Equal = gomega.Equal
var BeEquivalentTo = gomega.BeEquivalentTo
var BeIdenticalTo = gomega.BeIdenticalTo
var BeNil = gomega.BeNil
var BeTrue = gomega.BeTrue
var BeFalse = gomega.BeFalse
var HaveOccurred = gomega.HaveOccurred
var Succeed = gomega.Succeed
var MatchError = gomega.MatchError
var BeClosed = gomega.BeClosed
var Receive = gomega.Receive
var BeSent = gomega.BeSent
var MatchRegexp = gomega.MatchRegexp
var ContainSubstring = gomega.ContainSubstring
var HavePrefix = gomega.HavePrefix
var HaveSuffix = gomega.HaveSuffix
var MatchJSON = gomega.MatchJSON
var MatchXML = gomega.MatchXML
var MatchYAML = gomega.MatchYAML
var BeEmpty = gomega.BeEmpty
var HaveLen = gomega.HaveLen
var HaveCap = gomega.HaveCap
var BeZero = gomega.BeZero
var ContainElement = gomega.ContainElement
var BeElementOf = gomega.BeElementOf
var ConsistOf = gomega.ConsistOf
var HaveKey = gomega.HaveKey
var HaveKeyWithValue = gomega.HaveKeyWithValue
var BeNumerically = gomega.BeNumerically
var BeTemporally = gomega.BeTemporally
var BeAssignableToTypeOf = gomega.BeAssignableToTypeOf
var Panic = gomega.Panic
var BeAnExistingFile = gomega.BeAnExistingFile
var BeARegularFile = gomega.BeARegularFile
var BeADirectory = gomega.BeADirectory
var And = gomega.And
var SatisfyAll = gomega.SatisfyAll
var Or = gomega.Or
var SatisfyAny = gomega.SatisfyAny
var Not = gomega.Not
var WithTransform = gomega.WithTransform
