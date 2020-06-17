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

var _ = Describe("Unblocking an account #unblock", func() {
	var (
		unBlockReq                *account.ChangeAccountRequest
		ctx                       context.Context
		superAdminOwnerAndActive  string
		superAdminViewerAndActive string
		superAdminOwnerNotActive  string
	)

	BeforeEach(func() {
		unBlockReq = &account.ChangeAccountRequest{
			AccountId:    uuid.New().String(),
			SuperAdminId: uuid.New().String(),
		}
		ctx = auth.AddSuperAdminMD(context.Background())
	})

	Describe("Calling UnBlockAccount with nil or malformed request", func() {
		It("should fail when the request is nil", func() {
			unBlockReq = nil
			blockRes, err := AccountAPI.UnBlockAccount(ctx, unBlockReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
		})
		It("should fail when the super admin id is missing", func() {
			unBlockReq.SuperAdminId = ""
			blockRes, err := AccountAPI.UnBlockAccount(ctx, unBlockReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
		})
		It("should fail when the account id is missing", func() {
			unBlockReq.AccountId = ""
			blockRes, err := AccountAPI.UnBlockAccount(ctx, unBlockReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
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

	Describe("Calling unblock account on existing account", func() {

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

			Describe("Unblocking an account when the account is not blocked", func() {
				It("should fail to unblock the account because it is not blocked", func() {
					unBlockReq.AccountId = accountID
					unBlockReq.SuperAdminId = superAdminOwnerAndActive
					blockRes, err := AccountAPI.UnBlockAccount(ctx, unBlockReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
					Expect(blockRes).Should(BeNil())
				})

			})

			Describe("Blocking an account", func() {
				Context("Lets activate the account then block it, then unblock it", func() {
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

					Describe("Lets block the account for real", func() {
						It("should block the account", func() {
							unBlockReq.AccountId = accountID
							unBlockReq.SuperAdminId = superAdminOwnerAndActive
							blockRes, err := AccountAPI.BlockAccount(ctx, unBlockReq)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.OK))
							Expect(blockRes).ShouldNot(BeNil())
						})
					})

					Describe("Unblocking the account", func() {
						Context("Unblocking account when the the super admin type is not OWNER", func() {
							Describe("Lets try to unblock the account", func() {
								It("should fail because the super admin type is not OWNER", func() {
									unBlockReq.AccountId = accountID
									unBlockReq.SuperAdminId = superAdminViewerAndActive
									blockRes, err := AccountAPI.UnBlockAccount(ctx, unBlockReq)
									Expect(err).Should(HaveOccurred())
									Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
									Expect(blockRes).Should(BeNil())
								})

								Describe("Let's get the account", func() {
									It("should fail because the account is still blocked", func() {
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

						Context("When the the super admin account state is not ACTIVE", func() {
							Describe("Lets try to unblock the account", func() {
								It("should fail because the super admin state is INACTIVE", func() {
									unBlockReq.AccountId = accountID
									unBlockReq.SuperAdminId = superAdminOwnerNotActive
									blockRes, err := AccountAPI.UnBlockAccount(ctx, unBlockReq)
									Expect(err).Should(HaveOccurred())
									Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
									Expect(blockRes).Should(BeNil())
								})

								Describe("Let's get the account", func() {
									It("should fail because the account is still blocked", func() {
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

						Context("When the the super admin type is OWNER and account state ACTIVE", func() {
							Describe("Lets try to unblock the account", func() {
								It("should succeed because the super admin state is ACTIVE and type OWNER", func() {
									unBlockReq.AccountId = accountID
									unBlockReq.SuperAdminId = superAdminOwnerAndActive
									blockRes, err := AccountAPI.UnBlockAccount(ctx, unBlockReq)
									Expect(err).ShouldNot(HaveOccurred())
									Expect(status.Code(err)).Should(Equal(codes.OK))
									Expect(blockRes).ShouldNot(BeNil())
								})

								Describe("Let's get the account", func() {
									It("should succeed because the account has been unblocked", func() {
										getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
											AccountId: accountID,
										})
										Expect(err).ShouldNot(HaveOccurred())
										Expect(status.Code(err)).Should(Equal(codes.OK))
										Expect(getRes).ShouldNot(BeNil())
									})
								})
							})
						})
					})
				})
			})
		})
	})
})
