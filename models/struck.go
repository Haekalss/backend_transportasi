package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rute struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	KodeRute string             `json:"kode_rute" bson:"kode_rute"`
	NamaRute string             `json:"nama_rute" bson:"nama_rute"`
	Asal     string             `json:"asal" bson:"asal"`
	Tujuan   string             `json:"tujuan" bson:"tujuan"`
	JarakKM  int                `json:"jarak_km" bson:"jarak_km"`
}

type Kendaraan struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NomorPolisi string             `json:"nomor_polisi" bson:"nomor_polisi"`
	Jenis       string             `json:"jenis" bson:"jenis"`
	Kapasitas   int                `json:"kapasitas" bson:"kapasitas"`
	Status      string             `json:"status" bson:"status"`
}

type Jadwal struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Tanggal        string             `json:"tanggal" bson:"tanggal"`
	WaktuBerangkat string             `json:"waktu_berangkat" bson:"waktu_berangkat"`
	EstimasiTiba   string             `json:"estimasi_tiba" bson:"estimasi_tiba"`
	RuteID         primitive.ObjectID `json:"rute_id" bson:"rute_id"`
	KendaraanID    primitive.ObjectID `json:"kendaraan_id" bson:"kendaraan_id"`
}
