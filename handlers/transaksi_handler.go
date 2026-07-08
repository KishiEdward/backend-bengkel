package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/services"
	"github.com/gin-gonic/gin"
)

type TransaksiHandler struct {
	transaksiService *services.TransaksiService
}

func NewTransaksiHandler() *TransaksiHandler {
	return &TransaksiHandler{
		transaksiService: services.NewTransaksiService(),
	}
}

// CatatPembayaran menangani request multipart/form-data untuk pembayaran beserta upload bukti
func (h *TransaksiHandler) CatatPembayaran(c *gin.Context) {
	// 1. Ambil data teks dari form-data
	pesananIDStr := c.PostForm("pesanan_id")
	tipe := c.PostForm("tipe")
	jumlahStr := c.PostForm("jumlah")
	tglStr := c.PostForm("tgl")

	// Konversi string ke tipe data yang sesuai
	pesananID, err := strconv.ParseUint(pesananIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format pesanan_id tidak valid"})
		return
	}

	jumlah, err := strconv.ParseFloat(jumlahStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format jumlah tidak valid"})
		return
	}

	// Parsing tanggal (default ke waktu sekarang jika gagal)
	tgl, err := time.Parse(time.RFC3339, tglStr)
	if err != nil {
		tgl = time.Now()
	}

	// 2. Tangani Upload File Gambar (jika ada)
	var buktiBayarURL string
	file, err := c.FormFile("bukti_bayar") // Nama field file dari Flutter harus "bukti_bayar"
	
	if err == nil {
		// Buat nama file unik menggunakan timestamp agar tidak bentrok jika nama file sama
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		
		// Tentukan lokasi penyimpanan file fisik
		savePath := filepath.Join("uploads", "bukti_bayar", fileName)
		
		// Simpan file ke direktori lokal server
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan gambar bukti pembayaran"})
			return
		}
		
		// Simpan URL/path relatifnya ke database
		buktiBayarURL = "/uploads/bukti_bayar/" + fileName
	}

	// 3. Masukkan ke dalam struct model Pembayaran
	req := models.Pembayaran{
		PesananID:  uint(pesananID),
		Tipe:       tipe,
		Jumlah:     jumlah,
		Tgl:        tgl,
		BuktiBayar: buktiBayarURL,
	}

	// 4. Lanjutkan proses simpan ke database via Service
	if err := h.transaksiService.TambahPembayaran(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Pembayaran berhasil dicatat", "data": req})
}

// CatatBiayaTambahan menangani request JSON untuk penambahan biaya produksi
func (h *TransaksiHandler) CatatBiayaTambahan(c *gin.Context) {
	var req models.BiayaTambahan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.transaksiService.TambahBiayaTambahan(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Biaya tambahan berhasil dicatat", "data": req})
}