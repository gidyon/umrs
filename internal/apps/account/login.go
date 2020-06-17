package account

import (
	"context"
	"github.com/gidyon/umrs/internal/pkg/auth"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/account"
	"github.com/gidyon/umrs/pkg/api/ledger"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"strings"
	"time"
)

func (accountAPI *accountAPIServer) Login(
	ctx context.Context, loginReq *account.LoginRequest,
) (*account.LoginResponse, error) {
	// Request should not be nil
	if loginReq == nil {
		return nil, errs.NilObject("LoginRequest")
	}

	var (
		email         string
		phone         string
		nationalID    string
		password      string
		withCreds     bool
		socialProfile *account.SocialProfile
	)
	switch loginReq.Login.(type) {
	case nil:
		return nil, errs.NilObject("Login")
	case *account.LoginRequest_Creds:
		email = loginReq.GetCreds().Email
		phone = loginReq.GetCreds().Phone
		nationalID = loginReq.GetCreds().NationalId
		password = loginReq.GetCreds().Password
		if strings.TrimSpace(password) == "" {
			return nil, errs.MissingCredential("Password")
		}
		withCreds = true
	case *account.LoginRequest_Facebook:
		socialProfile = loginReq.GetFacebook().GetFbProfile()
		email = socialProfile.GetEmailAddress()
		phone = socialProfile.GetPhoneNumber()
	case *account.LoginRequest_Google:
		socialProfile = loginReq.GetGoogle().GetGoogleProfile()
		email = socialProfile.GetEmailAddress()
		phone = socialProfile.GetPhoneNumber()
	case *account.LoginRequest_Twitter:
		socialProfile = loginReq.GetTwitter().GetTwitterProfile()
		email = socialProfile.GetEmailAddress()
		phone = socialProfile.GetPhoneNumber()
	}

	// Validation
	if email == "" && phone == "" {
		return nil, errs.MissingCredential("Email or Phone")
	}

	accountDB := &Account{}

	// Query for user with email or phone as provided
	err := accountAPI.sqlDB.Select(
		"account_id,first_name,last_name,account_type,account_state,password,account_labels",
	).First(accountDB, "national_id=? OR phone=? OR email=?", nationalID, phone, email).Error
	switch {
	case err == nil:
	case gorm.IsRecordNotFoundError(err):
		emailOrPhone := func() string {
			if strings.Contains(email, "@") {
				return email
			}
			return phone
		}
		return nil, errs.AccountDoesntExist(emailOrPhone())
	default:
		return nil, errs.SQLQueryFailed(err, "Login")
	}

	if accountDB.Password == "" {
		return nil, errs.WrapMessage(
			codes.PermissionDenied, "account has no password; please request new password",
		)
	}
	accountPB, err := getAccountPB(accountDB)
	if err != nil {
		return nil, err
	}

	// Check that account is not blocked
	if accountPB.AccountState == account.AccountState_BLOCKED {
		return nil, errs.AccountBlocked()
	}

	// Check if password match if they logged in with Phone or Email
	if withCreds {
		err = compareHash(accountDB.Password, password)
		if err != nil {
			return nil, errs.WrongPassword()
		}
	}

	var token, loginGroup string
	if loginGroup = strings.TrimSpace(loginReq.GetGroup()); loginGroup != "" {
		for _, label := range accountPB.AccountLabels {
			label = strings.ToUpper(strings.TrimSpace(label))
			if label == loginGroup {
				// Generates the token with claims from profile object
				token, err = auth.GenToken(ctx, &auth.Payload{
					ID:           accountDB.AccountID,
					FirstName:    accountDB.FirstName,
					LastName:     accountDB.LastName,
					PhoneNumber:  accountDB.Phone,
					EmailAddress: accountDB.Email,
					Label:        loginGroup,
					Group:        ledger.Actor_value[loginGroup],
				}, ledger.Actor_value[loginGroup], 0)
				if err != nil {
					return nil, errs.FailedToGenToken(err)
				}
			}
		}
	}

	if token == "" {
		return nil, errs.AccountNotGroupMember(loginReq.GetGroup())
	}

	// Update account in case they logged in with twitter, google or facebook
	if !withCreds {
		go func(accountDB *Account) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_, err = accountAPI.Update(ctx, &account.UpdateRequest{
				AccountId: accountDB.AccountID,
				Account: &account.Account{
					AccountId:       accountDB.AccountID,
					Email:           socialProfile.GetEmailAddress(),
					Phone:           socialProfile.GetPhoneNumber(),
					FirstName:       socialProfile.GetFirstName(),
					LastName:        socialProfile.GetLastName(),
					BirthDate:       socialProfile.GetBirthDate(),
					Gender:          socialProfile.GetGender(),
					ProfileUrlThumb: socialProfile.GetProfileUrl(),
				},
			})
			if err != nil {
				errs.LogError("failed to update profile: %v", err)
			}
		}(accountDB)
	}

	// Return token
	return &account.LoginResponse{
		Token:        token,
		AccountId:    accountDB.AccountID,
		AccountState: accountDB.AccountState,
		AccountType:  accountDB.AccountType,
		AccountGroup: loginGroup,
	}, nil
}
