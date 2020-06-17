package treatment

import (
	"context"
	"github.com/gidyon/umrs/internal/chaincodes/treatment/mocks"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/gidyon/umrs/pkg/api/treatment"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestTreatment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Treatment Suite")
}

const dbAddress = "localhost"

var (
	TreatmentAPI       treatment.TreatmentAPIServer
	TreatmentAPIServer *treatmentAPIServer
	err                error
)

var _ = BeforeSuite(func() {
	var err error

	// Create ledger mock
	ledgerAPI := &mocks.ledgerAPIMock{}
	ledgerAPI.On("AddBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.AddBlockResponse{Hash: uuid.New().String()}, nil)
	ledgerAPI.On("GetBlock", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Block{Hash: uuid.New().String(), Payload: &ledger.Transaction{
			Operation: ledger.Operation_ADD_PATIENT_TREATMENT_RECORD,
		}}, nil)
	ledgerAPI.On("ListBlocks", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.Blocks{NextPageNumber: 3, Blocks: []*ledger.Block{}}, nil)
	ledgerAPI.On("RegisterContract", mock.Anything, mock.Anything, mock.Anything).
		Return(&ledger.RegisterContractResponse{ContractId: uuid.New().String()}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opt := &Options{
		ledgerClient: ledgerAPI,
		ContractID:       uuid.New().String(),
	}

	TreatmentAPI, err = NewTreatmentChaincode(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	TreatmentAPIServer, ok = TreatmentAPI.(*treatmentAPIServer)
	Expect(ok).Should(BeTrue())

	// When creating a treatment server with bad credentials
	_, err = NewTreatmentChaincode(ctx, nil)
	Expect(err).Should(HaveOccurred())

	_, err = NewTreatmentChaincode(nil, opt)
	Expect(err).Should(HaveOccurred())

	opt.ledgerClient = nil
	_, err = NewTreatmentChaincode(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.ledgerClient = ledgerAPI
	opt.ContractID = ""
	_, err = NewTreatmentChaincode(ctx, opt)
	Expect(err).Should(HaveOccurred())
})

var _ = AfterSuite(func() {
	// Allow for events to happen
	// Expect(DB.Close()).ShouldNot(HaveOccurred())
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
