package task

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tasker/database"
)

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("user").(database.User).ID
	var updatedTask database.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate deadline if provided
	if updatedTask.Deadline != nil && *updatedTask.Deadline != "" {
		if err := validateDatetime(*updatedTask.Deadline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339 format (e.g., 2023-12-31T23:59:59Z)"})
			return
		}
	}

	if err := database.DB.Model(&database.Task{}).Where("id = ? AND user_id = ?", id, userId).Updates(updatedTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}
