package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"learn-go-capstone/internal/auth"
	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/notifications"
	"learn-go-capstone/internal/task"
)

// setupTestDB creates a test database and returns the connection, repository, and cleanup function
func setupTestDB(t *testing.T) (*sql.DB, database.Repository, func()) {
	dbPath := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	cfg := &database.Config{
		Driver: "sqlite3",
		DSN:    dbPath,
	}
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	repo := database.NewSQLiteRepository(db)
	migrationManager := database.NewMigrationManager(db)
	err = migrationManager.Migrate()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db, repo, func() {
		database.Close(db)
	}
}

func setupTestAPI(t *testing.T) (*httptest.Server, *Server, func()) {
	// Setup test database
	_, repository, cleanup := setupTestDB(t)

	// Create managers
	taskManager := task.NewTaskManager()
	userManager := task.NewUserManager(repository)
	categoryManager := task.NewCategoryManager(repository)
	dependencyManager := task.NewDependencyManager(repository)
	searchManager := task.NewSearchManager(repository)
	exportManager := task.NewExportManager(repository)

	// Create auth service
	authService := auth.NewAuthService(repository)

	// Create notification service
	notifConfig := notifications.NotificationConfig{}
	notifService := notifications.NewNotificationService(repository, notifConfig)

	// Create notification manager
	notificationManager := task.NewNotificationManager(repository, notifService)

	// Create server
	server := NewServer(
		taskManager,
		userManager,
		categoryManager,
		dependencyManager,
		searchManager,
		exportManager,
		notificationManager,
		authService,
	)

	// Create test server
	testServer := httptest.NewServer(server.GetRouter())

	return testServer, server, cleanup
}

func TestHealthCheck(t *testing.T) {
	server, _, cleanup := setupTestAPI(t)
	defer cleanup()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Should make health check request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Health check should return 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Should decode health check response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Fatal("Health check should be successful")
	}
}

func TestUserRegistration(t *testing.T) {
	server, _, cleanup := setupTestAPI(t)
	defer cleanup()

	userData := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Should marshal user data: %v", err)
	}

	resp, err := http.Post(server.URL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Should make registration request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Registration should return 201, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Should decode registration response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Fatal("Registration should be successful")
	}
}

func TestUnauthorizedAccess(t *testing.T) {
	server, _, cleanup := setupTestAPI(t)
	defer cleanup()

	resp, err := http.Get(server.URL + "/api/v1/tasks")
	if err != nil {
		t.Fatalf("Should create request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Should return 401 for unauthorized access, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Should decode error response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Fatal("Response should indicate failure")
	}
}
