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

func getKendaraanCollection() *mongo.Collection {
	return config.GetCollection("kendaraan")
}

func GetAllKendaraan(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil semua kendaraan
	cursor, err := getKendaraanCollection().Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var kendaraanList []models.Kendaraan
	if err := cursor.All(ctx, &kendaraanList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(kendaraanList)
}

func GetKendaraanByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Konversi ID dari string ke ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var kendaraan models.Kendaraan
	err = getKendaraanCollection().FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&kendaraan)
	if err != nil {
		fmt.Println("❌ Kendaraan not found with ID:", id)
		return c.Status(404).JSON(fiber.Map{"error": "Kendaraan not found"})
	}

	return c.JSON(kendaraan)
}

func CreateKendaraan(c *fiber.Ctx) error {
	var input struct {
		NomorPolisi string `json:"nomor_polisi"`
		Jenis       string `json:"jenis"`
		Kapasitas   int    `json:"kapasitas"`
		Status      string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		fmt.Println("❌ Error parsing body:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Validasi field penting
	if input.NomorPolisi == "" || input.Jenis == "" || input.Kapasitas <= 0 || input.Status == "" {
		fmt.Println("❌ Field validation failed")
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field wajib diisi dan kapasitas harus lebih dari 0",
		})
	}

	// Buat kendaraan baru
	kendaraan := models.Kendaraan{
		ID:          primitive.NewObjectID(),
		NomorPolisi: input.NomorPolisi,
		Jenis:       input.Jenis,
		Kapasitas:   input.Kapasitas,
		Status:      input.Status,
	}

	_, err := getKendaraanCollection().InsertOne(context.TODO(), kendaraan)
	if err != nil {
		fmt.Println("❌ Error saat menyimpan kendaraan:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("✅ Kendaraan berhasil dibuat:", kendaraan)
	return c.Status(201).JSON(kendaraan)
}

func UpdateKendaraan(c *fiber.Ctx) error {
	id := c.Params("id")

	// Konversi ID dari string ke ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var kendaraan models.Kendaraan
	if err := c.BodyParser(&kendaraan); err != nil {
		fmt.Println("❌ Error parsing body:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Validasi field penting
	if kendaraan.NomorPolisi == "" || kendaraan.Jenis == "" || kendaraan.Kapasitas <= 0 || kendaraan.Status == "" {
		fmt.Println("❌ Field validation failed")
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field wajib diisi dan kapasitas harus lebih dari 0",
		})
	}

	update := bson.M{"$set": kendaraan}
	_, err = getKendaraanCollection().UpdateByID(context.TODO(), objID, update)
	if err != nil {
		fmt.Println("❌ Error saat mengupdate kendaraan:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("✅ Kendaraan berhasil diupdate:", kendaraan)
	return c.JSON(fiber.Map{"message": "Data kendaraan diupdate"})
}

func DeleteKendaraan(c *fiber.Ctx) error {
	id := c.Params("id")

	// Konversi ID dari string ke ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = getKendaraanCollection().DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println("❌ Error saat menghapus kendaraan:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("✅ Kendaraan berhasil dihapus:", id)
	return c.JSON(fiber.Map{"message": "Data kendaraan dihapus"})
}
