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
    // 1. Cari pesanan berdasarkan ID
    pesanan, err := s.pesananRepo.FindByID(pembayaran.PesananID)
    if err != nil {
        return errors.New("pesanan tidak ditemukan")
    }

    // 2. Simpan pembayaran ke database
    err = s.transaksiRepo.CreatePembayaran(pembayaran)
    if err != nil {
        return err // Jika gagal simpan pembayaran, hentikan proses
    }

    // =========================================================
    // 3. LOGIKA BARU: UPDATE STATUS PESANAN SETELAH BAYAR DP
    // =========================================================
    // Cek apakah tipe pembayarannya "DP" dan status pesanan masih "Menunggu DP"
    // (Pastikan variabel pembayaran.Tipe sesuai dengan nama di struct models Anda)
    if pembayaran.Tipe == "DP" && pesanan.Status == "Menunggu DP" {
        
        pesanan.Status = "WIP" // Atau ganti "Diproses", sesuaikan dengan istilah bengkel Anda
        
        // Panggil fungsi untuk meng-update pesanan ke database
        err = s.pesananRepo.Update(pesanan) 
        
        if err != nil {
            return errors.New("pembayaran berhasil, tapi gagal update status pesanan")
        }
    }

    return nil
}

func (s *TransaksiService) TambahBiayaTambahan(biaya *models.BiayaTambahan) error {
	// Validasi: Pastikan ID Pesanan valid
	_, err := s.pesananRepo.FindByID(biaya.PesananID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}
	return s.transaksiRepo.CreateBiayaTambahan(biaya)
}