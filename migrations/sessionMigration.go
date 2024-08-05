package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateSession() {
	err := d.DB.AutoMigrate(&models.Session{})
	if err != nil {
		return
	}
}
