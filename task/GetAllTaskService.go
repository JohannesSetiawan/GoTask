package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-tasker/database"
)

type PaginatedResponse struct {
	Data       []database.Task `json:"data"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
}

func GetTasks(c *gin.Context) {
	var tasks []database.Task
	var total int64
	user := c.MustGet("user").(database.User)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Ensure minimum values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Maximum limit to prevent abuse
	}

	// Parse filter parameters
	title := c.Query("title")
	description := c.Query("description")
	status := c.Query("status")

	// Parse sorting parameters
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate sort_by parameter
	validSortFields := map[string]bool{
		"created_at": true,
		"deadline":   true,
		"title":      true,
		"status":     true,
	}

	if !validSortFields[sortBy] {
		sortBy = "created_at" // Default to created_at if invalid
	}

	// Validate sort_order parameter
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc" // Default to desc if invalid
	}

	// Build query
	query := database.DB.Model(&database.Task{}).Where("user_id = ?", user.ID)

	// Apply filters
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if description != "" {
		query = query.Where("description ILIKE ?", "%"+description+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	query.Count(&total)

	// Calculate offset
	offset := (page - 1) * limit

	// Execute query with pagination and sorting
	orderClause := sortBy + " " + sortOrder

	// Handle deadline sorting with NULL values (put NULL values last)
	if sortBy == "deadline" {
		if sortOrder == "asc" {
			orderClause = "deadline ASC NULLS LAST"
		} else {
			orderClause = "deadline DESC NULLS LAST"
		}
	}

	query.Select("id", "title", "description", "status", "created_at", "deadline").
		Order(orderClause).
		Offset(offset).
		Limit(limit).
		Find(&tasks)

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := PaginatedResponse{
		Data:       tasks,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}
