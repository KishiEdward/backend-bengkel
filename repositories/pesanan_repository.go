package repositories

import (
	"strings"

	"github.com/KishiEdward/backend-bengkel.git/config"
	"github.com/KishiEdward/backend-bengkel.git/models"
)

type PesananRepository struct{}

func NewPesananRepository() *PesananRepository {
	return &PesananRepository{}
}

// Create: Menyimpan pesanan baru ke database
func (r *PesananRepository) Create(pesanan *models.Pesanan) error {
	return config.DB.Create(pesanan).Error
}

// FindAll: Mengambil semua data pesanan beserta data customernya
func (r *PesananRepository) FindAll() ([]models.Pesanan, error) {
	var pesanans []models.Pesanan
	// Preload digunakan untuk melakukan otomatis JOIN ke tabel Customer
	err := config.DB.Preload("Customer").Find(&pesanans).Error
	return pesanans, err
}

// FindByID: Mengambil satu pesanan spesifik beserta seluruh relasi biayanya
func (r *PesananRepository) FindByID(id uint) (*models.Pesanan, error) {
	var pesanan models.Pesanan
	// Preload semua tabel yang berelasi untuk kalkulasi margin nanti
	err := config.DB.
		Preload("Customer").
		Preload("Pembayarans").
		Preload("BiayaTambahans").
		Preload("PesananMaterials.Material").
		First(&pesanan, id).Error
	return &pesanan, err
}

// Update: Memperbarui status atau data pesanan (misal saat status jadi WIP)
func (r *PesananRepository) Update(pesanan *models.Pesanan) error {
	return config.DB.Save(pesanan).Error
}

// FindByStatus: Mengambil pesanan berdasarkan status untuk keperluan laporan
func (r *PesananRepository) FindByStatus(status string) ([]models.Pesanan, error) {
	var pesanans []models.Pesanan
	err := config.DB.
		Preload("Customer").
		Preload("Pembayarans").
		Preload("BiayaTambahans").
		Preload("PesananMaterials.Material").
		Where("LOWER(status) = ?", strings.ToLower(status)).
		Find(&pesanans).Error
	return pesanans, err
}