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

	// Ambil semua jadwal
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("❌ Error saat mengambil jadwal:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data jadwal"})
	}

	var jadwals []models.Jadwal
	if err := cursor.All(ctx, &jadwals); err != nil {
		fmt.Println("❌ Gagal parsing jadwal:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal parsing data jadwal"})
	}

	var result []JadwalWithRute

	for _, j := range jadwals {
		fmt.Println("🔍 Memproses jadwal dengan ID:", j.ID.Hex(), "dan RuteID:", j.RuteID)

		var rute models.Rute
		err = ruteCollection.FindOne(ctx, bson.M{"_id": j.RuteID}).Decode(&rute)
		if err != nil {
			fmt.Println("⚠️ Error saat mencari rute:", err)
			rute = models.Rute{}
		}

		result = append(result, JadwalWithRute{
			ID:             j.ID.Hex(),
			Tanggal:        j.Tanggal,
			WaktuBerangkat: j.WaktuBerangkat,
			EstimasiTiba:   j.EstimasiTiba,
			RuteID:         j.RuteID.Hex(),
			Rute:           rute,
		})
	}

	return c.JSON(result)
}

func GetJadwalByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var jadwal models.Jadwal
	err = getJadwalCollection().FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&jadwal)
	if err != nil {
		fmt.Println("❌ Jadwal not found with ID:", id)
		return c.Status(404).JSON(fiber.Map{"error": "Jadwal not found"})
	}

	return c.JSON(jadwal)
}

func CreateJadwal(c *fiber.Ctx) error {
	var input struct {
		Tanggal        string `json:"tanggal"`
		WaktuBerangkat string `json:"waktu_berangkat"`
		EstimasiTiba   string `json:"estimasi_tiba"`
		KodeRute       string `json:"kode_rute"`    // Input kode_rute dari frontend
		NomorPolisi    string `json:"nomor_polisi"` // Input nomor_polisi dari frontend
	}

	if err := c.BodyParser(&input); err != nil {
		fmt.Println("❌ Error parsing body:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Validasi field penting
	if input.Tanggal == "" || input.WaktuBerangkat == "" || input.EstimasiTiba == "" || input.KodeRute == "" || input.NomorPolisi == "" {
		fmt.Println("❌ Field validation failed")
		return c.Status(400).JSON(fiber.Map{
			"error": "Semua field wajib diisi",
		})
	}

	// Cari rute_id berdasarkan kode_rute
	var rute models.Rute
	err := config.GetCollection("rutes").FindOne(context.TODO(), bson.M{"kode_rute": input.KodeRute}).Decode(&rute)
	if err != nil {
		fmt.Println("❌ Rute not found with kode_rute:", input.KodeRute)
		return c.Status(404).JSON(fiber.Map{"error": "Rute not found"})
	}

	// Cari kendaraan_id berdasarkan nomor_polisi
	var kendaraan models.Kendaraan
	err = config.GetCollection("kendaraan").FindOne(context.TODO(), bson.M{"nomor_polisi": input.NomorPolisi}).Decode(&kendaraan)
	if err != nil {
		fmt.Println("❌ Kendaraan not found with nomor_polisi:", input.NomorPolisi)
		return c.Status(404).JSON(fiber.Map{"error": "Kendaraan not found"})
	}

	// Buat jadwal baru
	jadwal := models.Jadwal{
		ID:             primitive.NewObjectID(),
		Tanggal:        input.Tanggal,
		WaktuBerangkat: input.WaktuBerangkat,
		EstimasiTiba:   input.EstimasiTiba,
		RuteID:         rute.ID,      // Gunakan rute_id yang ditemukan
		KendaraanID:    kendaraan.ID, // Gunakan kendaraan_id yang ditemukan
	}

	_, err = getJadwalCollection().InsertOne(context.TODO(), jadwal)
	if err != nil {
		fmt.Println("❌ Error saat menyimpan jadwal:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("✅ Jadwal berhasil dibuat:", jadwal)
	return c.Status(201).JSON(jadwal)
}

func UpdateJadwal(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input struct {
		Tanggal        string `json:"tanggal"`
		WaktuBerangkat string `json:"waktu_berangkat"`
		EstimasiTiba   string `json:"estimasi_tiba"`
		KodeRute       string `json:"kode_rute"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Cari rute berdasarkan kode_rute
	var rute models.Rute
	err = getRuteCollection().FindOne(context.TODO(), bson.M{"kode_rute": input.KodeRute}).Decode(&rute)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Rute not found"})
	}

	// Update jadwal
	update := bson.M{
		"$set": bson.M{
			"tanggal":         input.Tanggal,
			"waktu_berangkat": input.WaktuBerangkat,
			"estimasi_tiba":   input.EstimasiTiba,
			"kode_rute":       input.KodeRute,
			"rute_id":         rute.ID,
		},
	}
	_, err = getJadwalCollection().UpdateByID(context.TODO(), objID, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Ambil ulang jadwal yang sudah diupdate
	var jadwal models.Jadwal
	err = getJadwalCollection().FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&jadwal)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return jadwal with rute info using a response struct
	response := struct {
		models.Jadwal
		Rute models.Rute `json:"rute"`
	}{
		Jadwal: jadwal,
		Rute:   rute,
	}

	return c.Status(200).JSON(response)
}

func DeleteJadwal(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Invalid ID format:", id)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = getJadwalCollection().DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println("❌ Error saat menghapus jadwal:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Jadwal berhasil dihapus"})
}
