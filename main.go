package main

import (
	"fmt"
	"log"
	"os"

	"transport-app/config"
	"transport-app/middleware"
	"transport-app/routes"

	_ "transport-app/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"

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

// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load .env file hanya jika TIDAK sedang di Railway
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Gagal memuat file .env")
		} else {
			fmt.Println("✅ File .env dimuat (local) ")
		}
	}

	config.ConnectDB()

	app := fiber.New()

	middleware.SetupCORS(app)
	middleware.SetupLogger(app)

	app.Get("/docs/*", swagger.HandlerDefault)

	routes.SetupRoutes(app)

	// WAJIB pakai PORT dari env agar Railway bisa deteksi ini web service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("✅ Server running on port:", port)
	log.Fatal(app.Listen(":" + port))
}
