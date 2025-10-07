package task

import (
	"sync"
	"time"
	
	"learn-go-capstone/internal/database"
)

// StorageType represents the type of storage to use
type StorageType int

const (
	MemoryStorage StorageType = iota
	DatabaseStorage
	HybridStorage // Uses database with memory fallback
)

// HybridTaskManager combines memory and database storage
type HybridTaskManager struct {
	// Original memory-based manager
	memoryManager *TaskManager
	
	// Database repository
	repository database.Repository
	
	// Storage configuration
	storageType StorageType
	
	// Synchronization
	mu sync.RWMutex
}

// NewHybridTaskManager creates a new hybrid task manager
func NewHybridTaskManager(repository database.Repository, storageType StorageType) *HybridTaskManager {
	return &HybridTaskManager{
		memoryManager: NewTaskManager(),
		repository:    repository,
		storageType:   storageType,
	}
}

// AddTask adds a new task using the configured storage
func (htm *HybridTaskManager) AddTask(title, description string, priority Priority, dueDate *time.Time) *Task {
	htm.mu.Lock()
	defer htm.mu.Unlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.AddTask(title, description, priority, dueDate)
		
	case DatabaseStorage:
		dbTask := &database.DatabaseTask{
			Title:       title,
			Description: description,
			Priority:    int(priority),
			Status:      int(Pending),
			DueDate:     dueDate,
		}
		
		if err := htm.repository.CreateTask(dbTask); err != nil {
			// Fallback to memory storage if database fails
			return htm.memoryManager.AddTask(title, description, priority, dueDate)
		}
		
		// Convert back to original Task struct
		task := convertFromDatabaseTask(dbTask)
		return &task
		
	case HybridStorage:
		// Try database first, fallback to memory
		dbTask := &database.DatabaseTask{
			Title:       title,
			Description: description,
			Priority:    int(priority),
			Status:      int(Pending),
			DueDate:     dueDate,
		}
		
		if err := htm.repository.CreateTask(dbTask); err != nil {
			// Fallback to memory storage
			return htm.memoryManager.AddTask(title, description, priority, dueDate)
		}
		
		// Also add to memory for fast access
		htm.memoryManager.AddTask(title, description, priority, dueDate)
		
		// Convert back to original Task struct
		task := convertFromDatabaseTask(dbTask)
		return &task
		
	default:
		return htm.memoryManager.AddTask(title, description, priority, dueDate)
	}
}

// GetTask retrieves a task by ID
func (htm *HybridTaskManager) GetTask(id int) (*Task, error) {
	htm.mu.RLock()
	defer htm.mu.RUnlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.GetTask(id)
		
	case DatabaseStorage:
		dbTask, err := htm.repository.GetTask(id)
		if err != nil {
			return nil, err
		}
		task := convertFromDatabaseTask(dbTask)
		return &task, nil
		
	case HybridStorage:
		// Try memory first for speed
		if task, err := htm.memoryManager.GetTask(id); err == nil {
			return task, nil
		}
		
		// Fallback to database
		dbTask, err := htm.repository.GetTask(id)
		if err != nil {
			return nil, err
		}
		task := convertFromDatabaseTask(dbTask)
		return &task, nil
		
	default:
		return htm.memoryManager.GetTask(id)
	}
}

// UpdateTaskStatus updates the status of a task
func (htm *HybridTaskManager) UpdateTaskStatus(id int, status Status) error {
	htm.mu.Lock()
	defer htm.mu.Unlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.UpdateTaskStatus(id, status)
		
	case DatabaseStorage:
		dbTask, err := htm.repository.GetTask(id)
		if err != nil {
			return err
		}
		dbTask.Status = int(status)
		return htm.repository.UpdateTask(dbTask)
		
	case HybridStorage:
		// Update both memory and database
		if err := htm.memoryManager.UpdateTaskStatus(id, status); err != nil {
			// If memory update fails, try database
			dbTask, err := htm.repository.GetTask(id)
			if err != nil {
				return err
			}
			dbTask.Status = int(status)
			return htm.repository.UpdateTask(dbTask)
		}
		
		// Also update database
		dbTask, err := htm.repository.GetTask(id)
		if err == nil {
			dbTask.Status = int(status)
			htm.repository.UpdateTask(dbTask)
		}
		
		return nil
		
	default:
		return htm.memoryManager.UpdateTaskStatus(id, status)
	}
}

// DeleteTask removes a task by ID
func (htm *HybridTaskManager) DeleteTask(id int) error {
	htm.mu.Lock()
	defer htm.mu.Unlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.DeleteTask(id)
		
	case DatabaseStorage:
		return htm.repository.DeleteTask(id)
		
	case HybridStorage:
		// Delete from both memory and database
		memoryErr := htm.memoryManager.DeleteTask(id)
		dbErr := htm.repository.DeleteTask(id)
		
		// Return error if both fail
		if memoryErr != nil && dbErr != nil {
			return dbErr
		}
		
		return nil
		
	default:
		return htm.memoryManager.DeleteTask(id)
	}
}

// GetAllTasks returns all tasks
func (htm *HybridTaskManager) GetAllTasks() []Task {
	htm.mu.RLock()
	defer htm.mu.RUnlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.GetAllTasks()
		
	case DatabaseStorage:
		dbTasks, err := htm.repository.GetAllTasks()
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetAllTasks()
		}
		return convertFromDatabaseTasks(dbTasks)
		
	case HybridStorage:
		// Try database first
		dbTasks, err := htm.repository.GetAllTasks()
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetAllTasks()
		}
		return convertFromDatabaseTasks(dbTasks)
		
	default:
		return htm.memoryManager.GetAllTasks()
	}
}

// GetTasksByStatus returns tasks filtered by status
func (htm *HybridTaskManager) GetTasksByStatus(status Status) []Task {
	htm.mu.RLock()
	defer htm.mu.RUnlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.GetTasksByStatus(status)
		
	case DatabaseStorage:
		dbTasks, err := htm.repository.GetTasksByStatus(int(status))
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetTasksByStatus(status)
		}
		return convertFromDatabaseTasks(dbTasks)
		
	case HybridStorage:
		// Try database first
		dbTasks, err := htm.repository.GetTasksByStatus(int(status))
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetTasksByStatus(status)
		}
		return convertFromDatabaseTasks(dbTasks)
		
	default:
		return htm.memoryManager.GetTasksByStatus(status)
	}
}

// GetTasksByPriority returns tasks filtered by priority
func (htm *HybridTaskManager) GetTasksByPriority(priority Priority) []Task {
	htm.mu.RLock()
	defer htm.mu.RUnlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.GetTasksByPriority(priority)
		
	case DatabaseStorage:
		dbTasks, err := htm.repository.GetTasksByPriority(int(priority))
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetTasksByPriority(priority)
		}
		return convertFromDatabaseTasks(dbTasks)
		
	case HybridStorage:
		// Try database first
		dbTasks, err := htm.repository.GetTasksByPriority(int(priority))
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetTasksByPriority(priority)
		}
		return convertFromDatabaseTasks(dbTasks)
		
	default:
		return htm.memoryManager.GetTasksByPriority(priority)
	}
}

// SortTasksByPriority sorts tasks by priority (highest first)
func (htm *HybridTaskManager) SortTasksByPriority() []Task {
	// This method doesn't need database-specific implementation
	// as it's just sorting the results from GetAllTasks
	return htm.memoryManager.SortTasksByPriority()
}

// GetOverdueTasks returns tasks that are past their due date
func (htm *HybridTaskManager) GetOverdueTasks() []Task {
	htm.mu.RLock()
	defer htm.mu.RUnlock()
	
	switch htm.storageType {
	case MemoryStorage:
		return htm.memoryManager.GetOverdueTasks()
		
	case DatabaseStorage:
		dbTasks, err := htm.repository.GetOverdueTasks()
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetOverdueTasks()
		}
		return convertFromDatabaseTasks(dbTasks)
		
	case HybridStorage:
		// Try database first
		dbTasks, err := htm.repository.GetOverdueTasks()
		if err != nil {
			// Fallback to memory
			return htm.memoryManager.GetOverdueTasks()
		}
		return convertFromDatabaseTasks(dbTasks)
		
	default:
		return htm.memoryManager.GetOverdueTasks()
	}
}

// SetStorageType changes the storage type at runtime
func (htm *HybridTaskManager) SetStorageType(storageType StorageType) {
	htm.mu.Lock()
	defer htm.mu.Unlock()
	htm.storageType = storageType
}

// GetStorageType returns the current storage type
func (htm *HybridTaskManager) GetStorageType() StorageType {
	htm.mu.RLock()
	defer htm.mu.RUnlock()
	return htm.storageType
}

// SyncToDatabase syncs memory data to database
func (htm *HybridTaskManager) SyncToDatabase() error {
	htm.mu.Lock()
	defer htm.mu.Unlock()
	
	if htm.storageType == MemoryStorage {
		return nil // Nothing to sync
	}
	
	tasks := htm.memoryManager.GetAllTasks()
	for _, task := range tasks {
		dbTask := convertToDatabaseTask(task)
		if err := htm.repository.CreateTask(dbTask); err != nil {
			// Log error but continue with other tasks
			continue
		}
	}
	
	return nil
}

// LoadFromDatabase loads all tasks from database into memory
func (htm *HybridTaskManager) LoadFromDatabase() error {
	htm.mu.Lock()
	defer htm.mu.Unlock()
	
	if htm.storageType == MemoryStorage {
		return nil // Nothing to load
	}
	
	dbTasks, err := htm.repository.GetAllTasks()
	if err != nil {
		return err
	}
	
	// Clear memory and load from database
	htm.memoryManager = NewTaskManager()
	for _, dbTask := range dbTasks {
		task := convertFromDatabaseTask(&dbTask)
		htm.memoryManager.AddTask(task.Title, task.Description, task.Priority, task.DueDate)
	}
	
	return nil
}

// Converter functions to avoid import cycles

// convertToDatabaseTask converts a task.Task to database.DatabaseTask
func convertToDatabaseTask(t Task) *database.DatabaseTask {
	return &database.DatabaseTask{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    int(t.Priority),
		Status:      int(t.Status),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DueDate:     t.DueDate,
		CategoryID:  nil, // Will be set when categories are implemented
		UserID:      nil, // Will be set when users are implemented
		IsArchived:  false,
	}
}

// convertFromDatabaseTask converts a database.DatabaseTask to task.Task
func convertFromDatabaseTask(dt *database.DatabaseTask) Task {
	return Task{
		ID:          dt.ID,
		Title:       dt.Title,
		Description: dt.Description,
		Priority:    Priority(dt.Priority),
		Status:      Status(dt.Status),
		CreatedAt:   dt.CreatedAt,
		UpdatedAt:   dt.UpdatedAt,
		DueDate:     dt.DueDate,
	}
}

// convertFromDatabaseTasks converts a slice of database.DatabaseTask to []task.Task
func convertFromDatabaseTasks(dt []database.DatabaseTask) []Task {
	result := make([]Task, len(dt))
	for i, t := range dt {
		result[i] = convertFromDatabaseTask(&t)
	}
	return result
}
