package main

import (
	"cartoon-network-go/backend/src/cache"
	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
	"cartoon-network-go/backend/src/router"
	"cartoon-network-go/backend/src/worker"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("backend/src/.env")

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
