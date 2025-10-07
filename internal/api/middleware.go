package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"learn-go-capstone/internal/auth"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Message: "Authorization header required",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Message: "Invalid authorization header format",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Message: "Token is required",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   err.Error(),
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Set user information in the context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// Continue to the next handler
		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware handles request logging
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// RateLimitMiddleware handles rate limiting (simple implementation)
func RateLimitMiddleware() gin.HandlerFunc {
	// This is a simple in-memory rate limiter
	// In production, you'd use Redis or a more sophisticated solution
	rateLimiter := make(map[string]int)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		// Simple rate limiting: max 100 requests per minute
		if rateLimiter[clientIP] > 100 {
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Success: false,
				Message: "Rate limit exceeded",
				Code:    http.StatusTooManyRequests,
			})
			c.Abort()
			return
		}
		
		rateLimiter[clientIP]++
		
		// Reset counter every minute (simplified)
		go func() {
			time.Sleep(time.Minute)
			rateLimiter[clientIP] = 0
		}()
		
		c.Next()
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Success: false,
					Message: "Internal server error",
					Error:   "An unexpected error occurred",
					Code:    http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()
		
		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	// Simple implementation - in production, use a proper UUID library
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
