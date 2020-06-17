package insurance

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/chaincodes/insurance/mocks"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"

	"github.com/gidyon/umrs/pkg/api/insurance"

	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestInsurance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Insurance Suite")
}

const (
	dbAddress           = "localhost"
	defualtTemplateFile = "/home/gideon/go/src/github.com/gidyon/umrs/internal/chaincodes/insurance/templates/add.html"
)

var (
	DB                 *gorm.DB
	InsuranceAPI       insurance.InsuranceAPIServer
	InsuranceAPIServer *insuranceAPIServer
	err                error
)

func initDB() (*gorm.DB, error) {
	param := "charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("root:hakty11@tcp(%s:3306)/umrs?%s", dbAddress, param)
	return gorm.Open("mysql", dsn)
}

var _ = BeforeSuite(func() {
	err = os.Setenv("ADD_INSURANCE_TEMPLATE_FILE", defualtTemplateFile)
	Expect(err).ShouldNot(HaveOccurred())

	err = os.Setenv("UPDATE_INSURANCE_TEMPLATE_FILE", defualtTemplateFile)
	Expect(err).ShouldNot(HaveOccurred())

	err = os.Setenv("DELETE_INSURANCE_TEMPLATE_FILE", defualtTemplateFile)
	Expect(err).ShouldNot(HaveOccurred())

	// Start real database
	DB, err = initDB()
	Expect(err).ShouldNot(HaveOccurred())

	DB.LogMode(true)

	// Create ledger mock
	ledgerAPI := &mocks.ledgerAPIMock{}
	ledgerAPI.On("AddBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.AddBlockResponse{Hash: uuid.New().String()}, nil)
	ledgerAPI.On("GetBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Block{Hash: uuid.New().String()}, nil)
	ledgerAPI.On("ListBlocks", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Blocks{NextPageNumber: 3, Blocks: []*ledger.Block{}}, nil)
	ledgerAPI.On("RegisterContract", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.RegisterContractResponse{ContractId: uuid.New().String()}, nil)

	// Notification mock
	notificationAPI := new(mocks.NotificationAPIMock)
	notificationAPI.On("CreateNotificationAccount", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	notificationAPI.On("Send", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	notificationAPI.On("ChannelSend", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)

	opt := &Options{
		ContractID:         uuid.New().String(),
		SQLDB:              DB,
		ledgerClient:   ledgerAPI,
		NotificationClient: notificationAPI,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	InsuranceAPI, err = NewInsuranceAPIServer(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	InsuranceAPIServer, ok = InsuranceAPI.(*insuranceAPIServer)
	Expect(ok).Should(BeTrue())

	// When creating a patient server with bad credentials
	_, err = NewInsuranceAPIServer(nil, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = nil
	_, err = NewInsuranceAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = DB
	opt.ledgerClient = nil
	_, err = NewInsuranceAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.ledgerClient = ledgerAPI
	opt.ContractID = ""
	_, err = NewInsuranceAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.ContractID = uuid.New().String()
	opt.NotificationClient = nil
	_, err = NewInsuranceAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.NotificationClient = notificationAPI
})

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
