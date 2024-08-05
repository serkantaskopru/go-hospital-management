package models

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	FirstName        string `gorm:"not null"`
	LastName         string `gorm:"not null"`
	IdentityNumber   string `gorm:"not null;unique"`
	Phone            string `gorm:"not null;unique"`
	JobGroupID       uint   `gorm:"not null"`
	TitleID          *uint
	HospitalClinicID *uint
	JobGroup         JobGroup        `gorm:"foreignKey:JobGroupID"`
	Title            *Title          `gorm:"foreignKey:TitleID"`
	HospitalClinic   *HospitalClinic `gorm:"foreignKey:HospitalClinicID"`
	HospitalID       int             `gorm:"not null"`
	Hospital         Hospital        `gorm:"foreignKey:HospitalID"`
}
