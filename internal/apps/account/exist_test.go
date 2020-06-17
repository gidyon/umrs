package account

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Exist Phone or Email #exist", func() {
	var (
		existReq *account.ExistRequest
		ctx      context.Context
	)

	BeforeEach(func() {
		existReq = &account.ExistRequest{
			Email:      randomdata.Email(),
			Phone:      randomdata.PhoneNumber(),
			NationalId: fmt.Sprintf("%d", randomdata.Number(2000000, 40000000)),
		}
		ctx = auth.AddAdminMD(context.Background())
	})

	Context("Failure scenarios", func() {
		When("Checking if email or phone exists with bad credentials", func() {
			It("should fail when the request is nil", func() {
				existReq = nil
				existRes, err := AccountAPI.Exist(ctx, existReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(existRes).Should(BeNil())
			})

			It("should fail when the national id and email and phone are missing", func() {
				existReq.Email = ""
				existReq.Phone = ""
				existReq.NationalId = ""
				existRes, err := AccountAPI.Exist(ctx, existReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(existRes).Should(BeNil())
			})
		})
	})

	Context("Success scenarios", func() {
		When("Checking existence with correct credentials", func() {
			It("should succeed with correct phone is given and email is not", func() {
				existReq.Phone = ""
				existReq.NationalId = ""
				existRes, err := AccountAPI.Exist(ctx, existReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(existRes).ShouldNot(BeNil())
			})
			It("should succeed with correct email is given and phone is not", func() {
				existReq.Email = ""
				existReq.NationalId = ""
				existRes, err := AccountAPI.Exist(ctx, existReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(existRes).ShouldNot(BeNil())
			})
			It("should succeed with correct national id is given", func() {
				existReq.Email = ""
				existReq.Phone = ""
				existRes, err := AccountAPI.Exist(ctx, existReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(existRes).ShouldNot(BeNil())
			})
			It("should succeed with correct val if national id, email and phone are given", func() {
				existRes, err := AccountAPI.Exist(ctx, existReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(existRes).ShouldNot(BeNil())
			})
		})

		Describe("Creating an account then checking if it exists", func() {
			var phone, email, nationalID string
			It("should succeed in creating an account", func() {
				createReq := &account.CreateRequest{
					AccountLabel:   randomdata.Adjective(),
					Account:        fakeAccount(),
					PrivateAccount: fakePrivateAccount(),
				}
				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
				phone = createReq.Account.Phone
				email = createReq.Account.Email
				nationalID = createReq.Account.NationalId
			})

			When("Checking existence with credentials that already exist", func() {
				It("should succeed if phone already exists and existence should be true", func() {
					existReq.Phone = phone
					existRes, err := AccountAPI.Exist(ctx, existReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(existRes).ShouldNot(BeNil())
					Expect(existRes.Exists).Should(BeTrue())
				})
				It("should succeed if email already exists and existence should be true", func() {
					existReq.Email = email
					existRes, err := AccountAPI.Exist(ctx, existReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(existRes).ShouldNot(BeNil())
					Expect(existRes.Exists).Should(BeTrue())
				})
				It("should succeed if national id already exists and existence should be true", func() {
					existReq.NationalId = nationalID
					existRes, err := AccountAPI.Exist(ctx, existReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(existRes).ShouldNot(BeNil())
					Expect(existRes.Exists).Should(BeTrue())
				})
			})
		})
	})
})
