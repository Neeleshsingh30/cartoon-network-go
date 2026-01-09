package controllers

import (
	"net/http"
	"strconv"

	"cartoon-network-go/backend/src/cache"
	"cartoon-network-go/backend/src/db"
	"cartoon-network-go/backend/src/models"
	"cartoon-network-go/backend/src/worker"

	"cartoon-network-go/backend/src/services"

	"github.com/gin-gonic/gin"
)

/* ===== HOMEPAGE CARTOONS (CACHED + CONCURRENT) ===== */
func GetAllCartoons(c *gin.Context) {
	cartoons := cache.GetHomeCache()
	c.JSON(http.StatusOK, cartoons)
}

func GetCartoonByID(c *gin.Context) {
	id := c.Param("id")

	var cartoon models.Cartoon
	db.DB.Preload("Images").Preload("Characters").First(&cartoon, id)

	// 🔥 only fetch if rating not stored
	if cartoon.ImdbRating == 0 {

		rating, err := services.FetchIMDBRating(cartoon.Name)
		if err == nil && rating > 0 {
			cartoon.ImdbRating = rating
			db.DB.Save(&cartoon)
		}
	}

	c.JSON(http.StatusOK, cartoon)
}

/* ===== LIKE CARTOON (ASYNC) ===== */
func LikeCartoon(c *gin.Context) {
	cartoonID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	worker.JobQueue <- worker.Job{
		Type: "LIKE",
		Data: &models.Like{CartoonID: uint(cartoonID), UserID: userID},
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like queued"})
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
	db.DB.Preload("Cartoon").Where("user_id = ?", userID).Find(&likes)

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
	cartoons := cache.GetShowTimeCache()
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
	var cartoons []models.Cartoon

	db.DB.Raw(`
		SELECT cartoons.*
		FROM cartoons
		JOIN cartoon_views ON cartoons.id = cartoon_views.cartoon_id
		GROUP BY cartoons.id
		ORDER BY COUNT(cartoon_views.id) DESC
		LIMIT 5
	`).Scan(&cartoons)

	c.JSON(http.StatusOK, cartoons)
}
func GetAgeGroups(c *gin.Context) {
	var ages []string

	db.DB.Model(&models.Cartoon{}).
		Distinct().
		Pluck("age_group", &ages)

	c.JSON(http.StatusOK, ages)
}

func GetGenres(c *gin.Context) {
	var genres []string

	db.DB.Model(&models.Cartoon{}).
		Distinct().
		Pluck("genre", &genres)

	c.JSON(http.StatusOK, genres)
}
