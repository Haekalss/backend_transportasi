package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// Fungsi untk melindungi rute yang membutuhkan autentikasi
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: jwtError,
	})
}

// hanya admin yang bisa akses
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)

		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Admin access required",
			})
		}
		return c.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}