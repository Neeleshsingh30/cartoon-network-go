package router

import (
	"cartoon-network-go/backend/src/controllers"
	"cartoon-network-go/backend/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
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
	r.GET("/cartoons/age-groups", controllers.GetAgeGroups)
	r.GET("/cartoons/genres", controllers.GetGenres)
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
