package task

import (
	"os"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
)

func TestCategoryManager(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_categories.db"
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
	
	// Create repository and category manager
	repository := database.NewSQLiteRepository(db)
	categoryManager := NewCategoryManager(repository)
	
	// Test creating a category
	category, err := categoryManager.CreateCategory("Work", "Work-related tasks", "#ff0000")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	
	if category.ID == 0 {
		t.Error("Expected category ID to be set after creation")
	}
	
	if category.Name != "Work" {
		t.Errorf("Expected category name 'Work', got '%s'", category.Name)
	}
	
	// Test retrieving the category
	retrievedCategory, err := categoryManager.GetCategory(category.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve category: %v", err)
	}
	
	if retrievedCategory.Name != category.Name {
		t.Errorf("Expected category name %s, got %s", category.Name, retrievedCategory.Name)
	}
	
	// Test creating a tag
	tag, err := categoryManager.CreateTag("urgent", "#ff0000")
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}
	
	if tag.ID == 0 {
		t.Error("Expected tag ID to be set after creation")
	}
	
	if tag.Name != "urgent" {
		t.Errorf("Expected tag name 'urgent', got '%s'", tag.Name)
	}
	
	// Test getting all categories
	categories, err := categoryManager.GetAllCategories()
	if err != nil {
		t.Fatalf("Failed to get all categories: %v", err)
	}
	
	if len(categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(categories))
	}
	
	// Test getting all tags
	tags, err := categoryManager.GetAllTags()
	if err != nil {
		t.Fatalf("Failed to get all tags: %v", err)
	}
	
	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}
	
	// Test updating category
	category.Description = "Updated work tasks"
	err = categoryManager.UpdateCategory(category)
	if err != nil {
		t.Fatalf("Failed to update category: %v", err)
	}
	
	// Test updating tag
	tag.Color = "#00ff00"
	err = categoryManager.UpdateTag(tag)
	if err != nil {
		t.Fatalf("Failed to update tag: %v", err)
	}
	
	// Test deleting category
	err = categoryManager.DeleteCategory(category.ID)
	if err != nil {
		t.Fatalf("Failed to delete category: %v", err)
	}
	
	// Test deleting tag
	err = categoryManager.DeleteTag(tag.ID)
	if err != nil {
		t.Fatalf("Failed to delete tag: %v", err)
	}
}

func TestTaskWithCategoryAndTags(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_task_categories.db"
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
	
	// Create repository and category manager
	repository := database.NewSQLiteRepository(db)
	categoryManager := NewCategoryManager(repository)
	
	// Create a category
	category, err := categoryManager.CreateCategory("Personal", "Personal tasks", "#00ff00")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	
	// Create a tag
	tag, err := categoryManager.CreateTag("important", "#ff0000")
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}
	
	// Create a task with category
	now := time.Now()
	task := Task{
		Title:       "Test Task with Category",
		Description: "Testing task with category",
		Priority:    High,
		Status:      Pending,
		CreatedAt:   now,
		UpdatedAt:   now,
		Category:    &Category{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			Color:       category.Color,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		},
		Tags: []Tag{},
	}
	
	// Test converting to database task
	dbTask := convertToDatabaseTask(task)
	if dbTask.CategoryID == nil {
		t.Error("Expected category ID to be set")
	}
	
	if *dbTask.CategoryID != category.ID {
		t.Errorf("Expected category ID %d, got %d", category.ID, *dbTask.CategoryID)
	}
	
	// Test converting back from database task
	convertedTask := convertFromDatabaseTask(dbTask)
	// Note: Category is not loaded in the converter to avoid circular dependencies
	// It will be loaded separately when needed
	if convertedTask.Title != task.Title {
		t.Errorf("Expected task title %s, got %s", task.Title, convertedTask.Title)
	}
	
	// Test task-tag relationship
	// First create a task in the database
	err = repository.CreateTask(dbTask)
	if err != nil {
		t.Fatalf("Failed to create task in database: %v", err)
	}
	
	// Add tag to task
	err = categoryManager.AddTagToTask(dbTask.ID, tag.ID)
	if err != nil {
		t.Fatalf("Failed to add tag to task: %v", err)
	}
	
	// Get task tags
	taskTags, err := categoryManager.GetTaskTags(dbTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task tags: %v", err)
	}
	
	if len(taskTags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(taskTags))
	}
	
	if taskTags[0].Name != tag.Name {
		t.Errorf("Expected tag name %s, got %s", tag.Name, taskTags[0].Name)
	}
	
	// Get tasks by tag
	tasksByTag, err := categoryManager.GetTasksByTag(tag.ID)
	if err != nil {
		t.Fatalf("Failed to get tasks by tag: %v", err)
	}
	
	if len(tasksByTag) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasksByTag))
	}
	
	if tasksByTag[0].Title != task.Title {
		t.Errorf("Expected task title %s, got %s", task.Title, tasksByTag[0].Title)
	}
	
	// Remove tag from task
	err = categoryManager.RemoveTagFromTask(dbTask.ID, tag.ID)
	if err != nil {
		t.Fatalf("Failed to remove tag from task: %v", err)
	}
	
	// Verify tag was removed
	taskTags, err = categoryManager.GetTaskTags(dbTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task tags after removal: %v", err)
	}
	
	if len(taskTags) != 0 {
		t.Errorf("Expected 0 tags after removal, got %d", len(taskTags))
	}
}
