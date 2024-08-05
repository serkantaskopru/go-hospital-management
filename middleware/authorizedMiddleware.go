package middleware

import (
	"hospital-management/models"

	"github.com/gofiber/fiber/v2"
)

func AuthorizedMiddleware(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Yetkisiz işlem",
		})
	}

	if user.Staff.Role != "authorized" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlemi gerçekleştirmeniz için yetkiniz yok",
		})
	}

	return c.Next()
}
