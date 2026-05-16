package services

import (
	"errors"

	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/repositories"
)

type TransaksiService struct {
	transaksiRepo *repositories.TransaksiRepository
	pesananRepo   *repositories.PesananRepository
}

func NewTransaksiService() *TransaksiService {
	return &TransaksiService{
		transaksiRepo: repositories.NewTransaksiRepository(),
		pesananRepo:   repositories.NewPesananRepository(), // Memanggil repo pesanan untuk validasi
	}
}

func (s *TransaksiService) TambahPembayaran(pembayaran *models.Pembayaran) error {
	// Validasi: Pastikan ID Pesanan yang diinput valid dan ada di database
	_, err := s.pesananRepo.FindByID(pembayaran.PesananID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}
	return s.transaksiRepo.CreatePembayaran(pembayaran)
}

func (s *TransaksiService) TambahBiayaTambahan(biaya *models.BiayaTambahan) error {
	// Validasi: Pastikan ID Pesanan valid
	_, err := s.pesananRepo.FindByID(biaya.PesananID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}
	return s.transaksiRepo.CreateBiayaTambahan(biaya)
}