package models

import (
	"time"

	"gorm.io/gorm"
)

type ResetCode struct {
	gorm.Model
	Phone     string
	Code      string
	ExpiresAt time.Time
}
