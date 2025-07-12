package main

import (
	"log"
	"os"

	"transport-app/config"
	"transport-app/middleware"
	"transport-app/routes"

	_ "transport-app/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// @title Transport App API
// @version 1.0
// @description This is a sample server for a transport app.

// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url https://github.com/Haekalss

// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
