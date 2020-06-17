package account

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/pkg/api/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Creating Account #create", func() {

	var (
		createReq *account.CreateRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &account.CreateRequest{
			AccountLabel:   randomdata.Adjective(),
			Account:        fakeAccount(),
			PrivateAccount: fakePrivateAccount(),
		}
		ctx = auth.AddHospitalMD(context.Background())
	})

	Context("Failing Scenarios", func() {
		When("Creating account with nil request", func() {
			It("should fail when request is nil", func() {
				createReq = nil
				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(createRes).Should(BeNil())
			})
		})

		When("Creating account with some missing request fields", func() {
			It("should fail when first name is missing", func() {
				createReq.Account.FirstName = ""
				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(createRes).Should(BeNil())
			})
			It("should fail when last name is missing", func() {
				createReq.Account.LastName = ""
				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).Should(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
				Expect(createRes).Should(BeNil())
			})
			It("should succeed when both phone and email are missing", func() {
				createReq.Account.Phone = ""
				createReq.Account.Email = ""
				createRes, err := AccountAPI.Create(ctx, createReq)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
			})
		})

		When("Creating account with email or phone that already exists in database", func() {
			var email, phone string
			Context("Lets create an account", func() {
				It("should create account in database without error", func() {
					createRes, err := AccountAPI.Create(ctx, createReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(createRes).ShouldNot(BeNil())
					email = createReq.Account.Email
					phone = createReq.Account.Phone
				})

				Describe("Creating account with existing email or phone", func() {
					It("should fail when email is already registered", func() {
						createReq.Account.Email = email
						createRes, err := AccountAPI.Create(ctx, createReq)
						Expect(err).Should(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.ResourceExhausted))
						Expect(createRes).Should(BeNil())
					})
					It("should fail when phone is already registered", func() {
						createReq.Account.Phone = phone
						createRes, err := AccountAPI.Create(ctx, createReq)
						Expect(err).Should(HaveOccurred())
						Expect(status.Code(err)).Should(Equal(codes.ResourceExhausted))
						Expect(createRes).Should(BeNil())
					})
				})
			})
		})
	})

	Context("Success scenarios", func() {
		When("Creating account with valid request", func() {
			It("should succeed when email is missing but phone is not", func() {
				createReq.Account.Email = ""
				createReq.Account.Phone = randomdata.PhoneNumber()[:10]
				createRes, err := AccountAPI.Create(ctx, createReq)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
			})

			It("should succeed when phone is missing but email is not", func() {
				createReq.Account.Phone = ""
				createReq.Account.Email = randomdata.Email()
				createRes, err := AccountAPI.Create(ctx, createReq)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(status.Code(err)).Should(Equal(codes.OK))
				Expect(createRes).ShouldNot(BeNil())
			})

			When("Creating account with all request fields provided", func() {
				It("should create account in database without error", func() {
					createRes, err := AccountAPI.Create(ctx, createReq)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(status.Code(err)).Should(Equal(codes.OK))
					Expect(createRes).ShouldNot(BeNil())
				})
			})
		})
	})
})
