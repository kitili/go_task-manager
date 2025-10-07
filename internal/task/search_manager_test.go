package task

import (
	"os"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/search"
)

func TestSearchManager(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_search_manager.db"
	defer os.Remove(tempDB)
	
	// Create database configuration
	config := &database.Config{
		Driver: "sqlite3",
		DSN:    tempDB,
	}
	
	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	
	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Create repository and search manager
	repository := database.NewSQLiteRepository(db)
	searchManager := NewSearchManager(repository)
	
	// Create test user
	user := &database.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		IsActive: true,
	}
	err = repository.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Create test category
	category := &database.Category{
		Name:        "Work",
		Description: "Work-related tasks",
		Color:       "#ff0000",
	}
	err = repository.CreateCategory(category)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	
	// Create test tag
	tag := &database.Tag{
		Name:  "urgent",
		Color: "#ff0000",
	}
	err = repository.CreateTag(tag)
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}
	
	// Create test tasks
	now := time.Now()
	task1 := &database.DatabaseTask{
		Title:       "Learn Go Programming",
		Description: "Study Go language fundamentals",
		Priority:    3,
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserID:      &user.ID,
		CategoryID:  &category.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(task1)
	if err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}
	
	task2 := &database.DatabaseTask{
		Title:       "Go Project Development",
		Description: "Build a comprehensive Go application",
		Priority:    2,
		Status:      1,
		CreatedAt:   now.Add(-24 * time.Hour),
		UpdatedAt:   now.Add(-24 * time.Hour),
		UserID:      &user.ID,
		CategoryID:  &category.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(task2)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}
	
	// Add tag to task
	err = repository.AddTagToTask(task1.ID, tag.ID)
	if err != nil {
		t.Fatalf("Failed to add tag to task: %v", err)
	}
	
	// Test text search
	results, err := searchManager.SearchTasksByText("Go")
	if err != nil {
		t.Fatalf("Failed to search tasks by text: %v", err)
	}
	
	if len(results) < 2 {
		t.Errorf("Expected at least 2 results for 'Go' search, got %d", len(results))
	}
	
	// Test user-specific search
	userResults, err := searchManager.SearchTasksByUser(user.ID, "Go")
	if err != nil {
		t.Fatalf("Failed to search user tasks: %v", err)
	}
	
	if len(userResults) < 2 {
		t.Errorf("Expected at least 2 results for user 'Go' search, got %d", len(userResults))
	}
	
	// Test tag search
	tagResults, err := searchManager.SearchTasksByTag("urgent")
	if err != nil {
		t.Fatalf("Failed to search tasks by tag: %v", err)
	}
	
	if len(tagResults) != 1 {
		t.Errorf("Expected 1 result for 'urgent' tag search, got %d", len(tagResults))
	}
	
	// Test category search
	categoryResults, err := searchManager.SearchTasksByCategory("Work")
	if err != nil {
		t.Fatalf("Failed to search tasks by category: %v", err)
	}
	
	if len(categoryResults) < 2 {
		t.Errorf("Expected at least 2 results for 'Work' category search, got %d", len(categoryResults))
	}
	
	// Test search with filters
	status := Pending
	priority := High
	results, err = searchManager.SearchTasksWithFilters(
		"Go",
		&user.ID,
		&status,
		&priority,
		&category.ID,
		[]string{"urgent"},
		10,
		0,
	)
	if err != nil {
		t.Fatalf("Failed to search tasks with filters: %v", err)
	}
	
	if len(results) < 1 {
		t.Errorf("Expected at least 1 result for filtered search, got %d", len(results))
	}
	
	// Test get tasks by multiple filters
	tasks, err := searchManager.GetTasksByMultipleFilters(
		&user.ID,
		&status,
		&priority,
		&category.ID,
		[]string{"urgent"},
		10,
		0,
	)
	if err != nil {
		t.Fatalf("Failed to get tasks by multiple filters: %v", err)
	}
	
	if len(tasks) < 1 {
		t.Errorf("Expected at least 1 task for filtered search, got %d", len(tasks))
	}
	
	// Test get recent tasks
	recentTasks, err := searchManager.GetRecentTasks(5)
	if err != nil {
		t.Fatalf("Failed to get recent tasks: %v", err)
	}
	
	if len(recentTasks) < 2 {
		t.Errorf("Expected at least 2 recent tasks, got %d", len(recentTasks))
	}
	
	// Test get tasks by date range
	dateFrom := now.Add(-48 * time.Hour)
	dateTo := now
	dateRangeTasks, err := searchManager.GetTasksByDateRange(
		&user.ID,
		&dateFrom,
		&dateTo,
		10,
		0,
	)
	if err != nil {
		t.Fatalf("Failed to get tasks by date range: %v", err)
	}
	
	if len(dateRangeTasks) < 2 {
		t.Errorf("Expected at least 2 tasks in date range, got %d", len(dateRangeTasks))
	}
	
	// Test get overdue tasks with filters
	_, err = searchManager.GetOverdueTasksWithFilters(
		&user.ID,
		&priority,
		&category.ID,
		10,
		0,
	)
	if err != nil {
		t.Fatalf("Failed to get overdue tasks with filters: %v", err)
	}
	
	// Note: No overdue tasks in this test, so we just check for no error
	
	// Test get filter options
	filterOptions, err := searchManager.GetFilterOptions()
	if err != nil {
		t.Fatalf("Failed to get filter options: %v", err)
	}
	
	if len(filterOptions.Statuses) == 0 {
		t.Error("Expected status options to be available")
	}
	
	if len(filterOptions.Priorities) == 0 {
		t.Error("Expected priority options to be available")
	}
	
	if len(filterOptions.Categories) == 0 {
		t.Error("Expected category options to be available")
	}
	
	if len(filterOptions.Tags) == 0 {
		t.Error("Expected tag options to be available")
	}
	
	// Test get task statistics
	stats, err := searchManager.GetTaskStatistics()
	if err != nil {
		t.Fatalf("Failed to get task statistics: %v", err)
	}
	
	if stats.TotalTasks < 2 {
		t.Errorf("Expected at least 2 total tasks, got %d", stats.TotalTasks)
	}
}

func TestSearchQueryBuilding(t *testing.T) {
	// Test default search query
	query := search.DefaultSearchQuery()
	if query.Limit != 20 {
		t.Errorf("Expected default limit 20, got %d", query.Limit)
	}
	
	if query.SortBy != "created_at" {
		t.Errorf("Expected default sort by 'created_at', got '%s'", query.SortBy)
	}
	
	if query.SortOrder != "desc" {
		t.Errorf("Expected default sort order 'desc', got '%s'", query.SortOrder)
	}
	
	// Test query validation
	err := query.Validate()
	if err != nil {
		t.Fatalf("Expected validation to pass, got error: %v", err)
	}
	
	// Test custom query
	customQuery := search.SearchQuery{
		Query:      "test",
		Limit:      10,
		Offset:     5,
		SortBy:     "title",
		SortOrder:  "asc",
	}
	
	err = customQuery.Validate()
	if err != nil {
		t.Fatalf("Expected custom query validation to pass, got error: %v", err)
	}
	
	if customQuery.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", customQuery.Limit)
	}
	
	if customQuery.Offset != 5 {
		t.Errorf("Expected offset 5, got %d", customQuery.Offset)
	}
	
	if customQuery.SortBy != "title" {
		t.Errorf("Expected sort by 'title', got '%s'", customQuery.SortBy)
	}
	
	if customQuery.SortOrder != "asc" {
		t.Errorf("Expected sort order 'asc', got '%s'", customQuery.SortOrder)
	}
}
