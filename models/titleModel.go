package models

import (
	"gorm.io/gorm"
)

type Title struct {
	gorm.Model
	Name       string `gorm:"not null"`
	JobGroupID int
	JobGroup   *JobGroup `gorm:"foreignKey:JobGroupID"`
	Persons    []Person  `gorm:"foreignKey:TitleID"`
}
