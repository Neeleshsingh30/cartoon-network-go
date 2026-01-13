package router

import (
	"cartoon-network-go/backend/src/controllers"
	"cartoon-network-go/backend/src/middlewares"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ==============================
	// Health Check
	// ==============================
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Backend running",
		})
	})

	// ==============================
	// Public User Routes
	// ==============================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middlewares.RateLimiter())

	// Public Routes
	r.POST("/signup", controllers.Signup)

	// ==============================
	// Admin Routes
	// ==============================
	admin := r.Group("/admin")

	// 🔓 Public admin route
	admin.POST("/login", controllers.AdminLogin)

	// 🔐 Protected admin routes (JWT required)
	admin.Use(middlewares.AdminAuth())
	{
		// Cartoon management
		admin.POST("/cartoon", controllers.AddCartoon)
		admin.DELETE("/cartoon/:id", controllers.DeleteCartoon)

		// Character management
		admin.POST("/cartoon/:cartoon_id/character", controllers.AddCharacter)
		admin.GET("/cartoon/:cartoon_id/characters", controllers.GetCharactersByCartoon)

		// Cartoon listing
		admin.GET("/cartoons", controllers.GetAllCartoons)

		// Image upload
		admin.POST("/upload-image", controllers.UploadImage)

		// Admin activity logs
		admin.GET("/logs", controllers.GetAdminLogs)

		// ==========================
		// Admin Management
		// ==========================
		// Create admin (admin + super admin both allowed)
		admin.POST("/create-admin", controllers.CreateAdmin)

		// Delete admin (SUPER ADMIN ONLY)
		admin.DELETE(
			"/delete-admin/:id",
			middlewares.SuperAdminOnly(),
			controllers.DeleteAdmin,
		)
	}
	r.POST("/login", controllers.Login)
	r.GET("/cartoons", controllers.GetHomeCartoons)
	r.GET("/cartoon/:id", controllers.GetCartoonByID)
	r.GET("/cartoon/:id/views", controllers.GetCartoonViewCount)
	r.GET("/cartoons/timings", controllers.GetShowTimings)
	r.GET("/cartoons/search", controllers.SearchCartoons)
	r.GET("/cartoons/paginate", controllers.GetPaginatedCartoons)
	r.GET("/cartoons/trending", controllers.GetTrendingCartoons)
	r.GET("/cartoons/by-age-group", controllers.GetCartoonsByAgeGroup)

	r.GET("/cartoons/by-genre", controllers.GetCartoonsByGenre)

	r.GET("/cartoon/:id/recommendations", controllers.GetRecommendedCartoons)

	// Protected Routes (JWT)
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.GET("/user/history", controllers.GetUserHistory)
		auth.GET("/user/favourites", controllers.GetUserFavourites)
		auth.POST("/cartoon/:id/like", controllers.LikeCartoon)
		auth.DELETE("/cartoon/:id/like", controllers.UnlikeCartoon)
		auth.POST("/cartoon/:id/view", controllers.AddView)
	}

	return r
}
