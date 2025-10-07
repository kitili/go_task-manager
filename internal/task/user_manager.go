package task

import (
	"errors"
	"time"

	"learn-go-capstone/internal/auth"
	"learn-go-capstone/internal/database"
)

// UserManager manages user operations and authentication
type UserManager struct {
	repository   database.Repository
	authService  *auth.AuthService
}

// NewUserManager creates a new user manager
func NewUserManager(repository database.Repository) *UserManager {
	return &UserManager{
		repository:  repository,
		authService: auth.NewAuthService(repository),
	}
}

// RegisterUser registers a new user
func (um *UserManager) RegisterUser(username, email, password string) (*database.User, error) {
	return um.authService.RegisterUser(username, email, password)
}

// LoginUser authenticates a user and returns a JWT token
func (um *UserManager) LoginUser(username, password string) (string, *database.User, error) {
	return um.authService.LoginUser(username, password)
}

// ValidateToken validates a JWT token
func (um *UserManager) ValidateToken(tokenString string) (*auth.Claims, error) {
	return um.authService.ValidateToken(tokenString)
}

// GetUserFromToken gets user information from a JWT token
func (um *UserManager) GetUserFromToken(tokenString string) (*database.User, error) {
	return um.authService.GetUserFromToken(tokenString)
}

// ChangePassword changes a user's password
func (um *UserManager) ChangePassword(userID int, oldPassword, newPassword string) error {
	return um.authService.ChangePassword(userID, oldPassword, newPassword)
}

// RefreshToken generates a new token for an existing user
func (um *UserManager) RefreshToken(tokenString string) (string, error) {
	return um.authService.RefreshToken(tokenString)
}

// DeactivateUser deactivates a user account
func (um *UserManager) DeactivateUser(userID int) error {
	return um.authService.DeactivateUser(userID)
}

// ActivateUser activates a user account
func (um *UserManager) ActivateUser(userID int) error {
	return um.authService.ActivateUser(userID)
}

// GetUser gets a user by ID
func (um *UserManager) GetUser(userID int) (*database.User, error) {
	return um.repository.GetUser(userID)
}

// GetUserByUsername gets a user by username
func (um *UserManager) GetUserByUsername(username string) (*database.User, error) {
	return um.repository.GetUserByUsername(username)
}

// GetUserByEmail gets a user by email
func (um *UserManager) GetUserByEmail(email string) (*database.User, error) {
	return um.repository.GetUserByEmail(email)
}

// UpdateUser updates a user
func (um *UserManager) UpdateUser(user *database.User) error {
	return um.repository.UpdateUser(user)
}

// DeleteUser deletes a user
func (um *UserManager) DeleteUser(userID int) error {
	return um.repository.DeleteUser(userID)
}

// GetAllUsers gets all users
func (um *UserManager) GetAllUsers() ([]database.User, error) {
	return um.repository.GetAllUsers()
}

// GetUserTasks gets all tasks for a specific user
func (um *UserManager) GetUserTasks(userID int) ([]Task, error) {
	dbTasks, err := um.repository.GetTasksByUser(userID)
	if err != nil {
		return nil, err
	}
	
	return convertFromDatabaseTasks(dbTasks), nil
}

// CreateUserTask creates a task for a specific user
func (um *UserManager) CreateUserTask(userID int, title, description string, priority Priority, dueDate *time.Time) (*Task, error) {
	now := time.Now()
	task := &database.DatabaseTask{
		Title:       title,
		Description: description,
		Priority:    int(priority),
		Status:      int(Pending),
		CreatedAt:   now,
		UpdatedAt:   now,
		DueDate:     dueDate,
		UserID:      &userID,
		IsArchived:  false,
	}
	
	err := um.repository.CreateTask(task)
	if err != nil {
		return nil, err
	}
	
	// Convert to task.Task
	convertedTask := convertFromDatabaseTask(task)
	return &convertedTask, nil
}

// UpdateUserTask updates a task for a specific user
func (um *UserManager) UpdateUserTask(userID, taskID int, title, description string, priority Priority, status Status, dueDate *time.Time) error {
	// Get the task first to verify ownership
	task, err := um.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	// Check if the task belongs to the user
	if task.UserID == nil || *task.UserID != userID {
		return errors.New("task not found or access denied")
	}
	
	// Update the task
	task.Title = title
	task.Description = description
	task.Priority = int(priority)
	task.Status = int(status)
	task.DueDate = dueDate
	task.UpdatedAt = time.Now()
	
	return um.repository.UpdateTask(task)
}

// DeleteUserTask deletes a task for a specific user
func (um *UserManager) DeleteUserTask(userID, taskID int) error {
	// Get the task first to verify ownership
	task, err := um.repository.GetTask(taskID)
	if err != nil {
		return err
	}
	
	// Check if the task belongs to the user
	if task.UserID == nil || *task.UserID != userID {
		return errors.New("task not found or access denied")
	}
	
	return um.repository.DeleteTask(taskID)
}

// GetUserTasksByStatus gets tasks for a user by status
func (um *UserManager) GetUserTasksByStatus(userID int, status Status) ([]Task, error) {
	dbTasks, err := um.repository.GetTasksByStatus(int(status))
	if err != nil {
		return nil, err
	}
	
	// Filter by user ID
	var userTasks []Task
	for _, dbTask := range dbTasks {
		if dbTask.UserID != nil && *dbTask.UserID == userID {
			userTasks = append(userTasks, convertFromDatabaseTask(&dbTask))
		}
	}
	
	return userTasks, nil
}

// GetUserTasksByPriority gets tasks for a user by priority
func (um *UserManager) GetUserTasksByPriority(userID int, priority Priority) ([]Task, error) {
	dbTasks, err := um.repository.GetTasksByPriority(int(priority))
	if err != nil {
		return nil, err
	}
	
	// Filter by user ID
	var userTasks []Task
	for _, dbTask := range dbTasks {
		if dbTask.UserID != nil && *dbTask.UserID == userID {
			userTasks = append(userTasks, convertFromDatabaseTask(&dbTask))
		}
	}
	
	return userTasks, nil
}

// GetUserOverdueTasks gets overdue tasks for a user
func (um *UserManager) GetUserOverdueTasks(userID int) ([]Task, error) {
	dbTasks, err := um.repository.GetOverdueTasks()
	if err != nil {
		return nil, err
	}
	
	// Filter by user ID
	var userTasks []Task
	for _, dbTask := range dbTasks {
		if dbTask.UserID != nil && *dbTask.UserID == userID {
			userTasks = append(userTasks, convertFromDatabaseTask(&dbTask))
		}
	}
	
	return userTasks, nil
}
