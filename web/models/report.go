package models

import (
	"psyWeb/utils"
	"time"
)

type UserReport struct {
	Name     string           `json:"Name"`
	Age      uint8            `json:"Age"`
	Gender   utils.UserGender `json:"Gender"`
	Time     string           `json:"Time"`
	SASScore float32          `json:"SAS"`
	ESSScore float32          `json:"ESS"`
	ISIScore float32          `json:"ISI"`
	SDSScore float32          `json:"SDS"`
}

func (r *UserReport) GetResult(phone_number string) error {
	r.Time = time.Now().Format("2006-01-02 15:04:05")
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
