package subtask

import (
	"net/http"

	"go-tasker/database"

	"github.com/gin-gonic/gin"
)

func GetSubtasks(c *gin.Context) {
	taskID := c.Param("taskId")
	user := c.MustGet("user").(database.User)

	var task database.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var subtasks []database.Subtask
	if err := database.DB.Where("task_id = ? AND user_id = ?", taskID, user.ID).
		Order("\"order\" ASC, created_at ASC").
		Find(&subtasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subtasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    subtasks,
		"count":   len(subtasks),
		"task_id": taskID,
		"message": "Subtasks retrieved successfully",
	})
}