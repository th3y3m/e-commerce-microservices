package logging

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Log request details
		log.Printf("Incoming request: %s %s from %s",
			c.Request.Method, c.Request.RequestURI, c.ClientIP())

		// Continue to the next handler
		c.Next()

		// Log response details after handling
		latency := time.Since(startTime)
		log.Printf("Response: %d, latency: %v, path: %s",
			c.Writer.Status(), latency, c.Request.RequestURI)
	}
}
