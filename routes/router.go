package routes

import (
	"net/http"

	"github.com/KishiEdward/backend-bengkel.git/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Endpoint Ping untuk test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Pong! Server Bengkel Bubut menyala."})
	})

	// Inisialisasi Handler
	pesananHandler := handlers.NewPesananHandler()
	transaksiHandler := handlers.NewTransaksiHandler()

	// Grouping API v1
	v1 := r.Group("/v1")
	{
		// Rute untuk pesanan
		v1.POST("/pesanan", pesananHandler.Create)
		v1.GET("/pesanan", pesananHandler.GetAll)
		v1.GET("/pesanan/laporan", pesananHandler.GetLaporan)
		v1.GET("/pesanan/:id", pesananHandler.GetDetail)
		v1.PUT("/pesanan/:id/status", pesananHandler.UpdateStatus)

		// Rute untuk Transaksi
		v1.POST("/pembayaran", transaksiHandler.CatatPembayaran)
		v1.POST("/biaya-tambahan", transaksiHandler.CatatBiayaTambahan)
	}

	return r
}