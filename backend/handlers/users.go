package handlers

import (
	"net/http"
	"securetask/database"
	"securetask/models"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

func sanitizeHTML(input string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(input)
}

func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id") // From AuthMiddleware
	paramID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userID != paramID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update other users' profiles"})
		return
	}

	if bio, ok := updates["bio"].(string); ok {
		updates["bio"] = sanitizeHTML(bio)
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Model(&user).Updates(updates)

	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	// Check if user has admin role
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Check admin API key header
	adminKey := c.GetHeader("X-Admin-Key")
	if adminKey != getAdminAPIKey() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin API key"})
		return
	}

	var users []models.User
	database.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}
