// internal/middleware/middleware.go
package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs request information
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Log request information
		if c.Writer.Status() >= 400 {
			// Log errors with more detail
			logMsg := fmt.Sprintf(
				"Request failed - Method: %s, Path: %s, Status: %d, Duration: %s, Client IP: %s, Error: %s",
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				duration.String(),
				c.ClientIP(),
				c.Errors.String(),
			)
			// Add a single error to the context
			_ = c.Error(fmt.Errorf(logMsg))
		} else {
			// Log successful requests
			c.Request.Context()
			gin.DefaultWriter.Write([]byte("Request completed\n"))
			gin.DefaultWriter.Write([]byte("Method: " + c.Request.Method + "\n"))
			gin.DefaultWriter.Write([]byte("Path: " + c.Request.URL.Path + "\n"))
			gin.DefaultWriter.Write([]byte("Status: " + string(rune(c.Writer.Status())) + "\n"))
			gin.DefaultWriter.Write([]byte("Duration: " + duration.String() + "\n"))
			gin.DefaultWriter.Write([]byte("Client IP: " + c.ClientIP() + "\n"))
		}
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
