package search

import (
	"fmt"
	"os"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
)

func TestSearchService(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_search.db"
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
	
	// Create repository and search service
	repository := database.NewSQLiteRepository(db)
	searchService := NewSearchService(repository)
	
	// Create test data
	user1 := &database.User{
		Username: "user1",
		Email:    "user1@example.com",
		Password: "password123",
		IsActive: true,
	}
	err = repository.CreateUser(user1)
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}
	
	user2 := &database.User{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "password123",
		IsActive: true,
	}
	err = repository.CreateUser(user2)
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}
	
	// Create categories
	category1 := &database.Category{
		Name:        "Work",
		Description: "Work-related tasks",
		Color:       "#ff0000",
	}
	err = repository.CreateCategory(category1)
	if err != nil {
		t.Fatalf("Failed to create category1: %v", err)
	}
	
	category2 := &database.Category{
		Name:        "Personal",
		Description: "Personal tasks",
		Color:       "#00ff00",
	}
	err = repository.CreateCategory(category2)
	if err != nil {
		t.Fatalf("Failed to create category2: %v", err)
	}
	
	// Create tags
	tag1 := &database.Tag{
		Name:  "urgent",
		Color: "#ff0000",
	}
	err = repository.CreateTag(tag1)
	if err != nil {
		t.Fatalf("Failed to create tag1: %v", err)
	}
	
	tag2 := &database.Tag{
		Name:  "important",
		Color: "#0000ff",
	}
	err = repository.CreateTag(tag2)
	if err != nil {
		t.Fatalf("Failed to create tag2: %v", err)
	}
	
	// Create tasks
	now := time.Now()
	task1 := &database.DatabaseTask{
		Title:       "Learn Go Programming",
		Description: "Study Go language fundamentals and best practices",
		Priority:    3,
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserID:      &user1.ID,
		CategoryID:  &category1.ID,
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
		UserID:      &user1.ID,
		CategoryID:  &category1.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(task2)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}
	
	task3 := &database.DatabaseTask{
		Title:       "Personal Reading",
		Description: "Read books about Go and software development",
		Priority:    1,
		Status:      0,
		CreatedAt:   now.Add(-48 * time.Hour),
		UpdatedAt:   now.Add(-48 * time.Hour),
		UserID:      &user2.ID,
		CategoryID:  &category2.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(task3)
	if err != nil {
		t.Fatalf("Failed to create task3: %v", err)
	}
	
	// Add tags to tasks
	err = repository.AddTagToTask(task1.ID, tag1.ID)
	if err != nil {
		t.Fatalf("Failed to add tag to task1: %v", err)
	}
	
	err = repository.AddTagToTask(task2.ID, tag2.ID)
	if err != nil {
		t.Fatalf("Failed to add tag to task2: %v", err)
	}
	
	// Test text search
	results, err := searchService.SearchTasksByText("Go")
	if err != nil {
		t.Fatalf("Failed to search tasks by text: %v", err)
	}
	
	if len(results) < 2 {
		t.Errorf("Expected at least 2 results for 'Go' search, got %d", len(results))
	}
	
	// Test user-specific search
	userResults, err := searchService.SearchTasksByUser(user1.ID, "Go")
	if err != nil {
		t.Fatalf("Failed to search user tasks: %v", err)
	}
	
	if len(userResults) < 2 {
		t.Errorf("Expected at least 2 results for user1 'Go' search, got %d", len(userResults))
	}
	
	// Test tag search
	tagResults, err := searchService.SearchTasksByTag("urgent")
	if err != nil {
		t.Fatalf("Failed to search tasks by tag: %v", err)
	}
	
	if len(tagResults) != 1 {
		t.Errorf("Expected 1 result for 'urgent' tag search, got %d", len(tagResults))
	}
	
	// Test category search
	categoryResults, err := searchService.SearchTasksByCategory("Work")
	if err != nil {
		t.Fatalf("Failed to search tasks by category: %v", err)
	}
	
	if len(categoryResults) < 2 {
		t.Errorf("Expected at least 2 results for 'Work' category search, got %d", len(categoryResults))
	}
	
	// Test comprehensive search
	searchQuery := SearchQuery{
		Query:      "Go",
		UserID:     &user1.ID,
		Status:     func() *int { s := 0; return &s }(),
		Priority:   func() *int { p := 3; return &p }(),
		CategoryID: &category1.ID,
		Limit:      10,
		Offset:     0,
		SortBy:     "title",
		SortOrder:  "asc",
	}
	
	comprehensiveResults, err := searchService.SearchTasks(searchQuery)
	if err != nil {
		t.Fatalf("Failed to perform comprehensive search: %v", err)
	}
	
	if comprehensiveResults.Total < 1 {
		t.Errorf("Expected at least 1 result for comprehensive search, got %d", comprehensiveResults.Total)
	}
	
	// Test filter options
	filterOptions, err := searchService.GetFilterOptions()
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
}

func TestSearchQueryValidation(t *testing.T) {
	// Test default query
	query := DefaultSearchQuery()
	if query.Limit != 20 {
		t.Errorf("Expected default limit 20, got %d", query.Limit)
	}
	
	if query.SortBy != "created_at" {
		t.Errorf("Expected default sort by 'created_at', got '%s'", query.SortBy)
	}
	
	if query.SortOrder != "desc" {
		t.Errorf("Expected default sort order 'desc', got '%s'", query.SortOrder)
	}
	
	// Test validation
	err := query.Validate()
	if err != nil {
		t.Fatalf("Expected validation to pass, got error: %v", err)
	}
	
	// Test invalid sort field
	query.SortBy = "invalid_field"
	err = query.Validate()
	if err != nil {
		t.Fatalf("Expected validation to handle invalid sort field, got error: %v", err)
	}
	
	if query.SortBy != "created_at" {
		t.Errorf("Expected invalid sort field to be reset to 'created_at', got '%s'", query.SortBy)
	}
	
	// Test invalid sort order
	query.SortOrder = "invalid_order"
	err = query.Validate()
	if err != nil {
		t.Fatalf("Expected validation to handle invalid sort order, got error: %v", err)
	}
	
	if query.SortOrder != "desc" {
		t.Errorf("Expected invalid sort order to be reset to 'desc', got '%s'", query.SortOrder)
	}
}

func TestSearchResultPagination(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_pagination.db"
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
	
	// Create repository and search service
	repository := database.NewSQLiteRepository(db)
	searchService := NewSearchService(repository)
	
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
	
	// Create multiple tasks
	now := time.Now()
	for i := 0; i < 25; i++ {
		task := &database.DatabaseTask{
			Title:       fmt.Sprintf("Task %d", i+1),
			Description: fmt.Sprintf("Description for task %d", i+1),
			Priority:    (i%4) + 1,
			Status:      i % 4,
			CreatedAt:   now.Add(-time.Duration(i) * time.Hour),
			UpdatedAt:   now.Add(-time.Duration(i) * time.Hour),
			UserID:      &user.ID,
			IsArchived:  false,
		}
		err = repository.CreateTask(task)
		if err != nil {
			t.Fatalf("Failed to create task %d: %v", i+1, err)
		}
	}
	
	// Test pagination
	query := SearchQuery{
		Limit:  10,
		Offset: 0,
	}
	
	result, err := searchService.SearchTasks(query)
	if err != nil {
		t.Fatalf("Failed to search tasks: %v", err)
	}
	
	if len(result.Tasks) != 10 {
		t.Errorf("Expected 10 tasks on first page, got %d", len(result.Tasks))
	}
	
	if result.Total != 25 {
		t.Errorf("Expected total 25 tasks, got %d", result.Total)
	}
	
	if result.Page != 1 {
		t.Errorf("Expected page 1, got %d", result.Page)
	}
	
	if result.TotalPages != 3 {
		t.Errorf("Expected 3 total pages, got %d", result.TotalPages)
	}
	
	// Test second page
	query.Offset = 10
	result, err = searchService.SearchTasks(query)
	if err != nil {
		t.Fatalf("Failed to search tasks on second page: %v", err)
	}
	
	if len(result.Tasks) != 10 {
		t.Errorf("Expected 10 tasks on second page, got %d", len(result.Tasks))
	}
	
	if result.Page != 2 {
		t.Errorf("Expected page 2, got %d", result.Page)
	}
	
	// Test third page
	query.Offset = 20
	result, err = searchService.SearchTasks(query)
	if err != nil {
		t.Fatalf("Failed to search tasks on third page: %v", err)
	}
	
	if len(result.Tasks) != 5 {
		t.Errorf("Expected 5 tasks on third page, got %d", len(result.Tasks))
	}
	
	if result.Page != 3 {
		t.Errorf("Expected page 3, got %d", result.Page)
	}
}
