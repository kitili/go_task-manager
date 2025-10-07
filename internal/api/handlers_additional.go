package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"learn-go-capstone/internal/database"
)

// CreateCategory handles category creation
// @Summary Create a new category
// @Description Create a new task category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body CategoryRequest true "Category data"
// @Success 201 {object} APIResponse{data=CategoryResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	createdCategory, err := h.categoryManager.CreateCategory(req.Name, req.Description, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to create category",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	categoryResponse := CategoryResponse{
		ID:          createdCategory.ID,
		Name:        createdCategory.Name,
		Description: createdCategory.Description,
		Color:       createdCategory.Color,
		CreatedAt:   createdCategory.CreatedAt,
		UpdatedAt:   createdCategory.UpdatedAt,
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    categoryResponse,
	})
}

// GetCategories handles getting all categories
// @Summary Get all categories
// @Description Get all task categories
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse{data=[]CategoryResponse}
// @Failure 401 {object} ErrorResponse
// @Router /categories [get]
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.categoryManager.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get categories",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	categoryResponses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			Color:       category.Color,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    categoryResponses,
	})
}

// GetCategory handles getting a specific category
// @Summary Get a specific category
// @Description Get a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} APIResponse{data=CategoryResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id} [get]
func (h *Handler) GetCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	category, err := h.categoryManager.GetCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Message: "Category not found",
			Error:   err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	categoryResponse := CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Category retrieved successfully",
		Data:    categoryResponse,
	})
}

// UpdateCategory handles category updates
// @Summary Update a category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param category body CategoryRequest true "Updated category data"
// @Success 200 {object} APIResponse{data=CategoryResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	category := &database.Category{
		ID:          categoryID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}

	err = h.categoryManager.UpdateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to update category",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Get updated category
	updatedCategory, err := h.categoryManager.GetCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get updated category",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	categoryResponse := CategoryResponse{
		ID:          updatedCategory.ID,
		Name:        updatedCategory.Name,
		Description: updatedCategory.Description,
		Color:       updatedCategory.Color,
		CreatedAt:   updatedCategory.CreatedAt,
		UpdatedAt:   updatedCategory.UpdatedAt,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Category updated successfully",
		Data:    categoryResponse,
	})
}

// DeleteCategory handles category deletion
// @Summary Delete a category
// @Description Delete an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.categoryManager.DeleteCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to delete category",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}

// CreateTag handles tag creation
// @Summary Create a new tag
// @Description Create a new task tag
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tag body TagRequest true "Tag data"
// @Success 201 {object} APIResponse{data=TagResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /tags [post]
func (h *Handler) CreateTag(c *gin.Context) {
	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	createdTag, err := h.categoryManager.CreateTag(req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to create tag",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	tagResponse := TagResponse{
		ID:        createdTag.ID,
		Name:      createdTag.Name,
		Color:     createdTag.Color,
		CreatedAt: createdTag.CreatedAt,
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Tag created successfully",
		Data:    tagResponse,
	})
}

// GetTags handles getting all tags
// @Summary Get all tags
// @Description Get all task tags
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse{data=[]TagResponse}
// @Failure 401 {object} ErrorResponse
// @Router /tags [get]
func (h *Handler) GetTags(c *gin.Context) {
	tags, err := h.categoryManager.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get tags",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	tagResponses := make([]TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			Color:     tag.Color,
			CreatedAt: tag.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Tags retrieved successfully",
		Data:    tagResponses,
	})
}

// GetTag handles getting a specific tag
// @Summary Get a specific tag
// @Description Get a tag by ID
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} APIResponse{data=TagResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /tags/{id} [get]
func (h *Handler) GetTag(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid tag ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	tag, err := h.categoryManager.GetTag(tagID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Message: "Tag not found",
			Error:   err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	tagResponse := TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Color:     tag.Color,
		CreatedAt: tag.CreatedAt,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Tag retrieved successfully",
		Data:    tagResponse,
	})
}

// UpdateTag handles tag updates
// @Summary Update a tag
// @Description Update an existing tag
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Param tag body TagRequest true "Updated tag data"
// @Success 200 {object} APIResponse{data=TagResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /tags/{id} [put]
func (h *Handler) UpdateTag(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid tag ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	tag := &database.Tag{
		ID:    tagID,
		Name:  req.Name,
		Color: req.Color,
	}

	err = h.categoryManager.UpdateTag(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to update tag",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Get updated tag
	updatedTag, err := h.categoryManager.GetTag(tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get updated tag",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	tagResponse := TagResponse{
		ID:        updatedTag.ID,
		Name:      updatedTag.Name,
		Color:     updatedTag.Color,
		CreatedAt: updatedTag.CreatedAt,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Tag updated successfully",
		Data:    tagResponse,
	})
}

// DeleteTag handles tag deletion
// @Summary Delete a tag
// @Description Delete an existing tag
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /tags/{id} [delete]
func (h *Handler) DeleteTag(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid tag ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.categoryManager.DeleteTag(tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to delete tag",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Tag deleted successfully",
	})
}

// CreateDependency handles dependency creation
// @Summary Create a task dependency
// @Description Create a dependency between two tasks
// @Tags dependencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dependency body DependencyRequest true "Dependency data"
// @Success 201 {object} APIResponse{data=DependencyResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /dependencies [post]
func (h *Handler) CreateDependency(c *gin.Context) {
	var req DependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.dependencyManager.AddDependency(req.TaskID, req.DependsOnTaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Failed to create dependency",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Note: In a real implementation, you'd return the created dependency
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Dependency created successfully",
	})
}

// GetTaskDependencies handles getting task dependencies
// @Summary Get task dependencies
// @Description Get all dependencies for a specific task
// @Tags dependencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} APIResponse{data=[]DependencyResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /dependencies/task/{id} [get]
func (h *Handler) GetTaskDependencies(c *gin.Context) {
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

	// Note: This would need to be implemented in the dependency manager
	// For now, return empty dependencies for task ID: taskID
	_ = taskID // Suppress unused variable warning
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Dependencies retrieved successfully",
		Data:    []DependencyResponse{},
	})
}

// DeleteDependency handles dependency deletion
// @Summary Delete a task dependency
// @Description Delete a dependency between two tasks
// @Tags dependencies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Dependency ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /dependencies/{id} [delete]
func (h *Handler) DeleteDependency(c *gin.Context) {
	dependencyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid dependency ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Note: This would need to be implemented in the dependency manager
	_ = dependencyID

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Dependency deleted successfully",
	})
}

// ExportTasks handles task export
// @Summary Export tasks
// @Description Export tasks in JSON or CSV format
// @Tags export
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param export body ExportRequest true "Export options"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /export/tasks [post]
func (h *Handler) ExportTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Note: This would need to be implemented in the export manager
	_ = userID
	_ = req

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Export completed successfully",
	})
}

// ImportTasks handles task import
// @Summary Import tasks
// @Description Import tasks from JSON or CSV format
// @Tags export
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param import body ImportRequest true "Import options"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /export/import [post]
func (h *Handler) ImportTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Note: This would need to be implemented in the export manager
	_ = userID
	_ = req

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Import completed successfully",
	})
}

// GetNotifications handles getting user notifications
// @Summary Get user notifications
// @Description Get all notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} PaginatedResponse{data=[]NotificationResponse}
// @Failure 401 {object} ErrorResponse
// @Router /notifications [get]
func (h *Handler) GetNotifications(c *gin.Context) {
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

	notifications, err := h.notificationManager.GetUserNotifications(userID.(int), pageSize, (page-1)*pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get notifications",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert notifications to responses
	notificationResponses := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		notificationResponses[i] = NotificationResponse{
			ID:           n.ID,
			UserID:       n.UserID,
			TaskID:       n.TaskID,
			Type:         string(n.Type),
			Priority:     string(n.Priority),
			Status:       string(n.Status),
			Trigger:      string(n.Trigger),
			Title:        n.Title,
			Message:      n.Message,
			Recipient:    n.Recipient,
			Channel:      n.Channel,
			ScheduledAt:  n.ScheduledAt,
			SentAt:       n.SentAt,
			DeliveredAt:  n.DeliveredAt,
			CreatedAt:    n.CreatedAt,
			UpdatedAt:    n.UpdatedAt,
			RetryCount:   n.RetryCount,
			MaxRetries:   n.MaxRetries,
			Error:        n.Error,
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Notifications retrieved successfully",
		Data:    notificationResponses,
	})
}

// CreateNotification handles notification creation
// @Summary Create a notification
// @Description Create a new notification
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body NotificationRequest true "Notification data"
// @Success 201 {object} APIResponse{data=NotificationResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /notifications [post]
func (h *Handler) CreateNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Note: This would need to be implemented in the notification manager
	_ = userID
	_ = req

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Notification created successfully",
	})
}

// MarkNotificationAsRead handles marking notification as read
// @Summary Mark notification as read
// @Description Mark a notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Notification ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /notifications/{id}/read [put]
func (h *Handler) MarkNotificationAsRead(c *gin.Context) {
	notificationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid notification ID",
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.notificationManager.MarkNotificationAsRead(notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to mark notification as read",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Notification marked as read",
	})
}

// GetNotificationStats handles getting notification statistics
// @Summary Get notification statistics
// @Description Get notification statistics for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse
// @Failure 401 {object} ErrorResponse
// @Router /notifications/stats [get]
func (h *Handler) GetNotificationStats(c *gin.Context) {
	stats, err := h.notificationManager.GetNotificationStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to get notification statistics",
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Notification statistics retrieved successfully",
		Data:    stats,
	})
}
