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

var _ = Describe("Login Account #activate", func() {
	var (
		activateReq *account.ActivateRequest
		activateRes *account.ActivateResponse
		ctx         context.Context
		err         error
		token       string
		accountID   string
		email       string
		password    string
	)

	BeforeEach(func() {
		activateReq = &account.ActivateRequest{
			AccountId: uuid.New().String(),
			Token:     randomdata.RandStringRunes(32),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	AfterEach(func() {
		activateRes = nil
		err = nil
	})

	When("Activating account with missing or nil request", func() {
		It("should definitely fail when request is nil", func() {
			activateReq = nil
			activateRes, err = AccountAPI.Activate(ctx, activateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(activateRes).Should(BeNil())
		})
		It("should definitely fail when token is missing", func() {
			activateReq.Token = ""
			activateRes, err = AccountAPI.Activate(ctx, activateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(activateRes).Should(BeNil())
		})
		It("should definitely fail when accountID is missing", func() {
			activateReq.AccountId = ""
			activateRes, err = AccountAPI.Activate(ctx, activateReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(activateRes).Should(BeNil())
		})
	})

	Context("Activating account => create account, login to get token and id, then activate", func() {
		var label string
		Describe("Creating an account", func() {
			It("should create account in database without error", func() {
				createReq := &account.CreateRequest{
					AccountLabel:   getLabel(),
					Account:        fakeAccount(),
					PrivateAccount: fakePrivateAccount(),
				}

				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
				Expect(createRes.AccountId).ShouldNot(BeZero())

				email = createReq.Account.Email
				password = createReq.PrivateAccount.Password
				accountID = createRes.AccountId
				label = createReq.AccountLabel
			})

			Describe("Login to the created account", func() {
				It("should login the account and return some data", func() {
					loginReq := &account.LoginRequest{
						Login: &account.LoginRequest_Creds{
							Creds: &account.Creds{
								Email:    email,
								Password: password,
							},
						},
						Group: label,
					}
					loginRes, err := AccountAPI.Login(ctx, loginReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(loginRes).ShouldNot(BeNil())
					Expect(loginRes.Token).ShouldNot(BeZero())
					Expect(loginRes.AccountId).ShouldNot(BeZero())

					token = loginRes.Token
				})

				Describe("Activating the account", func() {
					It("should activate the account in database", func() {
						activateReq.AccountId = accountID
						activateReq.Token = token
						activateRes, err = AccountAPI.Activate(ctx, activateReq)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.OK))
						Expect(activateRes).ShouldNot(BeNil())
					})

					Describe("Getting the account and checking if its activated", func() {
						It("should get the account", func() {
							activeState := account.AccountState_ACTIVE.String()
							getReq := &account.GetRequest{
								AccountId: accountID,
							}
							getRes, err := AccountAPI.Get(auth.AddAdminMD(ctx), getReq)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(status.Code(err)).Should(Equal(codes.OK))
							Expect(getRes).ShouldNot(BeNil())
							Expect(getRes.AccountState.String()).Should(BeEquivalentTo(activeState))
						})
					})
				})
			})
		})
	})
})
