package controllers

import (
	"os"
	"strconv"
	"time"

	"backend/db"
	"backend/models"

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

	type AddCartoonInput struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Genre       string  `json:"genre"`
		AgeGroup    string  `json:"age_group"`
		Universe    string  `json:"universe"`
		ShowTime    string  `json:"show_time"`
		ImdbRating  float32 `json:"imdb_rating"`
		AirDate     string  `json:"air_date"` // frontend se string
	}

	var input AddCartoonInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// ✅ Parse air_date safely
	var airDate time.Time
	if input.AirDate != "" {
		parsedDate, err := time.Parse("2006-01-02", input.AirDate)
		if err == nil {
			airDate = parsedDate
		}
	}

	cartoon := models.Cartoon{
		Name:        input.Name,
		Description: input.Description,
		Genre:       input.Genre,
		AgeGroup:    input.AgeGroup,
		Universe:    input.Universe,
		ShowTime:    input.ShowTime,
		ImdbRating:  input.ImdbRating,
		AirDate:     airDate,
	}

	if err := db.DB.Create(&cartoon).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to add cartoon"})
		return
	}

	db.DB.Create(&models.AdminActivityLog{
		Action: "Added Cartoon: " + cartoon.Name,
	})

	c.JSON(201, gin.H{
		"message": "Cartoon added successfully",
		"cartoon": cartoon,
	})
}

/* =========================================================
   DELETE CARTOON (FULL CLEANUP ✅)
========================================================= */

func DeleteCartoon(c *gin.Context) {
	id := c.Param("id")

	var cartoon models.Cartoon
	if err := db.DB.First(&cartoon, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Cartoon not found"})
		return
	}

	// 🔥 Delete characters
	db.DB.Where("cartoon_id = ?", cartoon.ID).Delete(&models.Character{})

	// 🔥 Delete images (thumbnail/banner/poster)
	db.DB.Where("cartoon_id = ?", cartoon.ID).Delete(&models.CartoonImage{})

	// 🔥 Delete cartoon
	db.DB.Delete(&cartoon)

	db.DB.Create(&models.AdminActivityLog{
		Action: "Deleted Cartoon: " + cartoon.Name,
	})

	c.JSON(200, gin.H{
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
   GET CHARACTERS BY CARTOON  ✅ (MISSING FUNCTION)
========================================================= */

func GetCharactersByCartoon(c *gin.Context) {
	cartoonID := c.Param("cartoon_id")

	var characters []models.Character
	if err := db.DB.
		Where("cartoon_id = ?", cartoonID).
		Find(&characters).Error; err != nil {

		c.JSON(500, gin.H{
			"error": "Failed to fetch characters",
		})
		return
	}

	c.JSON(200, gin.H{
		"characters": characters,
	})
}

/* =========================================================
   DELETE CHARACTER ✅ (NEW)
========================================================= */

func DeleteCharacter(c *gin.Context) {
	id := c.Param("id")

	if err := db.DB.Delete(&models.Character{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete character"})
		return
	}

	db.DB.Create(&models.AdminActivityLog{
		Action: "Deleted Character ID: " + id,
	})

	c.JSON(200, gin.H{
		"message": "Character deleted successfully",
	})
}

/* =========================================================
   GET CARTOONS WITH THUMBNAIL
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

	c.JSON(200, gin.H{"cartoons": result})
}

/*
	=========================================================
	  UPLOAD CARTOON IMAGE (thumbnail/banner/poster) ✅

=========================================================
*/
type UploadImageInput struct {
	CartoonID uint   `json:"cartoon_id"`
	ImageType string `json:"image_type"` // thumbnail | banner | poster
	ImageURL  string `json:"image_url"`
}

func UploadCartoonImage(c *gin.Context) {
	var input UploadImageInput

	// ✅ JSON read karo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if input.ImageType == "" {
		c.JSON(400, gin.H{"error": "image_type required"})
		return
	}

	image := models.CartoonImage{
		CartoonID: input.CartoonID,
		ImageURL:  input.ImageURL,
		ImageType: input.ImageType,
	}

	db.DB.Create(&image)

	c.JSON(200, gin.H{
		"message": "Image uploaded successfully",
	})
}

/* =========================================================
   DELETE CARTOON IMAGE (thumbnail/banner) ✅
========================================================= */

func DeleteCartoonImage(c *gin.Context) {
	id := c.Param("id")

	if err := db.DB.Delete(&models.CartoonImage{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete image"})
		return
	}

	db.DB.Create(&models.AdminActivityLog{
		Action: "Deleted Cartoon Image ID: " + id,
	})

	c.JSON(200, gin.H{"message": "Image deleted successfully"})
}

/* =========================================================
   ADMIN LOGS
========================================================= */

func GetAdminLogs(c *gin.Context) {
	var logs []models.AdminActivityLog

	db.DB.Preload("Admin").
		Order("created_at desc").
		Find(&logs)

	c.JSON(200, gin.H{"logs": logs})
}

/* ================= CREATE / DELETE ADMIN ================= */

type CreateAdminInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func CreateAdmin(c *gin.Context) {
	var input CreateAdminInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, _ := hashPassword(input.Password)

	admin := models.Admin{
		Username: input.Username,
		Password: hashedPassword,
		Role:     input.Role,
	}

	if err := db.DB.Create(&admin).Error; err != nil {
		c.JSON(400, gin.H{"error": "Admin already exists"})
		return
	}

	c.JSON(201, gin.H{"message": "Admin created successfully"})
}

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
	c.JSON(200, gin.H{"message": "Admin deleted successfully"})
}

func GetAllAdmins(c *gin.Context) {
	var admins []models.Admin

	db.DB.Select("id, username, role").Find(&admins)
	c.JSON(200, gin.H{"admins": admins})
}
