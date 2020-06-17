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

var _ = Describe("Restoring a deleted account #undelete", func() {
	var (
		undeleteReq               *account.ChangeAccountRequest
		ctx                       context.Context
		superAdminOwnerAndActive  string
		superAdminViewerAndActive string
		superAdminOwnerNotActive  string
	)

	BeforeEach(func() {
		undeleteReq = &account.ChangeAccountRequest{
			AccountId:    uuid.New().String(),
			SuperAdminId: uuid.New().String(),
		}
		ctx = auth.AddSuperAdminMD(context.Background())
	})

	Describe("Calling Undelete with nil or malformed request", func() {
		It("should fail when the request is nil", func() {
			undeleteReq = nil
			blockRes, err := AccountAPI.Undelete(ctx, undeleteReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
		})
		It("should fail when the super admin id is missing", func() {
			undeleteReq.SuperAdminId = ""
			blockRes, err := AccountAPI.Undelete(ctx, undeleteReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(blockRes).Should(BeNil())
		})
		It("should fail when the account id is missing", func() {
			undeleteReq.AccountId = ""
			blockRes, err := AccountAPI.Undelete(ctx, undeleteReq)
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

	Describe("Calling undelete account on existing account", func() {

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

			Describe("Deleting an account", func() {
				It("should delete the account", func() {
					ctx = auth.AddAdminMD(context.Background())
					undeleteReq.AccountId = accountID
					undeleteReq.SuperAdminId = superAdminOwnerAndActive
					blockRes, err := AccountAPI.Delete(ctx, &account.DeleteRequest{
						AccountId: accountID,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(blockRes).ShouldNot(BeNil())
				})

				Describe("Restoring deleted the account", func() {
					Context("Restoring deleted account when the the super admin type is not OWNER", func() {
						It("should fail because the super admin type is not OWNER", func() {
							undeleteReq.AccountId = accountID
							undeleteReq.SuperAdminId = superAdminViewerAndActive
							blockRes, err := AccountAPI.Undelete(ctx, undeleteReq)
							Expect(err).Should(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
							Expect(blockRes).Should(BeNil())
						})

						Describe("Let's get the account", func() {
							It("should fail because the account is deleted", func() {
								getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
									AccountId: accountID,
								})
								Expect(err).Should(HaveOccurred())
								Expect(status.Code(err)).Should(Equal(codes.NotFound))
								Expect(getRes).Should(BeNil())
							})
						})
					})

					Context("When the the super admin account state is not ACTIVE", func() {
						It("should fail because the super admin state is INACTIVE", func() {
							undeleteReq.AccountId = accountID
							undeleteReq.SuperAdminId = superAdminOwnerNotActive
							blockRes, err := AccountAPI.Undelete(ctx, undeleteReq)
							Expect(err).Should(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.PermissionDenied))
							Expect(blockRes).Should(BeNil())
						})

						Describe("Let's get the account", func() {
							It("should fail because the account is still deleted", func() {
								getRes, err := AccountAPI.Get(ctx, &account.GetRequest{
									AccountId: accountID,
								})
								Expect(err).Should(HaveOccurred())
								Expect(status.Code(err)).Should(Equal(codes.NotFound))
								Expect(getRes).Should(BeNil())
							})
						})
					})

					Context("When the the super admin type is OWNER and account state ACTIVE", func() {
						Describe("Lets try to undelete the account", func() {
							It("should succeed because the super admin state is ACTIVE and type OWNER", func() {
								undeleteReq.AccountId = accountID
								undeleteReq.SuperAdminId = superAdminOwnerAndActive
								blockRes, err := AccountAPI.Undelete(ctx, undeleteReq)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(status.Code(err)).Should(Equal(codes.OK))
								Expect(blockRes).ShouldNot(BeNil())
							})

							Describe("Let's get the account", func() {
								It("should succeed because the account has been undeleted", func() {
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
