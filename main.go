package main

import (
	"log"
	"os"

	"transport-app/config"
	"transport-app/middleware"
	"transport-app/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	app := fiber.New()

	middleware.SetupCORS(app)
	middleware.SetupLogger(app)

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
