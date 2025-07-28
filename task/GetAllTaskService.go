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
		sortBy = "created_at"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Build base query with optimized selection
	query := database.DB.Model(&database.Task{}).Where("user_id = ?", user.ID)

	// Apply filters with optimized queries
	if title != "" {
		// Use full-text search or trigram similarity for better performance
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if description != "" {
		query = query.Where("description ILIKE ?", "%"+description+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Execute optimized query with window function for total count
	// This reduces the query from 2 to 1

	orderClause := sortBy + " " + sortOrder
	if sortBy == "deadline" {
		if sortOrder == "asc" {
			orderClause = "deadline ASC NULLS LAST"
		} else {
			orderClause = "deadline DESC NULLS LAST"
		}
	}

	// Single query with COUNT OVER() window function
	var rawSQL string
	if title != "" || description != "" || status != "" {
		// For filtered queries, we still need separate count for accuracy
		var total int64
		query.Count(&total)

		query.Select("id", "title", "description", "status", "created_at", "deadline").
			Order(orderClause).
			Offset(offset).
			Limit(limit).
			Find(&tasks)

		totalPages := int((total + int64(limit) - 1) / int64(limit))

		response := PaginatedResponse{
			Data:       tasks,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// For non-filtered queries, use window function
	rawSQL = `
		SELECT id, title, description, status, created_at, deadline,
			   COUNT(*) OVER() as total_count
		FROM tasks 
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY ` + orderClause + `
		LIMIT ? OFFSET ?`

	rows, err := database.DB.Raw(rawSQL, user.ID, limit, offset).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer rows.Close()

	var total int64
	tasks = []database.Task{}

	for rows.Next() {
		var task database.Task
		var totalCount int64

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status,
			&task.CreatedAt, &task.Deadline, &totalCount)
		if err != nil {
			continue
		}

		tasks = append(tasks, task)
		total = totalCount
	}

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
