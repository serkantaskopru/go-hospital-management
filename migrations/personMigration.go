package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigratePerson() {
	err := d.DB.AutoMigrate(&models.Person{})
	if err != nil {
		return
	}
}
