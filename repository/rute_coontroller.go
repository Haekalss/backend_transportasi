package repository

import (
	"context"
	"fmt"
	"time"
	"transport-app/config"
	"transport-app/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
	var rute models.Rute

	err := ruteCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&rute)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Rute not found"})
	}

	return c.JSON(rute)
}
func CreateRute(c *fiber.Ctx) error {
	var rute models.Rute
	if err := c.BodyParser(&rute); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if rute.ID == "" || rute.KodeRute == "" || rute.NamaRute == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID, kode_rute, dan nama_rute wajib diisi"})
	}

	// Cek ID unik
	var existing models.Rute
	err := getRuteCollection().FindOne(context.TODO(), bson.M{"_id": rute.ID}).Decode(&existing)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID sudah terdaftar, gunakan ID lain"})
	}

	_, err = getRuteCollection().InsertOne(context.TODO(), rute)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(rute)
}

func UpdateRute(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")
	var rute models.Rute
	if err := c.BodyParser(&rute); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	update := bson.M{
		"$set": rute,
	}

	_, err := ruteCollection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data berhasil diupdate"})
}
func DeleteRute(c *fiber.Ctx) error {
	ruteCollection := getRuteCollection()
	id := c.Params("id")

	_, err := ruteCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data berhasil dihapus"})
}
