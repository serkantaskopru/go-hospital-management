package controllers

import (
	"hospital-management/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PolyclinicController struct {
	DB *gorm.DB
}

func (lc *PolyclinicController) GetPolyclinics(c *fiber.Ctx) error {
	var polyclinics []models.Polyclinic
	if err := lc.DB.Preload("Clinics").Find(&polyclinics).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(polyclinics)
}
