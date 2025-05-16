package repository

import (
	"context"
	"time"
	"transport-app/config"
	"transport-app/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getKendaraanCollection() *mongo.Collection {
	return config.GetCollection("kendaraan")
}

func GetAllKendaraan(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getKendaraanCollection().Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var kendaraan []models.Kendaraan
	if err := cursor.All(ctx, &kendaraan); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(kendaraan)
}

func GetKendaraanByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var kendaraan models.Kendaraan

	err := getKendaraanCollection().FindOne(context.TODO(), bson.M{"_id": id}).Decode(&kendaraan)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kendaraan not found"})
	}

	return c.JSON(kendaraan)
}

func CreateKendaraan(c *fiber.Ctx) error {
	var kendaraan models.Kendaraan
	if err := c.BodyParser(&kendaraan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if kendaraan.ID == "" || kendaraan.NomorPolisi == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID dan Nomor Polisi wajib diisi"})
	}

	// Cek ID unik
	var existing models.Kendaraan
	err := getKendaraanCollection().FindOne(context.TODO(), bson.M{"_id": kendaraan.ID}).Decode(&existing)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID sudah terdaftar, gunakan ID lain"})
	}

	_, err = getKendaraanCollection().InsertOne(context.TODO(), kendaraan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(kendaraan)
}

func UpdateKendaraan(c *fiber.Ctx) error {
	id := c.Params("id")
	var kendaraan models.Kendaraan
	if err := c.BodyParser(&kendaraan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	update := bson.M{"$set": kendaraan}
	_, err := getKendaraanCollection().UpdateByID(context.TODO(), id, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data kendaraan diupdate"})
}

func DeleteKendaraan(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := getKendaraanCollection().DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data kendaraan dihapus"})
}
