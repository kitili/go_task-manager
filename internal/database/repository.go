package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository defines the interface for task data operations
type Repository interface {
	// Task operations
	CreateTask(task *DatabaseTask) error
	GetTask(id int) (*DatabaseTask, error)
	UpdateTask(task *DatabaseTask) error
	DeleteTask(id int) error
	GetAllTasks() ([]DatabaseTask, error)
	GetTasksByStatus(status int) ([]DatabaseTask, error)
	GetTasksByPriority(priority int) ([]DatabaseTask, error)
	GetOverdueTasks() ([]DatabaseTask, error)
	
	// Future operations (will be implemented in later phases)
	GetTasksByUser(userID int) ([]DatabaseTask, error)
	GetTasksByCategory(categoryID int) ([]DatabaseTask, error)
	SearchTasks(query string) ([]DatabaseTask, error)
	
	// Category operations (Phase 2)
	CreateCategory(category *Category) error
	GetCategory(id int) (*Category, error)
	GetAllCategories() ([]Category, error)
	UpdateCategory(category *Category) error
	DeleteCategory(id int) error
	
	// Tag operations (Phase 2)
	CreateTag(tag *Tag) error
	GetTag(id int) (*Tag, error)
	GetAllTags() ([]Tag, error)
	UpdateTag(tag *Tag) error
	DeleteTag(id int) error
	
	// Task-Tag relationship operations (Phase 2)
	AddTagToTask(taskID, tagID int) error
	RemoveTagFromTask(taskID, tagID int) error
	GetTaskTags(taskID int) ([]Tag, error)
	GetTasksByTag(tagID int) ([]DatabaseTask, error)
	
	// User operations (Phase 5)
	CreateUser(user *User) error
	GetUser(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
	GetAllUsers() ([]User, error)
}

// SQLiteRepository implements Repository interface for SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(db *sql.DB) Repository {
	return &SQLiteRepository{db: db}
}

// Task operations implementation

func (r *SQLiteRepository) CreateTask(task *DatabaseTask) error {
	query := `
	INSERT INTO tasks (title, description, priority, status, due_date, user_id, category_id, is_archived)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, 
		task.Title, 
		task.Description, 
		task.Priority, 
		task.Status, 
		task.DueDate, 
		task.UserID, 
		task.CategoryID, 
		task.IsArchived)
	
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get task ID: %w", err)
	}
	
	task.ID = int(id)
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	
	return nil
}

func (r *SQLiteRepository) GetTask(id int) (*DatabaseTask, error) {
	query := `
	SELECT id, title, description, priority, status, created_at, updated_at, due_date, user_id, category_id, is_archived
	FROM tasks WHERE id = ?`
	
	task := &DatabaseTask{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DueDate,
		&task.UserID,
		&task.CategoryID,
		&task.IsArchived,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	
	return task, nil
}

func (r *SQLiteRepository) UpdateTask(task *DatabaseTask) error {
	query := `
	UPDATE tasks 
	SET title = ?, description = ?, priority = ?, status = ?, updated_at = ?, due_date = ?, user_id = ?, category_id = ?, is_archived = ?
	WHERE id = ?`
	
	task.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(query,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.UpdatedAt,
		task.DueDate,
		task.UserID,
		task.CategoryID,
		task.IsArchived,
		task.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", task.ID)
	}
	
	return nil
}

func (r *SQLiteRepository) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}
	
	return nil
}

func (r *SQLiteRepository) GetAllTasks() ([]DatabaseTask, error) {
	query := `
	SELECT id, title, description, priority, status, created_at, updated_at, due_date, user_id, category_id, is_archived
	FROM tasks 
	WHERE is_archived = FALSE
	ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()
	
	var tasks []DatabaseTask
	for rows.Next() {
		task := DatabaseTask{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DueDate,
			&task.UserID,
			&task.CategoryID,
			&task.IsArchived,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

func (r *SQLiteRepository) GetTasksByStatus(status int) ([]DatabaseTask, error) {
	query := `
	SELECT id, title, description, priority, status, created_at, updated_at, due_date, user_id, category_id, is_archived
	FROM tasks 
	WHERE status = ? AND is_archived = FALSE
	ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}
	defer rows.Close()
	
	var tasks []DatabaseTask
	for rows.Next() {
		task := DatabaseTask{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DueDate,
			&task.UserID,
			&task.CategoryID,
			&task.IsArchived,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

func (r *SQLiteRepository) GetTasksByPriority(priority int) ([]DatabaseTask, error) {
	query := `
	SELECT id, title, description, priority, status, created_at, updated_at, due_date, user_id, category_id, is_archived
	FROM tasks 
	WHERE priority = ? AND is_archived = FALSE
	ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, priority)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by priority: %w", err)
	}
	defer rows.Close()
	
	var tasks []DatabaseTask
	for rows.Next() {
		task := DatabaseTask{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DueDate,
			&task.UserID,
			&task.CategoryID,
			&task.IsArchived,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

func (r *SQLiteRepository) GetOverdueTasks() ([]DatabaseTask, error) {
	query := `
	SELECT id, title, description, priority, status, created_at, updated_at, due_date, user_id, category_id, is_archived
	FROM tasks 
	WHERE due_date < ? AND status != 2 AND is_archived = FALSE
	ORDER BY due_date ASC`
	
	now := time.Now()
	rows, err := r.db.Query(query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}
	defer rows.Close()
	
	var tasks []DatabaseTask
	for rows.Next() {
		task := DatabaseTask{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DueDate,
			&task.UserID,
			&task.CategoryID,
			&task.IsArchived,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// Placeholder implementations for future phases
// These will be implemented when we reach the respective phases

func (r *SQLiteRepository) GetTasksByUser(userID int) ([]DatabaseTask, error) {
	// TODO: Implement in Phase 5
	return nil, fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) GetTasksByCategory(categoryID int) ([]DatabaseTask, error) {
	// TODO: Implement in Phase 2
	return nil, fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) SearchTasks(query string) ([]DatabaseTask, error) {
	// TODO: Implement in Phase 6
	return nil, fmt.Errorf("not implemented yet")
}

// Category operations (Phase 2)
func (r *SQLiteRepository) CreateCategory(category *Category) error {
	query := `
	INSERT INTO categories (name, description, color)
	VALUES (?, ?, ?)`
	
	result, err := r.db.Exec(query, 
		category.Name, 
		category.Description, 
		category.Color)
	
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get category ID: %w", err)
	}
	
	category.ID = int(id)
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	
	return nil
}

func (r *SQLiteRepository) GetCategory(id int) (*Category, error) {
	query := `
	SELECT id, name, description, color, created_at, updated_at
	FROM categories WHERE id = ?`
	
	category := &Category{}
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Color,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	
	return category, nil
}

func (r *SQLiteRepository) GetAllCategories() ([]Category, error) {
	query := `
	SELECT id, name, description, color, created_at, updated_at
	FROM categories 
	ORDER BY name ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()
	
	var categories []Category
	for rows.Next() {
		category := Category{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	
	return categories, nil
}

func (r *SQLiteRepository) UpdateCategory(category *Category) error {
	query := `
	UPDATE categories 
	SET name = ?, description = ?, color = ?, updated_at = ?
	WHERE id = ?`
	
	category.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(query,
		category.Name,
		category.Description,
		category.Color,
		category.UpdatedAt,
		category.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", category.ID)
	}
	
	return nil
}

func (r *SQLiteRepository) DeleteCategory(id int) error {
	query := `DELETE FROM categories WHERE id = ?`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", id)
	}
	
	return nil
}

// Tag operations (Phase 2)
func (r *SQLiteRepository) CreateTag(tag *Tag) error {
	query := `
	INSERT INTO tags (name, color)
	VALUES (?, ?)`
	
	result, err := r.db.Exec(query, 
		tag.Name, 
		tag.Color)
	
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get tag ID: %w", err)
	}
	
	tag.ID = int(id)
	tag.CreatedAt = time.Now()
	
	return nil
}

func (r *SQLiteRepository) GetTag(id int) (*Tag, error) {
	query := `
	SELECT id, name, color, created_at
	FROM tags WHERE id = ?`
	
	tag := &Tag{}
	err := r.db.QueryRow(query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Color,
		&tag.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	
	return tag, nil
}

func (r *SQLiteRepository) GetAllTags() ([]Tag, error) {
	query := `
	SELECT id, name, color, created_at
	FROM tags 
	ORDER BY name ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()
	
	var tags []Tag
	for rows.Next() {
		tag := Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}
	
	return tags, nil
}

func (r *SQLiteRepository) UpdateTag(tag *Tag) error {
	query := `
	UPDATE tags 
	SET name = ?, color = ?
	WHERE id = ?`
	
	result, err := r.db.Exec(query,
		tag.Name,
		tag.Color,
		tag.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("tag with ID %d not found", tag.ID)
	}
	
	return nil
}

func (r *SQLiteRepository) DeleteTag(id int) error {
	query := `DELETE FROM tags WHERE id = ?`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("tag with ID %d not found", id)
	}
	
	return nil
}

// Task-Tag relationship operations (Phase 2)
func (r *SQLiteRepository) AddTagToTask(taskID, tagID int) error {
	query := `
	INSERT INTO task_tags (task_id, tag_id)
	VALUES (?, ?)`
	
	_, err := r.db.Exec(query, taskID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to task: %w", err)
	}
	
	return nil
}

func (r *SQLiteRepository) RemoveTagFromTask(taskID, tagID int) error {
	query := `DELETE FROM task_tags WHERE task_id = ? AND tag_id = ?`
	
	result, err := r.db.Exec(query, taskID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from task: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("tag-task relationship not found")
	}
	
	return nil
}

func (r *SQLiteRepository) GetTaskTags(taskID int) ([]Tag, error) {
	query := `
	SELECT t.id, t.name, t.color, t.created_at
	FROM tags t
	INNER JOIN task_tags tt ON t.id = tt.tag_id
	WHERE tt.task_id = ?
	ORDER BY t.name ASC`
	
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task tags: %w", err)
	}
	defer rows.Close()
	
	var tags []Tag
	for rows.Next() {
		tag := Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}
	
	return tags, nil
}

func (r *SQLiteRepository) GetTasksByTag(tagID int) ([]DatabaseTask, error) {
	query := `
	SELECT t.id, t.title, t.description, t.priority, t.status, t.created_at, t.updated_at, t.due_date, t.user_id, t.category_id, t.is_archived
	FROM tasks t
	INNER JOIN task_tags tt ON t.id = tt.task_id
	WHERE tt.tag_id = ? AND t.is_archived = FALSE
	ORDER BY t.created_at DESC`
	
	rows, err := r.db.Query(query, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by tag: %w", err)
	}
	defer rows.Close()
	
	var tasks []DatabaseTask
	for rows.Next() {
		task := DatabaseTask{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DueDate,
			&task.UserID,
			&task.CategoryID,
			&task.IsArchived,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// User operations (Phase 5)
func (r *SQLiteRepository) CreateUser(user *User) error {
	// TODO: Implement in Phase 5
	return fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) GetUser(id int) (*User, error) {
	// TODO: Implement in Phase 5
	return nil, fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) GetUserByUsername(username string) (*User, error) {
	// TODO: Implement in Phase 5
	return nil, fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) GetUserByEmail(email string) (*User, error) {
	// TODO: Implement in Phase 5
	return nil, fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) UpdateUser(user *User) error {
	// TODO: Implement in Phase 5
	return fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) DeleteUser(id int) error {
	// TODO: Implement in Phase 5
	return fmt.Errorf("not implemented yet")
}

func (r *SQLiteRepository) GetAllUsers() ([]User, error) {
	// TODO: Implement in Phase 5
	return nil, fmt.Errorf("not implemented yet")
}
