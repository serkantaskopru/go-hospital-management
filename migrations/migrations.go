package migrations

import (
	"hospital-management/models"
	"hospital-management/utils"
	"log"
)

func (d PostgreDB) CreateDefaultData() {
	defaultUser := models.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: utils.HashPassword("123456"),
		Phone:    "05123456789",
	}

	var existingUser models.User
	if err := d.DB.Where("email = ?", defaultUser.Email).First(&existingUser).Error; err != nil {
		if err := d.DB.Create(&defaultUser).Error; err != nil {
			log.Fatalf("Failed to create default user: %v", err)
		}
	} else {
		defaultUser = existingUser
	}

	defaultHospital := models.Hospital{
		Name:      "Default Hospital",
		Mail:      "mail@hospital.com",
		Phone:     "2268129900",
		TaxNumber: 1111111111,
		Address:   "Must be filled this area",
	}

	var existingHospital models.Hospital
	if err := d.DB.Where("mail = ?", defaultHospital.Mail).First(&existingHospital).Error; err != nil {
		if err := d.DB.Create(&defaultHospital).Error; err != nil {
			log.Fatalf("Failed to create default hospital: %v", err)
		}
	} else {
		defaultHospital = existingHospital
	}

	defaultStaff := models.Staff{
		UserID:     defaultUser.ID,
		HospitalID: defaultHospital.ID,
		Role:       "authorized",
	}

	var existingStaff models.Staff
	if err := d.DB.Where("user_id = ? AND hospital_id = ?", defaultStaff.UserID, defaultStaff.HospitalID).First(&existingStaff).Error; err != nil {
		if err := d.DB.Create(&defaultStaff).Error; err != nil {
			log.Fatalf("Failed to create default staff: %v", err)
		}
	} else {
		log.Println("Default staff already exists.")
	}

	polyclinics := []models.Polyclinic{
		{Name: "Genel Cerrahi"},
		{Name: "Dahiliye"},
		{Name: "Pediatri"},
		{Name: "Ortopedi"},
		{Name: "Kardiyoloji"},
		{Name: "Kadın Hastalıkları ve Doğum"},
		{Name: "Göz Hastalıkları"},
		{Name: "Kulak Burun Boğaz"},
		{Name: "Üroloji"},
		{Name: "Nöroloji"},
	}

	for _, polyclinic := range polyclinics {
		if err := d.DB.FirstOrCreate(&polyclinic, models.Polyclinic{Name: polyclinic.Name}).Error; err != nil {
			log.Fatalf("Failed to create default polyclinic: %v", err)
		}
	}

}
