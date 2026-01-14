package main

import (
	"backend/cache"
	"backend/db"
	"backend/models"
	"backend/router"
	"backend/worker"

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
