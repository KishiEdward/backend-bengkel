// models/pesanan.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// 1. Model Customer
type Customer struct {
	gorm.Model
	Nama    string `gorm:"size:100;not null" json:"nama"`
	Alamat  string `gorm:"type:text" json:"alamat"`
	NoTelp  string `gorm:"size:20" json:"no_telp"`
}

// ==========================================
// STRUCT BARU PENANGKAP JSON DARI FLUTTER
// ==========================================
type JasaLainInput struct {
	NamaJasa string  `json:"nama_jasa"`
	Biaya    float64 `json:"biaya"`
}

// 2. Model Pesanan (Tabel Sentral)
type Pesanan struct {
	gorm.Model
	CustomerID   uint       `json:"customer_id"`
	Customer     Customer   `gorm:"foreignKey:CustomerID" json:"customer"`
	TglOrder     time.Time  `json:"tgl_order"`
	TglDeadline  time.Time  `json:"tgl_deadline"`
	TglSelesai   *time.Time `json:"tgl_selesai"` 
	Status       string     `gorm:"size:50;default:'Menunggu DP'" json:"status"`
	HargaJual    float64    `json:"harga_jual"`
	JasaDesainer bool       `gorm:"-" json:"jasa_desainer"`
	BiayaDesain  float64    `gorm:"-" json:"biaya_desain"`
	
	JasaCNC      bool       `gorm:"-" json:"jasa_cnc"`
	BiayaCNC     float64    `gorm:"-" json:"biaya_cnc"`

	// ==========================================
	// TAMBAHKAN PENAMPUNG ARRAY JASA LAINNYA DI SINI
	// ==========================================
	JasaLainnya  []JasaLainInput `gorm:"-" json:"jasa_lainnya"`

	// Relasi (Has Many) ke tabel operasional
	Pembayarans      []Pembayaran      `gorm:"foreignKey:PesananID" json:"pembayaran"`
	BiayaTambahans   []BiayaTambahan   `gorm:"foreignKey:PesananID" json:"biaya_tambahan"`
	PesananMaterials []PesananMaterial `gorm:"foreignKey:PesananID" json:"pesanan_material"`
}

// 3. Model Pembayaran
type Pembayaran struct {
	gorm.Model
	PesananID  uint      `json:"pesanan_id"`
	Tipe       string    `gorm:"size:50" json:"tipe"` 
	Jumlah     float64   `json:"jumlah"`
	Tgl        time.Time `json:"tgl"`
	BuktiBayar string    `gorm:"type:varchar(255)" json:"bukti_bayar"` 
}

// 4. Model Biaya Tambahan 
type BiayaTambahan struct {
	gorm.Model
	PesananID  uint    `json:"pesanan_id"`
	Kategori   string  `gorm:"size:50" json:"kategori"`
	Keterangan string  `gorm:"size:200" json:"keterangan"`
	Nominal    float64 `json:"nominal"`
}

// 5. Model Material dan Pivot
type Material struct {
	gorm.Model
	Nama         string  `gorm:"size:100" json:"nama"`
	Satuan       string  `gorm:"size:30" json:"satuan"`
	HargaDefault float64 `json:"harga_default"`
}

type PesananMaterial struct {
	gorm.Model
	PesananID   uint     `json:"pesanan_id"`
	MaterialID  uint     `json:"material_id"`
	Material    Material `gorm:"foreignKey:MaterialID" json:"material"`
	Qty         int      `json:"qty"`
	HargaSatuan float64  `json:"harga_satuan"`
}