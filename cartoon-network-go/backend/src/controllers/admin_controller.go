package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

/* =========================================================
   JWT CONFIG & HELPERS
========================================================= */

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	if os.Getenv("JWT_SECRET") != "" {
		return os.Getenv("JWT_SECRET")
	}
	return "cartoon_network_secret"
}

func generateAdminToken(adminID uint) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": adminID,
		"role":     "admin",
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

/* =========================================================
   ADMIN LOGIN
   POST /admin/login
========================================================= */

type AdminLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *gin.Context) {
	var input AdminLoginInput
	var admin models.Admin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateAdminToken(admin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	db.DB.Create(&models.AdminActivityLog{
		AdminID: admin.ID,
		Action:  "Admin Logged In",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

/* =========================================================
   ADD CARTOON
   POST /admin/cartoon
========================================================= */

func AddCartoon(c *gin.Context) {
	var cartoon models.Cartoon

	if err := c.ShouldBindJSON(&cartoon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cartoon data"})
		return
	}

	if err := db.DB.Create(&cartoon).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cartoon already exists"})
		return
	}

	db.DB.Create(&models.AdminActivityLog{
		Action: "Added Cartoon: " + cartoon.Name,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Cartoon added successfully",
		"cartoon": cartoon,
	})
}

/* =========================================================
   DELETE CARTOON
   DELETE /admin/cartoon/:id
========================================================= */

func DeleteCartoon(c *gin.Context) {
	id := c.Param("id")

	cartoonID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cartoon ID"})
		return
	}

	if err := db.DB.Delete(&models.Cartoon{}, cartoonID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cartoon"})
		return
	}

	db.DB.Create(&models.AdminActivityLog{
		Action: "Deleted Cartoon ID: " + id,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Cartoon deleted successfully",
	})
}

/* =========================================================
   IMAGE UPLOAD (SUPABASE PLACEHOLDER)
   POST /admin/upload-image
========================================================= */

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file required"})
		return
	}

	// ⚠️ Supabase integration yahan connect hogi
	imageURL := "https://supabase.fake.url/" + file.Filename

	db.DB.Create(&models.AdminActivityLog{
		Action: "Uploaded Image: " + file.Filename,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":   "Image uploaded successfully",
		"image_url": imageURL,
	})
}

/* =========================================================
   ADMIN ACTIVITY LOGS
   GET /admin/logs
========================================================= */

func GetAdminLogs(c *gin.Context) {
	var logs []models.AdminActivityLog

	db.DB.Preload("Admin").
		Order("created_at desc").
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}

func AddCharacter(c *gin.Context) {

	// 1️⃣ Get cartoon_id from URL
	cartoonIDParam := c.Param("cartoon_id")

	cartoonID, err := strconv.Atoi(cartoonIDParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid cartoon ID"})
		return
	}

	// 2️⃣ Check cartoon exists or not
	var cartoon models.Cartoon
	if err := db.DB.First(&cartoon, cartoonID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cartoon not found"})
		return
	}

	// 3️⃣ Bind character input
	var character models.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(400, gin.H{"error": "Invalid character data"})
		return
	}

	// 4️⃣ Set CartoonID
	character.CartoonID = uint(cartoonID)

	// 5️⃣ Save character
	if err := db.DB.Create(&character).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add character"})
		return
	}

	// 6️⃣ Log admin activity
	db.DB.Create(&models.AdminActivityLog{
		Action: "Added Character: " + character.Name + " to Cartoon: " + cartoon.Name,
	})

	// 7️⃣ Success response
	c.JSON(200, gin.H{
		"message":   "Character added successfully",
		"character": character,
	})
}
