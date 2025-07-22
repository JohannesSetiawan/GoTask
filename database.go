package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
	
	dsn := os.Getenv("DB_STRING")
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	fmt.Println("Database connection successful.")

	database.AutoMigrate(&Task{})

	DB = database
}