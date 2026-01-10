package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	tokens     int
	lastAccess time.Time
}

var (
	clients = make(map[string]*client)
	mutex   sync.Mutex
	limit   = 30 // 30 requests
	window  = 60 * time.Second
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mutex.Lock()
		defer mutex.Unlock()

		cl, exists := clients[ip]

		if !exists {
			clients[ip] = &client{tokens: limit - 1, lastAccess: time.Now()}
			c.Next()
			return
		}

		if time.Since(cl.lastAccess) > window {
			cl.tokens = limit - 1
			cl.lastAccess = time.Now()
			c.Next()
			return
		}

		if cl.tokens <= 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, slow down!",
			})
			c.Abort()
			return
		}

		cl.tokens--
		c.Next()
	}
}
