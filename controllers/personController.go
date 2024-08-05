package controllers

import (
	"hospital-management/models"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PersonController struct {
	DB *gorm.DB
}

func (controller *PersonController) ListPersons(c *fiber.Ctx) error {
	var persons []models.Person
	var totalRecords int64
	user := c.Locals("user").(models.User)

	hospitalID := user.Staff.HospitalID

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	identityNumber := c.Query("identityNumber")
	jobGroupID := c.Query("jobGroupId")
	titleID := c.Query("titleId")

	query := controller.DB.Model(&models.Person{}).
		Where("hospital_id = ?", hospitalID)

	if firstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+firstName+"%")
	}
	if lastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+lastName+"%")
	}
	if identityNumber != "" {
		query = query.Where("identity_number ILIKE ?", "%"+identityNumber+"%")
	}
	if jobGroupID != "" {
		query = query.Where("job_group_id = ?", jobGroupID)
	}
	if titleID != "" {
		query = query.Where("title_id = ?", titleID)
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Personel sayısı alma hatası"})
	}

	if err := query.
		Preload("JobGroup").
		Preload("Title").
		Preload("HospitalClinic").
		Preload("HospitalClinic.Hospital").
		Preload("HospitalClinic.Polyclinic").
		Offset(offset).
		Limit(limit).
		Find(&persons).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Personel verisi alma hatası"})
	}

	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	return c.JSON(fiber.Map{
		"data":        persons,
		"totalPages":  totalPages,
		"currentPage": page,
	})
}

func (controller *PersonController) CreatePerson(c *fiber.Ctx) error {
	type PersonInput struct {
		FirstName        string `json:"firstName" validate:"required"`
		LastName         string `json:"lastName" validate:"required"`
		IdentityNumber   string `json:"identityNumber" validate:"required"`
		Phone            string `json:"phone" validate:"required"`
		JobGroupID       string `json:"jobGroupId" validate:"required"`
		TitleID          string `json:"titleId"`
		HospitalClinicID string `json:"hospitalClinicId"`
	}

	user := c.Locals("user").(models.User)
	hospitalID := user.Staff.HospitalID

	var input PersonInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	JobGroupID, err := strconv.Atoi(input.JobGroupID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz meslek grubu verisi"})
	}

	var person models.Person
	person = models.Person{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		IdentityNumber: input.IdentityNumber,
		Phone:          input.Phone,
		JobGroupID:     uint(JobGroupID),
	}

	if input.TitleID != "" {
		titleID, err := strconv.Atoi(input.TitleID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz ünvan verisi"})
		}
		titleIDUint := uint(titleID)
		person.TitleID = &titleIDUint
	}

	if input.HospitalClinicID != "" {
		hospitalClinicID, err := strconv.Atoi(input.HospitalClinicID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz hastane poliklinik verisi"})
		}
		hospitalClinicIDUint := uint(hospitalClinicID)
		person.HospitalClinicID = &hospitalClinicIDUint
	}

	person.HospitalID = hospitalID

	if person.TitleID != nil {
		var titleName string
		if err := controller.DB.Model(&models.Title{}).
			Where("id = ?", *person.TitleID).
			Pluck("name", &titleName).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ünvan verisi alma hatası"})
		}

		if titleName == "Başhekim" {
			var existingChiefCount int64
			if err := controller.DB.Model(&models.Person{}).
				Where("hospital_id = ? AND title_id = ?", hospitalID, person.TitleID).
				Count(&existingChiefCount).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Hastane başhekimlik kontrolü hatası"})
			}

			if existingChiefCount > 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hastane de zaten bir başhekim mevcut"})
			}
		}
	}

	if err := controller.DB.Create(&person).Error; err != nil {
		log.Printf("Failed to create person: %v", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "people_phone_key") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bu telefon numarası kullanımda"})
			}
			if strings.Contains(err.Error(), "people_identity_number_key") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bu TC kimlik numarası kullanımda"})
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Personel oluşturma hatası"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Personel başarıyla oluşturuldu"})
}

func (controller *PersonController) DeletePerson(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := controller.DB.Delete(&models.Person{}, id).Error; err != nil {
		log.Printf("Failed to delete person: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Personel silme hatası"})
	}

	return c.Status(fiber.StatusNoContent).SendString("Personel başarıyla silindi")
}

func (controller *PersonController) UpdatePerson(c *fiber.Ctx) error {
	id := c.Params("id")

	type PersonInput struct {
		FirstName        string `json:"firstName"`
		LastName         string `json:"lastName"`
		IdentityNumber   string `json:"identityNumber"`
		Phone            string `json:"phone"`
		JobGroupID       int    `json:"jobGroupId"`
		TitleID          int    `json:"titleId"`
		HospitalClinicID int    `json:"hospitalClinicId"`
	}

	var input PersonInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	updateData := map[string]interface{}{
		"first_name":         input.FirstName,
		"last_name":          input.LastName,
		"identity_number":    input.IdentityNumber,
		"phone":              input.Phone,
		"job_group_id":       uint(input.JobGroupID),
		"title_id":           uint(input.TitleID),
		"hospital_clinic_id": uint(input.HospitalClinicID),
	}

	if err := controller.DB.Model(&models.Person{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		log.Printf("Failed to update person: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Personel güncelleme hatası"})
	}

	return c.JSON(fiber.Map{"message": "Personel başarıyla güncellendi"})
}
