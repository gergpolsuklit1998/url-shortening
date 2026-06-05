package routes

import (
	"net/http"

	"github.com/gergpolsuklit1998/url-shortening/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes and applies middleware.
// It organizes routes by version for API versioning support.
func SetupRoutes(router *gin.Engine, shortUrlHandler *handlers.ShortUrlHandler) {
	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		shortens := v1.Group("/shortens")
		{
			// Short URL endpoints
			shortens.POST("", shortUrlHandler.CreateShortUrl)
			shortens.GET("/:shortCode", shortUrlHandler.RedirectShortUrl)
		}
	}
}
