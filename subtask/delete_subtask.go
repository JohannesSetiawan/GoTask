package subtask

import (
	"net/http"

	"go-tasker/database"

	"github.com/gin-gonic/gin"
)

// DeleteSubtask deletes a subtask
func DeleteSubtask(c *gin.Context) {
	taskID := c.Param("taskId")
	subtaskID := c.Param("subtaskId")
	user := c.MustGet("user").(database.User)

	var task database.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var subtask database.Subtask
	if err := database.DB.Where("id = ? AND task_id = ? AND user_id = ?", subtaskID, taskID, user.ID).First(&subtask).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subtask not found"})
		return
	}

	if err := database.DB.Delete(&subtask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subtask"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Subtask deleted successfully",
		"deleted_id": subtaskID,
	})
}
