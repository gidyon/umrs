package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

// Notification is an object that contain requires users attention
type Notification struct {
	NotificationID string `gorm:"primary_key;varchar(50);"`
	AccountID      string `gorm:"index:query_index;type:varchar(50);not null"`
	Notification   []byte `gorm:"type:json;not null"`
	Seen           bool   `gorm:"type:tinyint(1);default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// BeforeCreate is a hook that is set before creating object
func (n *Notification) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("NotificationID", uuid.New().String())
}
