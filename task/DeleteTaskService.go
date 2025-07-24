package task

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tasker/database"
)

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("user").(database.User).ID
	if err := database.DB.Delete(&database.Task{}, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
