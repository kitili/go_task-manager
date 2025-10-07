package notifications

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"learn-go-capstone/internal/database"
)

// NotificationService handles notification operations
type NotificationService struct {
	repository database.Repository
	config     NotificationConfig
	queue      chan *Notification
	workers    []*NotificationWorker
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NotificationWorker handles processing notifications
type NotificationWorker struct {
	id       int
	service  *NotificationService
	queue    chan *Notification
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewNotificationService creates a new notification service
func NewNotificationService(repository database.Repository, config NotificationConfig) *NotificationService {
	ctx, cancel := context.WithCancel(context.Background())
	
	service := &NotificationService{
		repository: repository,
		config:     config,
		queue:      make(chan *Notification, config.BatchSize),
		ctx:        ctx,
		cancel:     cancel,
	}
	
	// Start workers
	service.startWorkers()
	
	return service
}

// startWorkers starts the notification workers
func (ns *NotificationService) startWorkers() {
	ns.workers = make([]*NotificationWorker, ns.config.WorkerCount)
	
	for i := 0; i < ns.config.WorkerCount; i++ {
		workerCtx, workerCancel := context.WithCancel(ns.ctx)
		worker := &NotificationWorker{
			id:      i,
			service: ns,
			queue:   ns.queue,
			ctx:     workerCtx,
			cancel:  workerCancel,
		}
		
		ns.workers[i] = worker
		ns.wg.Add(1)
		go worker.start()
	}
}

// start starts the worker
func (nw *NotificationWorker) start() {
	defer nw.service.wg.Done()
	
	log.Printf("Notification worker %d started", nw.id)
	
	for {
		select {
		case notification := <-nw.queue:
			nw.processNotification(notification)
		case <-nw.ctx.Done():
			log.Printf("Notification worker %d stopped", nw.id)
			return
		}
	}
}

// processNotification processes a single notification
func (nw *NotificationWorker) processNotification(notification *Notification) {
	log.Printf("Worker %d processing notification %d", nw.id, notification.ID)
	
	// Update status to pending
	notification.Status = StatusPending
	notification.UpdatedAt = time.Now()
	
	// Send notification based on type
	var err error
	switch notification.Type {
	case TypeEmail:
		err = nw.sendEmail(notification)
	case TypeInApp:
		err = nw.sendInApp(notification)
	case TypeSMS:
		err = nw.sendSMS(notification)
	case TypeWebhook:
		err = nw.sendWebhook(notification)
	case TypeSlack:
		err = nw.sendSlack(notification)
	case TypeDiscord:
		err = nw.sendDiscord(notification)
	default:
		err = fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
	
	// Update notification status
	if err != nil {
		notification.Status = StatusFailed
		notification.Error = err.Error()
		notification.RetryCount++
		log.Printf("Failed to send notification %d: %v", notification.ID, err)
		
		// Retry if under max retries
		if notification.RetryCount < notification.MaxRetries {
			notification.Status = StatusPending
			notification.ScheduledAt = &[]time.Time{time.Now().Add(time.Duration(nw.service.config.RetryDelay) * time.Second)}[0]
			nw.queue <- notification
		}
	} else {
		notification.Status = StatusSent
		now := time.Now()
		notification.SentAt = &now
		log.Printf("Successfully sent notification %d", notification.ID)
	}
	
	notification.UpdatedAt = time.Now()
	// Note: In a real implementation, you'd update the database here
}

// sendEmail sends an email notification
func (nw *NotificationWorker) sendEmail(notification *Notification) error {
	// This is a placeholder implementation
	// In a real system, you'd use an email service like SendGrid, AWS SES, etc.
	log.Printf("Sending email to %s: %s", notification.Recipient, notification.Title)
	
	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)
	
	// Simulate occasional failures for testing
	if notification.RetryCount > 0 && notification.RetryCount%3 == 0 {
		return fmt.Errorf("simulated email sending failure")
	}
	
	return nil
}

// sendInApp sends an in-app notification
func (nw *NotificationWorker) sendInApp(notification *Notification) error {
	// This is a placeholder implementation
	// In a real system, you'd store this in a database for the user to see
	log.Printf("Sending in-app notification to user %d: %s", notification.UserID, notification.Title)
	
	// Simulate in-app notification delay
	time.Sleep(50 * time.Millisecond)
	
	return nil
}

// sendSMS sends an SMS notification
func (nw *NotificationWorker) sendSMS(notification *Notification) error {
	// This is a placeholder implementation
	// In a real system, you'd use Twilio, AWS SNS, etc.
	log.Printf("Sending SMS to %s: %s", notification.Recipient, notification.Title)
	
	// Simulate SMS sending delay
	time.Sleep(200 * time.Millisecond)
	
	return nil
}

// sendWebhook sends a webhook notification
func (nw *NotificationWorker) sendWebhook(notification *Notification) error {
	// This is a placeholder implementation
	// In a real system, you'd make an HTTP POST request
	log.Printf("Sending webhook to %s: %s", notification.Channel, notification.Title)
	
	// Simulate webhook delay
	time.Sleep(150 * time.Millisecond)
	
	return nil
}

// sendSlack sends a Slack notification
func (nw *NotificationWorker) sendSlack(notification *Notification) error {
	// This is a placeholder implementation
	// In a real system, you'd use Slack's webhook API
	log.Printf("Sending Slack notification to %s: %s", notification.Channel, notification.Title)
	
	// Simulate Slack delay
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// sendDiscord sends a Discord notification
func (nw *NotificationWorker) sendDiscord(notification *Notification) error {
	// This is a placeholder implementation
	// In a real system, you'd use Discord's webhook API
	log.Printf("Sending Discord notification to %s: %s", notification.Channel, notification.Title)
	
	// Simulate Discord delay
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// SendNotification sends a notification immediately
func (ns *NotificationService) SendNotification(notification *Notification) error {
	// Set default values
	if notification.MaxRetries == 0 {
		notification.MaxRetries = ns.config.MaxRetries
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}
	notification.UpdatedAt = time.Now()
	
	// Queue the notification
	select {
	case ns.queue <- notification:
		return nil
	case <-ns.ctx.Done():
		return fmt.Errorf("notification service is shutting down")
	default:
		return fmt.Errorf("notification queue is full")
	}
}

// ScheduleNotification schedules a notification for later delivery
func (ns *NotificationService) ScheduleNotification(notification *Notification, scheduledAt time.Time) error {
	notification.ScheduledAt = &scheduledAt
	notification.Status = StatusPending
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	
	// In a real implementation, you'd store this in a database
	// and have a scheduler pick it up at the scheduled time
	go func() {
		time.Sleep(time.Until(scheduledAt))
		ns.SendNotification(notification)
	}()
	
	return nil
}

// CreateTaskReminder creates a reminder for a task
func (ns *NotificationService) CreateTaskReminder(userID, taskID int, taskTitle string, dueDate time.Time, reminderMinutes int) error {
	reminderTime := dueDate.Add(-time.Duration(reminderMinutes) * time.Minute)
	
	notification := &Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        TypeEmail,
		Priority:    PriorityNormal,
		Trigger:     TriggerDueDate,
		Title:       fmt.Sprintf("Task Reminder: %s", taskTitle),
		Message:     fmt.Sprintf("Your task '%s' is due in %d minutes.", taskTitle, reminderMinutes),
		Recipient:   "", // Would be populated from user email
		MaxRetries:  ns.config.MaxRetries,
	}
	
	return ns.ScheduleNotification(notification, reminderTime)
}

// CreateOverdueReminder creates an overdue reminder for a task
func (ns *NotificationService) CreateOverdueReminder(userID, taskID int, taskTitle string, overdueHours int) error {
	notification := &Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        TypeEmail,
		Priority:    PriorityHigh,
		Trigger:     TriggerOverdue,
		Title:       fmt.Sprintf("Overdue Task: %s", taskTitle),
		Message:     fmt.Sprintf("Your task '%s' is %d hours overdue.", taskTitle, overdueHours),
		Recipient:   "", // Would be populated from user email
		MaxRetries:  ns.config.MaxRetries,
	}
	
	return ns.SendNotification(notification)
}

// CreateStatusChangeNotification creates a notification for task status changes
func (ns *NotificationService) CreateStatusChangeNotification(userID, taskID int, taskTitle, oldStatus, newStatus string) error {
	notification := &Notification{
		UserID:      userID,
		TaskID:      taskID,
		Type:        TypeInApp,
		Priority:    PriorityNormal,
		Trigger:     TriggerStatusChange,
		Title:       "Task Status Updated",
		Message:     fmt.Sprintf("Task '%s' status changed from %s to %s", taskTitle, oldStatus, newStatus),
		Recipient:   "", // In-app notifications don't need recipient
		MaxRetries:  ns.config.MaxRetries,
	}
	
	return ns.SendNotification(notification)
}

// GetUserNotifications gets notifications for a user
func (ns *NotificationService) GetUserNotifications(userID int, limit, offset int) ([]*Notification, error) {
	// This is a placeholder implementation
	// In a real system, you'd query the database
	return []*Notification{}, nil
}

// MarkNotificationAsRead marks a notification as read
func (ns *NotificationService) MarkNotificationAsRead(notificationID int) error {
	// This is a placeholder implementation
	// In a real system, you'd update the database
	return nil
}

// GetNotificationStats gets notification statistics
func (ns *NotificationService) GetNotificationStats() (*NotificationStats, error) {
	// This is a placeholder implementation
	// In a real system, you'd query the database for statistics
	return &NotificationStats{
		TotalSent:       0,
		TotalDelivered:  0,
		TotalFailed:     0,
		TotalPending:    0,
		SuccessRate:     0.0,
		AverageDeliveryTime: 0,
		ByType:         make(map[NotificationType]int),
		ByPriority:     make(map[NotificationPriority]int),
		ByTrigger:      make(map[NotificationTrigger]int),
	}, nil
}

// Stop stops the notification service
func (ns *NotificationService) Stop() {
	log.Println("Stopping notification service...")
	
	// Cancel context to stop workers
	ns.cancel()
	
	// Wait for workers to finish
	ns.wg.Wait()
	
	// Close queue
	close(ns.queue)
	
	log.Println("Notification service stopped")
}

// GetQueueStatus gets the current queue status
func (ns *NotificationService) GetQueueStatus() map[string]interface{} {
	return map[string]interface{}{
		"queue_length": len(ns.queue),
		"worker_count": len(ns.workers),
		"is_running":   ns.ctx.Err() == nil,
	}
}
