package models

import "gorm.io/gorm"

type Polyclinic struct {
	gorm.Model
	Name    string           `gorm:"unique;not null"`
	Clinics []HospitalClinic `gorm:"foreignKey:PolyclinicID"`
}
