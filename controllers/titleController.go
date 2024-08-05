package controllers

import (
	"hospital-management/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TitleController struct {
	DB *gorm.DB
}

func (controller *TitleController) GetTitles(c *fiber.Ctx) error {
	var titles []models.Title
	jobGroupId := c.Params("jobGroupId")

	if err := controller.DB.Where("job_group_id = ?", jobGroupId).Find(&titles).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Ünvan bilgilerini alma hatası"})
	}

	return c.JSON(titles)
}
