package task

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tasker/database"
)


func GetTasks(c *gin.Context) {
	var tasks []database.Task
	user := c.MustGet("user").(database.User)
	database.DB.Select("id", "title", "description", "status").Where("user_id = ?", user.ID).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

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
	if err := database.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, newTask.ID)
}

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

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("user").(database.User).ID
	var updatedTask database.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Model(&database.Task{}).Where("id = ? AND user_id = ?", id, userId).Updates(updatedTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	userId := c.MustGet("user").(database.User).ID
	if err := database.DB.Delete(&database.Task{}, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}