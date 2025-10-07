package task

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Priority represents task priority levels
type Priority int

const (
	Low Priority = iota + 1
	Medium
	High
	Urgent
)

// Status represents task status
type Status int

const (
	Pending Status = iota
	InProgress
	Completed
	Cancelled
)

// Task represents a single task in our task manager
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	// Phase 2: Enhanced data model
	Category    *Category `json:"category,omitempty"`
	Tags        []Tag     `json:"tags,omitempty"`
}

// Category represents a task category
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag represents a task tag
type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskManager manages a collection of tasks
type TaskManager struct {
	tasks []Task
	nextID int
	mu     sync.RWMutex
}

// NewTaskManager creates a new task manager instance
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

// AddTask adds a new task to the manager
func (tm *TaskManager) AddTask(title, description string, priority Priority, dueDate *time.Time) *Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	now := time.Now()
	task := Task{
		ID:          tm.nextID,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      Pending,
		CreatedAt:   now,
		UpdatedAt:   now,
		DueDate:     dueDate,
	}
	
	tm.tasks = append(tm.tasks, task)
	tm.nextID++
	
	return &task
}

// GetTask retrieves a task by ID
func (tm *TaskManager) GetTask(id int) (*Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			return &tm.tasks[i], nil
		}
	}
	
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// UpdateTaskStatus updates the status of a task
func (tm *TaskManager) UpdateTaskStatus(id int, status Status) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks[i].Status = status
			tm.tasks[i].UpdatedAt = time.Now()
			return nil
		}
	}
	
	return fmt.Errorf("task with ID %d not found", id)
}

// DeleteTask removes a task by ID
func (tm *TaskManager) DeleteTask(id int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("task with ID %d not found", id)
}

// GetAllTasks returns all tasks
func (tm *TaskManager) GetAllTasks() []Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	// Return a copy to prevent external modifications
	tasks := make([]Task, len(tm.tasks))
	copy(tasks, tm.tasks)
	return tasks
}

// GetTasksByStatus returns tasks filtered by status
func (tm *TaskManager) GetTasksByStatus(status Status) []Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	var filtered []Task
	for _, task := range tm.tasks {
		if task.Status == status {
			filtered = append(filtered, task)
		}
	}
	
	return filtered
}

// GetTasksByPriority returns tasks filtered by priority
func (tm *TaskManager) GetTasksByPriority(priority Priority) []Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	var filtered []Task
	for _, task := range tm.tasks {
		if task.Priority == priority {
			filtered = append(filtered, task)
		}
	}
	
	return filtered
}

// SortTasksByPriority sorts tasks by priority (highest first)
func (tm *TaskManager) SortTasksByPriority() []Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	tasks := make([]Task, len(tm.tasks))
	copy(tasks, tm.tasks)
	
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})
	
	return tasks
}

// GetOverdueTasks returns tasks that are past their due date
func (tm *TaskManager) GetOverdueTasks() []Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	now := time.Now()
	var overdue []Task
	
	for _, task := range tm.tasks {
		if task.DueDate != nil && task.DueDate.Before(now) && task.Status != Completed {
			overdue = append(overdue, task)
		}
	}
	
	return overdue
}

// String methods for better display
func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Urgent:
		return "Urgent"
	default:
		return "Unknown"
	}
}

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case InProgress:
		return "In Progress"
	case Completed:
		return "Completed"
	case Cancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

func (t Task) String() string {
	dueDateStr := "No due date"
	if t.DueDate != nil {
		dueDateStr = t.DueDate.Format("2006-01-02 15:04")
	}
	
	return fmt.Sprintf("ID: %d | %s | %s | %s | %s | Due: %s", 
		t.ID, t.Title, t.Priority.String(), t.Status.String(), 
		t.CreatedAt.Format("2006-01-02"), dueDateStr)
}
