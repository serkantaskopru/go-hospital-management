package migrations

import (
	"hospital-management/models"

	"gorm.io/gorm"
)

type PostgreDB struct {
	DB *gorm.DB
}

func (d PostgreDB) MigratePolyclinic() {
	err := d.DB.AutoMigrate(&models.Polyclinic{})
	if err != nil {
		return
	}
}
