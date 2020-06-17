package employment

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/gidyon/umrs/pkg/api/employment"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Checking employment status #status", func() {
	var (
		checkReq *employment.CheckEmploymentStatusRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		checkReq = &employment.CheckEmploymentStatusRequest{
			AccountId: uuid.New().String(),
			Actor:     newActor(),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	Describe("Checking employment status wil malformed request", func() {
		It("should fail when the request is nil", func() {
			checkReq = nil
			checkRes, err := EmploymentAPI.CheckEmploymentStatus(ctx, checkReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(checkRes).Should(BeNil())
		})
		It("should fail when the actor is changed", func() {
			checkReq.Actor.Actor = int32(ledger.Actor_INSURANCE)
			checkRes, err := EmploymentAPI.CheckEmploymentStatus(ctx, checkReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(checkRes).Should(BeNil())
		})
		It("should fail when account id is missing", func() {
			checkReq.AccountId = ""
			checkRes, err := EmploymentAPI.CheckEmploymentStatus(ctx, checkReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(checkRes).Should(BeNil())
		})
	})

	Describe("Checking employment status with well formed request", func() {
		var accountID string
		Context("Set employment status and verified as true", func() {
			It("should succeed when request is valid", func() {
				employmentPB := newEmployment()
				employmentDB, err := getEmploymentDB(employmentPB)
				Expect(err).ShouldNot(HaveOccurred())

				employmentDB.IsRecent = true
				employmentDB.EmploymentVerified = true
				employmentDB.StillEmployed = true
				err = DB.Create(employmentDB).Error
				Expect(err).ShouldNot(HaveOccurred())

				accountID = employmentDB.AccountID
			})
			Describe("Checking employment status", func() {
				It("should succeed when request is valid", func() {
					checkReq.AccountId = accountID
					checkRes, err := EmploymentAPI.CheckEmploymentStatus(ctx, checkReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(checkRes).ShouldNot(BeNil())
					Expect(checkRes.IsEmployed).Should(BeTrue())
					Expect(checkRes.IsVerified).Should(BeTrue())
				})
			})
		})
		Context("Set employment status as true and verified as false", func() {
			It("should succeed when request is valid", func() {
				employmentPB := newEmployment()
				employmentDB, err := getEmploymentDB(employmentPB)
				Expect(err).ShouldNot(HaveOccurred())

				employmentDB.IsRecent = true
				employmentDB.EmploymentVerified = false
				employmentDB.StillEmployed = true
				err = DB.Create(employmentDB).Error
				Expect(err).ShouldNot(HaveOccurred())

				accountID = employmentDB.AccountID
			})
			Describe("Checking employment status", func() {
				It("should succeed when request is valid", func() {
					checkReq.AccountId = accountID
					checkRes, err := EmploymentAPI.CheckEmploymentStatus(ctx, checkReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(checkRes).ShouldNot(BeNil())
					Expect(checkRes.IsEmployed).Should(BeTrue())
					Expect(checkRes.IsVerified).Should(BeFalse())
				})
			})
		})

		It("should succeed", func() {
			checkReq.AccountId = uuid.New().String()
			checkRes, err := EmploymentAPI.CheckEmploymentStatus(ctx, checkReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(checkRes).ShouldNot(BeNil())
			Expect(checkRes.IsEmployed).Should(BeFalse())
		})
	})
})
