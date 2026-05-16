package main

import (
	"log"

	"github.com/KishiEdward/backend-bengkel.git/config"
	"github.com/KishiEdward/backend-bengkel.git/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	config.InitDatabase()

	// Menggunakan SetupRouter dari folder routes
	r := routes.SetupRouter()

	log.Println("Server berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}