package controllers

import (
	"context"
	"encoding/json"
	"hospital-management/models"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type JobGroupController struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func NewJobGroupController(db *gorm.DB, redisClient *redis.Client) *JobGroupController {
	return &JobGroupController{
		DB:          db,
		RedisClient: redisClient,
	}
}

func (controller *JobGroupController) GetJobGroups(c *fiber.Ctx) error {
	ctx := context.Background()

	// RedisClient ve DB'nin nil olmadığını kontrol et
	if controller.RedisClient == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Redis client is nil"})
	}
	if controller.DB == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection is nil"})
	}

	jobGroupsData, err := controller.RedisClient.Get(ctx, "job_groups").Result()
	if err == redis.Nil {
		// Redis'te veri yok, veritabanından alalım
		var jobGroups []models.JobGroup
		if err := controller.DB.Find(&jobGroups).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Meslek grupları alma hatası"})
		}

		jobGroupsDataBytes, err := json.Marshal(jobGroups)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Veri dönüştürme hatası"})
		}
		jD := string(jobGroupsDataBytes)

		if err := controller.RedisClient.Set(ctx, "job_groups", jD, 10*time.Minute).Err(); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Önbellek veri yazma hatası"})
		}

		return c.JSON(jobGroups)
	} else if err != nil {
		// Redis bağlantı hatası
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Redis bağlantı hatası"})
	}

	// Redis'te veri bulundu, deserialize edelim
	var jobGroups []models.JobGroup
	if err := json.Unmarshal([]byte(jobGroupsData), &jobGroups); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Önbellek veri alma hatası"})
	}
	return c.JSON(jobGroups)
}
