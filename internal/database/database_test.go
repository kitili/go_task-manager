package database

import (
	"os"
	"testing"
	"time"
)

func TestDatabaseOperations(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_tasks.db"
	defer os.Remove(tempDB) // Clean up after test
	
	// Create database configuration
	config := &Config{
		Driver: "sqlite3",
		DSN:    tempDB,
	}
	
	// Connect to database
	db, err := Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer Close(db)
	
	// Run migrations
	migrationManager := NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Create repository
	repo := NewSQLiteRepository(db)
	
	// Test creating a task
	now := time.Now()
	task := &DatabaseTask{
		Title:       "Test Task",
		Description: "This is a test task",
		Priority:    2,
		Status:      0,
		DueDate:     &now,
	}
	
	err = repo.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	
	if task.ID == 0 {
		t.Error("Expected task ID to be set after creation")
	}
	
	// Test retrieving the task
	retrievedTask, err := repo.GetTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve task: %v", err)
	}
	
	if retrievedTask.Title != task.Title {
		t.Errorf("Expected title %s, got %s", task.Title, retrievedTask.Title)
	}
	
	// Test updating the task
	retrievedTask.Status = 2 // Completed
	err = repo.UpdateTask(retrievedTask)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}
	
	// Test getting all tasks
	allTasks, err := repo.GetAllTasks()
	if err != nil {
		t.Fatalf("Failed to get all tasks: %v", err)
	}
	
	if len(allTasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(allTasks))
	}
	
	// Test getting tasks by status
	completedTasks, err := repo.GetTasksByStatus(2)
	if err != nil {
		t.Fatalf("Failed to get tasks by status: %v", err)
	}
	
	if len(completedTasks) != 1 {
		t.Errorf("Expected 1 completed task, got %d", len(completedTasks))
	}
	
	// Test deleting the task
	err = repo.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
	
	// Verify task is deleted
	_, err = repo.GetTask(task.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted task")
	}
}

func TestMigrationSystem(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_migrations.db"
	defer os.Remove(tempDB)
	
	// Create database configuration
	config := &Config{
		Driver: "sqlite3",
		DSN:    tempDB,
	}
	
	// Connect to database
	db, err := Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer Close(db)
	
	// Run migrations
	migrationManager := NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Check if migrations table exists and has records
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count migrations: %v", err)
	}
	
	if count == 0 {
		t.Error("Expected migrations to be recorded")
	}
	
	// Run migrations again (should be idempotent)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations again: %v", err)
	}
	
	// Check that we didn't create duplicate migration records
	var newCount int
	err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&newCount)
	if err != nil {
		t.Fatalf("Failed to count migrations after second run: %v", err)
	}
	
	if newCount != count {
		t.Error("Expected migration count to remain the same after second run")
	}
}
