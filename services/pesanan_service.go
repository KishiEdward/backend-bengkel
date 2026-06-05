package services

import (
    "errors"
    "fmt" // Tambahan baru untuk memformat tanggal
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
            PesananID:  pesanan.ID,
            Kategori:   "Jasa",
            Keterangan: "Jasa Desainer",
            Nominal:    pesanan.BiayaDesain,
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
    // 1. Ambil SELURUH pesanan (Bukan cuma yang selesai, agar kita bisa hitung Piutang dari pesanan WIP)
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

    // 2. Lakukan iterasi dan filter di dalam memori
    for _, p := range pesanans {
        // --- FILTER TANGGAL ---
        // Ubah TglOrder ke bentuk string (mengamankan format, baik itu string atau time.Time dari model)
        tglStr := fmt.Sprintf("%v", p.TglOrder) 
        
        if tahun != "Semua" && tahun != "" {
            if !strings.HasPrefix(tglStr, tahun) {
                continue // Skip pesanan ini jika tahunnya beda
            }
        }
        
        if bulan != "Semua" && bulan != "" {
            // Format tanggal umumnya "YYYY-MM-DD", index 5 sampai 7 adalah letak bulannya
            if len(tglStr) >= 7 {
                bulanPesanan := tglStr[5:7]
                if bulanPesanan != bulan {
                    continue // Skip pesanan ini jika bulannya beda
                }
            }
        }

        // --- KALKULASI RINCIAN HPP ---
        var totalMaterial float64 = 0
        for _, m := range p.PesananMaterials {
            totalMaterial += float64(m.Qty) * m.HargaSatuan
        }

        var totalJasa float64 = 0
        var totalTambahan float64 = 0
        for _, b := range p.BiayaTambahans {
            if strings.ToLower(b.Kategori) == "jasa" {
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
    }

    return listLaporan, ringkasanGlobal, nil
}