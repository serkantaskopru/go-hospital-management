package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateHospital() {
	err := d.DB.AutoMigrate(&models.Hospital{})
	if err != nil {
		return
	}
}
