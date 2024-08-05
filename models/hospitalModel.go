package models

import "gorm.io/gorm"

type Hospital struct {
	gorm.Model

	ID        int `json:"id" gorm:"primaryKey"`
	Name      string
	Mail      string `gorm:"unique"`
	Phone     string `gorm:"unique"`
	TaxNumber uint   `gorm:"unique"`
	Address   string
	Staffs    []Staff          `gorm:"foreignKey:HospitalID"`
	Clinics   []HospitalClinic `gorm:"foreignKey:HospitalID"`
}
