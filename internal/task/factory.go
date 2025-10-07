package task

import (
	"database/sql"
	
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


// Ensure TaskManager implements TaskManagerInterface
var _ TaskManagerInterface = (*TaskManager)(nil)

// Ensure HybridTaskManager implements TaskManagerInterface
var _ TaskManagerInterface = (*HybridTaskManager)(nil)
