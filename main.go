package main

import (
	"go-tasker/api"
	"go-tasker/database"
	"go-tasker/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to GoTasker API!",
			"status":  "running",
		})
	})

	// Auth routes (no versioning needed for basic auth)
	router.POST("/register", user.RegisterUser)
	router.POST("/login", user.LoginUser)

	// Setup API versioning
	api.SetupAPIInfo(router)
	api.SetupV1Routes(router)
	api.SetupLegacyRoutes(router)

	router.Run()
}
