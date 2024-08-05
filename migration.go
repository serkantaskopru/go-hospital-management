package main

import (
	"hospital-management/configs"
	"hospital-management/migrations"
)

func _main() {
	dbClient := configs.ConnectPostgreSQL()

	rm := migrations.PostgreDB{DB: dbClient}

	rm.MigrateHospital()
	rm.MigrateStaff()
}
