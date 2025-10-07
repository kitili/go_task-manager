package notifications

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	TypeEmail     NotificationType = "email"
	TypeInApp     NotificationType = "in_app"
	TypeSMS       NotificationType = "sms"
	TypeWebhook   NotificationType = "webhook"
	TypeSlack     NotificationType = "slack"
	TypeDiscord   NotificationType = "discord"
)

// NotificationPriority represents the priority level
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityNormal   NotificationPriority = "normal"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// NotificationStatus represents the delivery status
type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusSent      NotificationStatus = "sent"
	StatusDelivered NotificationStatus = "delivered"
	StatusFailed    NotificationStatus = "failed"
	StatusCancelled NotificationStatus = "cancelled"
)

// NotificationTrigger represents what triggers the notification
type NotificationTrigger string

const (
	TriggerDueDate      NotificationTrigger = "due_date"
	TriggerOverdue      NotificationTrigger = "overdue"
	TriggerStatusChange NotificationTrigger = "status_change"
	TriggerCreated      NotificationTrigger = "created"
	TriggerUpdated      NotificationTrigger = "updated"
	TriggerCustom       NotificationTrigger = "custom"
)

// Notification represents a notification in the system
type Notification struct {
	ID           int                `json:"id" db:"id"`
	UserID       int                `json:"user_id" db:"user_id"`
	TaskID       int                `json:"task_id" db:"task_id"`
	Type         NotificationType   `json:"type" db:"type"`
	Priority     NotificationPriority `json:"priority" db:"priority"`
	Status       NotificationStatus `json:"status" db:"status"`
	Trigger      NotificationTrigger `json:"trigger" db:"trigger"`
	Title        string             `json:"title" db:"title"`
	Message      string             `json:"message" db:"message"`
	Recipient    string             `json:"recipient" db:"recipient"`
	Channel      string             `json:"channel" db:"channel"`
	ScheduledAt  *time.Time         `json:"scheduled_at" db:"scheduled_at"`
	SentAt       *time.Time         `json:"sent_at" db:"sent_at"`
	DeliveredAt  *time.Time         `json:"delivered_at" db:"delivered_at"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`
	RetryCount   int                `json:"retry_count" db:"retry_count"`
	MaxRetries   int                `json:"max_retries" db:"max_retries"`
	Error        string             `json:"error" db:"error"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
}

// NotificationTemplate represents a reusable notification template
type NotificationTemplate struct {
	ID          int                `json:"id" db:"id"`
	Name        string             `json:"name" db:"name"`
	Type        NotificationType   `json:"type" db:"type"`
	Trigger     NotificationTrigger `json:"trigger" db:"trigger"`
	Subject     string             `json:"subject" db:"subject"`
	Body        string             `json:"body" db:"body"`
	IsActive    bool               `json:"is_active" db:"is_active"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// NotificationSettings represents user notification preferences
type NotificationSettings struct {
	ID                    int                `json:"id" db:"id"`
	UserID                int                `json:"user_id" db:"user_id"`
	EmailEnabled          bool               `json:"email_enabled" db:"email_enabled"`
	InAppEnabled          bool               `json:"in_app_enabled" db:"in_app_enabled"`
	SMSEnabled            bool               `json:"sms_enabled" db:"sms_enabled"`
	WebhookEnabled        bool               `json:"webhook_enabled" db:"webhook_enabled"`
	SlackEnabled          bool               `json:"slack_enabled" db:"slack_enabled"`
	DiscordEnabled        bool               `json:"discord_enabled" db:"discord_enabled"`
	DueDateReminder       bool               `json:"due_date_reminder" db:"due_date_reminder"`
	OverdueReminder       bool               `json:"overdue_reminder" db:"overdue_reminder"`
	StatusChangeReminder  bool               `json:"status_change_reminder" db:"status_change_reminder"`
	CreatedReminder       bool               `json:"created_reminder" db:"created_reminder"`
	UpdatedReminder       bool               `json:"updated_reminder" db:"updated_reminder"`
	ReminderTime          int                `json:"reminder_time" db:"reminder_time"` // minutes before due date
	OverdueInterval       int                `json:"overdue_interval" db:"overdue_interval"` // hours between overdue reminders
	QuietHoursStart       string             `json:"quiet_hours_start" db:"quiet_hours_start"` // HH:MM format
	QuietHoursEnd         string             `json:"quiet_hours_end" db:"quiet_hours_end"` // HH:MM format
	Timezone              string             `json:"timezone" db:"timezone"`
	CreatedAt             time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at" db:"updated_at"`
}

// NotificationChannel represents a notification delivery channel
type NotificationChannel struct {
	ID          int                `json:"id" db:"id"`
	UserID      int                `json:"user_id" db:"user_id"`
	Type        NotificationType   `json:"type" db:"type"`
	Name        string             `json:"name" db:"name"`
	Config      map[string]interface{} `json:"config" db:"config"`
	IsActive    bool               `json:"is_active" db:"is_active"`
	IsVerified  bool               `json:"is_verified" db:"is_verified"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// NotificationQueue represents a queued notification
type NotificationQueue struct {
	ID           int                `json:"id" db:"id"`
	NotificationID int              `json:"notification_id" db:"notification_id"`
	Priority     NotificationPriority `json:"priority" db:"priority"`
	ScheduledAt  time.Time          `json:"scheduled_at" db:"scheduled_at"`
	RetryCount   int                `json:"retry_count" db:"retry_count"`
	MaxRetries   int                `json:"max_retries" db:"max_retries"`
	Status       NotificationStatus `json:"status" db:"status"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalSent       int `json:"total_sent"`
	TotalDelivered  int `json:"total_delivered"`
	TotalFailed     int `json:"total_failed"`
	TotalPending    int `json:"total_pending"`
	SuccessRate     float64 `json:"success_rate"`
	AverageDeliveryTime int `json:"average_delivery_time_ms"`
	ByType          map[NotificationType]int `json:"by_type"`
	ByPriority      map[NotificationPriority]int `json:"by_priority"`
	ByTrigger       map[NotificationTrigger]int `json:"by_trigger"`
}

// NotificationConfig represents system-wide notification configuration
type NotificationConfig struct {
	SMTPHost         string `json:"smtp_host"`
	SMTPPort         int    `json:"smtp_port"`
	SMTPUsername     string `json:"smtp_username"`
	SMTPPassword     string `json:"smtp_password"`
	SMTPFromEmail    string `json:"smtp_from_email"`
	SMTPFromName     string `json:"smtp_from_name"`
	TwilioSID        string `json:"twilio_sid"`
	TwilioToken      string `json:"twilio_token"`
	TwilioFromNumber string `json:"twilio_from_number"`
	SlackWebhookURL  string `json:"slack_webhook_url"`
	DiscordWebhookURL string `json:"discord_webhook_url"`
	MaxRetries       int    `json:"max_retries"`
	RetryDelay       int    `json:"retry_delay_seconds"`
	BatchSize        int    `json:"batch_size"`
	WorkerCount      int    `json:"worker_count"`
}

// DefaultNotificationSettings returns default notification settings
func DefaultNotificationSettings() NotificationSettings {
	return NotificationSettings{
		EmailEnabled:          true,
		InAppEnabled:          true,
		SMSEnabled:            false,
		WebhookEnabled:        false,
		SlackEnabled:          false,
		DiscordEnabled:        false,
		DueDateReminder:       true,
		OverdueReminder:       true,
		StatusChangeReminder:  false,
		CreatedReminder:       false,
		UpdatedReminder:       false,
		ReminderTime:          60, // 1 hour before due date
		OverdueInterval:       24, // 24 hours between overdue reminders
		QuietHoursStart:       "22:00",
		QuietHoursEnd:         "08:00",
		Timezone:              "UTC",
	}
}

// DefaultNotificationConfig returns default notification configuration
func DefaultNotificationConfig() NotificationConfig {
	return NotificationConfig{
		MaxRetries:     3,
		RetryDelay:     300, // 5 minutes
		BatchSize:      100,
		WorkerCount:    5,
	}
}
