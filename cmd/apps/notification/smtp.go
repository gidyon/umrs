package main

import (
	"bytes"
	"github.com/Pallinder/go-randomdata"
	app_mail "github.com/gidyon/umrs/internal/pkg/mail"
	"gopkg.in/gomail.v2"

	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func emailDialer() (*gomail.Dialer, error) {
	var err error

	smtpPort := 585
	if portStr := strings.TrimPrefix(os.Getenv("SMTP_PORT"), ":"); portStr != "" {
		smtpPort, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	smtpPassword := setIfEmpty(os.Getenv("SMTP_PASSWORD"), randomdata.RandStringRunes(20))
	// Read password from file
	if smtpPasswordFile := strings.TrimSpace(os.Getenv("SMTP_PASSWORD_FILE")); smtpPasswordFile != "" {
		passwordBytes, err := ioutil.ReadFile(smtpPasswordFile)
		if err != nil {
			return nil, err
		}
		smtpPassword = setIfEmpty(string(bytes.TrimSpace(passwordBytes)), smtpPassword)
	}

	smtpUsername := setIfEmpty(os.Getenv("SMTP_USERNAME"), randomdata.Email())
	// Read username from file
	if smtpUsernameFile := strings.TrimSpace(os.Getenv("SMTP_USERNAME_FILE")); smtpUsernameFile != "" {
		usernameBytes, err := ioutil.ReadFile(smtpUsernameFile)
		if err != nil {
			return nil, err
		}
		smtpUsername = setIfEmpty(string(bytes.TrimSpace(usernameBytes)), smtpUsername)
	}

	smtpHost := setIfEmpty(os.Getenv("SMTP_HOST"), "localhost")

	return app_mail.NewDialer(smtpHost, smtpUsername, smtpPassword, smtpPort)
}

func setIfEmpty(cur, to string) string {
	if cur == "" {
		return to
	}
	return cur
}
