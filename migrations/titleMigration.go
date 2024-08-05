package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateTitle() {
	err := d.DB.AutoMigrate(&models.Title{})
	if err != nil {
		return
	}
}
