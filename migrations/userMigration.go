package migrations

import (
	"hospital-management/models"
	"hospital-management/utils"
	"log"
)

func (d PostgreDB) MigrateUser() {
	err := d.DB.AutoMigrate(&models.User{})
	if err != nil {
		return
	}
}
func (d PostgreDB) CreateDefaultUser() {
	var user models.User

	if err := d.DB.Where("email = ?", "admin@example.com").First(&user).Error; err == nil {
		log.Println("Default user already exists")
		return
	}

	defaultUser := models.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: utils.HashPassword("123456"),
		Phone:    "05123456789",
	}

	if err := d.DB.Create(&defaultUser).Error; err != nil {
		log.Fatalf("Failed to create default user: %v", err)
	}

	log.Println("Default user created successfully")
}
