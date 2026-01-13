package main

import (
	"fmt"

	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// ✅ VERY IMPORTANT: DB CONNECT FIRST
	db.ConnectDB() // 👈 YE LINE MISSING THI

	newPassword := "super123"

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		panic(err)
	}

	result := db.DB.Model(&models.Admin{}).
		Where("username = ?", "superadmin").
		Update("password", string(hash))

	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println("✅ Superadmin password reset successful")
}
