package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateHospitalClinic() {
	err := d.DB.AutoMigrate(&models.HospitalClinic{})
	if err != nil {
		return
	}
}
