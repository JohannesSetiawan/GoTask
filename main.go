package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


func main() {
	ConnectDatabase()
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	router.POST("/register", registerUser)

	router.GET("/tasks", getTasks)
	router.POST("/tasks", createTask)
	router.GET("/tasks/:id", getTasksById)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)

	router.Run()
}