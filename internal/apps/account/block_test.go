package account

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/account"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Blocking account #block", func() {
	var (
		blockReq *account.ChangeAccountRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		blockReq = &account.ChangeAccountRequest{
			AccountId:    uuid.New().String(),
			SuperAdminId: uuid.New().String(),
		}
		ctx = auth.AddSuperAdminMD(context.Background())
	})

	Describe("Calling BlockAccount with nil or malformed request", func() {
		It("should fail when the request is nil", func() {
			blockReq = nil
			blockRes, err := AccountAPI.BlockAccount(ctx, blockReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
		})
		It("should fail when the super admin id is missing", func() {
			blockReq.SuperAdminId = ""
			blockRes, err := AccountAPI.BlockAccount(ctx, blockReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
		})
		It("should fail when the account id is missing", func() {
			blockReq.AccountId = ""
			blockRes, err := AccountAPI.BlockAccount(ctx, blockReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
		})
	})

	Describe("Calling block account on existing account", func() {
		var superAdminViewerAndActive, superAdminOwnerNotActive, superAdminOwnerAndActive, accountID string
		Describe("Creating the super admin and user account", func() {
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
				accountID = createRes.AccountId
			})
			Describe("Creating super admins", func() {
				It("should create viewer admin", func() {
					var err error
					superAdminViewerAndActive, err = createAdmin(
						account.AccountType_ADMIN_VIEWER, account.AccountState_ACTIVE,
					)
					Expect(err).ShouldNot(HaveOccurred())
				})
				It("should create inactive super admin", func() {
					var err error
					superAdminOwnerNotActive, err = createAdmin(
						account.AccountType_ADMIN_OWNER, account.AccountState_INACTIVE,
					)
					Expect(err).ShouldNot(HaveOccurred())
				})
				It("should create active super admin", func() {
					var err error
					superAdminOwnerAndActive, err = createAdmin(
						account.AccountType_ADMIN_OWNER, account.AccountState_ACTIVE,
					)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		Context("Blocking account when the the super admin type is not OWNER", func() {
			It("should fail because the super admin type is not OWNER", func() {
				blockReq.AccountId = accountID
				blockReq.SuperAdminId = superAdminViewerAndActive
				blockRes, err := AccountAPI.BlockAccount(ctx, blockReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
				Expect(blockRes).Should(BeNil())
			})

			Describe("Let's get the account", func() {
				It("should succeed because the account is not blocked", func() {
					getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
						AccountId: accountID,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())
				})
			})
		})

		Context("When the the super admin account state is not ACTIVE", func() {
			It("should fail because the super admin state is INACTIVE", func() {
				blockReq.AccountId = accountID
				blockReq.SuperAdminId = superAdminOwnerNotActive
				blockRes, err := AccountAPI.BlockAccount(ctx, blockReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
				Expect(blockRes).Should(BeNil())
			})

			Describe("Let's get the account", func() {
				It("should succeed because the account is not blocked", func() {
					getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
						AccountId: accountID,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())
				})
			})
		})

		Context("When the the super admin type is OWNER and account state ACTIVE", func() {
			It("should succeed because account state is ACTIVE", func() {
				blockReq.AccountId = accountID
				blockReq.SuperAdminId = superAdminOwnerAndActive
				blockRes, err := AccountAPI.BlockAccount(ctx, blockReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
				Expect(blockRes).Should(BeNil())
			})

			Describe("Lets activate the account in order to block it", func() {
				It("should succeed because the super admin state is ACTIVE and type OWNER", func() {
					adminActivateReq := &account.ChangeAccountRequest{
						AccountId:    accountID,
						SuperAdminId: superAdminOwnerAndActive,
					}
					adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(adminActivateRes).ShouldNot(BeNil())
				})
			})

			It("should succeed because the account state is ACTIVE", func() {
				blockReq.AccountId = accountID
				blockReq.SuperAdminId = superAdminOwnerAndActive
				blockRes, err := AccountAPI.BlockAccount(ctx, blockReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(blockRes).ShouldNot(BeNil())
			})

			Describe("Let's get the account", func() {
				It("should fail because the account is blocked", func() {
					getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
						AccountId: accountID,
					})
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
					Expect(getRes).Should(BeNil())
				})
			})
		})
	})
})
