package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-tasker/api"
	"go-tasker/database"
	"go-tasker/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database with optimized settings
	database.ConnectDatabase()

	// Configure Gin for production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware for better performance
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add custom middleware for performance monitoring
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Log slow requests (>1 second)
		if duration > time.Second {
			log.Printf("SLOW REQUEST: %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
		}
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to GoTasker API!",
			"status":  "running",
		})
	})

	// Auth routes (no versioning needed for basic auth)
	router.POST("/register", user.RegisterUser)
	router.POST("/login", user.LoginUser)

	// Setup API versioning
	api.SetupAPIInfo(router)
	api.SetupV1Routes(router)
	api.SetupLegacyRoutes(router)

	// Configure HTTP server for high concurrency
	srv := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start server in goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
