package router

import (
	"time"

	"cartoon-network-go/backend/src/controllers"
	"cartoon-network-go/backend/src/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// =========================================
	// 🔥 STATIC FILE SERVING (UPLOADS)
	// =========================================
	// Folder must exist: backend/uploads
	r.Static("/uploads", "./uploads")

	// =========================================
	// CORS CONFIG
	// =========================================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate Limiter
	r.Use(middlewares.RateLimiter())

	// =========================================
	// HEALTH CHECK
	// =========================================
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Backend running",
		})
	})

	// =========================================
	// PUBLIC USER ROUTES
	// =========================================
	r.POST("/signup", controllers.Signup)
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

	// =========================================
	// ADMIN ROUTES
	// =========================================
	admin := r.Group("/admin")

	// 🔓 Admin Login (Public)
	admin.POST("/login", controllers.AdminLogin)

	// 🔐 Protected Admin Routes
	admin.Use(middlewares.AdminAuth())
	{
		// ================= CARTOONS =================
		admin.POST("/cartoon", controllers.AddCartoon)
		admin.GET("/cartoons", controllers.GetAllCartoons)
		admin.DELETE("/cartoon/:id", controllers.DeleteCartoon)

		// ================= CHARACTERS =================
		admin.POST("/cartoon/:cartoon_id/character", controllers.AddCharacter)
		admin.GET("/cartoon/:cartoon_id/characters", controllers.GetCharactersByCartoon)

		// ================= IMAGE UPLOAD =================
		admin.POST("/upload-image", controllers.UploadImage)

		// ================= LOGS =================
		admin.GET("/logs", controllers.GetAdminLogs)

		// ================= ADMIN MANAGEMENT =================

		// 🔐 SUPER ADMIN ONLY → LIST ADMINS
		admin.GET(
			"/list",
			middlewares.SuperAdminOnly(),
			controllers.GetAllAdmins,
		)

		// 🔐 CREATE ADMIN (role check handled in controller)
		admin.POST("/create-admin", controllers.CreateAdmin)

		// 🔐 SUPER ADMIN ONLY → DELETE ADMIN
		admin.DELETE(
			"/delete-admin/:id",
			middlewares.SuperAdminOnly(),
			controllers.DeleteAdmin,
		)
	}

	// =========================================
	// AUTHENTICATED USER ROUTES
	// =========================================
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.GET("/user/history", controllers.GetUserHistory)
		auth.GET("/user/favourites", controllers.GetUserFavourites)
		auth.POST("/cartoon/:id/like", controllers.LikeCartoon)
		auth.POST("/cartoon/:id/view", controllers.AddView)
	}

	return r
}
