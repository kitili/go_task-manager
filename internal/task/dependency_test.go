package task

import (
	"os"
	"testing"

	"learn-go-capstone/internal/database"
)

func TestDependencyManager(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_dependencies.db"
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
	
	// Create repository and dependency manager
	repository := database.NewSQLiteRepository(db)
	dependencyManager := NewDependencyManager(repository)
	
	// Create test tasks
	task1 := &database.DatabaseTask{
		Title:       "Task 1",
		Description: "First task",
		Priority:    1,
		Status:      0,
	}
	err = repository.CreateTask(task1)
	if err != nil {
		t.Fatalf("Failed to create task 1: %v", err)
	}
	
	task2 := &database.DatabaseTask{
		Title:       "Task 2",
		Description: "Second task",
		Priority:    2,
		Status:      0,
	}
	err = repository.CreateTask(task2)
	if err != nil {
		t.Fatalf("Failed to create task 2: %v", err)
	}
	
	task3 := &database.DatabaseTask{
		Title:       "Task 3",
		Description: "Third task",
		Priority:    3,
		Status:      0,
	}
	err = repository.CreateTask(task3)
	if err != nil {
		t.Fatalf("Failed to create task 3: %v", err)
	}
	
	// Test adding dependencies
	err = dependencyManager.AddDependency(task2.ID, task1.ID)
	if err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	
	err = dependencyManager.AddDependency(task3.ID, task2.ID)
	if err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	
	// Test circular dependency detection
	err = dependencyManager.AddDependency(task1.ID, task3.ID)
	if err == nil {
		t.Error("Expected circular dependency error, but got none")
	}
	
	// Test getting dependencies
	dependencies, err := dependencyManager.GetDependencies(task2.ID)
	if err != nil {
		t.Fatalf("Failed to get dependencies: %v", err)
	}
	
	if len(dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(dependencies))
	}
	
	if dependencies[0].DependsOnTaskID != task1.ID {
		t.Errorf("Expected dependency on task %d, got %d", task1.ID, dependencies[0].DependsOnTaskID)
	}
	
	// Test getting tasks that depend on a task
	dependentTasks, err := dependencyManager.GetTasksThatDependOn(task1.ID)
	if err != nil {
		t.Fatalf("Failed to get dependent tasks: %v", err)
	}
	
	if len(dependentTasks) != 1 {
		t.Errorf("Expected 1 dependent task, got %d", len(dependentTasks))
	}
	
	if dependentTasks[0].Title != "Task 2" {
		t.Errorf("Expected dependent task 'Task 2', got '%s'", dependentTasks[0].Title)
	}
	
	// Test getting tasks that a task depends on
	dependsOnTasks, err := dependencyManager.GetTasksThatTaskDependsOn(task3.ID)
	if err != nil {
		t.Fatalf("Failed to get tasks that task depends on: %v", err)
	}
	
	if len(dependsOnTasks) != 1 {
		t.Errorf("Expected 1 task that task depends on, got %d", len(dependsOnTasks))
	}
	
	if dependsOnTasks[0].Title != "Task 2" {
		t.Errorf("Expected task that depends on 'Task 2', got '%s'", dependsOnTasks[0].Title)
	}
	
	// Test dependency chain
	chain, err := dependencyManager.GetDependencyChain(task3.ID)
	if err != nil {
		t.Fatalf("Failed to get dependency chain: %v", err)
	}
	
	if len(chain) != 2 {
		t.Errorf("Expected dependency chain of length 2, got %d", len(chain))
	}
	
	// Test removing dependency
	err = dependencyManager.RemoveDependency(task2.ID, task1.ID)
	if err != nil {
		t.Fatalf("Failed to remove dependency: %v", err)
	}
	
	// Verify dependency was removed
	dependencies, err = dependencyManager.GetDependencies(task2.ID)
	if err != nil {
		t.Fatalf("Failed to get dependencies after removal: %v", err)
	}
	
	if len(dependencies) != 0 {
		t.Errorf("Expected 0 dependencies after removal, got %d", len(dependencies))
	}
}

func TestDependencyCompletion(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_dependency_completion.db"
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
	
	// Create repository and dependency manager
	repository := database.NewSQLiteRepository(db)
	dependencyManager := NewDependencyManager(repository)
	
	// Create test tasks
	task1 := &database.DatabaseTask{
		Title:       "Prerequisite Task",
		Description: "Task that must be completed first",
		Priority:    1,
		Status:      0, // Pending
	}
	err = repository.CreateTask(task1)
	if err != nil {
		t.Fatalf("Failed to create task 1: %v", err)
	}
	
	task2 := &database.DatabaseTask{
		Title:       "Dependent Task",
		Description: "Task that depends on task 1",
		Priority:    2,
		Status:      0, // Pending
	}
	err = repository.CreateTask(task2)
	if err != nil {
		t.Fatalf("Failed to create task 2: %v", err)
	}
	
	// Add dependency
	err = dependencyManager.AddDependency(task2.ID, task1.ID)
	if err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	
	// Test that task 2 cannot be completed yet
	canComplete, err := dependencyManager.CanCompleteTask(task2.ID)
	if err != nil {
		t.Fatalf("Failed to check if task can be completed: %v", err)
	}
	
	if canComplete {
		t.Error("Expected task 2 to not be completable yet")
	}
	
	// Test getting blocked tasks
	blockedTasks, err := dependencyManager.GetBlockedTasks()
	if err != nil {
		t.Fatalf("Failed to get blocked tasks: %v", err)
	}
	
	if len(blockedTasks) != 1 {
		t.Errorf("Expected 1 blocked task, got %d", len(blockedTasks))
	}
	
	if blockedTasks[0].Title != "Dependent Task" {
		t.Errorf("Expected blocked task 'Dependent Task', got '%s'", blockedTasks[0].Title)
	}
	
	// Test getting ready tasks
	readyTasks, err := dependencyManager.GetReadyTasks()
	if err != nil {
		t.Fatalf("Failed to get ready tasks: %v", err)
	}
	
	if len(readyTasks) != 1 {
		t.Errorf("Expected 1 ready task, got %d", len(readyTasks))
	}
	
	if readyTasks[0].Title != "Prerequisite Task" {
		t.Errorf("Expected ready task 'Prerequisite Task', got '%s'", readyTasks[0].Title)
	}
	
	// Complete task 1
	task1.Status = 2 // Completed
	err = repository.UpdateTask(task1)
	if err != nil {
		t.Fatalf("Failed to complete task 1: %v", err)
	}
	
	// Now task 2 should be completable
	canComplete, err = dependencyManager.CanCompleteTask(task2.ID)
	if err != nil {
		t.Fatalf("Failed to check if task can be completed: %v", err)
	}
	
	if !canComplete {
		t.Error("Expected task 2 to be completable now")
	}
	
	// Test getting ready tasks again
	readyTasks, err = dependencyManager.GetReadyTasks()
	if err != nil {
		t.Fatalf("Failed to get ready tasks: %v", err)
	}
	
	if len(readyTasks) != 1 {
		t.Errorf("Expected 1 ready task (task2), got %d", len(readyTasks))
	}
	
	if readyTasks[0].Title != "Dependent Task" {
		t.Errorf("Expected ready task 'Dependent Task', got '%s'", readyTasks[0].Title)
	}
}

func TestCircularDependencyDetection(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_circular_deps.db"
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
	
	// Create repository and dependency manager
	repository := database.NewSQLiteRepository(db)
	dependencyManager := NewDependencyManager(repository)
	
	// Create test tasks
	task1 := &database.DatabaseTask{
		Title:       "Task 1",
		Description: "First task",
		Priority:    1,
		Status:      0,
	}
	err = repository.CreateTask(task1)
	if err != nil {
		t.Fatalf("Failed to create task 1: %v", err)
	}
	
	task2 := &database.DatabaseTask{
		Title:       "Task 2",
		Description: "Second task",
		Priority:    2,
		Status:      0,
	}
	err = repository.CreateTask(task2)
	if err != nil {
		t.Fatalf("Failed to create task 2: %v", err)
	}
	
	// Test direct circular dependency (task depends on itself)
	err = dependencyManager.AddDependency(task1.ID, task1.ID)
	if err == nil {
		t.Error("Expected error for self-dependency, but got none")
	}
	
	// Test indirect circular dependency
	err = dependencyManager.AddDependency(task1.ID, task2.ID)
	if err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	
	err = dependencyManager.AddDependency(task2.ID, task1.ID)
	if err == nil {
		t.Error("Expected error for circular dependency, but got none")
	}
}
