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

var _ = Describe("Change Account Type #change", func() {
	var (
		changeReq                 *account.ChangeAccountTypeRequest
		ctx                       context.Context
		superAdminOwnerAndActive  string
		superAdminViewerAndActive string
		superAdminOwnerNotActive  string
	)

	BeforeEach(func() {
		changeReq = &account.ChangeAccountTypeRequest{
			AccountId:    uuid.New().String(),
			SuperAdminId: uuid.New().String(),
			Type:         account.AccountType_ADMIN_EDITOR,
		}
		ctx = auth.AddSuperAdminMD(context.Background())
	})

	Context("Failure scenarios", func() {
		When("Changing account type", func() {
			It("should fail when the request is nil", func() {
				changeReq = nil
				changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(changeRes).Should(BeNil())
			})
			It("should fail when the account id is missing", func() {
				changeReq.AccountId = ""
				changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(changeRes).Should(BeNil())
			})
			It("should fail when the super admin id is missing", func() {
				changeReq.SuperAdminId = ""
				changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(changeRes).Should(BeNil())
			})
			It("should fail when the account id does not exist", func() {
				changeReq.AccountId = "dontexist"
				changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.NotFound))
				Expect(changeRes).Should(BeNil())
			})
			It("should fail when the super admin doesn't exist", func() {
				changeReq.SuperAdminId = "dontexist"
				changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.NotFound))
				Expect(changeRes).Should(BeNil())
			})
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

	Describe("Calling ChangeAccountType account on existing account", func() {

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

			Describe("Admin changing the account type", func() {
				Context("Admin changing account type when the the super admin type is not OWNER", func() {
					It("should fail because the super admin type is not OWNER", func() {
						changeReq.AccountId = accountID
						changeReq.SuperAdminId = superAdminViewerAndActive
						changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
						Expect(err).Should(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
						Expect(changeRes).Should(BeNil())
					})
				})

				Context("When the the super admin account state is not ACTIVE", func() {
					It("should fail because the super admin state is INACTIVE", func() {
						changeReq.AccountId = accountID
						changeReq.SuperAdminId = superAdminOwnerNotActive
						changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
						Expect(err).Should(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
						Expect(changeRes).Should(BeNil())
					})
				})

				Context("When the the super admin type is OWNER and account state ACTIVE", func() {
					var accountType account.AccountType
					It("should succeed because the super admin state is ACTIVE and type OWNER", func() {
						changeReq.AccountId = accountID
						changeReq.SuperAdminId = superAdminOwnerAndActive
						changeRes, err := AccountAPI.ChangeAccountType(ctx, changeReq)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.OK))
						Expect(changeRes).ShouldNot(BeNil())
						accountType = changeReq.Type
					})

					Context("Let's get the account", func() {
						It("should succeed and account type changed", func() {
							getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
								AccountId: accountID,
							})
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.OK))
							Expect(getRes).ShouldNot(BeNil())
							Expect(getRes.AccountType).Should(Equal(accountType))
						})
					})
				})
			})
		})
	})
})
