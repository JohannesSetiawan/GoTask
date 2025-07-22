package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"go-tasker/database"
	"go-tasker/user"
	"go-tasker/task"
)


func main() {
	database.ConnectDatabase()
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	router.POST("/register", user.RegisterUser)
	router.POST("/login", user.LoginUser)

	taskApi := router.Group("/api")
	taskApi.Use(user.AuthMiddleware())

	taskApi.GET("/tasks", task.GetTasks)
	taskApi.POST("/tasks", task.CreateTask)
	taskApi.GET("/tasks/:id", task.GetTasksById)
	taskApi.PUT("/tasks/:id", task.UpdateTask)
	taskApi.DELETE("/tasks/:id", task.DeleteTask)

	router.Run()
}