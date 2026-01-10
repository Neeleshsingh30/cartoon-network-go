// package router

// import (
// 	"cartoon-network-go/backend/src/controllers"

// 	"github.com/gin-gonic/gin"
// )

// func SetupRouter() *gin.Engine {
// 	r := gin.Default()

// 	// ✅ Health check
// 	r.GET("/", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"status": "Backend running",
// 		})
// 	})

// 	// ✅ User routes
// 	r.POST("/signup", controllers.Signup)

// 	// 🔐 ADMIN ROUTES
// 	admin := r.Group("/admin")
// 	{
// 		admin.POST("/login", controllers.AdminLogin)
// 		admin.POST("/cartoon", controllers.AddCartoon)
// 		admin.DELETE("/cartoon/:id", controllers.DeleteCartoon)
// 		admin.POST("/cartoon/:cartoon_id/character", controllers.AddCharacter)
// 		admin.POST("/upload-image", controllers.UploadImage)
// 		admin.GET("/logs", controllers.GetAdminLogs)
// 		admin.GET("/cartoons", controllers.GetAllCartoons)
// 		admin.GET("/cartoon/:cartoon_id/characters", controllers.GetCharactersByCartoon)
// 	}

// 	return r
// }

package router

import (
	"cartoon-network-go/backend/src/controllers"
	"cartoon-network-go/backend/src/middlewares"

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

	return r
}
