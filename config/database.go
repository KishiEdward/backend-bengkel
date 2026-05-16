package config

import (
	"fmt"
	"log"
	"os"

	"github.com/KishiEdward/backend-bengkel.git/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Format Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Menampilkan log query di terminal
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// Setup connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	// AutoMigrate membaca struct di folder models dan membuat tabelnya
	err = DB.AutoMigrate(
		&models.Customer{},
		&models.Pesanan{},
		&models.Pembayaran{},
		&models.BiayaTambahan{},
		&models.Material{},
		&models.PesananMaterial{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate gagal: %v", err)
	}

	log.Println("Database berhasil terhubung dan tabel sudah di-migrate!")
}