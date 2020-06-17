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

var _ = Describe("Deleting Account #delete", func() {
	var (
		delReq *account.DeleteRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		delReq = &account.DeleteRequest{
			AccountId: uuid.New().String(),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Describe("Deleting account with nil request", func() {
		It("should fail when request is nil", func() {
			delReq = nil
			delRes, err := AccountAPI.Delete(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
	})

	When("Deleting account with missing account id", func() {
		It("should fail when account id is missing", func() {
			delReq.AccountId = ""
			delRes, err := AccountAPI.Delete(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(delRes).Should(BeNil())
		})
	})

	When("Deleting account with unknown account id", func() {
		It("should fail when account id doesn't exist", func() {
			delReq.AccountId = randomdata.RandStringRunes(20)
			delRes, err := AccountAPI.Delete(ctx, delReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
			Expect(delRes).Should(BeNil())
		})
	})

	When("Deleting account with correct account id", func() {
		var accountID string
		Describe("Creating the account first", func() {
			It("should create account in database without error", func() {
				createRes, err := AccountAPI.Create(ctx, &account.CreateRequest{
					AccountLabel:   randomdata.Adjective(),
					Account:        fakeAccount(),
					PrivateAccount: fakePrivateAccount(),
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
				accountID = createRes.AccountId
			})
			Describe("Deleting the account", func() {
				It("should delete account in database without error", func() {
					delReq.AccountId = accountID
					delRes, err := AccountAPI.Delete(ctx, delReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(delRes).ShouldNot(BeNil())
				})

				Describe("Getting the account", func() {
					It("should fail because it has been deleted", func() {
						getReq := &account.GetRequest{
							AccountId: accountID,
						}
						getRes, err := AccountAPI.Get(ctx, getReq)
						Expect(err).Should(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.NotFound))
						Expect(getRes).Should(BeNil())
					})
				})
			})
		})
	})
})
