package controllers

import (
	"net/http"
	"strconv"

	"cartoon-network-go/backend/src/cache"
	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
	"cartoon-network-go/backend/src/worker"

	"strings"

	"github.com/gin-gonic/gin"
)

/* ===== HOMEPAGE CARTOONS (CACHED + CONCURRENT) ===== */
func GetHomeCartoons(c *gin.Context) {
	cartoons := cache.GetHomeCache()
	c.JSON(http.StatusOK, cartoons)
}

func GetCartoonByID(c *gin.Context) {
	id := c.Param("id")

	var cartoon models.Cartoon
	db.DB.Preload("Images").Preload("Characters").First(&cartoon, id)

	c.JSON(http.StatusOK, cartoon)
}

/* ===== LIKE CARTOON (SYNC) ===== */
func LikeCartoon(c *gin.Context) {
	cartoonID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	like := models.Like{
		CartoonID: uint(cartoonID),
		UserID:    userID,
	}

	// ignore if already exists
	db.DB.Where("cartoon_id = ? AND user_id = ?", cartoonID, userID).
		FirstOrCreate(&like)

	c.JSON(http.StatusOK, gin.H{"message": "Liked"})
}

func UnlikeCartoon(c *gin.Context) {
	cartoonID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	db.DB.
		Where("cartoon_id = ? AND user_id = ?", cartoonID, userID).
		Delete(&models.Like{})

	c.JSON(http.StatusOK, gin.H{"message": "Unliked"})
}

/* ===== ADD VIEW (ASYNC) ===== */
func AddView(c *gin.Context) {
	cartoonID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	worker.JobQueue <- worker.Job{
		Type: "VIEW",
		Data: &models.CartoonView{CartoonID: uint(cartoonID), UserID: userID},
	}

	c.JSON(http.StatusOK, gin.H{"message": "View queued"})
}

/* ===== USER WATCH HISTORY ===== */
func GetUserHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var history []models.CartoonView
	db.DB.Preload("Cartoon").
		Where("user_id = ?", userID).
		Order("viewed_at DESC").
		Find(&history)

	c.JSON(http.StatusOK, history)
}

/* ===== USER FAVOURITES ===== */
func GetUserFavourites(c *gin.Context) {
	userID := c.GetUint("user_id")

	var likes []models.Like
	db.DB.
		Preload("Cartoon").
		Preload("Cartoon.Images").
		Where("user_id = ?", userID).
		Find(&likes)

	c.JSON(http.StatusOK, likes)
}

/* ===== CARTOON VIEW COUNT ===== */
func GetCartoonViewCount(c *gin.Context) {
	id := c.Param("id")

	var count int64
	db.DB.Model(&models.CartoonView{}).
		Where("cartoon_id = ?", id).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{"views": count})
}

func GetShowTimings(c *gin.Context) {
	var cartoons []models.Cartoon
	db.DB.Preload("Images").Order("show_time").Find(&cartoons)
	c.JSON(http.StatusOK, cartoons)
}

func SearchCartoons(c *gin.Context) {
	query := c.Query("q")

	var cartoons []models.Cartoon
	db.DB.Preload("Images").
		Where("name ILIKE ?", "%"+query+"%").
		Find(&cartoons)

	c.JSON(http.StatusOK, cartoons)
}
func GetPaginatedCartoons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	var cartoons []models.Cartoon
	db.DB.Preload("Images").
		Limit(limit).
		Offset(offset).
		Find(&cartoons)

	c.JSON(http.StatusOK, cartoons)
}

func GetTrendingCartoons(c *gin.Context) {

	type Result struct {
		models.Cartoon
		ViewCount int `json:"view_count"`
	}

	var results []Result

	db.DB.Raw(`
		SELECT cartoons.*, COUNT(cartoon_views.id) AS view_count
		FROM cartoons
		JOIN cartoon_views ON cartoons.id = cartoon_views.cartoon_id
		GROUP BY cartoons.id
		ORDER BY view_count DESC
		LIMIT 10
	`).Scan(&results)

	// preload images manually
	for i := range results {
		var images []models.CartoonImage
		db.DB.Where("cartoon_id = ?", results[i].ID).Find(&images)
		results[i].Images = images
	}

	c.JSON(http.StatusOK, results)
}
func GetCartoonsByAgeGroup(c *gin.Context) {

	var cartoons []models.Cartoon
	db.DB.Preload("Images").Find(&cartoons) // 👈 preload images

	result := make(map[string][]models.Cartoon)

	for _, ct := range cartoons {

		age := strings.Replace(ct.AgeGroup, "+", "", -1)
		val, _ := strconv.Atoi(age)

		var bucket string

		switch {
		case val <= 6:
			bucket = "5-6"
		case val <= 8:
			bucket = "7-8"
		case val <= 10:
			bucket = "9-10"
		default:
			bucket = "11-12"
		}

		result[bucket] = append(result[bucket], ct)
	}

	c.JSON(http.StatusOK, result)
}

func GetCartoonsByGenre(c *gin.Context) {

	var cartoons []models.Cartoon
	db.DB.Preload("Images").Find(&cartoons) // 👈 preload images

	result := make(map[string][]models.Cartoon)

	for _, ctoon := range cartoons {
		result[ctoon.Genre] = append(result[ctoon.Genre], ctoon)
	}

	c.JSON(http.StatusOK, result)
}

func GetRecommendedCartoons(c *gin.Context) {
	id := c.Param("id")

	var current models.Cartoon
	db.DB.First(&current, id)

	var cartoons []models.Cartoon
	db.DB.
		Preload("Images"). // 🔥 THIS WAS MISSING
		Where("id != ?", current.ID).
		Where("genre = ?", current.Genre).
		Order("imdb_rating DESC").
		Limit(10).
		Find(&cartoons)

	c.JSON(http.StatusOK, cartoons)
}
