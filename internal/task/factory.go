package task

import (
	"database/sql"
	"time"
	
	"learn-go-capstone/internal/database"
)

// TaskManagerFactory creates the appropriate TaskManager based on configuration
type TaskManagerFactory struct{}

// NewTaskManagerFactory creates a new factory
func NewTaskManagerFactory() *TaskManagerFactory {
	return &TaskManagerFactory{}
}

// CreateTaskManager creates a TaskManager based on the storage type
func (f *TaskManagerFactory) CreateTaskManager(storageType StorageType, db *sql.DB) TaskManagerInterface {
	switch storageType {
	case MemoryStorage:
		return NewTaskManager()
		
	case DatabaseStorage, HybridStorage:
		if db == nil {
			// Fallback to memory if no database connection
			return NewTaskManager()
		}
		
		repository := database.NewSQLiteRepository(db)
		return NewHybridTaskManager(repository, storageType)
		
	default:
		return NewTaskManager()
	}
}

// TaskManagerInterface defines the interface that all task managers must implement
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

// Ensure TaskManager implements TaskManagerInterface
var _ TaskManagerInterface = (*TaskManager)(nil)

// Ensure HybridTaskManager implements TaskManagerInterface
var _ TaskManagerInterface = (*HybridTaskManager)(nil)
