package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Task struct{
	ID          string    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var tasks = []Task{
	{ID: "1", Title: "Belajar Golang & Gin", Description: "", Status: "todo"},
	{ID: "2", Title: "Membangun API GoTasker", Description: "", Status: "todo"},
	{ID: "3", Title: "Deploy ke Docker", Description: "", Status: "todo"},
}

func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newID := len(tasks) + 1
	newTask.ID = strconv.Itoa(newID)
	tasks = append(tasks, newTask)
	c.JSON(http.StatusCreated, newTask)
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	router.GET("/tasks", getTasks)
	router.POST("/tasks", createTask)

	router.Run()
}