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

func getJadwalCollection() *mongo.Collection {
	return config.GetCollection("jadwal")
}

// JadwalWithRute is a struct to combine Jadwal with its related Rute
type JadwalWithRute struct {
	ID             string      `json:"id" bson:"_id"`
	Tanggal        string      `json:"tanggal"`
	WaktuBerangkat string      `json:"waktu_berangkat"`
	EstimasiTiba   string      `json:"estimasi_tiba"`
	RuteID         string      `json:"rute_id"`
	Rute           models.Rute `json:"rute"`
}

func GetAllJadwal(c *fiber.Ctx) error {
	collection := getJadwalCollection()
	ruteCollection := config.GetCollection("rutes")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var jadwals []models.Jadwal
	if err := cursor.All(ctx, &jadwals); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []JadwalWithRute

	for _, j := range jadwals {
		var rute models.Rute
		err := ruteCollection.FindOne(ctx, bson.M{"_id": j.RuteID}).Decode(&rute)
		if err != nil {
			fmt.Println("⚠️ Rute not found untuk jadwal:", j.ID)
		}

		result = append(result, JadwalWithRute{
			ID:             j.ID,
			Tanggal:        j.Tanggal,
			WaktuBerangkat: j.WaktuBerangkat,
			EstimasiTiba:   j.EstimasiTiba,
			RuteID:         j.RuteID,
			Rute:           rute,
		})
	}

	return c.JSON(result)
}

func GetJadwalByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var jadwal models.Jadwal

	err := getJadwalCollection().FindOne(context.TODO(), bson.M{"_id": id}).Decode(&jadwal)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}

	return c.JSON(jadwal)
}

func CreateJadwal(c *fiber.Ctx) error {
	var jadwal models.Jadwal
	if err := c.BodyParser(&jadwal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if jadwal.ID == "" || jadwal.RuteID == "" || jadwal.KendaraanID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID, rute_id, dan kendaraan_id wajib diisi"})
	}

	// Cek ID unik
	var existing models.Jadwal
	err := getJadwalCollection().FindOne(context.TODO(), bson.M{"_id": jadwal.ID}).Decode(&existing)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID sudah terdaftar, gunakan ID lain"})
	}

	_, err = getJadwalCollection().InsertOne(context.TODO(), jadwal)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(jadwal)
}

func UpdateJadwal(c *fiber.Ctx) error {
	id := c.Params("id")
	var jadwal models.Jadwal
	if err := c.BodyParser(&jadwal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	update := bson.M{"$set": jadwal}
	_, err := getJadwalCollection().UpdateByID(context.TODO(), id, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Jadwal berhasil diupdate"})
}

func DeleteJadwal(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := getJadwalCollection().DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Jadwal berhasil dihapus"})
}
