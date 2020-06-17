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

var _ = Describe("Get Account #get", func() {
	var (
		getReq *account.GetRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &account.GetRequest{
			AccountId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Context("Get account with nil request", func() {
		It("should fail when request is nil", func() {
			getReq = nil
			getRes, err := AccountAPI.Get(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
	})

	Context("Get account with missing/incorrect account id", func() {
		It("should fail when account id is missing", func() {
			getReq.AccountId = ""
			getRes, err := AccountAPI.Get(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(getRes).Should(BeNil())
		})
		It("should fail when account id is incorrect", func() {
			getReq.AccountId = "knowledge"
			getRes, err := AccountAPI.Get(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(getRes).Should(BeNil())
		})
	})

	Describe("Creating an account and getting it", func() {
		var accountID, nationalID string
		It("should succeed in creating account", func() {
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
			nationalID = createReq.Account.NationalId
		})

		Context("Get account with valid request", func() {
			It("should get the account", func() {
				getReq.AccountId = accountID
				getRes, err := AccountAPI.Get(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
			})
		})

		Context("Get account with valid request", func() {
			It("should get the account using national id true", func() {
				getReq.AccountId = nationalID
				getReq.WithNationalId = true
				getRes, err := AccountAPI.Get(ctx, getReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(getRes).ShouldNot(BeNil())
			})
		})
	})
})
