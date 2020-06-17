package patient

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/chaincodes/patient/mocks"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"

	"github.com/gidyon/umrs/pkg/api/patient"

	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestPatient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Patient Suite")
}

const dbAddress = "localhost"

var (
	DB               *gorm.DB
	PatientAPI       patient.PatientAPIServer
	PatientAPIServer *patientAPIServer
	err              error
)

func initDB() (*gorm.DB, error) {
	param := "charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("root:hakty11@tcp(%s:3306)/umrs?%s", dbAddress, param)
	return gorm.Open("mysql", dsn)
}

var _ = BeforeSuite(func() {
	var err error
	// Start real database
	DB, err = initDB()
	Expect(err).ShouldNot(HaveOccurred())

	// Create ledger mock
	ledgerAPI := &mocks.ledgerAPIMock{}
	ledgerAPI.On("AddBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.AddBlockResponse{Hash: uuid.New().String()}, nil)
	ledgerAPI.On("GetBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Block{Hash: uuid.New().String()}, nil)
	ledgerAPI.On("ListBlocks", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Blocks{NextPageNumber: 3, Blocks: []*ledger.Block{
			newBlock(), newBlock(), newBlock(), newBlock(),
		}}, nil)
	ledgerAPI.On("RegisterContract", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.RegisterContractResponse{ContractId: uuid.New().String()}, nil)

	opt := &Options{
		ContractID:       uuid.New().String(),
		SQLDB:            DB,
		ledgerClient: ledgerAPI,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	PatientAPI, err = NewPatientChaincode(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	PatientAPIServer, ok = PatientAPI.(*patientAPIServer)
	Expect(ok).Should(BeTrue())

	// When creating a patient server with bad credentials
	_, err = NewPatientChaincode(nil, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = nil
	_, err = NewPatientChaincode(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = DB
	opt.ledgerClient = nil
	_, err = NewPatientChaincode(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.ledgerClient = ledgerAPI
	opt.ContractID = ""
	_, err = NewPatientChaincode(ctx, opt)
	Expect(err).Should(HaveOccurred())

	rand.Seed(time.Now().UnixNano())
})

func newBlock() *ledger.Block {
	return &ledger.Block{
		Timestamp: time.Now().Unix(),
		Hash:      uuid.New().String(),
		PrevHash:  uuid.New().String(),
		Payload: &ledger.Transaction{
			Operation:    newOperation(),
			Creator:      newActorPayload(),
			Patient:      newActorPayload(),
			Organization: newActorPayload(),
		},
	}
}

var _ = AfterSuite(func() {
	// Allow for events to happen
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
