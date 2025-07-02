package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thefaheem/JEE-Leetcode/server/internal/auth"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // For development; restrict in production
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin", "Content-Type", "Accept", "Authorization",
		"X-Requested-With", "X-CSRF-Token",
	}
	router.Use(cors.New(config))

	// Set up routes
	setupRoutes(router)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("Server running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func setupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Initialize user service
	userService, err := auth.NewUserService()
	if err != nil {
		log.Fatalf("Failed to initialize user service: %v", err)
	}

	// Initialize and register user handler routes
	userHandler := auth.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		api.GET("/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "This is a public endpoint",
			})
		})

		// Protected routes
		protected := api.Group("/")
		protected.Use(auth.SupabaseAuthMiddleware())
		{
			protected.GET("/protected", func(c *gin.Context) {
				// Get the user ID from the verified token
				userId, exists := c.Get("userId")
				if !exists {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"message": "This is a protected endpoint",
					"userId":  userId,
				})
			})

			// Add more protected routes here
		}
	}
}
