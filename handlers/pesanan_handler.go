package handlers

import (
	"net/http"
	"strconv"

	"github.com/KishiEdward/backend-bengkel.git/models"
	"github.com/KishiEdward/backend-bengkel.git/services"
	"github.com/gin-gonic/gin"
)

type PesananHandler struct {
    pesananService *services.PesananService
}

func NewPesananHandler() *PesananHandler {
    return &PesananHandler{
        pesananService: services.NewPesananService(),
    }
}

// Create menerima data JSON dari Flutter untuk membuat pesanan baru
func (h *PesananHandler) Create(c *gin.Context) {
    var pesanan models.Pesanan

    // Bind JSON request ke struct Pesanan
    if err := c.ShouldBindJSON(&pesanan); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
        return
    }

    if err := h.pesananService.CreatePesanan(&pesanan); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat pesanan"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Pesanan berhasil dibuat", "data": pesanan})
}

// GetAll mengembalikan semua data pesanan
func (h *PesananHandler) GetAll(c *gin.Context) {
    pesanans, err := h.pesananService.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "data": pesanans})
}

// GetDetail mengembalikan detail pesanan beserta kalkulasi margin dinamisnya
func (h *PesananHandler) GetDetail(c *gin.Context) {
    // Ambil ID dari parameter URL (misal: /v1/pesanan/1)
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID Pesanan tidak valid"})
        return
    }

    pesanan, kalkulasi, err := h.pesananService.GetDetailDanHitungMargin(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Pesanan tidak ditemukan"})
        return
    }

    // Menggabungkan data pesanan dan kalkulasi keuangan dalam satu response JSON
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "pesanan":  pesanan,
            "keuangan": kalkulasi,
        },
    })
}

// UpdateStatus menangani request JSON untuk mengubah status pesanan
func (h *PesananHandler) UpdateStatus(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID Pesanan tidak valid"})
        return
    }

    // Kita bikin struct kecil sementara hanya untuk menangkap input "status"
    var req struct {
        Status string `json:"status" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Kolom status wajib diisi"})
        return
    }

    pesanan, err := h.pesananService.UpdateStatus(uint(id), req.Status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status pesanan berhasil diperbarui", "data": pesanan})
}

// GetLaporan menangani request JSON untuk halaman laporan keuangan bengkel
func (h *PesananHandler) GetLaporan(c *gin.Context) {
    // Tangkap parameter filter dari URL (misal: /v1/laporan?bulan=05&tahun=2026)
    // Jika tidak ada yang dikirim, default-nya adalah "Semua"
    bulan := c.DefaultQuery("bulan", "Semua")
    tahun := c.DefaultQuery("tahun", "Semua")

    // Masukkan parameter tersebut ke dalam pemanggilan fungsi service
    listLaporan, ringkasan, err := h.pesananService.GetLaporan(bulan, tahun)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menghasilkan laporan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Laporan keuangan berhasil ditarik",
        "data": gin.H{
            "ringkasan_global": ringkasan,
            "detail_pesanan":   listLaporan,
        },
    })
}