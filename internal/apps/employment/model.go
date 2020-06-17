package employment

import (
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/employment"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

const employmentTable = "employments"

type employmentModel struct {
	AccountID          string `gorm:"index:query_index;type:varchar(50)"`
	FullName           string `gorm:"type:varchar(50);not null"`
	ProfileThumbURL    string `gorm:"type:text(256)"`
	EmploymentID       string `gorm:"primary_key;type:varchar(50)"`
	EmploymentType     int8   `gorm:"type:tinyint(2);not null"`
	JoinedDate         string `gorm:"type:varchar(30);not null"`
	OrganizationType   string `gorm:"type:varchar(50);not null"`
	OrganizationName   string `gorm:"type:varchar(50);not null"`
	OrganizationID     string `gorm:"type:varchar(50);not null"`
	RoleAtOrganization string `gorm:"type:varchar(50);not null"`
	WorkEmail          string `gorm:"type:varchar(50);not null"`
	StillEmployed      bool   `gorm:"type:tinyint(1);default:0"`
	IsRecent           bool   `gorm:"index:query_index;type:tinyint(1);default:0"`
	EmploymentVerified bool   `gorm:"index:query_index;type:tinyint(1);default:0"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}

// BeforeCreate is a hook that is set before creating object
func (*employmentModel) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("EmploymentID", uuid.New().String())
}

// TableName is the table name
func (*employmentModel) TableName() string {
	return employmentTable
}

func getEmploymentDB(employmentPB *employment.Employment) (*employmentModel, error) {
	if employmentPB == nil {
		return nil, errs.NilObject("Employment")
	}
	employmentDB := &employmentModel{
		AccountID:          employmentPB.AccountId,
		FullName:           employmentPB.FullName,
		ProfileThumbURL:    employmentPB.ProfileThumbUrl,
		EmploymentID:       employmentPB.EmploymentId,
		EmploymentType:     int8(employmentPB.EmploymentType),
		JoinedDate:         employmentPB.JoinedDate,
		OrganizationType:   employmentPB.OrganizationType,
		OrganizationName:   employmentPB.OrganizationName,
		OrganizationID:     employmentPB.OrganizationId,
		RoleAtOrganization: employmentPB.RoleAtOrganization,
		WorkEmail:          employmentPB.WorkEmail,
		EmploymentVerified: employmentPB.EmploymentVerified,
		StillEmployed:      employmentPB.StillEmployed,
		IsRecent:           employmentPB.IsRecent,
	}
	return employmentDB, nil
}

func getEmploymentPB(employmentDB *employmentModel) (*employment.Employment, error) {
	if employmentDB == nil {
		return nil, errs.NilObject("employmentModel")
	}
	employmentPB := &employment.Employment{
		AccountId:          employmentDB.AccountID,
		FullName:           employmentDB.FullName,
		ProfileThumbUrl:    employmentDB.ProfileThumbURL,
		EmploymentId:       employmentDB.EmploymentID,
		EmploymentType:     employment.EmploymentType(employmentDB.EmploymentType),
		JoinedDate:         employmentDB.JoinedDate,
		OrganizationType:   employmentDB.OrganizationType,
		OrganizationName:   employmentDB.OrganizationName,
		OrganizationId:     employmentDB.OrganizationID,
		RoleAtOrganization: employmentDB.RoleAtOrganization,
		WorkEmail:          employmentDB.WorkEmail,
		EmploymentVerified: employmentDB.EmploymentVerified,
		StillEmployed:      employmentDB.StillEmployed,
		IsRecent:           employmentDB.IsRecent,
	}
	return employmentPB, nil
}
