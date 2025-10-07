package task

import (
	"fmt"
	"time"

	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/notifications"
)

// NotificationManager manages notifications for tasks
type NotificationManager struct {
	repository         database.Repository
	notificationService *notifications.NotificationService
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(repository database.Repository, notificationService *notifications.NotificationService) *NotificationManager {
	return &NotificationManager{
		repository:         repository,
		notificationService: notificationService,
	}
}

// CreateTaskReminder creates a reminder for a task
func (nm *NotificationManager) CreateTaskReminder(userID, taskID int, reminderMinutes int) error {
	// Get task details
	task, err := nm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	// Check if task has a due date
	if task.DueDate == nil {
		return fmt.Errorf("task does not have a due date")
	}
	
	// Create reminder
	return nm.notificationService.CreateTaskReminder(
		userID,
		taskID,
		task.Title,
		*task.DueDate,
		reminderMinutes,
	)
}

// CreateOverdueReminder creates an overdue reminder for a task
func (nm *NotificationManager) CreateOverdueReminder(userID, taskID int, overdueHours int) error {
	// Get task details
	task, err := nm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	// Create overdue reminder
	return nm.notificationService.CreateOverdueReminder(
		userID,
		taskID,
		task.Title,
		overdueHours,
	)
}

// CreateStatusChangeNotification creates a notification for task status changes
func (nm *NotificationManager) CreateStatusChangeNotification(userID, taskID int, oldStatus, newStatus Status) error {
	// Get task details
	task, err := nm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	// Create status change notification
	return nm.notificationService.CreateStatusChangeNotification(
		userID,
		taskID,
		task.Title,
		oldStatus.String(),
		newStatus.String(),
	)
}

// CreateTaskCreatedNotification creates a notification when a task is created
func (nm *NotificationManager) CreateTaskCreatedNotification(userID, taskID int) error {
	// Get task details
	task, err := nm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	notification := &notifications.Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        notifications.TypeInApp,
		Priority:    notifications.PriorityNormal,
		Trigger:     notifications.TriggerCreated,
		Title:       "New Task Created",
		Message:     fmt.Sprintf("Task '%s' has been created", task.Title),
		Recipient:   "",
		MaxRetries:  3,
	}
	
	return nm.notificationService.SendNotification(notification)
}

// CreateTaskUpdatedNotification creates a notification when a task is updated
func (nm *NotificationManager) CreateTaskUpdatedNotification(userID, taskID int) error {
	// Get task details
	task, err := nm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	notification := &notifications.Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        notifications.TypeInApp,
		Priority:    notifications.PriorityNormal,
		Trigger:     notifications.TriggerUpdated,
		Title:       "Task Updated",
		Message:     fmt.Sprintf("Task '%s' has been updated", task.Title),
		Recipient:   "",
		MaxRetries:  3,
	}
	
	return nm.notificationService.SendNotification(notification)
}

// CreateTaskDeletedNotification creates a notification when a task is deleted
func (nm *NotificationManager) CreateTaskDeletedNotification(userID, taskID int, taskTitle string) error {
	notification := &notifications.Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        notifications.TypeInApp,
		Priority:    notifications.PriorityNormal,
		Trigger:     notifications.TriggerCustom,
		Title:       "Task Deleted",
		Message:     fmt.Sprintf("Task '%s' has been deleted", taskTitle),
		Recipient:   "",
		MaxRetries:  3,
	}
	
	return nm.notificationService.SendNotification(notification)
}

// CheckOverdueTasks checks for overdue tasks and creates reminders
func (nm *NotificationManager) CheckOverdueTasks() error {
	// Get all overdue tasks
	overdueTasks, err := nm.repository.GetOverdueTasks()
	if err != nil {
		return err
	}
	
	// Create overdue reminders for each task
	for _, task := range overdueTasks {
		if task.UserID == nil {
			continue // Skip tasks without user ID
		}
		
		// Calculate how many hours overdue
		overdueHours := int(time.Since(*task.DueDate).Hours())
		
		// Create overdue reminder
		err = nm.CreateOverdueReminder(*task.UserID, task.ID, overdueHours)
		if err != nil {
			// Log error but continue with other tasks
			continue
		}
	}
	
	return nil
}

// CheckDueSoonTasks checks for tasks due soon and creates reminders
func (nm *NotificationManager) CheckDueSoonTasks(reminderMinutes int) error {
	// Get all tasks with due dates
	allTasks, err := nm.repository.GetAllTasks()
	if err != nil {
		return err
	}
	
	now := time.Now()
	reminderTime := now.Add(time.Duration(reminderMinutes) * time.Minute)
	
	// Find tasks due within the reminder time
	for _, task := range allTasks {
		if task.DueDate == nil || task.UserID == nil {
			continue
		}
		
		// Check if task is due within the reminder time
		if task.DueDate.After(now) && task.DueDate.Before(reminderTime) {
			// Create reminder
			err = nm.CreateTaskReminder(*task.UserID, task.ID, reminderMinutes)
			if err != nil {
				// Log error but continue with other tasks
				continue
			}
		}
	}
	
	return nil
}

// GetUserNotifications gets notifications for a user
func (nm *NotificationManager) GetUserNotifications(userID int, limit, offset int) ([]*notifications.Notification, error) {
	return nm.notificationService.GetUserNotifications(userID, limit, offset)
}

// MarkNotificationAsRead marks a notification as read
func (nm *NotificationManager) MarkNotificationAsRead(notificationID int) error {
	return nm.notificationService.MarkNotificationAsRead(notificationID)
}

// GetNotificationStats gets notification statistics
func (nm *NotificationManager) GetNotificationStats() (*notifications.NotificationStats, error) {
	return nm.notificationService.GetNotificationStats()
}

// GetQueueStatus gets the current queue status
func (nm *NotificationManager) GetQueueStatus() map[string]interface{} {
	return nm.notificationService.GetQueueStatus()
}

// CreateBulkTaskReminders creates reminders for multiple tasks
func (nm *NotificationManager) CreateBulkTaskReminders(userID int, taskIDs []int, reminderMinutes int) error {
	for _, taskID := range taskIDs {
		err := nm.CreateTaskReminder(userID, taskID, reminderMinutes)
		if err != nil {
			// Log error but continue with other tasks
			continue
		}
	}
	return nil
}

// CreatePriorityTaskReminder creates a high-priority reminder for a task
func (nm *NotificationManager) CreatePriorityTaskReminder(userID, taskID int, reminderMinutes int) error {
	// Get task details
	task, err := nm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	// Check if task has a due date
	if task.DueDate == nil {
		return fmt.Errorf("task does not have a due date")
	}
	
	// Create high-priority reminder
	reminderTime := task.DueDate.Add(-time.Duration(reminderMinutes) * time.Minute)
	
	notification := &notifications.Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        notifications.TypeEmail,
		Priority:    notifications.PriorityHigh,
		Trigger:     notifications.TriggerDueDate,
		Title:       fmt.Sprintf("URGENT: Task Reminder - %s", task.Title),
		Message:     fmt.Sprintf("Your high-priority task '%s' is due in %d minutes.", task.Title, reminderMinutes),
		Recipient:   "", // Would be populated from user email
		MaxRetries:  5,  // More retries for high-priority
	}
	
	return nm.notificationService.ScheduleNotification(notification, reminderTime)
}

// CreateCustomNotification creates a custom notification
func (nm *NotificationManager) CreateCustomNotification(userID, taskID int, title, message string, notificationType notifications.NotificationType, priority notifications.NotificationPriority) error {
	notification := &notifications.Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        notificationType,
		Priority:    priority,
		Trigger:     notifications.TriggerCustom,
		Title:       title,
		Message:     message,
		Recipient:   "",
		MaxRetries:  3,
	}
	
	return nm.notificationService.SendNotification(notification)
}

// ScheduleRecurringReminder schedules a recurring reminder for a task
func (nm *NotificationManager) ScheduleRecurringReminder(userID, taskID int, intervalMinutes int, maxReminders int) error {
	// Get task details
	task, err := nm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	// Check if task has a due date
	if task.DueDate == nil {
		return fmt.Errorf("task does not have a due date")
	}
	
	// Schedule multiple reminders
	for i := 0; i < maxReminders; i++ {
		reminderTime := task.DueDate.Add(-time.Duration((i+1)*intervalMinutes) * time.Minute)
		
		// Only schedule if the reminder time is in the future
		if reminderTime.After(time.Now()) {
			notification := &notifications.Notification{
				UserID:      userID,
				TaskID:      taskID,
				Type:        notifications.TypeEmail,
				Priority:    notifications.PriorityNormal,
				Trigger:     notifications.TriggerDueDate,
				Title:       fmt.Sprintf("Task Reminder %d/%d: %s", i+1, maxReminders, task.Title),
				Message:     fmt.Sprintf("Your task '%s' is due in %d minutes.", task.Title, (i+1)*intervalMinutes),
				Recipient:   "",
				MaxRetries:  3,
			}
			
			err = nm.notificationService.ScheduleNotification(notification, reminderTime)
			if err != nil {
				// Log error but continue with other reminders
				continue
			}
		}
	}
	
	return nil
}
