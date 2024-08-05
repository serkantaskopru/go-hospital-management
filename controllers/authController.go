package controllers

import (
	"fmt"
	"hospital-management/models"
	"hospital-management/utils"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func (controller *AuthController) Login(c *fiber.Ctx) error {
	var input struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	var user models.User

	if err := controller.DB.Where("email = ? OR phone = ?", input.Identifier, input.Identifier).Preload("Staff").First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Eposta veya telefon hatalı"})
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Şifre hatalı"})
	}

	if user.Staff == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgileri alınamadı"})
	}

	var activeSessions []models.Session
	if err := controller.DB.Where("hospital_id = ? AND expires_at > ?", user.Staff.HospitalID, time.Now()).Find(&activeSessions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Aktif oturumları okuma hatası"})
	}

	if len(activeSessions) >= 2 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Hastaneye ait 2 aktif bağlantı var, oturum açabilmek için diğer kullanıcıların çıkış yapmasını bekleyin."})
	}

	newSession := models.Session{
		UserID:     uint(user.ID),
		HospitalID: uint(user.Staff.HospitalID),
	}

	if err := controller.DB.Create(&newSession).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Oturum başlatma hatası"})
	}

	token, err := createJWT(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token oluşturma hatası"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (controller *AuthController) Register(c *fiber.Ctx) error {
	type RegisterInput struct {
		Name           string `json:"name" validate:"required"`
		Email          string `json:"email" validate:"required,email"`
		IdentityNumber string `json:"identityNumber" validate:"required"`
		Password       string `json:"password" validate:"required"`
		Phone          string `json:"phone" validate:"required"`
		TaxNumber      string `json:"taxNumber" validate:"required"`
		HospitalMail   string `json:"hospitalMail" validate:"required,email"`
		HospitalPhone  string `json:"hospitalPhone" validate:"required"`
		HospitalName   string `json:"hospitalName" validate:"required"`
		Address        string `json:"address" validate:"required"`
		City           string `json:"city" validate:"required"`
		District       string `json:"district" validate:"required"`
	}

	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	defaultUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: utils.HashPassword(input.Password),
		Phone:    input.Phone,
	}

	if err := controller.DB.Create(&defaultUser).Error; err != nil {
		log.Fatalf("Failed to create default user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı oluşturma hatası"})
	}

	var city models.City
	var district models.District

	cityId, err := strconv.Atoi(input.City)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz şehir formatı"})
	}

	districtId, err := strconv.Atoi(input.District)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz ilçe formatı"})
	}

	if err := controller.DB.First(&city, uint(cityId)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Şehir bulunamadı"})
	}

	if err := controller.DB.First(&district, uint(districtId)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "İlçe bulunamadı"})
	}

	taxNumber, err := strconv.Atoi(input.TaxNumber)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz vergi numarası formatı"})
	}

	defaultHospital := models.Hospital{
		Name:      input.HospitalName,
		Mail:      input.HospitalMail,
		Phone:     input.HospitalPhone,
		TaxNumber: uint(taxNumber),
		Address:   input.Address + ", " + city.Name + ", " + district.Name,
	}

	if err := controller.DB.Create(&defaultHospital).Error; err != nil {
		log.Fatalf("Failed to create default hospital: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Hastane bilgileri oluşturma hatası"})
	}

	defaultStaff := models.Staff{
		UserID:         defaultUser.ID,
		HospitalID:     defaultHospital.ID,
		Role:           "Staff",
		IdentityNumber: input.IdentityNumber,
	}

	if err := controller.DB.Create(&defaultStaff).Error; err != nil {
		log.Fatalf("Failed to create default staff: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Yetkili oluşturma hatası"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Kullanıcı başarıyla oluşturuldu"})
}

func (controller *AuthController) Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		log.Println("Authorized user not found or invalid type")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Kullanıcı verisi alınamadı",
		})
	}

	if err := controller.DB.Where("user_id = ?", user.ID).Delete(&models.Session{}).Error; err != nil {
		log.Printf("Failed to delete session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Oturum kapatma hatası"})
	}

	return c.JSON(fiber.Map{"message": "Oturum başarıyla kapatıldı"})
}

func createJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Minute * 45).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func (controller *AuthController) VerifyToken(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Token boş"})
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz token"})
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Token geçerli"})
	} else {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz token"})
	}
}
func (controller *AuthController) RequestPasswordReset(c *fiber.Ctx) error {
	type Request struct {
		Phone string `json:"phone"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	resetCode := models.ResetCode{
		Phone:     req.Phone,
		Code:      generateRandomCode(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	var existingCode models.ResetCode
	err := controller.DB.Where("phone = ?", req.Phone).First(&existingCode).Error
	if err == nil {
		resetCode.ID = existingCode.ID
		if err := controller.DB.Save(&resetCode).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sıfırlama kodu güncelleme hatası"})
		}
	} else if err == gorm.ErrRecordNotFound {
		if err := controller.DB.Create(&resetCode).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sıfırlama kodu oluşturma hatası"})
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sıfırlama kodu kontrol hatası"})
	}

	return c.JSON(fiber.Map{"reset_code": resetCode.Code})
}
func generateRandomCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
func (controller *AuthController) ResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Phone           string `json:"phone"`
		Code            string `json:"code"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	var resetCode models.ResetCode
	if err := controller.DB.Where("phone = ? AND code = ?", req.Phone, req.Code).First(&resetCode).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz sıfırlama kodu"})
	}

	if time.Now().After(resetCode.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Sıfırlama kodu kullanım süresi geçersiz"})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Şifreler eşleşmiyor"})
	}

	var user models.User
	if err := controller.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
	}

	user.Password = utils.HashPassword(req.NewPassword)
	if err := controller.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Şifre güncelleme hatası"})
	}

	if err := controller.DB.Delete(&resetCode).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sıfırlama kodu silme hatası"})
	}

	return c.JSON(fiber.Map{"message": "Şifre başarıyla güncellendi"})
}
