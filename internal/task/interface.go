package task

import "time"

// TaskManagerInterface defines the interface for task management operations
type TaskManagerInterface interface {
	AddTask(title, description string, priority Priority, dueDate *time.Time) *Task
	GetTask(id int) (*Task, error)
	UpdateTaskStatus(id int, status Status) error
	DeleteTask(id int) error
	GetAllTasks() []Task
	GetTasksByStatus(status Status) []Task
	GetTasksByPriority(priority Priority) []Task
	SortTasksByPriority() []Task
	GetOverdueTasks() []Task
}
