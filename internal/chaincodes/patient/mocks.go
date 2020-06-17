package patient

import (
	"github.com/gidyon/umrs/pkg/api/ledger"
)

// ledgerAPIMock is mock interface for ledger API
type ledgerAPIMock interface {
	ledger.ledgerClient
}
