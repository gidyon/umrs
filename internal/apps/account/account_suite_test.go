package account

import (
	"context"
	"fmt"
	"github.com/gidyon/umrs/internal/apps/account/mocks"
	"github.com/gidyon/umrs/pkg/api/account"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/gidyon/umrs/pkg/api/subscriber"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestAccount(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Suite")
}

var (
	DB               *gorm.DB
	RedisDB          *redis.Client
	AccountAPI       account.AccountAPIServer
	AccountAPIServer *accountAPIServer
	err              error
)

const (
	accountActivationURL = ""
	dbHost               = "localhost"
)

func initDB() (*gorm.DB, error) {
	param := "charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("root:hakty11@tcp(%s:3306)/umrs?%s", dbHost, param)
	return gorm.Open("mysql", dsn)
}

var _ = BeforeSuite(func() {
	rand.Seed(time.Now().UnixNano())

	// Testing templates
	os.Setenv(createAccTplFileEnv, "templates/email.html")
	os.Setenv(activateAccTplFileEnv, "templates/email.html")
	os.Setenv(blockAccTplFileEnv, "templates/email.html")
	os.Setenv(unBlockAccTplFileEnv, "templates/email.html")
	os.Setenv(deleteAccTplFileEnv, "templates/email.html")
	os.Setenv(unDeleteAccTplFileEnv, "templates/email.html")
	os.Setenv(changeAccTplFileEnv, "templates/email.html")
	os.Setenv(resetPasswordTplFileEnv, "templates/email.html")

	// Start real database
	DB, err = initDB()
	Expect(err).ShouldNot(HaveOccurred())

	DB.LogMode(true)

	RedisDB = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})

	// Mock notification API
	NotificationAPI := new(mocks.NotificationAPIMock)
	NotificationAPI.On("CreateNotificationAccount", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	NotificationAPI.On("Send", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)
	NotificationAPI.On("ChannelSend", mock.Anything, mock.Anything, mock.Anything).
		Return(&empty.Empty{}, nil)

	// Mock subscriber API
	SubscriberAPI := new(mocks.SubscriberAPIMock)
	SubscriberAPI.On("GetSendMethod", mock.Anything, mock.Anything, mock.Anything).
		Return(&subscriber.GetSendMethodResponse{
			SendMethod: notification.SendMethod_EMAIL_AND_SMS,
		}, nil)
	SubscriberAPI.On("GetSubscriber", mock.Anything, mock.Anything, mock.Anything).
		Return(&subscriber.Subscriber{AccountId: "me"}, nil)

	opt := &Options{
		SQLDB:              DB,
		RedisDB:            RedisDB,
		NotificationClient: NotificationAPI,
		SubscriberClient:   SubscriberAPI,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Inject stubs to the service
	AccountAPI, err = NewAccountAPI(ctx, "", opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	AccountAPIServer, ok = AccountAPI.(*accountAPIServer)
	Expect(ok).Should(BeTrue())

	// When creating singleton with bad credentialsopt.NotificationServer = NotificationAPI
	opt.SQLDB = nil
	_, err = NewAccountAPI(ctx, "", opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = DB
	opt.NotificationClient = nil
	_, err = NewAccountAPI(ctx, "", opt)
	Expect(err).Should(HaveOccurred())

	opt.NotificationClient = NotificationAPI
	opt.SubscriberClient = nil
	_, err = NewAccountAPI(ctx, "", opt)
	Expect(err).Should(HaveOccurred())

	opt.SubscriberClient = SubscriberAPI
})

var _ = AfterSuite(func() {
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
