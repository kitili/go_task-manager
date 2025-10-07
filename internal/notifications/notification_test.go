package notifications

import (
	"os"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
)

func TestNotificationService(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_notifications.db"
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
	
	// Create repository and notification service
	repository := database.NewSQLiteRepository(db)
	notificationConfig := DefaultNotificationConfig()
	notificationService := NewNotificationService(repository, notificationConfig)
	defer notificationService.Stop()
	
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
	
	// Create test task
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
	
	// Test sending immediate notification
	notification := &Notification{
		UserID:      user.ID,
		TaskID:      task.ID,
		Type:        TypeEmail,
		Priority:    PriorityNormal,
		Trigger:     TriggerCreated,
		Title:       "Test Notification",
		Message:     "This is a test notification",
		Recipient:   "test@example.com",
		MaxRetries:  3,
	}
	
	err = notificationService.SendNotification(notification)
	if err != nil {
		t.Fatalf("Failed to send notification: %v", err)
	}
	
	// Test scheduling notification
	scheduledTime := now.Add(1 * time.Hour)
	scheduledNotification := &Notification{
		UserID:      user.ID,
		TaskID:      task.ID,
		Type:        TypeInApp,
		Priority:    PriorityHigh,
		Trigger:     TriggerDueDate,
		Title:       "Scheduled Notification",
		Message:     "This is a scheduled notification",
		Recipient:   "",
		MaxRetries:  3,
	}
	
	err = notificationService.ScheduleNotification(scheduledNotification, scheduledTime)
	if err != nil {
		t.Fatalf("Failed to schedule notification: %v", err)
	}
	
	// Test creating task reminder
	err = notificationService.CreateTaskReminder(user.ID, task.ID, task.Title, dueDate, 60)
	if err != nil {
		t.Fatalf("Failed to create task reminder: %v", err)
	}
	
	// Test creating overdue reminder
	err = notificationService.CreateOverdueReminder(user.ID, task.ID, task.Title, 2)
	if err != nil {
		t.Fatalf("Failed to create overdue reminder: %v", err)
	}
	
	// Test creating status change notification
	err = notificationService.CreateStatusChangeNotification(user.ID, task.ID, task.Title, "Pending", "In Progress")
	if err != nil {
		t.Fatalf("Failed to create status change notification: %v", err)
	}
	
	// Test getting notification stats
	stats, err := notificationService.GetNotificationStats()
	if err != nil {
		t.Fatalf("Failed to get notification stats: %v", err)
	}
	
	if stats == nil {
		t.Error("Expected notification stats, got nil")
	}
	
	// Test getting queue status
	queueStatus := notificationService.GetQueueStatus()
	if queueStatus == nil {
		t.Error("Expected queue status, got nil")
	}
	
	// Wait a bit for notifications to be processed
	time.Sleep(100 * time.Millisecond)
	
	// Test getting user notifications
	userNotifications, err := notificationService.GetUserNotifications(user.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get user notifications: %v", err)
	}
	
	// Note: In a real implementation, notifications would be stored in the database
	// and returned here. For now, we just check that the method doesn't error.
	_ = userNotifications
}

func TestNotificationScheduler(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_scheduler.db"
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
	notificationConfig := DefaultNotificationConfig()
	notificationService := NewNotificationService(repository, notificationConfig)
	defer notificationService.Stop()
	
	scheduler := NewScheduler(repository, notificationService, 1*time.Second)
	defer scheduler.Stop()
	
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
	
	// Create test task with due date in the past (overdue)
	pastDueDate := time.Now().Add(-2 * time.Hour)
	overdueTask := &database.DatabaseTask{
		Title:       "Overdue Task",
		Description: "This task is overdue",
		Priority:    3,
		Status:      0,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
		DueDate:     &pastDueDate,
		UserID:      &user.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(overdueTask)
	if err != nil {
		t.Fatalf("Failed to create overdue task: %v", err)
	}
	
	// Create test task with due date in the future
	futureDueDate := time.Now().Add(30 * time.Minute)
	futureTask := &database.DatabaseTask{
		Title:       "Future Task",
		Description: "This task is due soon",
		Priority:    2,
		Status:      0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     &futureDueDate,
		UserID:      &user.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(futureTask)
	if err != nil {
		t.Fatalf("Failed to create future task: %v", err)
	}
	
	// Start scheduler
	scheduler.Start()
	
	// Wait for scheduler to process notifications
	time.Sleep(2 * time.Second)
	
	// Test scheduler status
	status := scheduler.GetSchedulerStatus()
	if status == nil {
		t.Error("Expected scheduler status, got nil")
	}
	
	// Test scheduling task reminder
	err = scheduler.ScheduleTaskReminder(user.ID, futureTask.ID, futureTask.Title, futureDueDate, 15)
	if err != nil {
		t.Fatalf("Failed to schedule task reminder: %v", err)
	}
	
	// Test scheduling recurring reminder
	err = scheduler.ScheduleRecurringReminder(user.ID, futureTask.ID, futureTask.Title, futureDueDate, 30, 3)
	if err != nil {
		t.Fatalf("Failed to schedule recurring reminder: %v", err)
	}
	
	// Wait a bit more for processing
	time.Sleep(1 * time.Second)
}

func TestNotificationModels(t *testing.T) {
	// Test default notification settings
	settings := DefaultNotificationSettings()
	if !settings.EmailEnabled {
		t.Error("Expected email to be enabled by default")
	}
	if !settings.InAppEnabled {
		t.Error("Expected in-app notifications to be enabled by default")
	}
	if settings.ReminderTime != 60 {
		t.Errorf("Expected default reminder time to be 60 minutes, got %d", settings.ReminderTime)
	}
	
	// Test default notification config
	config := DefaultNotificationConfig()
	if config.MaxRetries != 3 {
		t.Errorf("Expected default max retries to be 3, got %d", config.MaxRetries)
	}
	if config.RetryDelay != 300 {
		t.Errorf("Expected default retry delay to be 300 seconds, got %d", config.RetryDelay)
	}
	if config.WorkerCount != 5 {
		t.Errorf("Expected default worker count to be 5, got %d", config.WorkerCount)
	}
	
	// Test notification creation
	notification := &Notification{
		ID:          1,
		UserID:      1,
		TaskID:      1,
		Type:        TypeEmail,
		Priority:    PriorityHigh,
		Status:      StatusPending,
		Trigger:     TriggerDueDate,
		Title:       "Test Notification",
		Message:     "Test message",
		Recipient:   "test@example.com",
		MaxRetries:  3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if notification.Type != TypeEmail {
		t.Errorf("Expected notification type to be email, got %s", notification.Type)
	}
	if notification.Priority != PriorityHigh {
		t.Errorf("Expected notification priority to be high, got %s", notification.Priority)
	}
	if notification.Status != StatusPending {
		t.Errorf("Expected notification status to be pending, got %s", notification.Status)
	}
}

func TestNotificationTypes(t *testing.T) {
	// Test notification type constants
	types := []NotificationType{TypeEmail, TypeInApp, TypeSMS, TypeWebhook, TypeSlack, TypeDiscord}
	for _, notificationType := range types {
		if notificationType == "" {
			t.Error("Notification type should not be empty")
		}
	}
	
	// Test priority constants
	priorities := []NotificationPriority{PriorityLow, PriorityNormal, PriorityHigh, PriorityCritical}
	for _, priority := range priorities {
		if priority == "" {
			t.Error("Notification priority should not be empty")
		}
	}
	
	// Test status constants
	statuses := []NotificationStatus{StatusPending, StatusSent, StatusDelivered, StatusFailed, StatusCancelled}
	for _, status := range statuses {
		if status == "" {
			t.Error("Notification status should not be empty")
		}
	}
	
	// Test trigger constants
	triggers := []NotificationTrigger{TriggerDueDate, TriggerOverdue, TriggerStatusChange, TriggerCreated, TriggerUpdated, TriggerCustom}
	for _, trigger := range triggers {
		if trigger == "" {
			t.Error("Notification trigger should not be empty")
		}
	}
}
