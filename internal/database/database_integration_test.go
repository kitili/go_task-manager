package database

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

// setupTestDB creates a test database and returns the connection, repository, and cleanup function
func setupTestDB(t *testing.T) (*sql.DB, Repository, func()) {
	dbPath := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	cfg := &Config{
		Driver: "sqlite3",
		DSN:    dbPath,
	}
	db, err := Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	repo := NewSQLiteRepository(db)
	migrationManager := NewMigrationManager(db)
	err = migrationManager.Migrate()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db, repo, func() {
		Close(db)
	}
}

func TestDatabaseConnection(t *testing.T) {
	db, repository, cleanup := setupTestDB(t)
	defer cleanup()

	if db == nil {
		t.Fatal("Database connection should not be nil")
	}
	if repository == nil {
		t.Fatal("Repository should not be nil")
	}
}

func TestUserCRUD(t *testing.T) {
	_, repository, cleanup := setupTestDB(t)
	defer cleanup()

	// Test CreateUser
	user := &User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repository.CreateUser(user)
	if err != nil {
		t.Fatalf("Should create user successfully: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("User ID should be set")
	}

	// Test GetUser
	retrievedUser, err := repository.GetUser(user.ID)
	if err != nil {
		t.Fatalf("Should retrieve user successfully: %v", err)
	}
	if user.Username != retrievedUser.Username {
		t.Fatalf("Username should match: expected %s, got %s", user.Username, retrievedUser.Username)
	}

	// Test UpdateUser
	user.Username = "updateduser"
	user.UpdatedAt = time.Now()
	err = repository.UpdateUser(user)
	if err != nil {
		t.Fatalf("Should update user successfully: %v", err)
	}

	// Test DeleteUser
	err = repository.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("Should delete user successfully: %v", err)
	}
}

func TestTaskCRUD(t *testing.T) {
	_, repository, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	user := &User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repository.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test CreateTask
	task := &DatabaseTask{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    1,
		Status:      0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      &user.ID,
		IsArchived:  false,
	}

	err = repository.CreateTask(task)
	if err != nil {
		t.Fatalf("Should create task successfully: %v", err)
	}
	if task.ID == 0 {
		t.Fatal("Task ID should be set")
	}

	// Test GetTask
	retrievedTask, err := repository.GetTask(task.ID)
	if err != nil {
		t.Fatalf("Should retrieve task successfully: %v", err)
	}
	if task.Title != retrievedTask.Title {
		t.Fatalf("Task title should match: expected %s, got %s", task.Title, retrievedTask.Title)
	}

	// Test UpdateTask
	task.Title = "Updated Task"
	task.UpdatedAt = time.Now()
	err = repository.UpdateTask(task)
	if err != nil {
		t.Fatalf("Should update task successfully: %v", err)
	}

	// Test DeleteTask
	err = repository.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("Should delete task successfully: %v", err)
	}
}
