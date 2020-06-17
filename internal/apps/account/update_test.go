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

var _ = Describe("Update Account #update", func() {
	Describe("Update Method #updateaccount", func() {
		var (
			updateReq *account.UpdateRequest
			ctx       context.Context
		)

		BeforeEach(func() {
			updateReq = &account.UpdateRequest{
				AccountId: uuid.New().String(),
				Account:   fakeAccount(),
			}
			// Disable some fields
			updateReq.Account.Nationality = ""
			updateReq.Account.BirthDate = ""

			ctx = auth.AddAdminMD(context.Background())
		})

		Describe("Updating account with malformed request", func() {
			It("should fail when request is nil", func() {
				updateReq = nil
				updateRes, err := AccountAPI.Update(ctx, updateReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(updateRes).Should(BeNil())
			})
			It("should definitely fail when account id is missing", func() {
				updateReq.AccountId = ""
				updateRes, err := AccountAPI.Update(ctx, updateReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(updateRes).Should(BeNil())
			})
			It("should definitely fail when account id is incorrect", func() {
				updateReq.AccountId = "omen"
				updateRes, err := AccountAPI.Update(ctx, updateReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.NotFound))
				Expect(updateRes).Should(BeNil())
			})
			It("should fail when account is nil", func() {
				updateReq.Account = nil
				updateRes, err := AccountAPI.Update(ctx, updateReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(updateRes).Should(BeNil())
			})
		})

		Context("Updating account with valid request", func() {
			var (
				accountID   string
				accountType account.AccountType
			)

			Describe("Create account first", func() {
				It("should create account without error", func() {
					createReq := &account.CreateRequest{
						AccountLabel:   randomdata.Adjective(),
						Account:        fakeAccount(),
						PrivateAccount: fakePrivateAccount(),
					}
					// Create user account
					createReq.Account.AccountType = account.AccountType_USER_OWNER
					createRes, err := AccountAPI.Create(ctx, createReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(createRes).ShouldNot(BeNil())
					accountID = createRes.AccountId
					accountType = createReq.Account.AccountType
				})
			})

			It("should update account in database without error", func() {
				updateReq.AccountId = accountID
				// Set the account state to active
				updateReq.Account.AccountState = account.AccountState_ACTIVE
				updateReq.Account.AccountType = account.AccountType_ADMIN_EDITOR
				updateRes, err := AccountAPI.Update(ctx, updateReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(updateRes).ShouldNot(BeNil())
			})

			Describe("Getting the updated account", func() {
				It("should not be possible to have updated account state or account type", func() {
					getReq := &account.GetRequest{
						AccountId: accountID,
					}
					getReq.AccountId = accountID
					getRes, err := AccountAPI.Get(ctx, getReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(getRes).ShouldNot(BeNil())

					// The account state should remain unchanged
					Expect(getRes.AccountState).Should(Equal(account.AccountState_INACTIVE))
					// The account type should remain unchanged
					Expect(getRes.AccountType).Should(Equal(accountType))
				})
			})
		})
	})
})

var _ = Describe("UpdatePrivate Method #updateprivate", func() {
	var (
		updateReq *account.UpdatePrivateRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &account.UpdatePrivateRequest{
			AccountId:      uuid.New().String(),
			PrivateAccount: fakePrivateAccount(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Updating account private profile with nil request", func() {
		It("should fail when request is nil", func() {
			updateReq = nil
			updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should definitely fail when account id is missing", func() {
			updateReq.AccountId = ""
			updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should definitely fail when account id is incorrect", func() {
			updateReq.AccountId = "omen"
			updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when account is nil", func() {
			updateReq.PrivateAccount = nil
			updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
	})

	Describe("Create account first", func() {
		var (
			accountPB *account.Account
			accountID string
			token     string
		)
		It("should create account without error", func() {
			accountPB = fakeAccount()
			createRes, err := AccountAPI.Create(ctx, &account.CreateRequest{
				AccountLabel:   randomdata.Adjective(),
				Account:        accountPB,
				PrivateAccount: fakePrivateAccount(),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.OK))
			Expect(createRes).ShouldNot(BeNil())
			accountID = createRes.AccountId
		})

		BeforeEach(func() {
			updateReq.AccountId = accountID
		})

		Context("Updating private account without token", func() {
			Context("Updating private account without token", func() {
				It("should succeed asking for update token", func() {
					updateReq.PrivateAccount = nil
					updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
					Expect(updateRes).Should(BeNil())
				})
			})

		})

		Context("Asking for update token", func() {
			It("should request for token", func() {
				reqReq := &account.RequestChangePrivateAccountRequest{
					Payload: accountPB.Email,
				}
				updateRes, err := AccountAPI.RequestChangePrivateAccount(ctx, reqReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(updateRes).ShouldNot(BeNil())

				v, err := RedisDB.Get(accountID).Result()
				Expect(err).ShouldNot(HaveOccurred())
				token = v
			})
		})

		Describe("Updating account with update token", func() {
			BeforeEach(func() {
				updateReq.ChangeToken = token
			})

			Context("Updating account private profile with bad private payload", func() {
				It("should fail when private profile is nil", func() {
					updateReq.PrivateAccount = nil
					updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
					Expect(updateRes).Should(BeNil())
				})
				It("should fail when passwords do not match", func() {
					updateReq.PrivateAccount.Password = "we dont match"
					updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
					Expect(err).Should(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
					Expect(updateRes).Should(BeNil())
				})
			})

			Context("Updating account private profile with valid request", func() {
				It("should update account in database without error", func() {
					updateReq.PrivateAccount.Password = "hakty11"
					updateReq.PrivateAccount.ConfirmPassword = "hakty11"
					updateRes, err := AccountAPI.UpdatePrivate(ctx, updateReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(updateRes).ShouldNot(BeNil())
				})
			})
		})
	})
})
