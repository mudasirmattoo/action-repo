package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService *UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers the user-related HTTP routes
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		// Public routes
		auth.GET("/verify", h.VerifySession)

		// Protected routes (require authentication)
		protected := auth.Group("/")
		protected.Use(SupabaseAuthMiddleware())
		{
			protected.GET("/me", h.GetCurrentUser)
			protected.PUT("/me", h.UpdateCurrentUser)
		}
	}
}

// VerifySession verifies if the current session is valid
func (h *UserHandler) VerifySession(c *gin.Context) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}

	// Check if the Authorization header starts with "Bearer "
	if _, err := validateToken(authHeader[7:]); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}

// GetCurrentUser gets the current user's profile
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	profile, err := h.userService.GetUser(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateCurrentUser updates the current user's profile
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var profile UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the profile
	if err := h.userService.UpdateUserProfile(userID.(string), &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
