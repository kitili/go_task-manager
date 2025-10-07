package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"learn-go-capstone/internal/database"
)

// ImportService handles data import operations
type ImportService struct {
	repository database.Repository
}

// NewImportService creates a new import service
func NewImportService(repository database.Repository) *ImportService {
	return &ImportService{
		repository: repository,
	}
}

// ImportTasks imports tasks from a file
func (is *ImportService) ImportTasks(filePath string, options ImportOptions) (*ImportResult, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Determine format from file extension
	format := is.detectFormat(filePath)
	if format == "" {
		return nil, fmt.Errorf("unsupported file format")
	}

	// Import based on format
	var result *ImportResult
	switch format {
	case FormatJSON:
		result, err = is.importFromJSON(file, options)
	case FormatCSV:
		result, err = is.importFromCSV(file, options)
	default:
		return nil, fmt.Errorf("unsupported import format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to import from %s: %w", format, err)
	}

	result.FileName = filepath.Base(filePath)
	result.Format = format
	result.ImportedAt = time.Now()
	result.Options = options

	return result, nil
}

// detectFormat detects the file format from the file extension
func (is *ImportService) detectFormat(filePath string) ExportFormat {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return FormatJSON
	case ".csv":
		return FormatCSV
	default:
		return ""
	}
}

// importFromJSON imports tasks from JSON format
func (is *ImportService) importFromJSON(file io.Reader, options ImportOptions) (*ImportResult, error) {
	var exportData ExportData
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&exportData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	result := &ImportResult{
		TotalRecords: len(exportData.Tasks),
		Imported:     0,
		Skipped:      0,
		Errors:       0,
		ErrorDetails: []ImportError{},
	}

	// Import categories first if they exist
	categoryMap := make(map[int]int) // old ID -> new ID
	if len(exportData.Categories) > 0 {
		for _, catExport := range exportData.Categories {
			category := &database.Category{
				Name:        catExport.Name,
				Description: catExport.Description,
				Color:       catExport.Color,
			}
			
			if err := is.repository.CreateCategory(category); err != nil {
				result.ErrorDetails = append(result.ErrorDetails, ImportError{
					Row:     0,
					Field:   "category",
					Value:   catExport.Name,
					Message: fmt.Sprintf("failed to create category: %v", err),
				})
				result.Errors++
				continue
			}
			
			categoryMap[catExport.ID] = category.ID
		}
	}

	// Import tags if they exist
	tagMap := make(map[int]int) // old ID -> new ID
	if len(exportData.Tags) > 0 {
		for _, tagExport := range exportData.Tags {
			tag := &database.Tag{
				Name:  tagExport.Name,
				Color: tagExport.Color,
			}
			
			if err := is.repository.CreateTag(tag); err != nil {
				result.ErrorDetails = append(result.ErrorDetails, ImportError{
					Row:     0,
					Field:   "tag",
					Value:   tagExport.Name,
					Message: fmt.Sprintf("failed to create tag: %v", err),
				})
				result.Errors++
				continue
			}
			
			tagMap[tagExport.ID] = tag.ID
		}
	}

	// Import users if they exist
	userMap := make(map[int]int) // old ID -> new ID
	if len(exportData.Users) > 0 {
		for _, userExport := range exportData.Users {
			user := &database.User{
				Username: userExport.Username,
				Email:    userExport.Email,
				Password: "imported_user_password", // Default password for imported users
				IsActive: userExport.IsActive,
			}
			
			if err := is.repository.CreateUser(user); err != nil {
				result.ErrorDetails = append(result.ErrorDetails, ImportError{
					Row:     0,
					Field:   "user",
					Value:   userExport.Username,
					Message: fmt.Sprintf("failed to create user: %v", err),
				})
				result.Errors++
				continue
			}
			
			userMap[userExport.ID] = user.ID
		}
	}

	// Import tasks
	for i, taskExport := range exportData.Tasks {
		// Validate task data
		if options.ValidateData {
			if err := is.validateTaskExport(taskExport); err != nil {
				result.ErrorDetails = append(result.ErrorDetails, ImportError{
					Row:     i + 1,
					Field:   "task",
					Value:   taskExport.Title,
					Message: fmt.Sprintf("validation failed: %v", err),
				})
				result.Errors++
				continue
			}
		}

		// Check for duplicates if skip duplicates is enabled
		if options.SkipDuplicates {
			// This would require a more sophisticated duplicate detection
			// For now, we'll skip this check
		}

		// Create task
		task := &database.DatabaseTask{
			Title:       taskExport.Title,
			Description: taskExport.Description,
			Priority:    taskExport.Priority,
			Status:      taskExport.Status,
			CreatedAt:   taskExport.CreatedAt,
			UpdatedAt:   taskExport.UpdatedAt,
			DueDate:     taskExport.DueDate,
			IsArchived:  taskExport.IsArchived,
		}

		// Map category ID
		if taskExport.CategoryID != nil {
			if newCategoryID, exists := categoryMap[*taskExport.CategoryID]; exists {
				task.CategoryID = &newCategoryID
			}
		}

		// Map user ID
		if taskExport.UserID != nil {
			if newUserID, exists := userMap[*taskExport.UserID]; exists {
				task.UserID = &newUserID
			}
		}

		// Create task
		if err := is.repository.CreateTask(task); err != nil {
			result.ErrorDetails = append(result.ErrorDetails, ImportError{
				Row:     i + 1,
				Field:   "task",
				Value:   taskExport.Title,
				Message: fmt.Sprintf("failed to create task: %v", err),
			})
			result.Errors++
			continue
		}

		// Import task tags if they exist
		if len(taskExport.Tags) > 0 {
			for _, tagExport := range taskExport.Tags {
				if newTagID, exists := tagMap[tagExport.ID]; exists {
					if err := is.repository.AddTagToTask(task.ID, newTagID); err != nil {
						// Log error but don't fail the entire import
						result.ErrorDetails = append(result.ErrorDetails, ImportError{
							Row:     i + 1,
							Field:   "tag",
							Value:   tagExport.Name,
							Message: fmt.Sprintf("failed to add tag to task: %v", err),
						})
						result.Errors++
					}
				}
			}
		}

		result.Imported++
	}

	return result, nil
}

// importFromCSV imports tasks from CSV format
func (is *ImportService) importFromCSV(file io.Reader, options ImportOptions) (*ImportResult, error) {
	reader := csv.NewReader(file)
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Validate header
	expectedHeader := []string{
		"id", "title", "description", "priority", "status", "created_at", "updated_at",
		"due_date", "user_id", "category_id", "is_archived",
	}
	
	if len(header) < len(expectedHeader) {
		return nil, fmt.Errorf("invalid CSV header: expected at least %d columns, got %d", len(expectedHeader), len(header))
	}

	result := &ImportResult{
		TotalRecords: 0,
		Imported:     0,
		Skipped:      0,
		Errors:       0,
		ErrorDetails: []ImportError{},
	}

	// Read data rows
	rowNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.ErrorDetails = append(result.ErrorDetails, ImportError{
				Row:     rowNum,
				Message: fmt.Sprintf("failed to read row: %v", err),
			})
			result.Errors++
			rowNum++
			continue
		}

		result.TotalRecords++

		// Parse task from CSV record
		task, parseErr := is.parseTaskFromCSV(record, rowNum)
		if parseErr != nil {
			result.ErrorDetails = append(result.ErrorDetails, *parseErr)
			result.Errors++
			rowNum++
			continue
		}

		// Validate task data
		if options.ValidateData {
			if err := is.validateTask(task); err != nil {
				result.ErrorDetails = append(result.ErrorDetails, ImportError{
					Row:     rowNum,
					Field:   "task",
					Value:   task.Title,
					Message: fmt.Sprintf("validation failed: %v", err),
				})
				result.Errors++
				rowNum++
				continue
			}
		}

		// Create task
		if err := is.repository.CreateTask(task); err != nil {
			result.ErrorDetails = append(result.ErrorDetails, ImportError{
				Row:     rowNum,
				Field:   "task",
				Value:   task.Title,
				Message: fmt.Sprintf("failed to create task: %v", err),
			})
			result.Errors++
			rowNum++
			continue
		}

		result.Imported++
		rowNum++
	}

	return result, nil
}

// parseTaskFromCSV parses a task from CSV record
func (is *ImportService) parseTaskFromCSV(record []string, rowNum int) (*database.DatabaseTask, *ImportError) {
	if len(record) < 11 {
		return nil, &ImportError{
			Row:     rowNum,
			Message: "insufficient columns in CSV record",
		}
	}

	task := &database.DatabaseTask{}

	// Parse ID (skip for new tasks)
	// task.ID will be set by the database

	// Parse title
	task.Title = record[1]
	if task.Title == "" {
		return nil, &ImportError{
			Row:     rowNum,
			Field:   "title",
			Message: "title cannot be empty",
		}
	}

	// Parse description
	task.Description = record[2]

	// Parse priority
	priority, err := strconv.Atoi(record[3])
	if err != nil {
		return nil, &ImportError{
			Row:     rowNum,
			Field:   "priority",
			Value:   record[3],
			Message: "invalid priority value",
		}
	}
	task.Priority = priority

	// Parse status
	status, err := strconv.Atoi(record[4])
	if err != nil {
		return nil, &ImportError{
			Row:     rowNum,
			Field:   "status",
			Value:   record[4],
			Message: "invalid status value",
		}
	}
	task.Status = status

	// Parse created_at
	createdAt, err := time.Parse(time.RFC3339, record[5])
	if err != nil {
		return nil, &ImportError{
			Row:     rowNum,
			Field:   "created_at",
			Value:   record[5],
			Message: "invalid created_at format",
		}
	}
	task.CreatedAt = createdAt

	// Parse updated_at
	updatedAt, err := time.Parse(time.RFC3339, record[6])
	if err != nil {
		return nil, &ImportError{
			Row:     rowNum,
			Field:   "updated_at",
			Value:   record[6],
			Message: "invalid updated_at format",
		}
	}
	task.UpdatedAt = updatedAt

	// Parse due_date (optional)
	if record[7] != "" {
		dueDate, err := time.Parse(time.RFC3339, record[7])
		if err != nil {
			return nil, &ImportError{
				Row:     rowNum,
				Field:   "due_date",
				Value:   record[7],
				Message: "invalid due_date format",
			}
		}
		task.DueDate = &dueDate
	}

	// Parse user_id (optional)
	if record[8] != "" {
		userID, err := strconv.Atoi(record[8])
		if err != nil {
			return nil, &ImportError{
				Row:     rowNum,
				Field:   "user_id",
				Value:   record[8],
				Message: "invalid user_id value",
			}
		}
		task.UserID = &userID
	}

	// Parse category_id (optional)
	if record[9] != "" {
		categoryID, err := strconv.Atoi(record[9])
		if err != nil {
			return nil, &ImportError{
				Row:     rowNum,
				Field:   "category_id",
				Value:   record[9],
				Message: "invalid category_id value",
			}
		}
		task.CategoryID = &categoryID
	}

	// Parse is_archived
	isArchived, err := strconv.ParseBool(record[10])
	if err != nil {
		return nil, &ImportError{
			Row:     rowNum,
			Field:   "is_archived",
			Value:   record[10],
			Message: "invalid is_archived value",
		}
	}
	task.IsArchived = isArchived

	return task, nil
}

// validateTaskExport validates a task export
func (is *ImportService) validateTaskExport(task TaskExport) error {
	if task.Title == "" {
		return fmt.Errorf("title is required")
	}
	if task.Priority < 1 || task.Priority > 4 {
		return fmt.Errorf("priority must be between 1 and 4")
	}
	if task.Status < 0 || task.Status > 3 {
		return fmt.Errorf("status must be between 0 and 3")
	}
	return nil
}

// validateTask validates a database task
func (is *ImportService) validateTask(task *database.DatabaseTask) error {
	if task.Title == "" {
		return fmt.Errorf("title is required")
	}
	if task.Priority < 1 || task.Priority > 4 {
		return fmt.Errorf("priority must be between 1 and 4")
	}
	if task.Status < 0 || task.Status > 3 {
		return fmt.Errorf("status must be between 0 and 3")
	}
	return nil
}
