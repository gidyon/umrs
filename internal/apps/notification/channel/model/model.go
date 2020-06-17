package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

// Channel is a bulk channel
type Channel struct {
	ChannelID   string `gorm:"type:varchar(50);primary_key"`
	ChannelName string `gorm:"type:varchar(50);not null;unique_index"`
	OwnerID     string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:text;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// BeforeCreate is a hook that is set before creating object
func (channel *Channel) BeforeCreate(scope *gorm.Scope) error {
	if channel.ChannelID != "" {
		return nil
	}
	return scope.SetColumn("ChannelID", uuid.New().String())
}
