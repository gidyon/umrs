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

var _ = Describe("Admin activating an account #admin_activate", func() {
	var (
		adminActivateReq          *account.ChangeAccountRequest
		ctx                       context.Context
		superAdminOwnerAndActive  string
		superAdminViewerAndActive string
		superAdminOwnerNotActive  string
	)

	BeforeEach(func() {
		adminActivateReq = &account.ChangeAccountRequest{
			AccountId:    uuid.New().String(),
			SuperAdminId: uuid.New().String(),
		}
		ctx = auth.AddSuperAdminMD(context.Background())
	})

	Describe("Calling AdminActivate with nil or malformed request", func() {
		It("should fail when the request is nil", func() {
			adminActivateReq = nil
			adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(adminActivateRes).Should(BeNil())
		})
		It("should fail when the super admin id is missing", func() {
			adminActivateReq.SuperAdminId = ""
			adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(adminActivateRes).Should(BeNil())
		})
		It("should fail when the account id is missing", func() {
			adminActivateReq.AccountId = ""
			adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(adminActivateRes).Should(BeNil())
		})
	})

	Describe("Creating Super Admins", func() {
		It("should create super admin", func() {
			var err error
			superAdminViewerAndActive, err = createAdmin(
				account.AccountType_ADMIN_VIEWER, account.AccountState_ACTIVE,
			)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should create super admin", func() {
			var err error
			superAdminOwnerNotActive, err = createAdmin(
				account.AccountType_ADMIN_OWNER, account.AccountState_INACTIVE,
			)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should create super admin", func() {
			var err error
			superAdminOwnerAndActive, err = createAdmin(
				account.AccountType_ADMIN_OWNER, account.AccountState_ACTIVE,
			)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("Calling AdminActivate with incorrect super admin id or account id", func() {
		Context("Lets create an account first", func() {
			var accountID string
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

			Context("When the super admin id is incorrect", func() {
				It("should fail when the super admin id is incorrect", func() {
					adminActivateReq.AccountId = accountID
					adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.NotFound))
					Expect(adminActivateRes).Should(BeNil())
				})
				It("should fail when the account id is incorrect", func() {
					adminActivateReq.SuperAdminId = superAdminOwnerAndActive
					adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.NotFound))
					Expect(adminActivateRes).Should(BeNil())
				})
				It("should succeed when the account id and super admin id is correct", func() {
					adminActivateReq.SuperAdminId = superAdminOwnerAndActive
					adminActivateReq.AccountId = accountID
					adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(adminActivateRes).ShouldNot(BeNil())
				})
			})
		})
	})

	Describe("Calling AdminActivate account on existing account", func() {

		Context("Lets create an account first", func() {
			var accountID string
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

			Describe("Admin activating the account", func() {
				Context("Admin activating account when the the super admin type is not OWNER", func() {
					It("should fail because the super admin type is not OWNER", func() {
						adminActivateReq.AccountId = accountID
						adminActivateReq.SuperAdminId = superAdminViewerAndActive
						adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
						Expect(err).Should(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
						Expect(adminActivateRes).Should(BeNil())
					})

					Describe("Let's get the account", func() {
						It("should succeed but account state still INACTIVE", func() {
							getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
								AccountId: accountID,
							})
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.OK))
							Expect(getRes).ShouldNot(BeNil())
							Expect(getRes.AccountState).Should(Equal(account.AccountState_INACTIVE))
						})
					})
				})

				Context("When the the super admin account state is not ACTIVE", func() {
					It("should fail because the super admin state is INACTIVE", func() {
						adminActivateReq.AccountId = accountID
						adminActivateReq.SuperAdminId = superAdminOwnerNotActive
						adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
						Expect(err).Should(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
						Expect(adminActivateRes).Should(BeNil())
					})

					Describe("Let's get the account", func() {
						It("should succeed but account state still INACTIVE", func() {
							getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
								AccountId: accountID,
							})
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.OK))
							Expect(getRes).ShouldNot(BeNil())
							Expect(getRes.AccountState).Should(Equal(account.AccountState_INACTIVE))
						})
					})
				})

				Context("When the the super admin type is OWNER and account state ACTIVE", func() {
					It("should succeed because the super admin state is ACTIVE and type OWNER", func() {
						adminActivateReq.AccountId = accountID
						adminActivateReq.SuperAdminId = superAdminOwnerAndActive
						adminActivateRes, err := AccountAPI.AdminActivate(ctx, adminActivateReq)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.OK))
						Expect(adminActivateRes).ShouldNot(BeNil())
					})

					Describe("Let's get the account", func() {
						It("should succeed and account state is ACTIVE", func() {
							getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
								AccountId: accountID,
							})
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.OK))
							Expect(getRes).ShouldNot(BeNil())
							Expect(getRes.AccountState).Should(Equal(account.AccountState_ACTIVE))
						})
					})
				})
			})
		})
	})
})
