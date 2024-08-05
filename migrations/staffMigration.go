package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateStaff() {
	err := d.DB.AutoMigrate(&models.Staff{})
	if err != nil {
		return
	}
}
