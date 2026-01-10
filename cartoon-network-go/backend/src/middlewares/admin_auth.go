package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ✅ MUST MATCH admin_controller.go
var JwtSecret = []byte("cartoon_network_secret")

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1️⃣ Read Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// 2️⃣ Expect: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3️⃣ Parse & validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 4️⃣ Extract claims safely
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// 5️⃣ Store in context (request scoped)
		c.Set("admin_id", claims["admin_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists || role != "super_admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Only super admin can perform this action",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetAdminIDFromContext(c *gin.Context) uint {
	adminID, exists := c.Get("admin_id")
	if !exists {
		return 0
	}
	return adminID.(uint)
}
