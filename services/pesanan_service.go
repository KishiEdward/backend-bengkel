package services

import (
	"errors"
	"fmt"
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
func (s *PesananService) CreatePesanan(pesanan *models.Pesanan) error {

	// =====================
	// SIMPAN PESANAN UTAMA
	// =====================
	if err := s.pesananRepo.Create(pesanan); err != nil {
		return err
	}

	// =====================
	// AUTO INPUT BIAYA DESIGNER
	// =====================
	if pesanan.JasaDesainer && pesanan.BiayaDesain > 0 {
		biayaDesigner := models.BiayaTambahan{
			PesananID:  pesanan.ID,
			Kategori:   "Jasa",
			Keterangan: "Jasa Desainer",
			Nominal:    pesanan.BiayaDesain,
		}

		if err := s.pesananRepo.CreateBiayaTambahan(&biayaDesigner); err != nil {
			return err
		}
	}

	// =====================
	// AUTO INPUT JASA CNC
	// =====================
	if pesanan.JasaCNC && pesanan.BiayaCNC > 0 {
		biayaCNC := models.BiayaTambahan{
			PesananID:  pesanan.ID,
			Kategori:   "jasa",
			Keterangan: "Jasa CNC",
			Nominal:    pesanan.BiayaCNC,
		}

		if err := s.pesananRepo.CreateBiayaTambahan(&biayaCNC); err != nil {
			return err
		}
	}

	// ==========================================
	// AUTO INPUT JASA LAINNYA (YANG DINAMIS)
	// ==========================================
	for _, jasa := range pesanan.JasaLainnya {
		if jasa.Biaya > 0 && jasa.NamaJasa != "" {
			biayaLain := models.BiayaTambahan{
				PesananID:  pesanan.ID,
				Kategori:   "jasa", // Masuk ke kategori Jasa karena disepakati di awal
				Keterangan: jasa.NamaJasa,
				Nominal:    jasa.Biaya,
			}

			if err := s.pesananRepo.CreateBiayaTambahan(&biayaLain); err != nil {
				return err
			}
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
		// SAMAKAN LOGIKA FILTER DENGAN FLUTTER
		kat := strings.ToLower(tambahan.Kategori)
		ket := strings.ToLower(tambahan.Keterangan)

		isJasa := kat == "jasa" || strings.Contains(ket, "jasa") || strings.Contains(ket, "designer") || strings.Contains(ket, "desainer") || strings.Contains(ket, "cnc")

		if isJasa {
			totalBiayaJasa += tambahan.Nominal
		} else {
			totalBiayaTambahan += tambahan.Nominal
		}
	}

	// 4. Hitung HPP Aktual dan Margin Dinamis
	hppAktual := totalBiayaMaterial + totalBiayaJasa + totalBiayaTambahan
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
		"total_biaya_jasa":     totalBiayaJasa,
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

// GetLaporan mengkalkulasi laporan keuangan secara dinamis berdasarkan filter bulan dan tahun
func (s *PesananService) GetLaporan(bulan string, tahun string) ([]map[string]interface{}, map[string]interface{}, error) {
	// 1. Ambil SELURUH pesanan
	pesanans, err := s.pesananRepo.FindAll()
	if err != nil {
		return nil, nil, err
	}

	var listLaporan []map[string]interface{}

	var grandTotalPendapatan float64 = 0
	var grandTotalHPP float64 = 0
	var grandTotalMargin float64 = 0

	var totalPiutang float64 = 0
	var jumlahPesanan int = 0
	var rincianMaterial float64 = 0
	var rincianJasa float64 = 0
	var rincianExtra float64 = 0

	// ==========================================
	// TAMBAHAN: Variabel Penghitung Status Pesanan
	// ==========================================
	var countSelesai int = 0
	var countMenungguPelunasan int = 0
	var countWIP int = 0
	var countMenungguDP int = 0
	var countBatal int = 0

	// 2. Lakukan iterasi dan filter di dalam memori
	for _, p := range pesanans {
		// --- FILTER TANGGAL ---
		// --- FILTER TANGGAL ---
		tglOrderStr := fmt.Sprintf("%v", p.TglOrder)
		tglDeadlineStr := fmt.Sprintf("%v", p.TglDeadline) // Mengambil tanggal deadline

		matchTahun := true
		matchBulan := true

		// Cek Tahun (Lolos jika Order atau Deadline sesuai)
		if tahun != "Semua" && tahun != "" {
			matchOrderTahun := strings.HasPrefix(tglOrderStr, tahun)
			matchDeadlineTahun := strings.HasPrefix(tglDeadlineStr, tahun)
			if !matchOrderTahun && !matchDeadlineTahun {
				matchTahun = false
			}
		}

		// Cek Bulan (Lolos jika Order atau Deadline sesuai)
		if bulan != "Semua" && bulan != "" {
			matchOrderBulan := false
			matchDeadlineBulan := false

			if len(tglOrderStr) >= 7 {
				if tglOrderStr[5:7] == bulan {
					matchOrderBulan = true
				}
			}
			if len(tglDeadlineStr) >= 7 {
				if tglDeadlineStr[5:7] == bulan {
					matchDeadlineBulan = true
				}
			}

			if !matchOrderBulan && !matchDeadlineBulan {
				matchBulan = false
			}
		}

		// Jika tidak cocok tahun atau bulan, skip iterasi ini
		if !matchTahun || !matchBulan {
			continue
		}

		// --- KALKULASI RINCIAN HPP ---
		var totalMaterial float64 = 0
		for _, m := range p.PesananMaterials {
			totalMaterial += float64(m.Qty) * m.HargaSatuan
		}

		var totalJasa float64 = 0
		var totalTambahan float64 = 0
		for _, b := range p.BiayaTambahans {
			// SAMAKAN LOGIKA FILTER DENGAN FLUTTER
			kat := strings.ToLower(b.Kategori)
			ket := strings.ToLower(b.Keterangan)

			isJasa := kat == "jasa" || strings.Contains(ket, "jasa") || strings.Contains(ket, "designer") || strings.Contains(ket, "desainer") || strings.Contains(ket, "cnc")

			if isJasa {
				totalJasa += b.Nominal
			} else {
				totalTambahan += b.Nominal
			}
		}

		hpp := totalMaterial + totalJasa + totalTambahan
		margin := p.HargaJual - hpp

		// --- KALKULASI PIUTANG (Sisa Tagihan) ---
		var totalTerbayar float64 = 0
		for _, bayar := range p.Pembayarans {
			totalTerbayar += bayar.Jumlah
		}
		sisaTagihan := p.HargaJual - totalTerbayar
		if sisaTagihan > 0 {
			totalPiutang += sisaTagihan
		}

		// ==========================================
		// TAMBAHAN: Hitung Status Pesanan yang Lolos Filter
		// ==========================================
		status := strings.ToLower(p.Status)
		switch status {
		case "selesai":
			countSelesai++
		case "menunggu pelunasan":
			countMenungguPelunasan++
		case "wip":
			countWIP++
		case "menunggu dp":
			countMenungguDP++
		case "batal":
			countBatal++
		}

		// --- REKAPITULASI (MASUK KE GRAND TOTAL) ---
		jumlahPesanan++
		grandTotalPendapatan += p.HargaJual
		grandTotalHPP += hpp
		grandTotalMargin += margin

		rincianMaterial += totalMaterial
		rincianJasa += totalJasa
		rincianExtra += totalTambahan

		// Masukkan data satuan ke list
		listLaporan = append(listLaporan, map[string]interface{}{
			"pesanan_id":    p.ID,
			"customer":      p.Customer.Nama,
			"tgl_order":     p.TglOrder,
			"harga_jual":    p.HargaJual,
			"hpp_aktual":    hpp,
			"margin_aktual": margin,
			"sisa_tagihan":  sisaTagihan,
		})
	}

	// 3. Bungkus semua variabel ke dalam map interface{}
	ringkasanGlobal := map[string]interface{}{
		"total_pendapatan": grandTotalPendapatan,
		"total_hpp":        grandTotalHPP,
		"total_margin":     grandTotalMargin,
		"total_piutang":    totalPiutang,
		"jumlah_pesanan":   jumlahPesanan,
		"rincian_material": rincianMaterial,
		"rincian_jasa":     rincianJasa,
		"rincian_extra":    rincianExtra,

		// ==========================================
		// TAMBAHAN: Kirim data status ke Flutter
		// ==========================================
		"count_selesai":            countSelesai,
		"count_menunggu_pelunasan": countMenungguPelunasan,
		"count_wip":                countWIP,
		"count_menunggu_dp":        countMenungguDP,
		"count_batal":              countBatal,
	}

	return listLaporan, ringkasanGlobal, nil
}