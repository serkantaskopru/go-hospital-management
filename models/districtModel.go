package models

type District struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"not null"`
	CityID uint
	City   City
}
