package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DB_STRING")

	// Configure GORM with optimized settings for concurrency
	config := &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Error), // Only log errors in production
		PrepareStmt:                              true,                                 // Cache prepared statements
		DisableForeignKeyConstraintWhenMigrating: false,
	}

	database, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}

	// Configure connection pool for high concurrency
	sqlDB.SetMaxIdleConns(20)                  // Keep 20 idle connections
	sqlDB.SetMaxOpenConns(100)                 // Allow up to 100 concurrent connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Recycle connections every hour
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Close idle connections after 30 minutes

	fmt.Println("Database connection successful with optimized pool settings.")

	database.AutoMigrate(&Task{}, &User{})

	sqlBytes, err := os.ReadFile("./database/post_migrations.sql")
    if err != nil {
        log.Fatalf("Could not read SQL file: %v", err)
    }

    // 3. Execute raw SQL
    statements := strings.Split(string(sqlBytes), ";")
    for _, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt == "" {
            continue
        }

        fmt.Println("Executing:", stmt)
        if err := database.Exec(stmt).Error; err != nil {
            log.Fatalf("Error executing SQL: %v\nStatement: %s", err, stmt)
        }
    }

	DB = database
}
