package handlers

import (
	"net/http"

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

// CatatPembayaran menangani request JSON untuk pembayaran DP/Lunas
func (h *TransaksiHandler) CatatPembayaran(c *gin.Context) {
	var req models.Pembayaran
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

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