package controllers

import (
	"hospital-management/models"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HospitalClinicController struct {
	DB *gorm.DB
}

func (controller *HospitalClinicController) GetHospitalClinics(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	staff := user.Staff

	var hospital models.Hospital
	if err := controller.DB.
		Preload("Clinics.Persons.JobGroup").
		Preload("Clinics.Polyclinic").
		First(&hospital, staff.HospitalID).
		Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	clinicsData := []map[string]interface{}{}

	for _, clinic := range hospital.Clinics {
		personCount := len(clinic.Persons)

		jobGroupCounts := make(map[uint]int)
		for _, person := range clinic.Persons {
			if person.JobGroupID > 0 {
				jobGroupCounts[person.JobGroupID]++
			}
		}

		var jobGroups []map[string]interface{}
		for jobGroupID, count := range jobGroupCounts {
			var jobGroup models.JobGroup
			if err := controller.DB.First(&jobGroup, jobGroupID).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			jobGroups = append(jobGroups, map[string]interface{}{
				"id":    jobGroup.ID,
				"name":  jobGroup.Name,
				"count": count,
			})
		}

		clinicsData = append(clinicsData, map[string]interface{}{
			"id":          clinic.ID,
			"polyclinic":  clinic.Polyclinic.Name,
			"personCount": personCount,
			"jobGroups":   jobGroups,
		})
	}

	response := map[string]interface{}{
		"clinics": clinicsData,
	}

	return c.JSON(response)
}

func (controller *HospitalClinicController) CreateHospitalClinic(c *fiber.Ctx) error {
	type ClinicInput struct {
		PolyclinicId string `json:"polyclinicId" validate:"required"`
	}

	var input ClinicInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	PolyclinicId, err := strconv.Atoi(input.PolyclinicId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz poliklinik"})
	}

	user := c.Locals("user").(models.User)

	staff := user.Staff

	hospital := staff.Hospital

	var existingHospitalClinic models.HospitalClinic
	if err := controller.DB.Where("hospital_id = ? AND polyclinic_id = ?", hospital.ID, PolyclinicId).First(&existingHospitalClinic).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Bu poliklinik zaten hastanede mevcut"})
	}

	hospitalClinic := models.HospitalClinic{
		HospitalID:   hospital.ID,
		PolyclinicID: PolyclinicId,
	}

	if err := controller.DB.Create(&hospitalClinic).Error; err != nil {
		log.Fatalf("Failed to create hospital clinic: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Hastane polikliniği eklenemedi"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Klinik hastaneye başarıyla eklendi"})
}
