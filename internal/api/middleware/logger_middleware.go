package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()

		log.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		timeDuration := time.Since(start)

		log.Printf("Completed %s %s with status %d in %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), timeDuration)
	}
}
