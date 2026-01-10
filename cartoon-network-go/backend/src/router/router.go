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
	r.POST("/login", controllers.Login)
	r.GET("/cartoons", controllers.GetAllCartoons)
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
		auth.POST("/cartoon/:id/view", controllers.AddView)
	}

	return r
}
