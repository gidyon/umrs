package account

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

const tableName = "accounts"

// Account contains profile information stored in the database
type Account struct {
	AccountID        string `gorm:"primary_key;type:varchar(50);not null"`
	NationalID       string `gorm:"index:query_index;type:varchar(50);unique_index;not null"`
	Email            string `gorm:"index:query_index;type:varchar(50);unique_index;not null"`
	Phone            string `gorm:"index:query_index;type:varchar(50);unique_index;not null"`
	FirstName        string `gorm:"type:varchar(20);not null"`
	LastName         string `gorm:"type:varchar(20);not null"`
	BirthDate        string `gorm:"type:varchar(30);"`
	Gender           string `gorm:"type:varchar(10);default:'unknown'"`
	Nationality      string `gorm:"type:varchar(50);default:'Kenyan'"`
	ProfileURLThumb  string `gorm:"type:varchar(256)"`
	ProfileURLNormal string `gorm:"type:varchar(256)"`
	AccountType      string `gorm:"index:query_index;type:enum('ADMIN_VIEWER','ADMIN_OWNER','ADMIN_EDITOR','USER_OWNER');not null"`
	AccountState     string `gorm:"index:query_index;type:enum('BLOCKED','ACTIVE','INACTIVE');not null;default:'INACTIVE'"`
	SecurityQuestion string `gorm:"type:text"`
	SecurityAnswer   string `gorm:"type:varchar(50)"`
	Password         string `gorm:"type:text"`
	TrustedDevices   string `gorm:"type:text"`
	AccountLabels    string `gorm:"index:query_index;type:varchar(256)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// TableName is the name of the tables
func (u *Account) TableName() string {
	return tableName
}

// BeforeCreate is a hook that is set before creating object
func (u *Account) BeforeCreate(scope *gorm.Scope) error {
	accountID := uuid.New().String()

	fu := func(str string) string {
		if str == "" {
			return accountID
		}
		return str
	}

	// Update email and phone
	err := scope.SetColumn("Email", fu(u.Email))
	if err != nil {
		return err
	}
	err = scope.SetColumn("Phone", fu(u.Phone))
	if err != nil {
		return err
	}

	return scope.SetColumn("AccountID", accountID)
}

// AfterFind will reset email and phone to their zero value if they equal the user id
func (u *Account) AfterFind() (err error) {
	// Reset email
	if u.Email == u.AccountID {
		u.Email = ""
	}
	if u.Phone == u.AccountID {
		u.Phone = ""
	}
	return
}
