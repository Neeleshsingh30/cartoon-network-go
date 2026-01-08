package main

import (
	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
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

}
