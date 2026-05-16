package repositories

import (
	"github.com/KishiEdward/backend-bengkel.git/config"
	"github.com/KishiEdward/backend-bengkel.git/models"
)

type TransaksiRepository struct{}

func NewTransaksiRepository() *TransaksiRepository {
	return &TransaksiRepository{}
}

// CreatePembayaran menyimpan data DP atau Pelunasan
func (r *TransaksiRepository) CreatePembayaran(pembayaran *models.Pembayaran) error {
	return config.DB.Create(pembayaran).Error
}

// CreateBiayaTambahan menyimpan pengeluaran ekstra saat WIP
func (r *TransaksiRepository) CreateBiayaTambahan(biaya *models.BiayaTambahan) error {
	return config.DB.Create(biaya).Error
}