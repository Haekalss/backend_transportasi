package models

type Kendaraan struct {
	ID          string `json:"_id,omitempty" bson:"_id,omitempty"`
	NomorPolisi string `json:"nomor_polisi" bson:"nomor_polisi"`
	Jenis       string `json:"jenis" bson:"jenis"`
	Kapasitas   int    `json:"kapasitas" bson:"kapasitas"`
	Status      string `json:"status" bson:"status"`
	RuteID      string `json:"rute_id" bson:"rute_id"`
}
