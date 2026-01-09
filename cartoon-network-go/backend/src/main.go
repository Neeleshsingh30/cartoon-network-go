package main

import (
	"log"

	"cartoon-network-go/backend/src/cache"
	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
	"cartoon-network-go/backend/src/router"
	"cartoon-network-go/backend/src/worker"

	"github.com/joho/godotenv"
)

func main() {

	// 🔥 sabse pehle env load
	err := godotenv.Load("backend/src/.env")
	if err != nil {
		log.Println("⚠️ .env file not found, using system env variables")
	}

	db.ConnectDB()
	worker.StartWorkerPool()

	go cache.RefreshHomeCache()
	go cache.RefreshShowTimeCache()

	db.DB.AutoMigrate(
		&models.User{},
		&models.Cartoon{},
		&models.CartoonImage{},
		&models.Character{},
		&models.Like{},
		&models.CartoonView{},
		&models.Admin{},
		&models.AdminActivityLog{},
	)

	r := router.SetupRouter()
	r.Run(":8000")
}
