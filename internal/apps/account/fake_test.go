package account

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/pkg/api/account"
)

func createAdmin(
	accountType account.AccountType, accountState account.AccountState,
) (string, error) {
	accountPB := fakeAccount()
	accountPB.AccountType = accountType
	accountPB.AccountState = accountState

	// Get admin model
	accountDB, err := getAccountDB(accountPB)
	if err != nil {
		return "", err
	}

	// Save to database
	err = DB.Create(accountDB).Error
	if err != nil {
		return "", err
	}

	// Return account ID
	return accountDB.AccountID, nil
}

// creates a fake account
func fakeAccount() *account.Account {
	return &account.Account{
		AccountId:        randomdata.RandStringRunes(32),
		NationalId:       fmt.Sprintf("%d", randomdata.Number(2000000, 40000000)),
		Email:            randomdata.Email(),
		Phone:            randomdata.PhoneNumber()[:10],
		FirstName:        randomdata.FirstName(randomdata.Male),
		LastName:         randomdata.LastName(),
		BirthDate:        randomdata.FullDate(),
		Gender:           "male",
		Nationality:      randomdata.Country(randomdata.FullCountry),
		ProfileUrlThumb:  randomdata.MacAddress(),
		ProfileUrlNormal: randomdata.MacAddress(),
		AccountState:     account.AccountState_ACTIVE,
		AccountType:      account.AccountType_USER_OWNER,
		TrustedDevices: []string{
			randomdata.MacAddress(), randomdata.MacAddress(), randomdata.MacAddress(),
		},
	}
}

// create a fake account private profile
func fakePrivateAccount() *account.PrivateAccount {
	return &account.PrivateAccount{
		Password:         "hakty11",
		ConfirmPassword:  "hakty11",
		SecurityQuestion: "What is your pets name",
		SecurityAnswer:   randomdata.SillyName(),
	}
}
