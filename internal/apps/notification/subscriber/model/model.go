package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Subscriber is user who is subscribed to channels
type Subscriber struct {
	AccountID  string `gorm:"type:varchar(50);primary_key"`
	Email      string `gorm:"type:varchar(50);unique_index"`
	Phone      string `gorm:"type:varchar(50);unique_index"`
	SendMethod string `gorm:"type:enum('EMAIL','SMS','USSD', 'CALL', 'EMAIL_AND_SMS');default:'EMAIL';not null"`
	Channels   []byte `gorm:"type:json"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

// BeforeCreate will update email and phone to make sure they are not empty
func (subscriber *Subscriber) BeforeCreate(scope *gorm.Scope) error {
	updateField := func(str string) string {
		if str == "" {
			return subscriber.AccountID
		}
		return str
	}

	// Update email and phone
	err := scope.SetColumn("Email", updateField(subscriber.Email))
	if err != nil {
		return err
	}
	err = scope.SetColumn("Phone", updateField(subscriber.Phone))
	if err != nil {
		return err
	}

	return nil
}

// AfterFind runs logic after query
func (subscriber *Subscriber) AfterFind() error {
	if subscriber.Email == subscriber.AccountID {
		subscriber.Email = ""
	}
	if subscriber.Phone == subscriber.AccountID {
		subscriber.Phone = ""
	}
	return nil
}
