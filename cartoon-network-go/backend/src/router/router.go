// package router

// import (
// 	"cartoon-network-go/backend/src/controllers"

// 	"github.com/gin-gonic/gin"
// )

// func SetupRouter() *gin.Engine {
// 	r := gin.Default()

// 	// Auth Routes
// 	r.POST("/signup", controllers.Signup)

//		return r
//	}
package router

import (
	"cartoon-network-go/backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ Health check (optional but helpful)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Backend running",
		})
	})

	// ✅ User routes
	r.POST("/signup", controllers.Signup)

	// 🔐 ✅ ADMIN ROUTES (THIS WAS MISSING)
	admin := r.Group("/admin")
	{
		admin.POST("/login", controllers.AdminLogin)
		admin.POST("/cartoon", controllers.AddCartoon)
		admin.DELETE("/cartoon/:id", controllers.DeleteCartoon)
		admin.POST("/cartoon/:cartoon_id/character", controllers.AddCharacter)
		admin.POST("/upload-image", controllers.UploadImage)
		admin.GET("/logs", controllers.GetAdminLogs)
	}

	return r
}
