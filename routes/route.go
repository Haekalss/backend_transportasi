package routes

import (
	"transport-app/middleware"
	"transport-app/repository"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth --- Rute Publik ---
	api.Post("/register", repository.Register)
	api.Post("/login", repository.Login)

	// --- Rute untuk Semua User (user & admin) ---
	// Endpoint GET All bisa diakses oleh semua yang sudah login
	api.Get("/rutes", middleware.Protected(), repository.GetAllRute)
	api.Get("/kendaraans", middleware.Protected(), repository.GetAllKendaraan)
	api.Get("/jadwals", middleware.Protected(), repository.GetAllJadwal)

	// Endpoint GET by ID juga bisa diakses oleh semua yang sudah login
	api.Get("/rutes/:id", middleware.Protected(), repository.GetRuteByID)
	api.Get("/kendaraans/:id", middleware.Protected(), repository.GetKendaraanByID)
	api.Get("/jadwals/:id", middleware.Protected(), repository.GetJadwalByID)


	// --- Rute admin ---
	// Rute
	api.Post("/rutes",middleware.Protected(), middleware.AdminOnly(), repository.CreateRute)
	api.Put("/rutes/:id",middleware.Protected(), middleware.AdminOnly(), repository.UpdateRute)
	api.Delete("/rutes/:id",middleware.Protected(), middleware.AdminOnly(), repository.DeleteRute)

	// Kendaraan
	api.Post("/kendaraans",middleware.Protected(), middleware.AdminOnly(), repository.CreateKendaraan)
	api.Put("/kendaraans/:id",middleware.Protected(), middleware.AdminOnly(), repository.UpdateKendaraan)
	api.Delete("/kendaraans/:id",middleware.Protected(), middleware.AdminOnly(), repository.DeleteKendaraan)

	// import repository jadwal
	api.Post("/jadwals",middleware.Protected(), middleware.AdminOnly(), repository.CreateJadwal)
	api.Put("/jadwals/:id",middleware.Protected(), middleware.AdminOnly(), repository.UpdateJadwal)
	api.Delete("/jadwals/:id",middleware.Protected(), middleware.AdminOnly(), repository.DeleteJadwal)

}
