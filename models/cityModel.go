package models

import (
	"gorm.io/gorm"
)

type City struct {
	gorm.Model

	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"unique;not null"`
	Plate     int        `gorm:"unique;not null"`
	Districts []District `gorm:"foreignKey:CityID"`
}
