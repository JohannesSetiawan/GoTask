package task

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-tasker/database"
)

func CreateTask(c *gin.Context) {
	var newTask database.Task
	user := c.MustGet("user").(database.User)
	newTask.UserID = user.ID
	newTask.User = user
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newTask.Title == "" || newTask.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Description are required"})
		return
	}
	if newTask.Status == "" {
		newTask.Status = "pending"
	}

	// Validate deadline if provided
	if newTask.Deadline != nil && *newTask.Deadline != "" {
		if err := validateDatetime(*newTask.Deadline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339 format (e.g., 2023-12-31T23:59:59Z)"})
			return
		}
	}

	if err := database.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, newTask.ID)
}

// validateDatetime validates if the provided string is a valid datetime in RFC3339 format
func validateDatetime(dateStr string) error {
	// Try to parse as RFC3339 (standard format: 2006-01-02T15:04:05Z07:00)
	_, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// If RFC3339 fails, try RFC3339Nano for more precision
		_, err = time.Parse(time.RFC3339Nano, dateStr)
	}
	return err
}
