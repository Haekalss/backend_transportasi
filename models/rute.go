package models

type Rute struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	KodeRute string `json:"kode_rute" bson:"kode_rute"`
	NamaRute string `json:"nama_rute" bson:"nama_rute"`
	Asal     string `json:"asal" bson:"asal"`
	Tujuan   string `json:"tujuan" bson:"tujuan"`
	JarakKM  int    `json:"jarak_km" bson:"jarak_km"`
}
