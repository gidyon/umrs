package account

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Requesting for an account reset #reset", func() {
	var (
		resetReq *account.RequestChangePrivateAccountRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		resetReq = &account.RequestChangePrivateAccountRequest{
			Payload: randomdata.Email(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Requesting for change token with malformed request", func() {
		It("should fail when the request is nil", func() {
			resetReq = nil
			updateRes, err := AccountAPI.RequestChangePrivateAccount(ctx, resetReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
		It("should fail when the payload in request is missing", func() {
			resetReq.Payload = ""
			updateRes, err := AccountAPI.RequestChangePrivateAccount(ctx, resetReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(updateRes).Should(BeNil())
		})
	})

	Describe("Requesting for change token with well-formed request", func() {
		It("should fail when the account does not exist", func() {
			reqReq := &account.RequestChangePrivateAccountRequest{
				Payload: randomdata.Email(),
			}
			updateRes, err := AccountAPI.RequestChangePrivateAccount(ctx, reqReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(updateRes).Should(BeNil())
		})

		Context("Requesting for token with an existing account", func() {
			var payload string
			Context("Lets create an account first", func() {
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
					payload = createReq.Account.Email
				})
			})

			It("should succeed in requesting the token", func() {
				reqReq := &account.RequestChangePrivateAccountRequest{
					Payload: payload,
				}
				updateRes, err := AccountAPI.RequestChangePrivateAccount(ctx, reqReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(updateRes).ShouldNot(BeNil())
			})
		})
	})
})
