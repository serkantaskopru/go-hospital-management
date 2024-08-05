package models

import "gorm.io/gorm"

type HospitalClinic struct {
	gorm.Model
	HospitalID   int        `gorm:"not null"`
	PolyclinicID int        `gorm:"not null"`
	Hospital     Hospital   `gorm:"foreignKey:HospitalID"`
	Polyclinic   Polyclinic `gorm:"foreignKey:PolyclinicID"`
	Persons      []Person   `gorm:"foreignKey:HospitalClinicID"`
}
