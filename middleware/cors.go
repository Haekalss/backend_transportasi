package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupCORS(app *fiber.App) {
	app.Use(cors.New(cors.Config{

		AllowOrigins: "http://localhost:5173",

		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}
