package notification

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/account/pkg/api"
	"github.com/gidyon/umrs/internal/apps/notification/notification/mocks"
	"github.com/gidyon/umrs/pkg/api/notification"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestNotification(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Notification Suite")
}

var (
	DB              *gorm.DB
	RedisDB         *redis.Client
	NotificationAPI notification.NotificationServiceServer
	err             error
)

func startDB() (*gorm.DB, error) {
	param := "charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("root:hakty11@tcp(localhost:3306)/umrs?%s", param)
	return gorm.Open("mysql", dsn)
}

var _ = BeforeSuite(func() {
	// Start real database
	DB, err = startDB()
	Expect(err).ShouldNot(HaveOccurred())

	DB.LogMode(true)

	RedisDB = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})

	// Mock email dialer
	mailAPIMock := new(mocks.Dialer)
	mailAPIMock.On("DialAndSend", mock.Anything).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opt := &Options{
		SQLDB:       DB,
		RedisClient: RedisDB,
		SMTPDialer:  mailAPIMock,
	}
	// Inject stubs to the service
	NotificationAPI, err = NewNotificationServiceServer(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	// Failing code
	_, err = NewNotificationServiceServer(nil, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = nil
	_, err = NewNotificationServiceServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = DB
	opt.RedisClient = nil
	_, err = NewNotificationServiceServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.RedisClient = RedisDB
	opt.SMTPDialer = nil
	_, err = NewNotificationServiceServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SMTPDialer = mailAPIMock
})

var _ = AfterSuite(func() {
	// The worker should exit by then
	time.Sleep(2 * time.Second)
	Expect(DB.Close()).ShouldNot(HaveOccurred())
})

// creates a new admin creds
func mockeAdminCreds() *account.AdminCreds {
	return &account.AdminCreds{
		AdminId: randomdata.RandStringRunes(32),
		Level:   account.AdminLevel_OWNER,
	}
}

func mockNotification() *notification.Notification {

	content := &notification.NotificationContent{
		Subject: "Testing",
		Data:    randomdata.Paragraph(),
	}

	return &notification.Notification{
		NotificationId: uuid.New().String(),
		OwnerIds:       []string{uuid.New().String()},
		Priority:       notification.Priority_MEDIUM,
		SendMethod:     notification.SendMethod_EMAIL,
		Content:        content,
		CreateTimeSec:  int64(time.Now().Unix()),
		Payload: &notification.Notification_EmailNotification{
			EmailNotification: mockEmailNotification(),
		},
		Bulk: false,
	}
}

func mockEmailNotification() *notification.EmailNotification {
	content := &notification.NotificationContent{
		Subject: "Testing",
		Data:    randomdata.Paragraph(),
	}
	return &notification.EmailNotification{
		From: "awesomeness@smtp.com",
		To: []string{
			randomdata.Email(), randomdata.Email(), randomdata.Email(),
		},
		Subject:         content.Subject,
		BodyContentType: "text/html",
		Body:            content.Data,
	}
}

func mockSMSNotification() *notification.SMSNotification {
	content := &notification.NotificationContent{
		Subject: "Testing",
		Data:    randomdata.Paragraph(),
	}
	return &notification.SMSNotification{
		Keyword: content.Subject,
		DestinationPhone: []string{
			randomdata.PhoneNumber(), randomdata.PhoneNumber(), randomdata.PhoneNumber(),
		},
		Message: content.Data,
	}
}

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
