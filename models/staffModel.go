package models

import (
	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	ID             int `json:"id" gorm:"primaryKey"`
	UserID         int
	HospitalID     int
	User           *User     `gorm:"foreignKey:UserID"`
	Hospital       *Hospital `gorm:"foreignKey:HospitalID"`
	Role           string
	IdentityNumber string
}
