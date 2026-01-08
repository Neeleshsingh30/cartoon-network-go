package main

import (
	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
	"cartoon-network-go/backend/src/router"
)

func main() {
	db.ConnectDB()

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
