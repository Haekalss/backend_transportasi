package models

type Rute struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	KodeRute string `json:"kode_rute" bson:"kode_rute"`
	NamaRute string `json:"nama_rute" bson:"nama_rute"`
	Asal     string `json:"asal" bson:"asal"`
	Tujuan   string `json:"tujuan" bson:"tujuan"`
	JarakKM  int    `json:"jarak_km" bson:"jarak_km"`
}

type Kendaraan struct {
	ID          string `json:"_id,omitempty" bson:"_id,omitempty"`
	NomorPolisi string `json:"nomor_polisi" bson:"nomor_polisi"`
	Jenis       string `json:"jenis" bson:"jenis"`
	Kapasitas   int    `json:"kapasitas" bson:"kapasitas"`
	Status      string `json:"status" bson:"status"`
	RuteID      string `json:"rute_id" bson:"rute_id"`
}

type Jadwal struct {
	ID             string `json:"_id,omitempty" bson:"_id,omitempty"`
	Tanggal        string `json:"tanggal" bson:"tanggal"`
	WaktuBerangkat string `json:"waktu_berangkat" bson:"waktu_berangkat"`
	EstimasiTiba   string `json:"estimasi_tiba" bson:"estimasi_tiba"`
	RuteID         string `json:"rute_id" bson:"rute_id"`
	KendaraanID    string `json:"kendaraan_id" bson:"kendaraan_id"`
}
