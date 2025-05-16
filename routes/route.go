package routes

import (
	"transport-app/repository"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Rute
	api.Get("/rutes", repository.GetAllRute)
	api.Get("/rutes/:id", repository.GetRuteByID)
	api.Post("/rutes", repository.CreateRute)
	api.Put("/rutes/:id", repository.UpdateRute)
	api.Delete("/rutes/:id", repository.DeleteRute)

	// Kendaraan
	api.Get("/kendaraans", repository.GetAllKendaraan)
	api.Get("/kendaraans/:id", repository.GetKendaraanByID)
	api.Post("/kendaraans", repository.CreateKendaraan)
	api.Put("/kendaraans/:id", repository.UpdateKendaraan)
	api.Delete("/kendaraans/:id", repository.DeleteKendaraan)

	// import repository jadwal
	api.Get("/jadwals", repository.GetAllJadwal)
	api.Get("/jadwals/:id", repository.GetJadwalByID)
	api.Post("/jadwals", repository.CreateJadwal)
	api.Put("/jadwals/:id", repository.UpdateJadwal)
	api.Delete("/jadwals/:id", repository.DeleteJadwal)

}
