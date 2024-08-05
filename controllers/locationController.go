package controllers

import (
	"context"
	"encoding/json"
	"hospital-management/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LocationController struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func (lc *LocationController) GetCities(c *fiber.Ctx) error {
	citiesData, err := lc.RedisClient.Get(context.Background(), "cities").Result()
	if err == nil && citiesData != "" {
		var cities []models.City
		if err := json.Unmarshal([]byte(citiesData), &cities); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Önbellek veri alma hatası"})
		}
		return c.JSON(cities)
	}

	var cities []models.City
	if err := lc.DB.Preload("Districts").Find(&cities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	citiesDataBytes, err := json.Marshal(cities)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Veri dönüştürme hatası"})
	}
	cD := string(citiesDataBytes)

	if err := lc.RedisClient.Set(context.Background(), "cities", cD, 10*time.Minute).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Önbellek veri yazma hatası"})
	}

	return c.JSON(cities)
}

func (lc *LocationController) GetDistrictsByCity(c *fiber.Ctx) error {
	cityID := c.Params("cityId")

	districtsData, err := lc.RedisClient.Get(context.Background(), "districts:"+cityID).Result()
	if err == nil && districtsData != "" {
		var districts []models.District
		if err := json.Unmarshal([]byte(districtsData), &districts); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Önbellek veri alma hatası"})
		}
		return c.JSON(districts)
	}

	var districts []models.District
	if err := lc.DB.Where("city_id = ?", cityID).Find(&districts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	districtsDataBytes, err := json.Marshal(districts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Veri dönüştürme hatası"})
	}
	dD := string(districtsDataBytes)

	if err := lc.RedisClient.Set(context.Background(), "districts:"+cityID, dD, 10*time.Minute).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Önbellek veri yazma hatası"})
	}

	return c.JSON(districts)
}
