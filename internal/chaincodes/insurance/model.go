package insurance

import (
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/insurance"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"strings"
)

const tableName = "insurances"

type insuranceModel struct {
	InsuranceID      string `gorm:"primary_key;type:varchar(50)"`
	InsuranceName    string `gorm:"unique_index;type:varchar(50);not null"`
	WebsiteURL       string `gorm:"type:varchar(100);default:'NA'"`
	LogoURL          string `gorm:"type:varchar(256);default:'NA'"`
	About            string `gorm:"type:varchar(256);not null"`
	SupportEmail     string `gorm:"type:varchar(50);not null"`
	SupportTelNumber string `gorm:"type:varchar(50);not null"`
	AdminEmails      string `gorm:"type:text(500);not null"`
	Permission       bool   `gorm:"type:tinyint(1);default:0;not null"`
}

func (*insuranceModel) TableName() string {
	return tableName
}

// BeforeCreate is a hook that is set before creating object
func (*insuranceModel) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("InsuranceID", uuid.New().String())
}

func getInsurancePB(insuranceDB *insuranceModel) (*insurance.Insurance, error) {
	if insuranceDB == nil {
		return nil, errs.NilObject("insuranceModel")
	}
	insurancePB := &insurance.Insurance{
		InsuranceId:      insuranceDB.InsuranceID,
		InsuranceName:    insuranceDB.InsuranceName,
		WebsiteUrl:       insuranceDB.WebsiteURL,
		LogoUrl:          insuranceDB.LogoURL,
		About:            insuranceDB.About,
		SupportEmail:     insuranceDB.SupportEmail,
		SupportTelNumber: insuranceDB.SupportTelNumber,
		AdminEmails:      strings.Split(insuranceDB.AdminEmails, ","),
		Permission:       insuranceDB.Permission,
	}
	return insurancePB, nil
}

func getInsuranceDB(insurancePB *insurance.Insurance) (*insuranceModel, error) {
	if insurancePB == nil {
		return nil, errs.NilObject("Insurance")
	}
	insuranceDB := &insuranceModel{
		InsuranceName:    insurancePB.InsuranceName,
		InsuranceID:      insurancePB.InsuranceId,
		WebsiteURL:       insurancePB.WebsiteUrl,
		LogoURL:          insurancePB.LogoUrl,
		About:            insurancePB.About,
		SupportTelNumber: insurancePB.SupportTelNumber,
		SupportEmail:     insurancePB.SupportEmail,
		AdminEmails:      strings.Join(insurancePB.AdminEmails, ","),
		Permission:       insurancePB.Permission,
	}
	return insuranceDB, nil
}
