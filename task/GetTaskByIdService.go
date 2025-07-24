package task

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tasker/database"
)


func GetTasksById(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("user").(database.User).ID

	var task database.Task
	if err := database.DB.Select("id", "title", "description", "status").First(&task, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}