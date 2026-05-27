package services

import (
	"errors"
	"strings"
	"time"

	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/repositories"
)

type PesananService struct {
	pesananRepo *repositories.PesananRepository
}

func NewPesananService() *PesananService {
	return &PesananService{
		pesananRepo: repositories.NewPesananRepository(),
	}
}

// CreatePesanan menangani logika pembuatan order baru
func (s *PesananService) CreatePesanan(
	pesanan *models.Pesanan,
) error {

	// =====================
	// SIMPAN PESANAN
	// =====================
	if err := s.pesananRepo.Create(pesanan); err != nil {
		return err
	}

	// =====================
	// AUTO INPUT BIAYA DESIGNER
	// =====================
	if pesanan.JasaDesainer &&
		pesanan.BiayaDesain > 0 {

		biayaDesigner := models.BiayaTambahan{
			PesananID: pesanan.ID,
			Kategori: "Jasa",
			Keterangan: "Jasa Desainer",
			Nominal: pesanan.BiayaDesain,
		}

		if err := s.pesananRepo.CreateBiayaTambahan(
			&biayaDesigner,
		); err != nil {

			return err
		}
	}

	// =====================
	// AUTO INPUT JASA CNC
	// =====================
	if pesanan.JasaCNC &&
		pesanan.BiayaCNC > 0 {

		biayaCNC :=
			models.BiayaTambahan{

				PesananID:
					pesanan.ID,

				Kategori:
					"jasa",

				Keterangan:
					"Jasa CNC",

				Nominal:
					pesanan.BiayaCNC,
			}

		if err := s.pesananRepo.
			CreateBiayaTambahan(
				&biayaCNC,
			); err != nil {

			return err
		}
	}

	return nil
}

// GetAll mengambil ringkasan seluruh pesanan
func (s *PesananService) GetAll() ([]models.Pesanan, error) {
	return s.pesananRepo.FindAll()
}

// GetDetailDanHitungMargin adalah fitur utama untuk pelacakan margin dinamis
func (s *PesananService) GetDetailDanHitungMargin(id uint) (*models.Pesanan, map[string]float64, error) {
	// 1. Ambil data pesanan beserta seluruh relasinya dari Repository
	pesanan, err := s.pesananRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	// 2. Kalkulasi Total Biaya Material
	var totalBiayaMaterial float64 = 0
	for _, material := range pesanan.PesananMaterials {
		totalBiayaMaterial += float64(material.Qty) * material.HargaSatuan
	}

	// 3. Kalkulasi Total Biaya Tambahan (Overhead, Revisi, dll) saat proses WIP
	var totalBiayaJasa float64 = 0

	var totalBiayaTambahan float64 = 0
	
	for _, tambahan := range pesanan.BiayaTambahans {
	
		if tambahan.Kategori == "jasa" {
		
			totalBiayaJasa += tambahan.Nominal
		
		} else {
		
			totalBiayaTambahan += tambahan.Nominal
		}
	}

	// 4. Hitung HPP Aktual dan Margin Dinamis
	hppAktual :=
	totalBiayaMaterial +
		totalBiayaJasa +
		totalBiayaTambahan
	marginAktual := pesanan.HargaJual - hppAktual

	// 5. Hitung Sisa Tagihan (Total Pembayaran yang sudah masuk vs Harga Jual)
	var totalTerbayar float64 = 0
	for _, pembayaran := range pesanan.Pembayarans {
		totalTerbayar += pembayaran.Jumlah
	}
	sisaTagihan := pesanan.HargaJual - totalTerbayar

	// Bungkus hasil kalkulasi finansial ke dalam map
	kalkulasiKeuangan := map[string]float64{
		"total_biaya_material": totalBiayaMaterial,
		"total_biaya_jasa": totalBiayaJasa,
		"total_biaya_tambahan": totalBiayaTambahan,
		"hpp_aktual":           hppAktual,
		"margin_aktual":        marginAktual,
		"total_terbayar":       totalTerbayar,
		"sisa_tagihan":         sisaTagihan,
	}

	return pesanan, kalkulasiKeuangan, nil
}

// UpdateStatus memperbarui status pesanan
func (s *PesananService) UpdateStatus(id uint, status string) (*models.Pesanan, error) {
	// 1. Cari pesanannya dulu
	pesanan, err := s.pesananRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("pesanan tidak ditemukan")
	}

	// 2. Ubah statusnya
	pesanan.Status = status

	// 3. Logika Otomatis: Jika status diubah jadi "Selesai", catat tanggal selesainya hari ini
	if strings.ToLower(status) == "selesai" {
		now := time.Now()
		pesanan.TglSelesai = &now
	}

	// 4. Simpan perubahan ke database
	if err := s.pesananRepo.Update(pesanan); err != nil {
		return nil, err
	}

	return pesanan, nil
}

// GetLaporan mengkalkulasi total pendapatan, HPP, dan Margin dari pesanan yang sudah "selesai"
func (s *PesananService) GetLaporan() ([]map[string]interface{}, map[string]float64, error) {
	// 1. Ambil hanya pesanan yang sudah selesai
	pesanans, err := s.pesananRepo.FindByStatus("selesai")
	if err != nil {
		return nil, nil, err
	}

	var listLaporan []map[string]interface{}
	var grandTotalPendapatan float64 = 0
	var grandTotalHPP float64 = 0
	var grandTotalMargin float64 = 0

	// 2. Lakukan iterasi untuk menghitung HPP dan Margin tiap pesanan
	for _, p := range pesanans {
		var totalMaterial float64 = 0
		for _, m := range p.PesananMaterials {
			totalMaterial += float64(m.Qty) * m.HargaSatuan
		}

		var totalTambahan float64 = 0
		for _, b := range p.BiayaTambahans {
			totalTambahan += b.Nominal
		}

		hpp := totalMaterial + totalTambahan
		margin := p.HargaJual - hpp

		// Tambahkan ke akumulasi total keseluruhan (Grand Total)
		grandTotalPendapatan += p.HargaJual
		grandTotalHPP += hpp
		grandTotalMargin += margin

		// Masukkan ringkasan per pesanan ke dalam list laporan
		listLaporan = append(listLaporan, map[string]interface{}{
			"pesanan_id":    p.ID,
			"customer":      p.Customer.Nama,
			"tgl_selesai":   p.TglSelesai,
			"harga_jual":    p.HargaJual,
			"hpp_aktual":    hpp,
			"margin_aktual": margin,
		})
	}

	// 3. Bungkus Grand Total untuk mempermudah pembuatan grafik/ringkasan di Flutter nanti
	ringkasanGlobal := map[string]float64{
		"total_pendapatan": grandTotalPendapatan,
		"total_hpp":        grandTotalHPP,
		"total_margin":     grandTotalMargin,
	}

	return listLaporan, ringkasanGlobal, nil
}