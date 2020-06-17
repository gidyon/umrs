package hospital

import (
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/hospital"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"strings"
)

const tableName = "hospitals"

type hospitalModel struct {
	HospitalID   string `gorm:"primary_key;type:varchar(50)"`
	HospitalName string `gorm:"unique_index;type:varchar(50);not null"`
	WebsiteURL   string `gorm:"type:varchar(100);default:'NA'"`
	LogoURL      string `gorm:"type:varchar(256);default:'NA'"`
	County       string `gorm:"type:varchar(100);not null"`
	SubCounty    string `gorm:"type:varchar(100);not null"`
	AdminEmails  string `gorm:"type:text(500);not null"`
	Permission   bool   `gorm:"type:tinyint(1);default:0;not null"`
}

func (*hospitalModel) TableName() string {
	return tableName
}

// BeforeCreate is a hook that is set before creating object
func (*hospitalModel) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("HospitalID", uuid.New().String())
}

func getHospitalPB(hospitalDB *hospitalModel) (*hospital.Hospital, error) {
	if hospitalDB == nil {
		return nil, errs.NilObject("hospitalModel")
	}
	hospitalPB := &hospital.Hospital{
		HospitalName: hospitalDB.HospitalName,
		HospitalId:   hospitalDB.HospitalID,
		WebsiteUrl:   hospitalDB.WebsiteURL,
		LogoUrl:      hospitalDB.LogoURL,
		County:       hospitalDB.County,
		SubCounty:    hospitalDB.SubCounty,
		AdminEmails:  strings.Split(hospitalDB.AdminEmails, ","),
		Permission:   hospitalDB.Permission,
	}
	return hospitalPB, nil
}

func getHospitalDB(hospitalPB *hospital.Hospital) (*hospitalModel, error) {
	if hospitalPB == nil {
		return nil, errs.NilObject("Hospital")
	}
	hospitalDB := &hospitalModel{
		HospitalName: hospitalPB.HospitalName,
		HospitalID:   hospitalPB.HospitalId,
		WebsiteURL:   hospitalPB.WebsiteUrl,
		LogoURL:      hospitalPB.LogoUrl,
		County:       hospitalPB.County,
		SubCounty:    hospitalPB.SubCounty,
		AdminEmails:  strings.Join(hospitalPB.AdminEmails, ","),
		Permission:   hospitalPB.Permission,
	}
	return hospitalDB, nil
}
