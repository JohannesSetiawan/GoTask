package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-tasker/middleware"
	"go-tasker/subtask"
	"go-tasker/task"
)

func SetupAPIInfo(router *gin.Engine) {
	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":         "GoTasker API",
			"versions":        []string{"v1"},
			"current_version": "v1",
			"endpoints": gin.H{
				"v1":     "/api/v1",
				"legacy": "/api",
			},
			"documentation": gin.H{
				"auth": gin.H{
					"register": "POST /register",
					"login":    "POST /login",
				},
				"tasks": gin.H{
					"list":   "GET /api/v1/tasks",
					"create": "POST /api/v1/tasks",
					"get":    "GET /api/v1/tasks/:id",
					"update": "PUT /api/v1/tasks/:id",
					"delete": "DELETE /api/v1/tasks/:id",
				},
				"subtasks": gin.H{
					"list":   "GET /api/v1/tasks/:taskId/subtasks",
					"create": "POST /api/v1/tasks/:taskId/subtasks",
					"update": "PUT /api/v1/tasks/:taskId/subtasks/:subtaskId",
					"delete": "DELETE /api/v1/tasks/:taskId/subtasks/:subtaskId",
				},
			},
		})
	})
}

func SetupV1Routes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())

	v1.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":     "1.0.0",
			"api_version": "v1",
			"features": []string{
				"task_management",
				"subtask_management",
				"user_authentication",
				"pagination",
				"filtering",
				"sorting",
				"deadline_validation",
			},
		})
	})

	v1.GET("/tasks", task.GetTasks)
	v1.POST("/tasks", task.CreateTask)
	v1.GET("/tasks/:id", task.GetTasksById)
	v1.PUT("/tasks/:id", task.UpdateTask)
	v1.DELETE("/tasks/:id", task.DeleteTask)

	v1.GET("/task/:taskId/subtasks", subtask.GetSubtasks)
	v1.POST("/task/:taskId/subtasks", subtask.CreateSubtask)
	v1.GET("/task/:taskId/subtasks/:subtaskId", subtask.GetSubtaskById)
	v1.PUT("/task/:taskId/subtasks/:subtaskId", subtask.UpdateSubtask)
	v1.DELETE("/task/:taskId/subtasks/:subtaskId", subtask.DeleteSubtask)
}

func SetupLegacyRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.GET("/tasks", task.GetTasks)
	api.POST("/tasks", task.CreateTask)
	api.GET("/tasks/:id", task.GetTasksById)
	api.PUT("/tasks/:id", task.UpdateTask)
	api.DELETE("/tasks/:id", task.DeleteTask)
}
