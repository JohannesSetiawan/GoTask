package api

import (
	"net/http"

	"go-tasker/task"
	"go-tasker/user"

	"github.com/gin-gonic/gin"
)

// SetupV1Routes sets up all v1 API routes
func SetupV1Routes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.Use(user.AuthMiddleware())

	// Task routes
	v1.GET("/tasks", task.GetTasks)
	v1.POST("/tasks", task.CreateTask)
	v1.GET("/tasks/:id", task.GetTasksById)
	v1.PUT("/tasks/:id", task.UpdateTask)
	v1.DELETE("/tasks/:id", task.DeleteTask)

	// Version info endpoint
	v1.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":     "1.0.0",
			"api_version": "v1",
			"features": []string{
				"task_management",
				"user_authentication",
				"pagination",
				"filtering",
				"sorting",
				"deadline_validation",
			},
		})
	})
}

// SetupLegacyRoutes sets up legacy API routes for backward compatibility
func SetupLegacyRoutes(router *gin.Engine) {
	legacy := router.Group("/api")
	legacy.Use(user.AuthMiddleware())

	// Legacy task routes (same handlers as v1)
	legacy.GET("/tasks", task.GetTasks)
	legacy.POST("/tasks", task.CreateTask)
	legacy.GET("/tasks/:id", task.GetTasksById)
	legacy.PUT("/tasks/:id", task.UpdateTask)
	legacy.DELETE("/tasks/:id", task.DeleteTask)
}

// SetupAPIInfo sets up general API information endpoints
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
			},
		})
	})
}
