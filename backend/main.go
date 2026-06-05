// main.go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/gergpolsuklit1998/url-shortening/config"
    "github.com/gergpolsuklit1998/url-shortening/handlers"
    "github.com/gergpolsuklit1998/url-shortening/repository"
    "github.com/gergpolsuklit1998/url-shortening/routes"

    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load configuration from environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
    mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
    dbName := getEnv("DB_NAME", "gin_mongodb_api")
    port := getEnv("PORT", "8080")

    // Connect to MongoDB
    db, err := config.ConnectDB(mongoURI, dbName)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Initialize repository and handlers
    shortUrlCollection := db.GetCollection("short_urls")
    shortUrlRepo := repository.NewShortUrlRepository(shortUrlCollection)
    shortUrlHandler := handlers.NewShortUrlHandler(shortUrlRepo)

    // Create Gin router with default middleware (logger, recovery)
    router := gin.Default()

    // Configure routes
    routes.SetupRoutes(router, shortUrlHandler)

    // Start server in a goroutine
    go func() {
        log.Printf("Server starting on port %s", port)
        if err := router.Run(":" + port); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal for graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Disconnect from MongoDB
    if err := db.Disconnect(); err != nil {
        log.Printf("Error disconnecting from MongoDB: %v", err)
    }

    log.Println("Server stopped")
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
		log.Printf("Using environment variable %s = %s", key, value)
        return value
    }
    log.Printf("Using default value for %s = %s", key, defaultValue)
    return defaultValue
}