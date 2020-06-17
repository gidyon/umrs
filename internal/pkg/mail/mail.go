package mail

import (
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// NewDialer creates a smtp dialer
func NewDialer(
	host, username, password string, port int,
) (*gomail.Dialer, error) {

	switch {
	case host == "":
		return nil, errors.New("mail host is required")
	case username == "":
		return nil, errors.New("mail username is required")
	case password == "":
		return nil, errors.New("mail password is required")
	case port <= 0:
		return nil, errors.New("mail port is required")
	}

	return &gomail.Dialer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}, nil
}
