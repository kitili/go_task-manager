package testing

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"learn-go-capstone/internal/database"
)

// TestDBConfig creates a test database configuration
func TestDBConfig(t *testing.T) *database.Config {
	return &database.Config{
		Driver: "sqlite3",
		DSN:    fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name()),
	}
}

// SetupTestDB creates a test database with migrations
func SetupTestDB(t *testing.T) (*sql.DB, database.Repository, func()) {
	cfg := TestDBConfig(t)
	
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	repository := database.NewSQLiteRepository(db)

	cleanup := func() {
		database.Close(db)
	}

	return db, repository, cleanup
}

// SetupTestDBWithFile creates a test database with a temporary file
func SetupTestDBWithFile(t *testing.T) (*sql.DB, database.Repository, func()) {
	// Create a temporary file for the database
	tmpFile, err := os.CreateTemp("", "test_task_manager_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &database.Config{
		Driver: "sqlite3",
		DSN:    tmpFile.Name(),
	}
	
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	repository := database.NewSQLiteRepository(db)

	cleanup := func() {
		database.Close(db)
		os.Remove(tmpFile.Name())
	}

	return db, repository, cleanup
}

// AssertNoError checks that err is nil, failing the test if not
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError checks that err is not nil, failing the test if it is
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", message)
	}
}

// AssertEqual checks that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertNotEqual checks that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected == actual {
		t.Fatalf("%s: expected %v to not equal %v", message, expected, actual)
	}
}

// AssertTrue checks that a condition is true
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Fatalf("%s: expected true, got false", message)
	}
}

// AssertFalse checks that a condition is false
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Fatalf("%s: expected false, got true", message)
	}
}

// AssertNil checks that a value is nil
func AssertNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value != nil {
		t.Fatalf("%s: expected nil, got %v", message, value)
	}
}

// AssertNotNil checks that a value is not nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Fatalf("%s: expected non-nil value, got nil", message)
	}
}

// AssertContains checks that a slice contains a specific element
func AssertContains(t *testing.T, slice interface{}, element interface{}, message string) {
	t.Helper()
	// This is a simplified version - in practice, you'd use reflection or generics
	// For now, we'll implement basic string slice checking
	if strSlice, ok := slice.([]string); ok {
		if strElement, ok := element.(string); ok {
			for _, item := range strSlice {
				if item == strElement {
					return
				}
			}
		}
	}
	t.Fatalf("%s: expected slice to contain %v", message, element)
}

// AssertNotContains checks that a slice does not contain a specific element
func AssertNotContains(t *testing.T, slice interface{}, element interface{}, message string) {
	t.Helper()
	// This is a simplified version - in practice, you'd use reflection or generics
	if strSlice, ok := slice.([]string); ok {
		if strElement, ok := element.(string); ok {
			for _, item := range strSlice {
				if item == strElement {
					t.Fatalf("%s: expected slice to not contain %v", message, element)
				}
			}
		}
	}
}

// AssertLength checks that a slice has the expected length
func AssertLength(t *testing.T, slice interface{}, expectedLength int, message string) {
	t.Helper()
	// This is a simplified version - in practice, you'd use reflection
	if strSlice, ok := slice.([]string); ok {
		actualLength := len(strSlice)
		if actualLength != expectedLength {
			t.Fatalf("%s: expected length %d, got %d", message, expectedLength, actualLength)
		}
	}
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, repository database.Repository) *database.User {
	user := &database.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
	}
	
	err := repository.CreateUser(user)
	AssertNoError(t, err, "Failed to create test user")
	
	return user
}

// CreateTestCategory creates a test category in the database
func CreateTestCategory(t *testing.T, repository database.Repository) *database.Category {
	category := &database.Category{
		Name:        "Test Category",
		Description: "A test category",
		Color:       "#ff0000",
	}
	
	err := repository.CreateCategory(category)
	AssertNoError(t, err, "Failed to create test category")
	
	return category
}

// CreateTestTag creates a test tag in the database
func CreateTestTag(t *testing.T, repository database.Repository) *database.Tag {
	tag := &database.Tag{
		Name:  "test-tag",
		Color: "#00ff00",
	}
	
	err := repository.CreateTag(tag)
	AssertNoError(t, err, "Failed to create test tag")
	
	return tag
}

// CreateTestTask creates a test task in the database
func CreateTestTask(t *testing.T, repository database.Repository, userID int) *database.DatabaseTask {
	task := &database.DatabaseTask{
		Title:       "Test Task",
		Description: "A test task",
		Priority:    3,
		Status:      0,
		UserID:      &userID,
		IsArchived:  false,
	}
	
	err := repository.CreateTask(task)
	AssertNoError(t, err, "Failed to create test task")
	
	return task
}
