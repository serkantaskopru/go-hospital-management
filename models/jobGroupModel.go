package models

import (
	"gorm.io/gorm"
)

type JobGroup struct {
	gorm.Model
	Name    string   `gorm:"not null"`
	Titles  []Title  `gorm:"foreignKey:JobGroupID"`
	Persons []Person `gorm:"foreignKey:JobGroupID"`
}
