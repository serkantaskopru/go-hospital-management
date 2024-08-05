package controllers

import (
	"hospital-management/models"
	"hospital-management/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type StaffController struct {
	DB *gorm.DB
}

func (controller *StaffController) CreateStaff(c *fiber.Ctx) error {
	type StaffInput struct {
		Name           string `json:"name" validate:"required"`
		Email          string `json:"email" validate:"required,email"`
		Password       string `json:"password" validate:"required"`
		Phone          string `json:"phone" validate:"required"`
		Role           string `json:"role" validate:"required"`
		IdentityNumber string `json:"identityNumber"`
	}

	user := c.Locals("user").(models.User)

	if user.Staff.Role != "authorized" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Bir kullanıcı oluşturmaya yetkiniz yok"})
	}

	var input StaffInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	hospitalID := user.Staff.HospitalID

	newUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: utils.HashPassword(input.Password),
		Phone:    input.Phone,
	}

	if err := controller.DB.Create(&newUser).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı oluşturma hatası"})
	}

	staff := models.Staff{
		UserID:         newUser.ID,
		HospitalID:     hospitalID,
		Role:           input.Role,
		IdentityNumber: input.IdentityNumber,
	}

	if err := controller.DB.Create(&staff).Error; err != nil {
		log.Printf("Failed to create staff: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Yetkili oluşturma hatası"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Yetkili başarıyla oluşturuldu"})
}

func (controller *StaffController) UpdateStaff(c *fiber.Ctx) error {
	i := c.Params("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz kullanıcı verisi"})
	}

	type StaffInput struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		Role           string `json:"role"`
		IdentityNumber string `json:"identityNumber"`
	}

	var input StaffInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	user := c.Locals("user").(models.User)
	if user.Staff.Role != "authorized" && user.Staff.ID != id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Bu yetkiliyi güncellemek için izniniz yok"})
	}

	hospitalID := user.Staff.HospitalID

	staffUpdateData := map[string]interface{}{
		"role":            input.Role,
		"hospital_id":     hospitalID,
		"identity_number": input.IdentityNumber,
	}

	if err := controller.DB.Model(&models.Staff{}).Where("id = ?", id).Updates(staffUpdateData).Error; err != nil {
		log.Printf("Failed to update staff: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Yetkili güncelleme hatası"})
	}

	userUpdateData := map[string]interface{}{
		"name":  input.Name,
		"email": input.Email,
		"phone": input.Phone,
	}

	if err := controller.DB.Model(&models.User{}).Where("id = ?", user.Staff.UserID).Updates(userUpdateData).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı güncelleme hatası"})
	}

	return c.JSON(fiber.Map{"message": "Yetkili başarıyla güncellendi"})
}

func (controller *StaffController) DeleteStaff(c *fiber.Ctx) error {
	i := c.Params("id")

	id, err := strconv.Atoi(i)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz yetkili verisi"})
	}

	user := c.Locals("user").(models.User)
	if user.Staff.Role != "authorized" && user.Staff.ID != id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Bu yetkiliyi kaldırmanız için izniniz yok"})
	}

	if err := controller.DB.Delete(&models.Staff{}, id).Error; err != nil {
		log.Printf("Failed to delete staff: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Yetkili kaldırma hatası"})
	}

	return c.JSON(fiber.Map{"message": "Yetkili başarıyla sistemden kaldırıldı"})
}

func (controller *StaffController) ListStaff(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		log.Println("Authorized user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Yetkisiz işlem",
		})
	}

	if user.Staff == nil {
		log.Println("Staff information not loaded")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Yetkisiz işlem",
		})
	}

	var staff []models.Staff

	if err := controller.DB.
		Preload("User").
		Preload("Hospital").
		Where("hospital_id = ?", user.Staff.HospitalID).
		Find(&staff).Error; err != nil {
		log.Printf("Failed to fetch staff: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Yetkili listesini alma hatası"})
	}

	return c.JSON(fiber.Map{"data": staff})
}
