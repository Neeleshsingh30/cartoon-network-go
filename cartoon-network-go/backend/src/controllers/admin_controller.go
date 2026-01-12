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

/* ================= JWT ================= */

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	if os.Getenv("JWT_SECRET") != "" {
		return os.Getenv("JWT_SECRET")
	}
	return "cartoon_network_secret"
}

func generateAdminToken(admin models.Admin) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": admin.ID,
		"role":     admin.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

/* ================= PASSWORD HELPERS ================= */

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

/* ================= ADMIN LOGIN ================= */

type AdminLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *gin.Context) {
	var input AdminLoginInput
	var admin models.Admin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := db.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if !checkPassword(input.Password, admin.Password) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := generateAdminToken(admin)

	db.DB.Create(&models.AdminActivityLog{
		AdminID: admin.ID,
		Action:  "Admin Logged In",
	})

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
		"role":    admin.Role,
	})
}

/* =========================================================
   ADD CARTOON
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
   ADD CHARACTER
========================================================= */

func AddCharacter(c *gin.Context) {
	cartoonIDParam := c.Param("cartoon_id")

	cartoonID, err := strconv.Atoi(cartoonIDParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid cartoon ID"})
		return
	}

	var cartoon models.Cartoon
	if err := db.DB.First(&cartoon, cartoonID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cartoon not found"})
		return
	}

	var character models.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(400, gin.H{"error": "Invalid character data"})
		return
	}

	character.CartoonID = uint(cartoonID)

	if err := db.DB.Create(&character).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add character"})
		return
	}

	db.DB.Create(&models.AdminActivityLog{
		Action: "Added Character: " + character.Name + " to Cartoon: " + cartoon.Name,
	})

	c.JSON(200, gin.H{
		"message":   "Character added successfully",
		"character": character,
	})
}

/* =========================================================
   GET CARTOONS & CHARACTERS (UPDATED WITH IMAGES)
========================================================= */

func GetAllCartoons(c *gin.Context) {
	type CartoonResponse struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Thumbnail string `json:"thumbnail"`
	}

	var result []CartoonResponse

	err := db.DB.
		Table("cartoons").
		Select(`
			cartoons.id,
			cartoons.name,
			COALESCE(
				(
					SELECT image_url
					FROM cartoon_images
					WHERE cartoon_images.cartoon_id = cartoons.id
					AND LOWER(image_type) = 'thumbnail'
					LIMIT 1
				),
				(
					SELECT image_url
					FROM cartoon_images
					WHERE cartoon_images.cartoon_id = cartoons.id
					AND LOWER(image_type) = 'poster'
					LIMIT 1
				)
			) AS thumbnail
		`).
		Scan(&result).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch cartoons"})
		return
	}

	c.JSON(200, gin.H{
		"cartoons": result,
	})
}

func GetCharactersByCartoon(c *gin.Context) {
	cartoonID := c.Param("cartoon_id")

	var characters []models.Character
	if err := db.DB.
		Where("cartoon_id = ?", cartoonID).
		Find(&characters).Error; err != nil {

		c.JSON(500, gin.H{"error": "Failed to fetch characters"})
		return
	}

	c.JSON(200, gin.H{
		"characters": characters,
	})
}

/* =========================================================
   UPLOAD IMAGE
========================================================= */

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file required"})
		return
	}

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

/* ================= CREATE ADMIN ================= */

type CreateAdminInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // admin / super_admin
}

func CreateAdmin(c *gin.Context) {
	var input CreateAdminInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Password hashing failed"})
		return
	}

	role := "admin"
	if input.Role != "" {
		role = input.Role
	}

	admin := models.Admin{
		Username: input.Username,
		Password: hashedPassword,
		Role:     role,
	}

	if err := db.DB.Create(&admin).Error; err != nil {
		c.JSON(400, gin.H{"error": "Admin already exists"})
		return
	}

	creatorID := c.GetUint("admin_id")

	db.DB.Create(&models.AdminActivityLog{
		AdminID: creatorID,
		Action:  "Created Admin: " + admin.Username,
	})

	c.JSON(201, gin.H{
		"message": "Admin created successfully",
		"admin": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
}

/* ================= DELETE ADMIN ================= */

func DeleteAdmin(c *gin.Context) {
	id := c.Param("id")

	var admin models.Admin
	if err := db.DB.First(&admin, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Admin not found"})
		return
	}

	if admin.Role == "super_admin" {
		c.JSON(403, gin.H{"error": "Cannot delete super admin"})
		return
	}

	db.DB.Delete(&admin)

	db.DB.Create(&models.AdminActivityLog{
		Action: "Deleted Admin: " + admin.Username,
	})

	c.JSON(200, gin.H{"message": "Admin deleted successfully"})
}
