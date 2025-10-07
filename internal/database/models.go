package database

import (
	"time"
)

// DatabaseTask extends the original Task with database-specific fields
type DatabaseTask struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Status      int       `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	// Future extensions (will be added in later phases)
	CategoryID  *int      `json:"category_id,omitempty" db:"category_id"`
	UserID      *int      `json:"user_id,omitempty" db:"user_id"`
	IsArchived  bool      `json:"is_archived" db:"is_archived"`
}

// Category represents task categories (Phase 2)
type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Color       string    `json:"color" db:"color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Tag represents task tags (Phase 2)
type Tag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TaskTag represents many-to-many relationship between tasks and tags (Phase 2)
type TaskTag struct {
	TaskID int `json:"task_id" db:"task_id"`
	TagID  int `json:"tag_id" db:"tag_id"`
}

// User represents system users (Phase 5)
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Never expose password in JSON
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

