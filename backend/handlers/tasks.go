package handlers

import (
	"net/http"
	"securetask/database"
	"securetask/models"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
}

var validPriorities = map[string]bool{"low": true, "medium": true, "high": true}
var validStatuses = map[string]bool{"todo": true, "in_progress": true, "done": true}

func GetTasks(c *gin.Context) {
	userID := c.GetUint("user_id")

	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Preload("User").Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default and validate priority
	if req.Priority == "" {
		req.Priority = "medium"
	}
	if !validPriorities[req.Priority] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority: must be low, medium, or high"})
		return
	}

	userID := c.GetUint("user_id")

	// Sanitize user-supplied text to prevent stored XSS
	task := models.Task{
		Title:       sanitizeHTML(req.Title),
		Description: sanitizeHTML(req.Description),
		Priority:    req.Priority,
		Status:      "todo",
		UserID:      userID,
	}

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Bind to a typed struct (prevents mass assignment of arbitrary columns)
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build a whitelist of updatable fields, validating and sanitizing each
	updates := map[string]interface{}{}
	if req.Title != nil {
		title := sanitizeHTML(*req.Title)
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title cannot be empty"})
			return
		}
		updates["title"] = title
	}
	if req.Description != nil {
		updates["description"] = sanitizeHTML(*req.Description)
	}
	if req.Priority != nil {
		if !validPriorities[*req.Priority] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority: must be low, medium, or high"})
			return
		}
		updates["priority"] = *req.Priority
	}
	if req.Status != nil {
		if !validStatuses[*req.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status: must be todo, in_progress, or done"})
			return
		}
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	database.DB.Model(&task).Updates(updates)

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetUint("user_id")

	// Only allow deleting a task the authenticated user owns
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func SearchTasks(c *gin.Context) {
	searchTerm := c.Query("q")

	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term required"})
		return
	}

	userID := c.GetUint("user_id")

	// Parameterized query, scoped to the authenticated user's own tasks
	var tasks []models.Task
	like := "%" + searchTerm + "%"
	if err := database.DB.
		Where("user_id = ? AND (title LIKE ? OR description LIKE ?)", userID, like, like).
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
