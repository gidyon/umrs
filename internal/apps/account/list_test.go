package account

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func newCriteria() *account.Criteria {
	return &account.Criteria{
		Filter:               true,
		ShowAdmins:           false,
		ShowUsers:            false,
		ShowFemales:          false,
		ShowMales:            false,
		ShowActiveAccounts:   false,
		ShowBlockedAccounts:  false,
		ShowInactiveAccounts: false,
		FilterCreationDate:   false,
	}
}

var _ = Describe("Listing accounts #list", func() {
	var (
		listReq *account.ListAccountsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &account.ListAccountsRequest{
			PageToken:    0,
			ListCriteria: newCriteria(),
			View:         account.AccountView_FULL_VIEW,
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Calling ListAccounts with nil or malformed request", func() {
		It("should fail when request is nil", func() {
			listReq = nil
			listRes, err := AccountAPI.ListAccounts(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
		It("should fail when request does not originate from admin", func() {
			ctx = auth.AddHospitalMD(context.Background())
			listRes, err := AccountAPI.ListAccounts(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("Calling ListAccounts with correct request payload", func() {
		Context("Lets create one account first", func() {
			It("should create the account without any error", func() {
				createReq := &account.CreateRequest{
					AccountLabel:   randomdata.Adjective(),
					Account:        fakeAccount(),
					PrivateAccount: fakePrivateAccount(),
				}
				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
			})

			Describe("Calling ListAccounts", func() {
				It("should succeed", func() {
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
				})
			})

			Describe("Calling ListAccounts with date filter on", func() {
				It("should succeed", func() {
					listReq.ListCriteria.FilterCreationDate = true
					listReq.ListCriteria.CreatedFrom = time.Now().UnixNano()
					listReq.ListCriteria.CreatedUntil = time.Now().UnixNano() / 2
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
				})
			})

			Describe("Calling ListAccounts with show_admins = true", func() {
				It("should succeed and returns only admin users", func() {
					listReq.ListCriteria.ShowAdmins = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountType).ShouldNot(Equal(account.AccountType_USER_OWNER))
					}
				})
			})

			Describe("Calling ListAccounts with show_users = true", func() {
				It("should succeed and returns only normal users", func() {
					listReq.ListCriteria.ShowUsers = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountType).Should(Equal(account.AccountType_USER_OWNER))
					}
				})
			})

			Describe("Calling ListAccounts with show_males = true", func() {
				It("should succeed and returns only male users", func() {
					listReq.ListCriteria.ShowMales = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.Gender).Should(Equal("male"))
					}
				})
			})

			Describe("Calling ListAccounts with show_females = true", func() {
				It("should succeed and returns only female users", func() {
					listReq.ListCriteria.ShowFemales = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.Gender).Should(Equal("female"))
					}
				})
			})

			Describe("Calling ListAccounts with show_active_accounts true", func() {
				It("should succeed and returns only ACTIVE accounts", func() {
					listReq.ListCriteria.ShowActiveAccounts = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountState).Should(Equal(account.AccountState_ACTIVE))
					}
				})
			})

			Describe("Calling ListAccounts with show_inactive_accounts = true", func() {
				It("should succeed and returns only INACTIVE accounts", func() {
					listReq.ListCriteria.ShowInactiveAccounts = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountState).Should(Equal(account.AccountState_INACTIVE))
					}
				})
			})

			Describe("Calling ListAccounts with show_blocked_accounts = true", func() {
				It("should succeed and returns only BLOCKED accounts", func() {
					listReq.ListCriteria.ShowBlockedAccounts = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountState).Should(Equal(account.AccountState_BLOCKED))
					}
				})
			})

			Describe("Calling ListAccounts with show_blocked_accounts and show_active_accounts true", func() {
				It("should succeed and returns only BLOCKED accounts", func() {
					listReq.ListCriteria.ShowBlockedAccounts = true
					listReq.ListCriteria.ShowActiveAccounts = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					arr := []account.AccountState{
						account.AccountState_BLOCKED, account.AccountState_ACTIVE,
					}
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountState).Should(BeElementOf(arr))
					}
				})
			})

			Describe("Calling ListAccounts with show_inactive_accounts and show_active_accounts true", func() {
				It("should succeed and returns only ACTIBE or INACTIVE accounts", func() {
					listReq.ListCriteria.ShowInactiveAccounts = true
					listReq.ListCriteria.ShowActiveAccounts = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					arr := []account.AccountState{
						account.AccountState_INACTIVE, account.AccountState_ACTIVE,
					}
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountState).Should(BeElementOf(arr))
					}
				})
			})

			Describe("Calling ListAccounts with show_blocked_accounts and show_inactive_accounts true", func() {
				It("should succeed and returns only BLOCKED or INACTIVE accounts", func() {
					listReq.ListCriteria.ShowBlockedAccounts = true
					listReq.ListCriteria.ShowInactiveAccounts = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
					arr := []account.AccountState{
						account.AccountState_BLOCKED, account.AccountState_INACTIVE,
					}
					for _, adminPB := range listRes.Accounts {
						Expect(adminPB.AccountState).Should(BeElementOf(arr))
					}
				})
			})

			Describe("Calling ListAccounts where all filters is true", func() {
				It("should succeed", func() {
					listReq.ListCriteria.ShowAdmins = true
					listReq.ListCriteria.ShowUsers = true
					listReq.ListCriteria.ShowBlockedAccounts = true
					listReq.ListCriteria.ShowActiveAccounts = true
					listReq.ListCriteria.ShowInactiveAccounts = true
					listReq.ListCriteria.ShowFemales = true
					listReq.ListCriteria.ShowMales = true
					listRes, err := AccountAPI.ListAccounts(ctx, listReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(listRes).ShouldNot(BeNil())
				})
			})
		})
	})
})
