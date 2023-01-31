package models

import (
	"psyWeb/utils"
	"time"
)

type DiagnosticResult struct {
	PhoneNumber string
	Name   string
	Age    uint8
	Gender string
	Time   string
	SAS    string
	ESS    string
	ISI    string
	SDS    string
}

type UserReport struct {
	Name     string
	Age      uint8
	Gender   utils.UserGender
	Time     string
	SASScore float32
	ESSScore float32
	ISIScore float32
	SDSScore float32
}

func (r *UserReport) queryFromDB(phone_number string) error {
	query_str := `SELECT 
				 Name,Gender,Age,SASScore,ESSScore,ISIScore,SDSScore 
				 FROM user 
				 WHERE PhoneNumber=?`
	row := utils.GetPsyWebDataBaseInstance().Db.QueryRow(query_str, phone_number)
	return row.Scan(
		&r.Name,
		&r.Gender,
		&r.Age,
		&r.SASScore,
		&r.ESSScore,
		&r.ISIScore,
		&r.SDSScore,
	)
}

func (r *UserReport) sas() utils.DiseaseSeverity {
	if r.SASScore < 50.0 {
		return utils.Asymptomatic
	} else if r.SASScore < 60.0 {
		return utils.SlightSymptoms
	} else if r.SASScore < 70.0 {
		return utils.ModerateSymptoms
	} else {
		return utils.SeriousSymptoms
	}
}

func (r *UserReport) ess() utils.DiseaseSeverity {
	if r.ESSScore < 8.0 {
		return utils.Asymptomatic
	} else if r.ESSScore < 10.0 {
		return utils.SlightSymptoms
	} else if r.ESSScore < 12.0 {
		return utils.ModerateSymptoms
	} else {
		return utils.SeriousSymptoms
	}
}

func (r *UserReport) isi() utils.DiseaseSeverity {
	if r.ISIScore < 50.0 {
		return utils.Asymptomatic
	} else if r.ISIScore < 60.0 {
		return utils.SlightSymptoms
	} else if r.ISIScore < 70.0 {
		return utils.ModerateSymptoms
	} else {
		return utils.SeriousSymptoms
	}
}

func (r *UserReport) sds() utils.DiseaseSeverity {
	if r.SDSScore < 50.0 {
		return utils.Asymptomatic
	} else if r.SDSScore < 60.0 {
		return utils.SlightSymptoms
	} else if r.SDSScore < 70.0 {
		return utils.ModerateSymptoms
	} else {
		return utils.SeriousSymptoms
	}
}

func GetDiagnosticResult(phone_number string) (*DiagnosticResult, error) {
	r := UserReport{
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := r.queryFromDB(phone_number)
	if err != nil {
		return nil, err
	}
	d := &DiagnosticResult{
		PhoneNumber: phone_number,
		Name:   r.Name,
		Age:    r.Age,
		Gender: r.Gender.String(),
		Time:   r.Time,
		SAS:    r.sas().String(),
		ESS:    r.ess().String(),
		ISI:    r.isi().String(),
		SDS:    r.sds().String(),
	}
	return d, nil
}
