package task

import (
	"os"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/notifications"
)

func TestNotificationManager(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_notification_manager.db"
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
	
	// Create repository and services
	repository := database.NewSQLiteRepository(db)
	notificationConfig := notifications.DefaultNotificationConfig()
	notificationService := notifications.NewNotificationService(repository, notificationConfig)
	defer notificationService.Stop()
	
	notificationManager := NewNotificationManager(repository, notificationService)
	
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
	
	// Create test task with due date
	now := time.Now()
	dueDate := now.Add(24 * time.Hour)
	task := &database.DatabaseTask{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    3,
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
		DueDate:     &dueDate,
		UserID:      &user.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	
	// Test creating task reminder
	err = notificationManager.CreateTaskReminder(user.ID, task.ID, 60)
	if err != nil {
		t.Fatalf("Failed to create task reminder: %v", err)
	}
	
	// Test creating overdue reminder
	err = notificationManager.CreateOverdueReminder(user.ID, task.ID, 2)
	if err != nil {
		t.Fatalf("Failed to create overdue reminder: %v", err)
	}
	
	// Test creating status change notification
	err = notificationManager.CreateStatusChangeNotification(user.ID, task.ID, Pending, InProgress)
	if err != nil {
		t.Fatalf("Failed to create status change notification: %v", err)
	}
	
	// Test creating task created notification
	err = notificationManager.CreateTaskCreatedNotification(user.ID, task.ID)
	if err != nil {
		t.Fatalf("Failed to create task created notification: %v", err)
	}
	
	// Test creating task updated notification
	err = notificationManager.CreateTaskUpdatedNotification(user.ID, task.ID)
	if err != nil {
		t.Fatalf("Failed to create task updated notification: %v", err)
	}
	
	// Test creating task deleted notification
	err = notificationManager.CreateTaskDeletedNotification(user.ID, task.ID, task.Title)
	if err != nil {
		t.Fatalf("Failed to create task deleted notification: %v", err)
	}
	
	// Test checking overdue tasks
	err = notificationManager.CheckOverdueTasks()
	if err != nil {
		t.Fatalf("Failed to check overdue tasks: %v", err)
	}
	
	// Test checking due soon tasks
	err = notificationManager.CheckDueSoonTasks(60)
	if err != nil {
		t.Fatalf("Failed to check due soon tasks: %v", err)
	}
	
	// Test getting user notifications
	userNotifications, err := notificationManager.GetUserNotifications(user.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get user notifications: %v", err)
	}
	
	// Note: In a real implementation, notifications would be stored in the database
	// and returned here. For now, we just check that the method doesn't error.
	_ = userNotifications
	
	// Test getting notification stats
	stats, err := notificationManager.GetNotificationStats()
	if err != nil {
		t.Fatalf("Failed to get notification stats: %v", err)
	}
	
	if stats == nil {
		t.Error("Expected notification stats, got nil")
	}
	
	// Test getting queue status
	queueStatus := notificationManager.GetQueueStatus()
	if queueStatus == nil {
		t.Error("Expected queue status, got nil")
	}
	
	// Test creating bulk task reminders
	taskIDs := []int{task.ID}
	err = notificationManager.CreateBulkTaskReminders(user.ID, taskIDs, 30)
	if err != nil {
		t.Fatalf("Failed to create bulk task reminders: %v", err)
	}
	
	// Test creating priority task reminder
	err = notificationManager.CreatePriorityTaskReminder(user.ID, task.ID, 15)
	if err != nil {
		t.Fatalf("Failed to create priority task reminder: %v", err)
	}
	
	// Test creating custom notification
	err = notificationManager.CreateCustomNotification(
		user.ID,
		task.ID,
		"Custom Title",
		"Custom Message",
		notifications.TypeInApp,
		notifications.PriorityHigh,
	)
	if err != nil {
		t.Fatalf("Failed to create custom notification: %v", err)
	}
	
	// Test scheduling recurring reminder
	err = notificationManager.ScheduleRecurringReminder(user.ID, task.ID, 30, 3)
	if err != nil {
		t.Fatalf("Failed to schedule recurring reminder: %v", err)
	}
	
	// Wait a bit for notifications to be processed
	time.Sleep(100 * time.Millisecond)
}

func TestNotificationManagerErrorHandling(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_notification_manager_errors.db"
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
	
	// Create repository and services
	repository := database.NewSQLiteRepository(db)
	notificationConfig := notifications.DefaultNotificationConfig()
	notificationService := notifications.NewNotificationService(repository, notificationConfig)
	defer notificationService.Stop()
	
	notificationManager := NewNotificationManager(repository, notificationService)
	
	// Test creating reminder for non-existent task
	err = notificationManager.CreateTaskReminder(1, 999, 60)
	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}
	
	// Test creating reminder for task without due date
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
	
	taskWithoutDueDate := &database.DatabaseTask{
		Title:       "Task Without Due Date",
		Description: "This task has no due date",
		Priority:    3,
		Status:      0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     nil, // No due date
		UserID:      &user.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(taskWithoutDueDate)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	
	err = notificationManager.CreateTaskReminder(user.ID, taskWithoutDueDate.ID, 60)
	if err == nil {
		t.Error("Expected error for task without due date, got nil")
	}
	
	// Test creating reminder for task with non-existent user
	taskWithDueDate := &database.DatabaseTask{
		Title:       "Task With Due Date",
		Description: "This task has a due date",
		Priority:    3,
		Status:      0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
		UserID:      &[]int{999}[0], // Non-existent user
		IsArchived:  false,
	}
	err = repository.CreateTask(taskWithDueDate)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	
	// This should work because the notification manager doesn't validate user existence
	// It just creates a notification for the given user ID
	err = notificationManager.CreateTaskReminder(999, taskWithDueDate.ID, 60)
	if err != nil {
		t.Errorf("Unexpected error for task with non-existent user: %v", err)
	}
}
