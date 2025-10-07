package task

import (
	"os"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
)

func TestUserManager(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_user_manager.db"
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
	
	// Create repository and user manager
	repository := database.NewSQLiteRepository(db)
	userManager := NewUserManager(repository)
	
	// Test user registration
	user, err := userManager.RegisterUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	
	if user.ID == 0 {
		t.Error("Expected user ID to be set after registration")
	}
	
	// Test user login
	token, _, err := userManager.LoginUser("testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	
	if token == "" {
		t.Error("Expected token to be generated")
	}
	
	// Test getting user from token
	userFromToken, err := userManager.GetUserFromToken(token)
	if err != nil {
		t.Fatalf("Failed to get user from token: %v", err)
	}
	
	if userFromToken.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, userFromToken.ID)
	}
	
	// Test creating a task for the user
	dueDate := time.Now().Add(24 * time.Hour)
	task, err := userManager.CreateUserTask(user.ID, "Test Task", "Test Description", High, &dueDate)
	if err != nil {
		t.Fatalf("Failed to create user task: %v", err)
	}
	
	if task.Title != "Test Task" {
		t.Errorf("Expected task title 'Test Task', got '%s'", task.Title)
	}
	
	if task.Priority != High {
		t.Errorf("Expected task priority High, got %v", task.Priority)
	}
	
	// Test getting user tasks
	userTasks, err := userManager.GetUserTasks(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user tasks: %v", err)
	}
	
	if len(userTasks) != 1 {
		t.Errorf("Expected 1 user task, got %d", len(userTasks))
	}
	
	if userTasks[0].Title != "Test Task" {
		t.Errorf("Expected task title 'Test Task', got '%s'", userTasks[0].Title)
	}
	
	// Test updating user task
	err = userManager.UpdateUserTask(user.ID, task.ID, "Updated Task", "Updated Description", Medium, InProgress, &dueDate)
	if err != nil {
		t.Fatalf("Failed to update user task: %v", err)
	}
	
	// Verify task was updated
	updatedTasks, err := userManager.GetUserTasks(user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user tasks: %v", err)
	}
	
	if len(updatedTasks) != 1 {
		t.Errorf("Expected 1 updated user task, got %d", len(updatedTasks))
	}
	
	if updatedTasks[0].Title != "Updated Task" {
		t.Errorf("Expected updated task title 'Updated Task', got '%s'", updatedTasks[0].Title)
	}
	
	if updatedTasks[0].Status != InProgress {
		t.Errorf("Expected task status InProgress, got %v", updatedTasks[0].Status)
	}
	
	// Test getting tasks by status
	pendingTasks, err := userManager.GetUserTasksByStatus(user.ID, Pending)
	if err != nil {
		t.Fatalf("Failed to get pending tasks: %v", err)
	}
	
	if len(pendingTasks) != 0 {
		t.Errorf("Expected 0 pending tasks, got %d", len(pendingTasks))
	}
	
	inProgressTasks, err := userManager.GetUserTasksByStatus(user.ID, InProgress)
	if err != nil {
		t.Fatalf("Failed to get in-progress tasks: %v", err)
	}
	
	if len(inProgressTasks) != 1 {
		t.Errorf("Expected 1 in-progress task, got %d", len(inProgressTasks))
	}
	
	// Test getting tasks by priority
	highPriorityTasks, err := userManager.GetUserTasksByPriority(user.ID, High)
	if err != nil {
		t.Fatalf("Failed to get high priority tasks: %v", err)
	}
	
	if len(highPriorityTasks) != 0 {
		t.Errorf("Expected 0 high priority tasks, got %d", len(highPriorityTasks))
	}
	
	mediumPriorityTasks, err := userManager.GetUserTasksByPriority(user.ID, Medium)
	if err != nil {
		t.Fatalf("Failed to get medium priority tasks: %v", err)
	}
	
	if len(mediumPriorityTasks) != 1 {
		t.Errorf("Expected 1 medium priority task, got %d", len(mediumPriorityTasks))
	}
	
	// Test deleting user task
	err = userManager.DeleteUserTask(user.ID, task.ID)
	if err != nil {
		t.Fatalf("Failed to delete user task: %v", err)
	}
	
	// Verify task was deleted
	deletedTasks, err := userManager.GetUserTasks(user.ID)
	if err != nil {
		t.Fatalf("Failed to get tasks after deletion: %v", err)
	}
	
	if len(deletedTasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(deletedTasks))
	}
}

func TestUserManagerAccessControl(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_access_control.db"
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
	
	// Create repository and user manager
	repository := database.NewSQLiteRepository(db)
	userManager := NewUserManager(repository)
	
	// Create two users
	user1, err := userManager.RegisterUser("user1", "user1@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user1: %v", err)
	}
	
	user2, err := userManager.RegisterUser("user2", "user2@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user2: %v", err)
	}
	
	// Create a task for user1
	task, err := userManager.CreateUserTask(user1.ID, "User1 Task", "User1 Description", High, nil)
	if err != nil {
		t.Fatalf("Failed to create task for user1: %v", err)
	}
	
	// Test that user2 cannot update user1's task
	err = userManager.UpdateUserTask(user2.ID, task.ID, "Hacked Task", "Hacked Description", Low, Completed, nil)
	if err == nil {
		t.Error("Expected error when user2 tries to update user1's task")
	}
	
	// Test that user2 cannot delete user1's task
	err = userManager.DeleteUserTask(user2.ID, task.ID)
	if err == nil {
		t.Error("Expected error when user2 tries to delete user1's task")
	}
	
	// Test that user1 can still access their own task
	user1Tasks, err := userManager.GetUserTasks(user1.ID)
	if err != nil {
		t.Fatalf("Failed to get user1 tasks: %v", err)
	}
	
	if len(user1Tasks) != 1 {
		t.Errorf("Expected 1 task for user1, got %d", len(user1Tasks))
	}
	
	// Test that user2 cannot see user1's task
	user2Tasks, err := userManager.GetUserTasks(user2.ID)
	if err != nil {
		t.Fatalf("Failed to get user2 tasks: %v", err)
	}
	
	if len(user2Tasks) != 0 {
		t.Errorf("Expected 0 tasks for user2, got %d", len(user2Tasks))
	}
}

func TestUserManagerOverdueTasks(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_overdue_tasks.db"
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
	
	// Create repository and user manager
	repository := database.NewSQLiteRepository(db)
	userManager := NewUserManager(repository)
	
	// Create a user
	user, err := userManager.RegisterUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	
	// Create an overdue task (due yesterday)
	yesterday := time.Now().Add(-24 * time.Hour)
	_, err = userManager.CreateUserTask(user.ID, "Overdue Task", "This task is overdue", High, &yesterday)
	if err != nil {
		t.Fatalf("Failed to create overdue task: %v", err)
	}
	
	// Create a future task (due tomorrow)
	tomorrow := time.Now().Add(24 * time.Hour)
	futureTask, err := userManager.CreateUserTask(user.ID, "Future Task", "This task is due tomorrow", Medium, &tomorrow)
	if err != nil {
		t.Fatalf("Failed to create future task: %v", err)
	}
	
	// Test getting overdue tasks
	overdueTasks, err := userManager.GetUserOverdueTasks(user.ID)
	if err != nil {
		t.Fatalf("Failed to get overdue tasks: %v", err)
	}
	
	if len(overdueTasks) != 1 {
		t.Errorf("Expected 1 overdue task, got %d", len(overdueTasks))
	}
	
	if overdueTasks[0].Title != "Overdue Task" {
		t.Errorf("Expected overdue task title 'Overdue Task', got '%s'", overdueTasks[0].Title)
	}
	
	// Verify future task is not in overdue list
	for _, task := range overdueTasks {
		if task.ID == futureTask.ID {
			t.Error("Future task should not be in overdue tasks list")
		}
	}
}
