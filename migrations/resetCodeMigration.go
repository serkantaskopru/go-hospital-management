package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateResetCode() {
	err := d.DB.AutoMigrate(&models.ResetCode{})
	if err != nil {
		return
	}
}
