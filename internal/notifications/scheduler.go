package notifications

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"learn-go-capstone/internal/database"
)

// Scheduler handles scheduled notification processing
type Scheduler struct {
	repository         database.Repository
	notificationService *NotificationService
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
	ticker             *time.Ticker
	interval           time.Duration
}

// NewScheduler creates a new notification scheduler
func NewScheduler(repository database.Repository, notificationService *NotificationService, interval time.Duration) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Scheduler{
		repository:         repository,
		notificationService: notificationService,
		ctx:                ctx,
		cancel:             cancel,
		interval:           interval,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	log.Println("Starting notification scheduler...")
	
	s.ticker = time.NewTicker(s.interval)
	s.wg.Add(1)
	
	go func() {
		defer s.wg.Done()
		defer s.ticker.Stop()
		
		for {
			select {
			case <-s.ticker.C:
				s.processScheduledNotifications()
			case <-s.ctx.Done():
				log.Println("Notification scheduler stopped")
				return
			}
		}
	}()
	
	log.Printf("Notification scheduler started with interval: %v", s.interval)
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("Stopping notification scheduler...")
	s.cancel()
	s.wg.Wait()
	log.Println("Notification scheduler stopped")
}

// processScheduledNotifications processes all scheduled notifications
func (s *Scheduler) processScheduledNotifications() {
	log.Println("Processing scheduled notifications...")
	
	// Check for overdue tasks
	s.checkOverdueTasks()
	
	// Check for tasks due soon
	s.checkDueSoonTasks()
	
	// Check for status change notifications
	s.checkStatusChangeNotifications()
	
	// Process any other scheduled notifications
	s.processCustomScheduledNotifications()
}

// checkOverdueTasks checks for overdue tasks and creates reminders
func (s *Scheduler) checkOverdueTasks() {
	// Get all overdue tasks
	overdueTasks, err := s.repository.GetOverdueTasks()
	if err != nil {
		log.Printf("Error getting overdue tasks: %v", err)
		return
	}
	
	log.Printf("Found %d overdue tasks", len(overdueTasks))
	
	// Create overdue reminders for each task
	for _, task := range overdueTasks {
		if task.UserID == nil {
			continue // Skip tasks without user ID
		}
		
		// Calculate how many hours overdue
		overdueHours := int(time.Since(*task.DueDate).Hours())
		
		// Create overdue reminder
		notification := &Notification{
			UserID:      *task.UserID,
			TaskID:      task.ID,
			Type:        TypeEmail,
			Priority:    PriorityHigh,
			Trigger:     TriggerOverdue,
			Title:       fmt.Sprintf("Overdue Task: %s", task.Title),
			Message:     fmt.Sprintf("Your task '%s' is %d hours overdue.", task.Title, overdueHours),
			Recipient:   "", // Would be populated from user email
			MaxRetries:  3,
		}
		
		err = s.notificationService.SendNotification(notification)
		if err != nil {
			log.Printf("Error sending overdue reminder for task %d: %v", task.ID, err)
		}
	}
}

// checkDueSoonTasks checks for tasks due soon and creates reminders
func (s *Scheduler) checkDueSoonTasks() {
	// Get all tasks with due dates
	allTasks, err := s.repository.GetAllTasks()
	if err != nil {
		log.Printf("Error getting all tasks: %v", err)
		return
	}
	
	now := time.Now()
	reminderTime := now.Add(60 * time.Minute) // 1 hour from now
	
	// Find tasks due within the next hour
	dueSoonCount := 0
	for _, task := range allTasks {
		if task.DueDate == nil || task.UserID == nil {
			continue
		}
		
		// Check if task is due within the next hour
		if task.DueDate.After(now) && task.DueDate.Before(reminderTime) {
			dueSoonCount++
			
			// Calculate minutes until due
			minutesUntilDue := int(task.DueDate.Sub(now).Minutes())
			
			// Create reminder
			notification := &Notification{
				UserID:      *task.UserID,
				TaskID:      task.ID,
				Type:        TypeEmail,
				Priority:    PriorityNormal,
				Trigger:     TriggerDueDate,
				Title:       fmt.Sprintf("Task Reminder: %s", task.Title),
				Message:     fmt.Sprintf("Your task '%s' is due in %d minutes.", task.Title, minutesUntilDue),
				Recipient:   "", // Would be populated from user email
				MaxRetries:  3,
			}
			
			err = s.notificationService.SendNotification(notification)
			if err != nil {
				log.Printf("Error sending due soon reminder for task %d: %v", task.ID, err)
			}
		}
	}
	
	if dueSoonCount > 0 {
		log.Printf("Sent %d due soon reminders", dueSoonCount)
	}
}

// checkStatusChangeNotifications checks for status changes that need notifications
func (s *Scheduler) checkStatusChangeNotifications() {
	// This would typically check for recent status changes
	// and send notifications if the user has that setting enabled
	// For now, this is a placeholder
	log.Println("Checking status change notifications...")
}

// processCustomScheduledNotifications processes any custom scheduled notifications
func (s *Scheduler) processCustomScheduledNotifications() {
	// This would typically query a database for scheduled notifications
	// and process them if their scheduled time has arrived
	// For now, this is a placeholder
	log.Println("Processing custom scheduled notifications...")
}

// ScheduleTaskReminder schedules a reminder for a specific task
func (s *Scheduler) ScheduleTaskReminder(userID, taskID int, taskTitle string, dueDate time.Time, reminderMinutes int) error {
	reminderTime := dueDate.Add(-time.Duration(reminderMinutes) * time.Minute)
	
	// Only schedule if the reminder time is in the future
	if reminderTime.Before(time.Now()) {
		return fmt.Errorf("reminder time is in the past")
	}
	
	notification := &Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        TypeEmail,
		Priority:    PriorityNormal,
		Trigger:     TriggerDueDate,
		Title:       fmt.Sprintf("Task Reminder: %s", taskTitle),
		Message:     fmt.Sprintf("Your task '%s' is due in %d minutes.", taskTitle, reminderMinutes),
		Recipient:   "", // Would be populated from user email
		MaxRetries:  3,
	}
	
	return s.notificationService.ScheduleNotification(notification, reminderTime)
}

// ScheduleRecurringReminder schedules a recurring reminder for a task
func (s *Scheduler) ScheduleRecurringReminder(userID, taskID int, taskTitle string, dueDate time.Time, intervalMinutes int, maxReminders int) error {
	// Schedule multiple reminders
	for i := 0; i < maxReminders; i++ {
		reminderTime := dueDate.Add(-time.Duration((i+1)*intervalMinutes) * time.Minute)
		
		// Only schedule if the reminder time is in the future
		if reminderTime.After(time.Now()) {
			notification := &Notification{
				UserID:      userID,
				TaskID:      taskID,
				Type:        TypeEmail,
				Priority:    PriorityNormal,
				Trigger:     TriggerDueDate,
				Title:       fmt.Sprintf("Task Reminder %d/%d: %s", i+1, maxReminders, taskTitle),
				Message:     fmt.Sprintf("Your task '%s' is due in %d minutes.", taskTitle, (i+1)*intervalMinutes),
				Recipient:   "", // Would be populated from user email
				MaxRetries:  3,
			}
			
			err := s.notificationService.ScheduleNotification(notification, reminderTime)
			if err != nil {
				log.Printf("Error scheduling recurring reminder %d for task %d: %v", i+1, taskID, err)
				// Continue with other reminders
			}
		}
	}
	
	return nil
}

// GetSchedulerStatus gets the current scheduler status
func (s *Scheduler) GetSchedulerStatus() map[string]interface{} {
	return map[string]interface{}{
		"is_running": s.ctx.Err() == nil,
		"interval":   s.interval.String(),
		"ticker":     s.ticker != nil,
	}
}
