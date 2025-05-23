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

func GetRuteByID(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")

	// Konversi ID dari string ke ObjectID
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
func CreateRute(c *fiber.Ctx) error {
	var rute models.Rute
	if err := c.BodyParser(&rute); err != nil {
		fmt.Println("❌ Error parsing body:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Validasi field penting
	if rute.KodeRute == "" || rute.NamaRute == "" || rute.Asal == "" || rute.Tujuan == "" || rute.JarakKM <= 0 {
		fmt.Println("❌ Field validation failed")
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field wajib diisi dan jarak harus lebih dari 0",
		})
	}

	_, err := getRuteCollection().InsertOne(context.TODO(), rute)
	if err != nil {
		fmt.Println("❌ Error saat menyimpan rute:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("✅ Rute berhasil dibuat:", rute)
	return c.Status(201).JSON(rute)
}

func UpdateRute(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")

	// Konversi ID dari string ke ObjectID
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

	// Validasi field penting
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

func DeleteRute(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")

	// Konversi ID dari string ke ObjectID
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
