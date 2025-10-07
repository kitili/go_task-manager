package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"learn-go-capstone/internal/auth"
	"learn-go-capstone/internal/task"
)

// Handler contains all the handlers for the API
type Handler struct {
	taskManager        task.TaskManagerInterface
	userManager        *task.UserManager
	categoryManager    *task.CategoryManager
	dependencyManager  *task.DependencyManager
	searchManager      *task.SearchManager
	exportManager      *task.ExportManager
	notificationManager *task.NotificationManager
	authService        *auth.AuthService
}

// NewHandler creates a new API handler
func NewHandler(
	taskManager task.TaskManagerInterface,
	userManager *task.UserManager,
	categoryManager *task.CategoryManager,
	dependencyManager *task.DependencyManager,
	searchManager *task.SearchManager,
	exportManager *task.ExportManager,
	notificationManager *task.NotificationManager,
	authService *auth.AuthService,
) *Handler {
	return &Handler{
		taskManager:        taskManager,
		userManager:        userManager,
		categoryManager:    categoryManager,
		dependencyManager:  dependencyManager,
		searchManager:      searchManager,
		exportManager:      exportManager,
		notificationManager: notificationManager,
		authService:        authService,
	}
}

// HealthCheck handles health check requests
// @Summary Health check
// @Description Check the health status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    "2h30m15s", // This would be calculated from app start time
		Services: map[string]string{
			"database":      "healthy",
			"notifications": "healthy",
		},
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Service is healthy",
		Data:    response,
	})
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body UserRequest true "User registration data"
// @Success 201 {object} APIResponse{data=UserResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userManager.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, ErrorResponse{
			Success: false,
			Message: "Failed to register user",
			Error:   err.Error(),
			Code:    status,
		})
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    userResponse,
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} APIResponse{data=LoginResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	token, user, err := h.userManager.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "Invalid credentials",
			Error:   err.Error(),
			Code:    http.StatusUnauthorized,
		})
		return
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	loginResponse := LoginResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    loginResponse,
	})
}

// CreateTask handles task creation
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body TaskRequest true "Task data"
// @Success 201 {object} APIResponse{data=TaskResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /tasks [post]
func (h *Handler) CreateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	task := ConvertToTaskRequest(req)
	createdTask, err := h.userManager.CreateUserTask(
		userID.(int),
		task.Title,
		task.Description,
		task.Priority,
		task.DueDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to create task",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	taskResponse := ConvertToTaskResponse(*createdTask)
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Task created successfully",
		Data:    taskResponse,
	})
}

// GetTasks handles getting user tasks
// @Summary Get user tasks
// @Description Get all tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query int false "Filter by status"
// @Param priority query int false "Filter by priority"
// @Success 200 {object} PaginatedResponse{data=[]TaskResponse}
// @Failure 401 {object} ErrorResponse
// @Router /tasks [get]
func (h *Handler) GetTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	statusStr := c.Query("status")
	priorityStr := c.Query("priority")

	var tasks []task.Task
	var err error

	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		tasks, err = h.userManager.GetUserTasksByStatus(userID.(int), task.Status(status))
	} else if priorityStr != "" {
		priority, _ := strconv.Atoi(priorityStr)
		tasks, err = h.userManager.GetUserTasksByPriority(userID.(int), task.Priority(priority))
	} else {
		tasks, err = h.userManager.GetUserTasks(userID.(int))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get tasks",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert tasks to responses
	taskResponses := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		taskResponses[i] = ConvertToTaskResponse(t)
	}

	// Simple pagination (in a real app, this would be done at the database level)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(taskResponses) {
		taskResponses = []TaskResponse{}
	} else if end > len(taskResponses) {
		taskResponses = taskResponses[start:]
	} else {
		taskResponses = taskResponses[start:end]
	}

	totalPages := (len(tasks) + pageSize - 1) / pageSize
	pagination := Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: len(tasks),
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    "Tasks retrieved successfully",
		Data:       taskResponses,
		Pagination: pagination,
	})
}

// GetTask handles getting a specific task
// @Summary Get a specific task
// @Description Get a task by ID for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} APIResponse{data=TaskResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /tasks/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid task ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	task, err := h.userManager.GetUserTask(userID.(int), taskID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied: task does not belong to user or user ID is missing" {
			status = http.StatusNotFound
		}
		c.JSON(status, ErrorResponse{
			Success: false,
			Message: "Failed to get task",
			Error:   err.Error(),
			Code:    status,
		})
		return
	}

	taskResponse := ConvertToTaskResponse(*task)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Task retrieved successfully",
		Data:    taskResponse,
	})
}

// UpdateTaskStatus handles updating task status
// @Summary Update task status
// @Description Update the status of a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param status body map[string]int true "New status"
// @Success 200 {object} APIResponse{data=TaskResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /tasks/{id}/status [put]
func (h *Handler) UpdateTaskStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid task ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var statusUpdate struct {
		Status int `json:"status" binding:"required,min=0,max=3"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.userManager.UpdateUserTaskStatus(userID.(int), taskID, task.Status(statusUpdate.Status))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied: task does not belong to user or user ID is missing" {
			status = http.StatusNotFound
		}
		c.JSON(status, ErrorResponse{
			Success: false,
			Message: "Failed to update task status",
			Error:   err.Error(),
			Code:    status,
		})
		return
	}

	// Get updated task
	updatedTask, err := h.userManager.GetUserTask(userID.(int), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get updated task",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	taskResponse := ConvertToTaskResponse(*updatedTask)
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Task status updated successfully",
		Data:    taskResponse,
	})
}

// DeleteTask handles task deletion
// @Summary Delete a task
// @Description Delete a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid task ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.userManager.DeleteUserTask(userID.(int), taskID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied: task does not belong to user or user ID is missing" {
			status = http.StatusNotFound
		}
		c.JSON(status, ErrorResponse{
			Success: false,
			Message: "Failed to delete task",
			Error:   err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Task deleted successfully",
	})
}

// SearchTasks handles task search
// @Summary Search tasks
// @Description Search tasks with various filters
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search body SearchRequest true "Search criteria"
// @Success 200 {object} PaginatedResponse{data=[]TaskResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /tasks/search [post]
func (h *Handler) SearchTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Convert SearchRequest parameters
	var status *task.Status
	var priority *task.Priority
	if req.Status != nil {
		s := task.Status(*req.Status)
		status = &s
	}
	if req.Priority != nil {
		p := task.Priority(*req.Priority)
		priority = &p
	}

	userIDInt := userID.(int)
	searchResults, err := h.searchManager.SearchTasksWithFilters(
		req.Query,
		&userIDInt,
		status,
		priority,
		req.CategoryID,
		req.TagNames,
		req.PageSize,
		(req.Page-1)*req.PageSize,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to search tasks",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert tasks to responses
	taskResponses := make([]TaskResponse, len(searchResults))
	for i, result := range searchResults {
		// Convert TaskResult to Task
		task := task.Task{
			ID:          result.ID,
			Title:       result.Title,
			Description: result.Description,
			Priority:    task.Priority(result.Priority),
			Status:      task.Status(result.Status),
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
			DueDate:     result.DueDate,
		}
		taskResponses[i] = ConvertToTaskResponse(task)
	}

	// Simple pagination (in a real app, this would be done at the database level)
	totalItems := len(searchResults)
	totalPages := (totalItems + req.PageSize - 1) / req.PageSize
	pagination := Pagination{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    "Search completed successfully",
		Data:       taskResponses,
		Pagination: pagination,
	})
}

// GetStatistics handles getting application statistics
// @Summary Get statistics
// @Description Get application statistics for the authenticated user
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse{data=StatisticsResponse}
// @Failure 401 {object} ErrorResponse
// @Router /statistics [get]
func (h *Handler) GetStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Get user tasks
	tasks, err := h.userManager.GetUserTasks(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get statistics",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Calculate statistics
	stats := StatisticsResponse{
		TotalTasks: len(tasks),
	}

	for _, t := range tasks {
		switch t.Status {
		case task.Completed:
			stats.CompletedTasks++
		case task.Pending:
			stats.PendingTasks++
		}
	}

	// Count overdue tasks
	overdueTasks, err := h.userManager.GetUserOverdueTasks(userID.(int))
	if err == nil {
		stats.OverdueTasks = len(overdueTasks)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Statistics retrieved successfully",
		Data:    stats,
	})
}
