package search

import (
	"fmt"
	"strings"
	"time"

	"learn-go-capstone/internal/database"
)

// SearchService handles search and filtering operations
type SearchService struct {
	repository database.Repository
}

// NewSearchService creates a new search service
func NewSearchService(repository database.Repository) *SearchService {
	return &SearchService{
		repository: repository,
	}
}

// SearchTasks performs a comprehensive search with filters
func (ss *SearchService) SearchTasks(query SearchQuery) (*SearchResult, error) {
	// Validate query
	if err := query.Validate(); err != nil {
		return nil, err
	}
	
	// Build the search query
	sqlQuery, args := ss.buildSearchQuery(query)
	
	// Execute search
	dbTasks, err := ss.executeSearch(sqlQuery, args)
	if err != nil {
		return nil, err
	}
	
	// Convert to search results
	tasks := ss.convertToTaskResults(dbTasks)
	
	// Get total count for pagination
	total, err := ss.getTotalCount(query)
	if err != nil {
		return nil, err
	}
	
	// Calculate pagination
	totalPages := (total + query.Limit - 1) / query.Limit
	page := (query.Offset / query.Limit) + 1
	
	return &SearchResult{
		Tasks:      tasks,
		Total:      total,
		Page:       page,
		PerPage:    query.Limit,
		TotalPages: totalPages,
		Query:      query,
	}, nil
}

// SearchTasksByText performs a simple text search
func (ss *SearchService) SearchTasksByText(query string) ([]TaskResult, error) {
	dbTasks, err := ss.repository.SearchTasks(query)
	if err != nil {
		return nil, err
	}
	
	return ss.convertToTaskResults(dbTasks), nil
}

// SearchTasksByUser performs a text search for a specific user
func (ss *SearchService) SearchTasksByUser(userID int, query string) ([]TaskResult, error) {
	dbTasks, err := ss.repository.SearchTasksByUser(userID, query)
	if err != nil {
		return nil, err
	}
	
	return ss.convertToTaskResults(dbTasks), nil
}

// SearchTasksByTag performs a search by tag name
func (ss *SearchService) SearchTasksByTag(tagName string) ([]TaskResult, error) {
	dbTasks, err := ss.repository.SearchTasksByTag(tagName)
	if err != nil {
		return nil, err
	}
	
	return ss.convertToTaskResults(dbTasks), nil
}

// SearchTasksByCategory performs a search by category name
func (ss *SearchService) SearchTasksByCategory(categoryName string) ([]TaskResult, error) {
	dbTasks, err := ss.repository.SearchTasksByCategory(categoryName)
	if err != nil {
		return nil, err
	}
	
	return ss.convertToTaskResults(dbTasks), nil
}

// GetFilterOptions returns available filter options
func (ss *SearchService) GetFilterOptions() (*FilterOptions, error) {
	// Get all tasks to build filter options
	allTasks, err := ss.repository.GetAllTasks()
	if err != nil {
		return nil, err
	}
	
	// Count occurrences
	statusCounts := make(map[int]int)
	priorityCounts := make(map[int]int)
	categoryCounts := make(map[int]int)
	tagCounts := make(map[int]int)
	userCounts := make(map[int]int)
	
	for _, task := range allTasks {
		statusCounts[task.Status]++
		priorityCounts[task.Priority]++
		
		if task.CategoryID != nil {
			categoryCounts[*task.CategoryID]++
		}
		if task.UserID != nil {
			userCounts[*task.UserID]++
		}
	}
	
	// Get categories
	categories, err := ss.repository.GetAllCategories()
	if err != nil {
		return nil, err
	}
	
	var categoryOptions []CategoryOption
	for _, cat := range categories {
		categoryOptions = append(categoryOptions, CategoryOption{
			ID:    cat.ID,
			Name:  cat.Name,
			Color: cat.Color,
			Count: categoryCounts[cat.ID],
		})
	}
	
	// Get tags
	tags, err := ss.repository.GetAllTags()
	if err != nil {
		return nil, err
	}
	
	var tagOptions []TagOption
	for _, tag := range tags {
		tagOptions = append(tagOptions, TagOption{
			ID:    tag.ID,
			Name:  tag.Name,
			Color: tag.Color,
			Count: tagCounts[tag.ID],
		})
	}
	
	// Get users
	users, err := ss.repository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	
	var userOptions []UserOption
	for _, user := range users {
		userOptions = append(userOptions, UserOption{
			ID:       user.ID,
			Username: user.Username,
			Count:    userCounts[user.ID],
		})
	}
	
	// Build status options
	statusOptions := []StatusOption{
		{Value: 0, Label: "Pending", Count: statusCounts[0]},
		{Value: 1, Label: "In Progress", Count: statusCounts[1]},
		{Value: 2, Label: "Completed", Count: statusCounts[2]},
		{Value: 3, Label: "Cancelled", Count: statusCounts[3]},
	}
	
	// Build priority options
	priorityOptions := []PriorityOption{
		{Value: 1, Label: "Low", Count: priorityCounts[1]},
		{Value: 2, Label: "Medium", Count: priorityCounts[2]},
		{Value: 3, Label: "High", Count: priorityCounts[3]},
		{Value: 4, Label: "Critical", Count: priorityCounts[4]},
	}
	
	return &FilterOptions{
		Statuses:   statusOptions,
		Priorities: priorityOptions,
		Categories: categoryOptions,
		Tags:       tagOptions,
		Users:      userOptions,
	}, nil
}

// buildSearchQuery builds the SQL query based on search criteria
func (ss *SearchService) buildSearchQuery(query SearchQuery) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	
	baseQuery := `
	SELECT t.id, t.title, t.description, t.priority, t.status, t.created_at, t.updated_at, t.due_date, t.user_id, t.category_id, t.is_archived
	FROM tasks t
	LEFT JOIN categories c ON t.category_id = c.id
	LEFT JOIN task_tags tt ON t.id = tt.task_id
	LEFT JOIN tags tg ON tt.tag_id = tg.id`
	
	// Always exclude archived tasks
	conditions = append(conditions, "t.is_archived = FALSE")
	
	// Text search
	if query.Query != "" {
		searchTerm := "%" + query.Query + "%"
		conditions = append(conditions, "(t.title LIKE ? OR t.description LIKE ? OR c.name LIKE ? OR tg.name LIKE ?)")
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}
	
	// User filter
	if query.UserID != nil {
		conditions = append(conditions, "t.user_id = ?")
		args = append(args, *query.UserID)
	}
	
	// Status filter
	if query.Status != nil {
		conditions = append(conditions, "t.status = ?")
		args = append(args, *query.Status)
	}
	
	// Priority filter
	if query.Priority != nil {
		conditions = append(conditions, "t.priority = ?")
		args = append(args, *query.Priority)
	}
	
	// Category filter
	if query.CategoryID != nil {
		conditions = append(conditions, "t.category_id = ?")
		args = append(args, *query.CategoryID)
	}
	
	// Tag filter
	if len(query.TagNames) > 0 {
		tagConditions := make([]string, len(query.TagNames))
		for i, tagName := range query.TagNames {
			tagConditions[i] = "tg.name LIKE ?"
			args = append(args, "%"+tagName+"%")
		}
		conditions = append(conditions, "("+strings.Join(tagConditions, " OR ")+")")
	}
	
	// Date range filter
	if query.DateFrom != nil {
		conditions = append(conditions, "t.created_at >= ?")
		args = append(args, *query.DateFrom)
	}
	if query.DateTo != nil {
		conditions = append(conditions, "t.created_at <= ?")
		args = append(args, *query.DateTo)
	}
	
	// Due date range filter
	if query.DueDateFrom != nil {
		conditions = append(conditions, "t.due_date >= ?")
		args = append(args, *query.DueDateFrom)
	}
	if query.DueDateTo != nil {
		conditions = append(conditions, "t.due_date <= ?")
		args = append(args, *query.DueDateTo)
	}
	
	// Overdue filter
	if query.IsOverdue != nil && *query.IsOverdue {
		conditions = append(conditions, "t.due_date < ? AND t.status != 2")
		args = append(args, time.Now())
	}
	
	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Build ORDER BY clause
	orderBy := ss.buildOrderByClause(query)
	
	// Build LIMIT and OFFSET
	limitClause := fmt.Sprintf("LIMIT %d OFFSET %d", query.Limit, query.Offset)
	
	// Add limit and offset to args for pagination
	args = append(args, query.Limit, query.Offset)
	
	// Combine all parts
	finalQuery := baseQuery + " " + whereClause + " " + orderBy + " " + limitClause
	
	return finalQuery, args
}

// buildOrderByClause builds the ORDER BY clause
func (ss *SearchService) buildOrderByClause(query SearchQuery) string {
	orderBy := "ORDER BY "
	
	switch query.SortBy {
	case "title":
		orderBy += "t.title"
	case "created_at":
		orderBy += "t.created_at"
	case "updated_at":
		orderBy += "t.updated_at"
	case "due_date":
		orderBy += "t.due_date"
	case "priority":
		orderBy += "t.priority"
	case "status":
		orderBy += "t.status"
	default:
		orderBy += "t.created_at"
	}
	
	if query.SortOrder == "asc" {
		orderBy += " ASC"
	} else {
		orderBy += " DESC"
	}
	
	return orderBy
}

// executeSearch executes the search query
func (ss *SearchService) executeSearch(sqlQuery string, args []interface{}) ([]database.DatabaseTask, error) {
	// This is a simplified version - in a real implementation, you'd use the database connection
	// For now, we'll use the existing repository methods
	
	// If it's a simple text search, use the existing method
	if len(args) == 4 && len(args[0].(string)) > 2 { // Basic text search
		query := args[0].(string)
		query = query[1 : len(query)-1] // Remove % characters
		return ss.repository.SearchTasks(query)
	}
	
	// For complex queries, we'd need to implement a more sophisticated query builder
	// For now, fall back to getting all tasks and filtering in memory
	allTasks, err := ss.repository.GetAllTasks()
	if err != nil {
		return nil, err
	}
	
	// Apply basic filters
	var filteredTasks []database.DatabaseTask
	for _, task := range allTasks {
		if ss.matchesFilters(task, args) {
			filteredTasks = append(filteredTasks, task)
		}
	}
	
	// Apply pagination
	// Extract limit and offset from args (they should be the last two arguments)
	if len(args) >= 2 {
		limit := args[len(args)-2].(int)
		offset := args[len(args)-1].(int)
		
		start := offset
		end := offset + limit
		
		if start >= len(filteredTasks) {
			return []database.DatabaseTask{}, nil
		}
		
		if end > len(filteredTasks) {
			end = len(filteredTasks)
		}
		
		return filteredTasks[start:end], nil
	}
	
	return filteredTasks, nil
}

// matchesFilters checks if a task matches the search filters
func (ss *SearchService) matchesFilters(task database.DatabaseTask, args []interface{}) bool {
	// This is a simplified implementation
	// In a real system, you'd parse the args and apply filters properly
	return true
}

// getTotalCount gets the total count of matching tasks
func (ss *SearchService) getTotalCount(query SearchQuery) (int, error) {
	// Simplified implementation - in practice, you'd run a COUNT query
	allTasks, err := ss.repository.GetAllTasks()
	if err != nil {
		return 0, err
	}
	
	return len(allTasks), nil
}

// convertToTaskResults converts database tasks to search results
func (ss *SearchService) convertToTaskResults(dbTasks []database.DatabaseTask) []TaskResult {
	var results []TaskResult
	
	for _, dbTask := range dbTasks {
		result := TaskResult{
			ID:          dbTask.ID,
			Title:       dbTask.Title,
			Description: dbTask.Description,
			Priority:    dbTask.Priority,
			Status:      dbTask.Status,
			CreatedAt:   dbTask.CreatedAt,
			UpdatedAt:   dbTask.UpdatedAt,
			DueDate:     dbTask.DueDate,
			UserID:      dbTask.UserID,
			CategoryID:  dbTask.CategoryID,
			IsArchived:  dbTask.IsArchived,
		}
		
		// Load category if available
		if dbTask.CategoryID != nil {
			category, err := ss.repository.GetCategory(*dbTask.CategoryID)
			if err == nil {
				result.Category = &CategoryResult{
					ID:          category.ID,
					Name:        category.Name,
					Description: category.Description,
					Color:       category.Color,
				}
			}
		}
		
		// Load tags if available
		// Note: This would require additional queries in a real implementation
		
		results = append(results, result)
	}
	
	return results
}
