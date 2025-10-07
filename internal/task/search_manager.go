package task

import (
	"time"

	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/search"
)

// SearchManager manages search and filtering operations for tasks
type SearchManager struct {
	repository    database.Repository
	searchService *search.SearchService
}

// NewSearchManager creates a new search manager
func NewSearchManager(repository database.Repository) *SearchManager {
	return &SearchManager{
		repository:    repository,
		searchService: search.NewSearchService(repository),
	}
}

// SearchTasks performs a comprehensive search with filters
func (sm *SearchManager) SearchTasks(query search.SearchQuery) (*search.SearchResult, error) {
	return sm.searchService.SearchTasks(query)
}

// SearchTasksByText performs a simple text search
func (sm *SearchManager) SearchTasksByText(query string) ([]search.TaskResult, error) {
	return sm.searchService.SearchTasksByText(query)
}

// SearchTasksByUser performs a text search for a specific user
func (sm *SearchManager) SearchTasksByUser(userID int, query string) ([]search.TaskResult, error) {
	return sm.searchService.SearchTasksByUser(userID, query)
}

// SearchTasksByTag performs a search by tag name
func (sm *SearchManager) SearchTasksByTag(tagName string) ([]search.TaskResult, error) {
	return sm.searchService.SearchTasksByTag(tagName)
}

// SearchTasksByCategory performs a search by category name
func (sm *SearchManager) SearchTasksByCategory(categoryName string) ([]search.TaskResult, error) {
	return sm.searchService.SearchTasksByCategory(categoryName)
}

// GetFilterOptions returns available filter options
func (sm *SearchManager) GetFilterOptions() (*search.FilterOptions, error) {
	return sm.searchService.GetFilterOptions()
}

// SearchTasksWithFilters performs a search with multiple filters
func (sm *SearchManager) SearchTasksWithFilters(
	textQuery string,
	userID *int,
	status *Status,
	priority *Priority,
	categoryID *int,
	tagNames []string,
	limit, offset int,
) ([]search.TaskResult, error) {
	// Build search query
	searchQuery := search.SearchQuery{
		Query:      textQuery,
		UserID:     userID,
		Status:     (*int)(status),
		Priority:   (*int)(priority),
		CategoryID: categoryID,
		TagNames:   tagNames,
		Limit:      limit,
		Offset:     offset,
	}
	
	// Perform search
	result, err := sm.SearchTasks(searchQuery)
	if err != nil {
		return nil, err
	}
	
	return result.Tasks, nil
}

// GetTasksByMultipleFilters gets tasks using multiple filter criteria
func (sm *SearchManager) GetTasksByMultipleFilters(
	userID *int,
	status *Status,
	priority *Priority,
	categoryID *int,
	tagNames []string,
	limit, offset int,
) ([]Task, error) {
	// Build search query
	searchQuery := search.SearchQuery{
		UserID:     userID,
		Status:     (*int)(status),
		Priority:   (*int)(priority),
		CategoryID: categoryID,
		TagNames:   tagNames,
		Limit:      limit,
		Offset:     offset,
	}
	
	// Perform search
	result, err := sm.SearchTasks(searchQuery)
	if err != nil {
		return nil, err
	}
	
	// Convert to Task objects
	var tasks []Task
	for _, taskResult := range result.Tasks {
		task := Task{
			ID:          taskResult.ID,
			Title:       taskResult.Title,
			Description: taskResult.Description,
			Priority:    Priority(taskResult.Priority),
			Status:      Status(taskResult.Status),
			CreatedAt:   taskResult.CreatedAt,
			UpdatedAt:   taskResult.UpdatedAt,
			DueDate:     taskResult.DueDate,
		}
		
		// Add category if available
		if taskResult.Category != nil {
			task.Category = &Category{
				ID:          taskResult.Category.ID,
				Name:        taskResult.Category.Name,
				Description: taskResult.Category.Description,
				Color:       taskResult.Category.Color,
			}
		}
		
		// Add tags if available
		if len(taskResult.Tags) > 0 {
			task.Tags = make([]Tag, len(taskResult.Tags))
			for i, tagResult := range taskResult.Tags {
				task.Tags[i] = Tag{
					ID:    tagResult.ID,
					Name:  tagResult.Name,
					Color: tagResult.Color,
				}
			}
		}
		
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// GetRecentTasks gets recently created or updated tasks
func (sm *SearchManager) GetRecentTasks(limit int) ([]Task, error) {
	searchQuery := search.SearchQuery{
		Limit:     limit,
		Offset:    0,
		SortBy:    "updated_at",
		SortOrder: "desc",
	}
	
	result, err := sm.SearchTasks(searchQuery)
	if err != nil {
		return nil, err
	}
	
	// Convert to Task objects
	var tasks []Task
	for _, taskResult := range result.Tasks {
		task := Task{
			ID:          taskResult.ID,
			Title:       taskResult.Title,
			Description: taskResult.Description,
			Priority:    Priority(taskResult.Priority),
			Status:      Status(taskResult.Status),
			CreatedAt:   taskResult.CreatedAt,
			UpdatedAt:   taskResult.UpdatedAt,
			DueDate:     taskResult.DueDate,
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// GetTasksByDateRange gets tasks within a date range
func (sm *SearchManager) GetTasksByDateRange(
	userID *int,
	dateFrom, dateTo *time.Time,
	limit, offset int,
) ([]Task, error) {
	searchQuery := search.SearchQuery{
		UserID:    userID,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		Limit:     limit,
		Offset:    offset,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	
	result, err := sm.SearchTasks(searchQuery)
	if err != nil {
		return nil, err
	}
	
	// Convert to Task objects
	var tasks []Task
	for _, taskResult := range result.Tasks {
		task := Task{
			ID:          taskResult.ID,
			Title:       taskResult.Title,
			Description: taskResult.Description,
			Priority:    Priority(taskResult.Priority),
			Status:      Status(taskResult.Status),
			CreatedAt:   taskResult.CreatedAt,
			UpdatedAt:   taskResult.UpdatedAt,
			DueDate:     taskResult.DueDate,
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// GetOverdueTasksWithFilters gets overdue tasks with additional filters
func (sm *SearchManager) GetOverdueTasksWithFilters(
	userID *int,
	priority *Priority,
	categoryID *int,
	limit, offset int,
) ([]Task, error) {
	isOverdue := true
	searchQuery := search.SearchQuery{
		UserID:     userID,
		Priority:   (*int)(priority),
		CategoryID: categoryID,
		IsOverdue:  &isOverdue,
		Limit:      limit,
		Offset:     offset,
		SortBy:    "due_date",
		SortOrder: "asc",
	}
	
	result, err := sm.SearchTasks(searchQuery)
	if err != nil {
		return nil, err
	}
	
	// Convert to Task objects
	var tasks []Task
	for _, taskResult := range result.Tasks {
		task := Task{
			ID:          taskResult.ID,
			Title:       taskResult.Title,
			Description: taskResult.Description,
			Priority:    Priority(taskResult.Priority),
			Status:      Status(taskResult.Status),
			CreatedAt:   taskResult.CreatedAt,
			UpdatedAt:   taskResult.UpdatedAt,
			DueDate:     taskResult.DueDate,
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// GetTaskStatistics gets search statistics
func (sm *SearchManager) GetTaskStatistics() (*search.SearchStats, error) {
	// Get all tasks
	allTasks, err := sm.repository.GetAllTasks()
	if err != nil {
		return nil, err
	}
	
	// Count by status
	statusCounts := make(map[string]int)
	for _, task := range allTasks {
		switch task.Status {
		case 0:
			statusCounts["pending"]++
		case 1:
			statusCounts["in_progress"]++
		case 2:
			statusCounts["completed"]++
		case 3:
			statusCounts["cancelled"]++
		}
	}
	
	return &search.SearchStats{
		TotalTasks: len(allTasks),
		FilteredBy: statusCounts,
	}, nil
}
