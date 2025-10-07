package task

import (
	"learn-go-capstone/internal/database"
)

// CategoryManager manages categories and tags
type CategoryManager struct {
	repository database.Repository
}

// NewCategoryManager creates a new category manager
func NewCategoryManager(repository database.Repository) *CategoryManager {
	return &CategoryManager{
		repository: repository,
	}
}

// Category operations

// CreateCategory creates a new category
func (cm *CategoryManager) CreateCategory(name, description, color string) (*database.Category, error) {
	category := &database.Category{
		Name:        name,
		Description: description,
		Color:       color,
	}
	
	err := cm.repository.CreateCategory(category)
	if err != nil {
		return nil, err
	}
	
	return category, nil
}

// GetCategory retrieves a category by ID
func (cm *CategoryManager) GetCategory(id int) (*database.Category, error) {
	return cm.repository.GetCategory(id)
}

// GetAllCategories returns all categories
func (cm *CategoryManager) GetAllCategories() ([]database.Category, error) {
	return cm.repository.GetAllCategories()
}

// UpdateCategory updates a category
func (cm *CategoryManager) UpdateCategory(category *database.Category) error {
	return cm.repository.UpdateCategory(category)
}

// DeleteCategory deletes a category by ID
func (cm *CategoryManager) DeleteCategory(id int) error {
	return cm.repository.DeleteCategory(id)
}

// Tag operations

// CreateTag creates a new tag
func (cm *CategoryManager) CreateTag(name, color string) (*database.Tag, error) {
	tag := &database.Tag{
		Name:  name,
		Color: color,
	}
	
	err := cm.repository.CreateTag(tag)
	if err != nil {
		return nil, err
	}
	
	return tag, nil
}

// GetTag retrieves a tag by ID
func (cm *CategoryManager) GetTag(id int) (*database.Tag, error) {
	return cm.repository.GetTag(id)
}

// GetAllTags returns all tags
func (cm *CategoryManager) GetAllTags() ([]database.Tag, error) {
	return cm.repository.GetAllTags()
}

// UpdateTag updates a tag
func (cm *CategoryManager) UpdateTag(tag *database.Tag) error {
	return cm.repository.UpdateTag(tag)
}

// DeleteTag deletes a tag by ID
func (cm *CategoryManager) DeleteTag(id int) error {
	return cm.repository.DeleteTag(id)
}

// Task-Tag relationship operations

// AddTagToTask adds a tag to a task
func (cm *CategoryManager) AddTagToTask(taskID, tagID int) error {
	return cm.repository.AddTagToTask(taskID, tagID)
}

// RemoveTagFromTask removes a tag from a task
func (cm *CategoryManager) RemoveTagFromTask(taskID, tagID int) error {
	return cm.repository.RemoveTagFromTask(taskID, tagID)
}

// GetTaskTags returns all tags for a task
func (cm *CategoryManager) GetTaskTags(taskID int) ([]database.Tag, error) {
	return cm.repository.GetTaskTags(taskID)
}

// GetTasksByTag returns all tasks with a specific tag
func (cm *CategoryManager) GetTasksByTag(tagID int) ([]Task, error) {
	dbTasks, err := cm.repository.GetTasksByTag(tagID)
	if err != nil {
		return nil, err
	}
	
	return convertFromDatabaseTasks(dbTasks), nil
}
