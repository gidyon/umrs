package patient

import (
	"encoding/json"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/patient"
)

type patientMedicalData struct {
	PatientID    string `gorm:"primary_key;type:varchar(50)"`
	HospitalID   string `gorm:"type:varchar(50);not null"`
	PatientState string `gorm:"type:varchar(10);not null"`
	Details      []byte `gorm:"type:json;"`
}

func getPatientMedDataDB(medDataPB *patient.MedicalData) (*patientMedicalData, error) {
	if medDataPB == nil {
		return nil, errs.NilObject("MedicalData")
	}
	patientMed := &patientMedicalData{
		HospitalID:   medDataPB.HospitalId,
		PatientState: medDataPB.PatientState.String(),
		Details:      make([]byte, 0),
	}
	if len(medDataPB.GetDetails().GetDetails()) != 0 {
		bs, err := json.Marshal(medDataPB.Details)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Details")
		}
		patientMed.Details = bs
	}
	return patientMed, nil
}

func getPatientMedDataPB(medDataDB *patientMedicalData) (*patient.MedicalData, error) {
	if medDataDB == nil {
		return nil, errs.NilObject("patientMedicalData")
	}
	patientMed := &patient.MedicalData{
		HospitalId:   medDataDB.HospitalID,
		PatientState: patient.State(patient.State_value[medDataDB.PatientState]),
	}
	if len(medDataDB.Details) != 0 {
		err := json.Unmarshal(medDataDB.Details, &patientMed.Details)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Details")
		}
	}
	return patientMed, nil
}
