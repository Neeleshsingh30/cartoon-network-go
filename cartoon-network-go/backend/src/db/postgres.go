package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	// 1️⃣ Try environment variable first (production / server)
	dsn := os.Getenv("DB_URL")

	// 2️⃣ Fallback to local DB for development
	if dsn == "" {
		log.Println("⚠️ DB_URL not found, using local database config")
		dsn = "host=localhost user=postgres password=12345678 dbname=CN port=5432 sslmode=disable"
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	DB = database
	log.Println("✅ PostgreSQL Connected Successfully")
}
