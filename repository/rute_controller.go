package repository

import (
	"context"
	"fmt"
	"time"
	"transport-app/config"
	"transport-app/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getRuteCollection() *mongo.Collection {
	return config.GetCollection("rutes")
}

// GetAllRute godoc
// @Summary Get all rutes
// @Description Mengambil semua data rute
// @Tags Rute
// @Accept json
// @Produce json
// @Success 200 {array} models.Rute "Daftar semua rute"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/rutes [get]
// @Security BearerAuth
func GetAllRute(c *fiber.Ctx) error {
	fmt.Println("🔍 GET /api/rutes called")
	ruteCollection := getRuteCollection()

	if ruteCollection == nil {
		fmt.Println("❌ Collection is nil")
		return c.Status(500).JSON(fiber.Map{"error": "Rute collection is nil"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ruteCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("❌ Error saat Find:", err.Error())
		return c.Status(500).JSON(fiber.Map{"error": "Find error: " + err.Error()})
	}

	var rutes []models.Rute
	if err := cursor.All(ctx, &rutes); err != nil {
		fmt.Println("❌ Error parsing cursor:", err.Error())
		return c.Status(500).JSON(fiber.Map{"error": "Cursor decode error: " + err.Error()})
	}

	fmt.Println("✅ Data ditemukan:", rutes)
	return c.JSON(rutes)
}

// GetRuteByID godoc
// @Summary Get a rute by ID
// @Description Mengambil data rute berdasarkan ID
// @Tags Rute
// @Accept json
// @Produce json
// @Param id path string true "Rute ID"
// @Success 200 {object} models.Rute "Data rute yang ditemukan"
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 404 {object} models.ErrorResponse "Rute not found"
// @Router /api/rutes/{id} [get]
// @Security BearerAuth
func GetRuteByID(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var rute models.Rute
	err = ruteCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&rute)
	if err != nil {
		fmt.Println("❌ Rute not found with ID:", id)
		return c.Status(404).JSON(fiber.Map{"error": "Rute not found"})
	}

	return c.JSON(rute)
}

// CreateRute godoc
// @Summary Create a new rute
// @Description Membuat data rute baru
// @Tags Rute
// @Accept json
// @Produce json
// @Param rute body models.Rute true "Data rute baru"
// @Success 201 {object} models.Rute "Rute berhasil dibuat"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/rutes [post]
// @Security BearerAuth
func CreateRute(c *fiber.Ctx) error {
	var rute models.Rute
	if err := c.BodyParser(&rute); err != nil {
		fmt.Println("❌ Error parsing body:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if rute.KodeRute == "" || rute.NamaRute == "" || rute.Asal == "" || rute.Tujuan == "" || rute.JarakKM <= 0 {
		fmt.Println("❌ Field validation failed")
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field wajib diisi dan jarak harus lebih dari 0",
		})
	}

	// Set ID baru secara manual agar bisa dikembalikan di response
	rute.ID = primitive.NewObjectID()

	_, err := getRuteCollection().InsertOne(context.TODO(), rute)
	if err != nil {
		fmt.Println("❌ Error saat menyimpan rute:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("✅ Rute berhasil dibuat:", rute)
	return c.Status(201).JSON(rute)
}

// UpdateRute godoc
// @Summary Update an existing rute
// @Description Memperbarui data rute yang sudah ada berdasarkan ID
// @Tags Rute
// @Accept json
// @Produce json
// @Param id path string true "Rute ID"
// @Param rute body models.Rute true "Data rute yang akan diupdate"
// @Success 200 {object} models.SuccessResponse "Data berhasil diupdate"
// @Failure 400 {object} models.ErrorResponse "Invalid ID atau Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/rutes/{id} [put]
// @Security BearerAuth
func UpdateRute(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var rute models.Rute
	if err := c.BodyParser(&rute); err != nil {
		fmt.Println("❌ Error parsing body:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if rute.KodeRute == "" || rute.NamaRute == "" || rute.Asal == "" || rute.Tujuan == "" || rute.JarakKM <= 0 {
		fmt.Println("❌ Field validation failed")
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field wajib diisi dan jarak harus lebih dari 0",
		})
	}

	update := bson.M{"$set": bson.M{
		"kode_rute": rute.KodeRute,
		"nama_rute": rute.NamaRute,
		"asal":      rute.Asal,
		"tujuan":    rute.Tujuan,
		"jarak_km":  rute.JarakKM,
	}}

	_, err = ruteCollection.UpdateByID(context.TODO(), objID, update)
	if err != nil {
		fmt.Println("❌ Error saat mengupdate rute:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data berhasil diupdate"})
}

// DeleteRute godoc
// @Summary Delete a rute
// @Description Menghapus data rute berdasarkan ID
// @Tags Rute
// @Accept json
// @Produce json
// @Param id path string true "Rute ID"
// @Success 200 {object} models.SuccessResponse "Data berhasil dihapus"
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/rutes/{id} [delete]
// @Security BearerAuth
func DeleteRute(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = ruteCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println("❌ Error saat menghapus rute:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data berhasil dihapus"})
}
