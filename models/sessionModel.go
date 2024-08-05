package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	HospitalID uint      `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
}

func (session *Session) BeforeCreate(tx *gorm.DB) (err error) {
	location, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		return
	}

	now := time.Now().In(location)

	session.ExpiresAt = now.Add(45 * time.Minute)
	return
}
