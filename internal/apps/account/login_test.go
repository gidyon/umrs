package account

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

var labels = []string{"HOSPITAL", "INSURANCE", "GOVERNMENT", "ADMIN"}

func getLabel() string {
	return labels[rand.Intn(len(labels))]
}

var _ = Describe("Login Account #login", func() {
	var (
		loginReq *account.LoginRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		loginReq = &account.LoginRequest{
			Login: &account.LoginRequest_Creds{
				Creds: &account.Creds{
					Email:    randomdata.Email(),
					Password: "hakty11",
				},
			},
			Group: getLabel(),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	When("An account login with nil request", func() {
		It("should definitely fail", func() {
			loginReq = nil
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(loginRes).Should(BeNil())
		})
	})

	When("A account login with missing credentials", func() {
		It("should fail when email and phone is missing in the login credentials", func() {
			loginReq.GetCreds().Email = ""
			loginReq.GetCreds().Phone = ""
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(loginRes).Should(BeNil())
		})
		It("should fail when password is missing in the login credentials", func() {
			loginReq.GetCreds().Password = ""
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(loginRes).Should(BeNil())
		})
	})

	When("An account login with incorrect credentials", func() {
		It("should fail when email is incorrect", func() {
			loginReq.GetCreds().Email = "incorrect@gmail.com"
			loginReq.GetCreds().Phone = ""
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(loginRes).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("account"))
		})
		It("should fail when phone is incorrect", func() {
			loginReq.GetCreds().Phone = "07notexist"
			loginReq.GetCreds().Email = ""
			loginRes, err := AccountAPI.Login(ctx, loginReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(loginRes).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("account"))
		})
	})

	When("A account login with valid credentials", func() {
		var accountID, email, password, label string
		Context("Let's create an account first", func() {
			It("should create account without error", func() {
				createReq := &account.CreateRequest{
					AccountLabel:   getLabel(),
					Account:        fakeAccount(),
					PrivateAccount: fakePrivateAccount(),
				}
				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
				email = createReq.Account.Email
				password = createReq.PrivateAccount.Password
				accountID = createRes.AccountId
				label = createReq.AccountLabel
			})
			Context("Login into the account", func() {
				It("should login the account and return token and some data", func() {
					loginReq.GetCreds().Email = email
					loginReq.GetCreds().Password = password
					loginReq.Group = label
					loginRes, err := AccountAPI.Login(ctx, loginReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(loginRes).ShouldNot(BeNil())
					Expect(loginRes.Token).ShouldNot(BeZero())
					Expect(loginRes.AccountId).ShouldNot(BeZero())
					Expect(loginRes.AccountId).Should(Equal(accountID))
				})
			})
		})
	})
})
