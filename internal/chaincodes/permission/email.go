package permission

import (
	"github.com/gidyon/umrs/pkg/api/permission"
)

type emailData struct {
	RequesterProfile *permission.BasicProfile
	Reason           string
	GrantAccessURL   string
}
