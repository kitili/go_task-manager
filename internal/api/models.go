package api

import (
	"time"

	"learn-go-capstone/internal/task"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:""`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success" example:"true"`
	Message    string      `json:"message" example:"Data retrieved successfully"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Error      string      `json:"error,omitempty" example:""`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	TotalItems int   `json:"total_items" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
	HasNext    bool  `json:"has_next" example:"true"`
	HasPrev    bool  `json:"has_prev" example:"false"`
}

// TaskRequest represents a task creation/update request
type TaskRequest struct {
	Title       string     `json:"title" binding:"required" example:"Learn Go Programming"`
	Description string     `json:"description" example:"Study Go language fundamentals and best practices"`
	Priority    int        `json:"priority" binding:"required,min=1,max=5" example:"3"`
	Status      int        `json:"status" binding:"min=0,max=3" example:"0"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2024-12-31T23:59:59Z"`
	CategoryID  *int       `json:"category_id,omitempty" example:"1"`
	TagNames    []string   `json:"tag_names,omitempty" example:"[\"learning\", \"programming\"]"`
}

// TaskResponse represents a task response
type TaskResponse struct {
	ID          int                `json:"id" example:"1"`
	Title       string             `json:"title" example:"Learn Go Programming"`
	Description string             `json:"description" example:"Study Go language fundamentals and best practices"`
	Priority    int                `json:"priority" example:"3"`
	Status      int                `json:"status" example:"0"`
	CreatedAt   time.Time          `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time          `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DueDate     *time.Time         `json:"due_date,omitempty" example:"2024-12-31T23:59:59Z"`
	Category    *CategoryResponse  `json:"category,omitempty"`
	Tags        []TagResponse      `json:"tags,omitempty"`
	UserID      *int               `json:"user_id,omitempty" example:"1"`
	IsArchived  bool               `json:"is_archived" example:"false"`
}

// CategoryRequest represents a category creation/update request
type CategoryRequest struct {
	Name        string `json:"name" binding:"required" example:"Work"`
	Description string `json:"description" example:"Work-related tasks"`
	Color       string `json:"color" example:"#ff0000"`
}

// CategoryResponse represents a category response
type CategoryResponse struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"Work"`
	Description string    `json:"description" example:"Work-related tasks"`
	Color       string    `json:"color" example:"#ff0000"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// TagRequest represents a tag creation/update request
type TagRequest struct {
	Name  string `json:"name" binding:"required" example:"urgent"`
	Color string `json:"color" example:"#ff0000"`
}

// TagResponse represents a tag response
type TagResponse struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"urgent"`
	Color     string    `json:"color" example:"#ff0000"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// UserRequest represents a user registration request
type UserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID        int       `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query       string    `json:"query" example:"learn go"`
	Status      *int      `json:"status,omitempty" example:"0"`
	Priority    *int      `json:"priority,omitempty" example:"3"`
	CategoryID  *int      `json:"category_id,omitempty" example:"1"`
	TagNames    []string  `json:"tag_names,omitempty" example:"[\"learning\", \"programming\"]"`
	CreatedAfter *time.Time `json:"created_after,omitempty" example:"2024-01-01T00:00:00Z"`
	DueBefore   *time.Time `json:"due_before,omitempty" example:"2024-12-31T23:59:59Z"`
	SortBy      string    `json:"sort_by" example:"created_at"`
	SortOrder   string    `json:"sort_order" example:"desc"`
	Page        int       `json:"page" example:"1"`
	PageSize    int       `json:"page_size" example:"10"`
}

// ExportRequest represents an export request
type ExportRequest struct {
	Format      string    `json:"format" binding:"required,oneof=json csv" example:"json"`
	Status      *int      `json:"status,omitempty" example:"0"`
	Priority    *int      `json:"priority,omitempty" example:"3"`
	CategoryID  *int      `json:"category_id,omitempty" example:"1"`
	TagNames    []string  `json:"tag_names,omitempty" example:"[\"learning\", \"programming\"]"`
	CreatedAfter *time.Time `json:"created_after,omitempty" example:"2024-01-01T00:00:00Z"`
	DueBefore   *time.Time `json:"due_before,omitempty" example:"2024-12-31T23:59:59Z"`
	IncludeUsers bool     `json:"include_users" example:"true"`
	IncludeCategories bool `json:"include_categories" example:"true"`
	IncludeTags bool     `json:"include_tags" example:"true"`
}

// ImportRequest represents an import request
type ImportRequest struct {
	Format          string `json:"format" binding:"required,oneof=json csv" example:"json"`
	OverwriteExisting bool `json:"overwrite_existing" example:"false"`
	SkipDuplicates   bool `json:"skip_duplicates" example:"true"`
	DryRun          bool `json:"dry_run" example:"false"`
}

// NotificationRequest represents a notification creation request
type NotificationRequest struct {
	Type        string    `json:"type" binding:"required,oneof=email in_app sms webhook slack discord" example:"email"`
	Priority    string    `json:"priority" binding:"required,oneof=low normal high critical" example:"normal"`
	Title       string    `json:"title" binding:"required" example:"Task Reminder"`
	Message     string    `json:"message" binding:"required" example:"Your task is due soon"`
	Recipient   string    `json:"recipient" example:"user@example.com"`
	Channel     string    `json:"channel" example:"#general"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID           int       `json:"id" example:"1"`
	UserID       int       `json:"user_id" example:"1"`
	TaskID       int       `json:"task_id" example:"1"`
	Type         string    `json:"type" example:"email"`
	Priority     string    `json:"priority" example:"normal"`
	Status       string    `json:"status" example:"pending"`
	Trigger      string    `json:"trigger" example:"due_date"`
	Title        string    `json:"title" example:"Task Reminder"`
	Message      string    `json:"message" example:"Your task is due soon"`
	Recipient    string    `json:"recipient" example:"user@example.com"`
	Channel      string    `json:"channel" example:"#general"`
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty" example:"2024-12-31T23:59:59Z"`
	SentAt       *time.Time `json:"sent_at,omitempty" example:"2024-12-31T23:59:59Z"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty" example:"2024-12-31T23:59:59Z"`
	CreatedAt    time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	RetryCount   int       `json:"retry_count" example:"0"`
	MaxRetries   int       `json:"max_retries" example:"3"`
	Error        string    `json:"error,omitempty" example:""`
}

// DependencyRequest represents a task dependency creation request
type DependencyRequest struct {
	TaskID          int `json:"task_id" binding:"required" example:"2"`
	DependsOnTaskID int `json:"depends_on_task_id" binding:"required" example:"1"`
}

// DependencyResponse represents a task dependency response
type DependencyResponse struct {
	ID               int       `json:"id" example:"1"`
	TaskID           int       `json:"task_id" example:"2"`
	DependsOnTaskID  int       `json:"depends_on_task_id" example:"1"`
	CreatedAt        time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// StatisticsResponse represents application statistics
type StatisticsResponse struct {
	TotalTasks      int `json:"total_tasks" example:"100"`
	CompletedTasks  int `json:"completed_tasks" example:"75"`
	PendingTasks    int `json:"pending_tasks" example:"20"`
	OverdueTasks    int `json:"overdue_tasks" example:"5"`
	TotalUsers      int `json:"total_users" example:"10"`
	TotalCategories int `json:"total_categories" example:"5"`
	TotalTags       int `json:"total_tags" example:"15"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status" example:"healthy"`
	Timestamp time.Time         `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Version   string            `json:"version" example:"1.0.0"`
	Uptime    string            `json:"uptime" example:"2h30m15s"`
	Services  map[string]string `json:"services" example:"{\"database\":\"healthy\",\"notifications\":\"healthy\"}"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"An error occurred"`
	Error   string `json:"error" example:"detailed error message"`
	Code    int    `json:"code" example:"400"`
}

// ConvertToTaskResponse converts a task.Task to TaskResponse
func ConvertToTaskResponse(t task.Task) TaskResponse {
	response := TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    int(t.Priority),
		Status:      int(t.Status),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DueDate:     t.DueDate,
		IsArchived:  false, // Default value
	}

	if t.Category != nil {
		response.Category = &CategoryResponse{
			ID:          t.Category.ID,
			Name:        t.Category.Name,
			Description: t.Category.Description,
			Color:       t.Category.Color,
			CreatedAt:   t.Category.CreatedAt,
			UpdatedAt:   t.Category.UpdatedAt,
		}
	}

	if len(t.Tags) > 0 {
		response.Tags = make([]TagResponse, len(t.Tags))
		for i, tag := range t.Tags {
			response.Tags[i] = TagResponse{
				ID:        tag.ID,
				Name:      tag.Name,
				Color:     tag.Color,
				CreatedAt: tag.CreatedAt,
			}
		}
	}

	return response
}

// ConvertToTaskRequest converts a TaskRequest to task.Task
func ConvertToTaskRequest(req TaskRequest) task.Task {
	t := task.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    task.Priority(req.Priority),
		Status:      task.Status(req.Status),
		DueDate:     req.DueDate,
	}

	if req.CategoryID != nil {
		t.Category = &task.Category{
			ID: *req.CategoryID,
		}
	}

	if len(req.TagNames) > 0 {
		t.Tags = make([]task.Tag, len(req.TagNames))
		for i, tagName := range req.TagNames {
			t.Tags[i] = task.Tag{
				Name: tagName,
			}
		}
	}

	return t
}
