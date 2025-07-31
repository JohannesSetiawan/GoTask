package task

import (
	"net/http"

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
		newTask.Status = database.StatusCreated
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