package app

import (
	"hospital-management/controllers"
	"hospital-management/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAuthRoutes(c *fiber.App, db *gorm.DB) {
	authController := &controllers.AuthController{DB: db}

	auth := c.Group("/auth")
	{
		auth.Get("/verify-token", authController.VerifyToken)
		auth.Post("/login", authController.Login)
		auth.Post("/request-password-reset", authController.RequestPasswordReset)
		auth.Post("/reset-password", authController.ResetPassword)
		auth.Post("/register", authController.Register)
		auth.Post("/logout", middleware.AuthMiddleware, authController.Logout)
	}

}
