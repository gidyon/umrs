package account

import (
	"context"
	"github.com/gidyon/umrs/pkg/api/account"
)

// GetFullName gets the full name of account with given id
func GetFullName(
	ctx context.Context, accountID string, accountAPIServer account.AccountAPIServer,
) (string, error) {
	accountPB, err := accountAPIServer.Get(ctx, &account.GetRequest{
		AccountId: accountID,
	})
	if err != nil {
		return "", err
	}
	return accountPB.FirstName + " " + accountPB.LastName, nil
}
