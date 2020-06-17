package hospital

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/pkg/templateutil"
	"github.com/google/uuid"
	"os"
	"testing"

	"github.com/gidyon/umrs/pkg/api/hospital"

	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	// sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestHospital(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hospital Suite")
}

const (
	dbAddress           = "ec2-18-218-27-110.us-east-2.compute.amazonaws.com"
	dbName              = "umrs-testing"
	dbPort              = 30760
	defualtTemplateFile = "/home/gideon/go/src/github.com/gidyon/umrs/internal/chaincodes/hospital/templates/add.html"
)

var (
	DB                *gorm.DB
	HospitalAPI       hospital.HospitalAPIServer
	HospitalAPIServer *hospitalAPIServer
	err               error
)

func initDB() (*gorm.DB, error) {
	param := "charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("root:@@umrs2020@tcp(%s:%d)/%s?%s", dbAddress, dbPort, dbName, param)
	return gorm.Open("mysql", dsn)
}

var _ = BeforeSuite(func() {
	err = os.Setenv(templateutil.TemplateDirsEnv, "./templates")
	Expect(err).ShouldNot(HaveOccurred())

	// Start real database
	DB, err = initDB()
	Expect(err).ShouldNot(HaveOccurred())

	DB.LogMode(true)

	// Create ledger mock
	ledgerAPI := fakeledger()

	// Notification mock
	notificationAPI := fakeNotificationAPI()

	opt := &Options{
		ContractID:         uuid.New().String(),
		SQLDB:              DB,
		ledgerClient:   ledgerAPI,
		NotificationClient: notificationAPI,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	HospitalAPI, err = NewHospitalAPIServer(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	HospitalAPIServer, ok = HospitalAPI.(*hospitalAPIServer)
	Expect(ok).Should(BeTrue())

	// When creating a patient server with bad credentials
	_, err = NewHospitalAPIServer(nil, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = nil
	_, err = NewHospitalAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = DB
	opt.ledgerClient = nil
	_, err = NewHospitalAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.ledgerClient = ledgerAPI
	opt.ContractID = ""
	_, err = NewHospitalAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.ContractID = uuid.New().String()
	opt.NotificationClient = nil
	_, err = NewHospitalAPIServer(ctx, opt)
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
