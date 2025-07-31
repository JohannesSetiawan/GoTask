package subtask

import (
	"net/http"

	"go-tasker/database"

	"github.com/gin-gonic/gin"
)

func UpdateSubtask(c *gin.Context) {
	taskID := c.Param("taskId")
	subtaskID := c.Param("subtaskId")
	user := c.MustGet("user").(database.User)

	var task database.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Verify that the subtask exists and belongs to the user and task
	var existingSubtask database.Subtask
	if err := database.DB.Where("id = ? AND task_id = ? AND user_id = ?", subtaskID, taskID, user.ID).First(&existingSubtask).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subtask not found"})
		return
	}

	var updatedSubtask database.Subtask
	if err := c.ShouldBindJSON(&updatedSubtask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedSubtask.Deadline != nil && *updatedSubtask.Deadline != "" {
		if err := validateDatetime(*updatedSubtask.Deadline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339 format (e.g., 2023-12-31T23:59:59Z)"})
			return
		}
	}

	if updatedSubtask.Status != "" && !validateTaskStatus(updatedSubtask.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid values are: Created, In Progress, Done"})
		return
	}

	updatedSubtask.TaskID = existingSubtask.TaskID
	updatedSubtask.UserID = existingSubtask.UserID

	if err := database.DB.Model(&existingSubtask).Updates(updatedSubtask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subtask"})
		return
	}

	database.DB.Where("id = ?", subtaskID).First(&existingSubtask)

	c.JSON(http.StatusOK, gin.H{
		"data":    existingSubtask,
		"message": "Subtask updated successfully",
	})
}
