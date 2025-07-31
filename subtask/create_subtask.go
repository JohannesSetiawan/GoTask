package subtask

import (
	"net/http"
	"strconv"

	"go-tasker/database"

	"github.com/gin-gonic/gin"
)

func CreateSubtask(c *gin.Context) {
	taskID := c.Param("taskId")
	user := c.MustGet("user").(database.User)

	var task database.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var newSubtask database.Subtask
	if err := c.ShouldBindJSON(&newSubtask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newSubtask.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if newSubtask.Deadline != nil && *newSubtask.Deadline != "" {
		if err := validateDatetime(*newSubtask.Deadline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339 format (e.g., 2023-12-31T23:59:59Z)"})
			return
		}
	}

	if newSubtask.Status != "" && !validateTaskStatus(newSubtask.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid values are: Created, In Progress, Done"})
		return
	}

	if newSubtask.Status == "" {
		newSubtask.Status = database.StatusCreated
	}

	taskIDUint, _ := strconv.ParseUint(taskID, 10, 32)
	newSubtask.TaskID = uint(taskIDUint)
	newSubtask.UserID = user.ID

	if newSubtask.Order <= 0 {
		var maxOrder int
		database.DB.Model(&database.Subtask{}).
			Where("task_id = ? AND user_id = ?", taskID, user.ID).
			Select("COALESCE(MAX(\"order\"), 0)").
			Scan(&maxOrder)
		newSubtask.Order = maxOrder + 1
	}

	if err := database.DB.Create(&newSubtask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subtask"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    newSubtask,
		"message": "Subtask created successfully",
	})
}
