package routes

import (
	"net/http"

	"github.com/KishiEdward/backend-bengkel.git/handlers"
	"github.com/KishiEdward/backend-bengkel.git/middleware" // Pastikan import ini ada
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Endpoint Ping - Public
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Pong! Server Bengkel Bubut menyala."})
	})

	// Inisialisasi Handler
	pesananHandler := handlers.NewPesananHandler()
	transaksiHandler := handlers.NewTransaksiHandler()
	authHandler := handlers.NewAuthHandler()
	customerHandler := handlers.NewCustomerHandler()
	materialHandler := handlers.NewMaterialHandler()

	// Grouping API v1
	v1 := r.Group("/v1")
	{
		// 1. Rute Public (Tanpa Login)
		v1.POST("/login", authHandler.Login)

		// 2. Rute Protected (Wajib Login/Bearer Token)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware()) // Mengunci semua rute di dalam grup ini
		{
			// Rute untuk pesanan
			protected.POST("/pesanan", pesananHandler.Create)
			protected.GET("/pesanan", pesananHandler.GetAll)
			protected.GET("/pesanan/laporan", pesananHandler.GetLaporan)
			protected.GET("/pesanan/:id", pesananHandler.GetDetail)
			protected.PUT("/pesanan/:id/status", pesananHandler.UpdateStatus)

			// Rute untuk Transaksi
			protected.POST("/pembayaran", transaksiHandler.CatatPembayaran)
			protected.POST("/biaya-tambahan", transaksiHandler.CatatBiayaTambahan)

			// Rute untuk Customer
			protected.POST("/customers", customerHandler.Create)
			protected.GET("/customers", customerHandler.GetAll)

			// Rute untuk Material
			protected.POST("/materials", materialHandler.Create)
			protected.GET("/materials", materialHandler.GetAll)
		}
	}

	return r
}