package task

import (
	"learn-go-capstone/internal/database"
)

// DependencyManager manages task dependencies
type DependencyManager struct {
	repository database.Repository
}

// NewDependencyManager creates a new dependency manager
func NewDependencyManager(repository database.Repository) *DependencyManager {
	return &DependencyManager{
		repository: repository,
	}
}

// AddDependency adds a dependency between two tasks
func (dm *DependencyManager) AddDependency(taskID, dependsOnTaskID int) error {
	// Validate that both tasks exist
	_, err := dm.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	_, err = dm.repository.GetTask(dependsOnTaskID)
	if err != nil {
		return err
	}
	
	// Add the dependency (this will check for circular dependencies)
	return dm.repository.AddTaskDependency(taskID, dependsOnTaskID)
}

// RemoveDependency removes a dependency between two tasks
func (dm *DependencyManager) RemoveDependency(taskID, dependsOnTaskID int) error {
	return dm.repository.RemoveTaskDependency(taskID, dependsOnTaskID)
}

// GetDependencies returns all dependencies for a task
func (dm *DependencyManager) GetDependencies(taskID int) ([]database.TaskDependency, error) {
	return dm.repository.GetTaskDependencies(taskID)
}

// GetTasksThatDependOn returns all tasks that depend on the given task
func (dm *DependencyManager) GetTasksThatDependOn(taskID int) ([]Task, error) {
	dbTasks, err := dm.repository.GetTasksThatDependOn(taskID)
	if err != nil {
		return nil, err
	}
	
	return convertFromDatabaseTasks(dbTasks), nil
}

// GetTasksThatTaskDependsOn returns all tasks that the given task depends on
func (dm *DependencyManager) GetTasksThatTaskDependsOn(taskID int) ([]Task, error) {
	dbTasks, err := dm.repository.GetTasksThatTaskDependsOn(taskID)
	if err != nil {
		return nil, err
	}
	
	return convertFromDatabaseTasks(dbTasks), nil
}

// CheckCircularDependency checks if adding a dependency would create a circular dependency
func (dm *DependencyManager) CheckCircularDependency(taskID, dependsOnTaskID int) (bool, error) {
	return dm.repository.CheckCircularDependency(taskID, dependsOnTaskID)
}

// GetDependencyChain returns the full dependency chain for a task
func (dm *DependencyManager) GetDependencyChain(taskID int) ([]Task, error) {
	// Get all tasks that this task depends on (recursively)
	visited := make(map[int]bool)
	return dm.getDependencyChainRecursive(taskID, visited)
}

// getDependencyChainRecursive is a helper function for GetDependencyChain
func (dm *DependencyManager) getDependencyChainRecursive(taskID int, visited map[int]bool) ([]Task, error) {
	if visited[taskID] {
		return nil, nil // Avoid infinite recursion
	}
	
	visited[taskID] = true
	defer delete(visited, taskID)
	
	// Get direct dependencies
	dependencies, err := dm.GetTasksThatTaskDependsOn(taskID)
	if err != nil {
		return nil, err
	}
	
	var allDependencies []Task
	for _, dep := range dependencies {
		// Add this dependency
		allDependencies = append(allDependencies, dep)
		
		// Get dependencies of this dependency
		subDeps, err := dm.getDependencyChainRecursive(dep.ID, visited)
		if err != nil {
			return nil, err
		}
		
		// Add sub-dependencies
		allDependencies = append(allDependencies, subDeps...)
	}
	
	return allDependencies, nil
}

// CanCompleteTask checks if a task can be completed (all dependencies are completed)
func (dm *DependencyManager) CanCompleteTask(taskID int) (bool, error) {
	dependencies, err := dm.GetTasksThatTaskDependsOn(taskID)
	if err != nil {
		return false, err
	}
	
	// Check if all dependencies are completed
	for _, dep := range dependencies {
		if dep.Status != Completed {
			return false, nil
		}
	}
	
	return true, nil
}

// GetBlockedTasks returns all tasks that are blocked by incomplete dependencies
func (dm *DependencyManager) GetBlockedTasks() ([]Task, error) {
	// This would require a more complex query to find all tasks with incomplete dependencies
	// For now, we'll implement a simple version that checks all tasks
	allTasks, err := dm.repository.GetAllTasks()
	if err != nil {
		return nil, err
	}
	
	var blockedTasks []Task
	for _, dbTask := range allTasks {
		task := convertFromDatabaseTask(&dbTask)
		if task.Status == Pending || task.Status == InProgress {
			canComplete, err := dm.CanCompleteTask(task.ID)
			if err != nil {
				return nil, err
			}
			if !canComplete {
				blockedTasks = append(blockedTasks, task)
			}
		}
	}
	
	return blockedTasks, nil
}

// GetReadyTasks returns all tasks that are ready to be started (no incomplete dependencies)
func (dm *DependencyManager) GetReadyTasks() ([]Task, error) {
	allTasks, err := dm.repository.GetAllTasks()
	if err != nil {
		return nil, err
	}
	
	var readyTasks []Task
	for _, dbTask := range allTasks {
		task := convertFromDatabaseTask(&dbTask)
		if task.Status == Pending {
			canComplete, err := dm.CanCompleteTask(task.ID)
			if err != nil {
				return nil, err
			}
			if canComplete {
				readyTasks = append(readyTasks, task)
			}
		}
	}
	
	return readyTasks, nil
}
