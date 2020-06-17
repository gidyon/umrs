package account

import (
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/account"
	"strings"
)

func getAccountPB(accountDB *Account) (*account.Account, error) {
	if accountDB == nil {
		return nil, errs.NilObject("Account")
	}
	accountPB := &account.Account{
		AccountId:        accountDB.AccountID,
		NationalId:       accountDB.NationalID,
		Email:            accountDB.Email,
		Phone:            accountDB.Phone,
		FirstName:        accountDB.FirstName,
		LastName:         accountDB.LastName,
		BirthDate:        accountDB.BirthDate,
		Gender:           accountDB.Gender,
		Nationality:      accountDB.Nationality,
		ProfileUrlThumb:  accountDB.ProfileURLThumb,
		ProfileUrlNormal: accountDB.ProfileURLNormal,
		AccountType:      account.AccountType(account.AccountType_value[accountDB.AccountType]),
		AccountState:     account.AccountState(account.AccountState_value[accountDB.AccountState]),
		AccountLabels:    strings.Split(accountDB.AccountLabels, ","),
		TrustedDevices:   strings.Split(accountDB.TrustedDevices, ","),
	}

	return accountPB, nil
}

func getAccountDB(accountPB *account.Account) (*Account, error) {
	accountDB := &Account{
		AccountID:        accountPB.AccountId,
		NationalID:       accountPB.NationalId,
		Email:            accountPB.Email,
		Phone:            accountPB.Phone,
		FirstName:        accountPB.FirstName,
		LastName:         accountPB.LastName,
		BirthDate:        accountPB.BirthDate,
		Gender:           accountPB.Gender,
		Nationality:      accountPB.Nationality,
		ProfileURLThumb:  accountPB.ProfileUrlThumb,
		ProfileURLNormal: accountPB.ProfileUrlNormal,
		AccountType:      accountPB.AccountType.String(),
		AccountState:     accountPB.AccountState.String(),
		AccountLabels:    strings.Join(accountPB.AccountLabels, ","),
		TrustedDevices:   strings.Join(accountPB.TrustedDevices, ","),
	}

	return accountDB, nil
}

func getAccountPBView(accountPB *account.Account, view account.AccountView) *account.Account {
	if accountPB == nil {
		return accountPB
	}
	switch view {
	case account.AccountView_SEARCH_VIEW, account.AccountView_LIST_VIEW:
		return &account.Account{
			AccountId:    accountPB.AccountId,
			Email:        accountPB.Email,
			Phone:        accountPB.Phone,
			FirstName:    accountPB.FirstName,
			LastName:     accountPB.LastName,
			AccountType:  accountPB.AccountType,
			AccountState: accountPB.AccountState,
		}
	default:
		return accountPB
	}
}
