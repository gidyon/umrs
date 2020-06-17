package notification

import (
	"gopkg.in/gomail.v2"
)

// Dialer is the dialer interface
type Dialer interface {
	DialAndSend(...*gomail.Message) error
}
