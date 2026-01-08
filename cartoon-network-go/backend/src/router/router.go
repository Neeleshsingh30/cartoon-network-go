package router

import (
	"cartoon-network-go/backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Auth Routes
	r.POST("/signup", controllers.Signup)

	return r
}
