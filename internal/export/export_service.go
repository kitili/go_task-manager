package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"learn-go-capstone/internal/database"
)

// ExportService handles data export operations
type ExportService struct {
	repository database.Repository
}

// NewExportService creates a new export service
func NewExportService(repository database.Repository) *ExportService {
	return &ExportService{
		repository: repository,
	}
}

// ExportTasks exports tasks to a file
func (es *ExportService) ExportTasks(options ExportOptions) (*ExportResult, error) {
	// Get tasks based on filters
	tasks, err := es.getFilteredTasks(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Convert to export format
	exportData, err := es.convertToExportData(tasks, options)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to export format: %w", err)
	}

	// Generate filename
	fileName := es.generateFileName(options.Format)

	// Export based on format
	var filePath string
	var fileSize int64

	switch options.Format {
	case FormatJSON:
		filePath, fileSize, err = es.exportToJSON(exportData, fileName)
	case FormatCSV:
		filePath, fileSize, err = es.exportToCSV(exportData, fileName)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", options.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export to %s: %w", options.Format, err)
	}

	return &ExportResult{
		Format:      options.Format,
		FileName:    fileName,
		FilePath:    filePath,
		RecordCount: len(exportData.Tasks),
		FileSize:    fileSize,
		ExportedAt:  time.Now(),
		Options:     options,
	}, nil
}

// getFilteredTasks gets tasks based on export options
func (es *ExportService) getFilteredTasks(options ExportOptions) ([]database.DatabaseTask, error) {
	// For now, get all tasks and filter in memory
	// In a real implementation, you'd build a query based on options
	allTasks, err := es.repository.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var filteredTasks []database.DatabaseTask
	for _, task := range allTasks {
		// Apply filters
		if options.UserID != nil && (task.UserID == nil || *task.UserID != *options.UserID) {
			continue
		}
		if options.Status != nil && task.Status != *options.Status {
			continue
		}
		if options.Priority != nil && task.Priority != *options.Priority {
			continue
		}
		if options.DateFrom != nil && task.CreatedAt.Before(*options.DateFrom) {
			continue
		}
		if options.DateTo != nil && task.CreatedAt.After(*options.DateTo) {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, nil
}

// convertToExportData converts database tasks to export format
func (es *ExportService) convertToExportData(tasks []database.DatabaseTask, options ExportOptions) (*ExportData, error) {
	exportData := &ExportData{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Tasks:      make([]TaskExport, 0, len(tasks)),
		Metadata: ExportMetadata{
			TotalTasks: len(tasks),
			ExportOptions: options,
		},
	}

	// Convert tasks
	for _, task := range tasks {
		taskExport := TaskExport{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			DueDate:     task.DueDate,
			UserID:      task.UserID,
			CategoryID:  task.CategoryID,
			IsArchived:  task.IsArchived,
		}

		// Add category if requested and available
		if options.IncludeCategories && task.CategoryID != nil {
			category, err := es.repository.GetCategory(*task.CategoryID)
			if err == nil {
				taskExport.Category = &CategoryExport{
					ID:          category.ID,
					Name:        category.Name,
					Description: category.Description,
					Color:       category.Color,
					CreatedAt:   category.CreatedAt,
					UpdatedAt:   category.UpdatedAt,
				}
			}
		}

		// Add user if requested and available
		if options.IncludeUsers && task.UserID != nil {
			user, err := es.repository.GetUser(*task.UserID)
			if err == nil {
				taskExport.User = &UserExport{
					ID:        user.ID,
					Username:  user.Username,
					Email:     user.Email,
					IsActive:  user.IsActive,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				}
			}
		}

		// Add tags if requested
		if options.IncludeTags {
			// Note: This would require additional queries to get task tags
			// For now, we'll leave it empty
			taskExport.Tags = []TagExport{}
		}

		exportData.Tasks = append(exportData.Tasks, taskExport)
	}

	// Add categories if requested
	if options.IncludeCategories {
		categories, err := es.repository.GetAllCategories()
		if err == nil {
			exportData.Categories = make([]CategoryExport, 0, len(categories))
			for _, cat := range categories {
				exportData.Categories = append(exportData.Categories, CategoryExport{
					ID:          cat.ID,
					Name:        cat.Name,
					Description: cat.Description,
					Color:       cat.Color,
					CreatedAt:   cat.CreatedAt,
					UpdatedAt:   cat.UpdatedAt,
				})
			}
			exportData.Metadata.TotalCategories = len(categories)
		}
	}

	// Add tags if requested
	if options.IncludeTags {
		tags, err := es.repository.GetAllTags()
		if err == nil {
			exportData.Tags = make([]TagExport, 0, len(tags))
			for _, tag := range tags {
				exportData.Tags = append(exportData.Tags, TagExport{
					ID:        tag.ID,
					Name:      tag.Name,
					Color:     tag.Color,
					CreatedAt: tag.CreatedAt,
				})
			}
			exportData.Metadata.TotalTags = len(tags)
		}
	}

	// Add users if requested
	if options.IncludeUsers {
		users, err := es.repository.GetAllUsers()
		if err == nil {
			exportData.Users = make([]UserExport, 0, len(users))
			for _, user := range users {
				exportData.Users = append(exportData.Users, UserExport{
					ID:        user.ID,
					Username:  user.Username,
					Email:     user.Email,
					IsActive:  user.IsActive,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				})
			}
			exportData.Metadata.TotalUsers = len(users)
		}
	}

	return exportData, nil
}

// generateFileName generates a filename for the export
func (es *ExportService) generateFileName(format ExportFormat) string {
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("tasks_export_%s.%s", timestamp, format)
}

// exportToJSON exports data to JSON format
func (es *ExportService) exportToJSON(data *ExportData, fileName string) (string, int64, error) {
	filePath := filepath.Join("exports", fileName)
	
	// Create exports directory if it doesn't exist
	if err := os.MkdirAll("exports", 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create exports directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(data); err != nil {
		file.Close()
		return "", 0, fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Close file before getting size
	file.Close()

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return filePath, fileInfo.Size(), nil
}

// exportToCSV exports data to CSV format
func (es *ExportService) exportToCSV(data *ExportData, fileName string) (string, int64, error) {
	filePath := filepath.Join("exports", fileName)
	
	// Create exports directory if it doesn't exist
	if err := os.MkdirAll("exports", 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create exports directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create file: %w", err)
	}

	writer := csv.NewWriter(file)

	// Write header
	header := []string{
		"id", "title", "description", "priority", "status", "created_at", "updated_at",
		"due_date", "user_id", "category_id", "is_archived",
	}
	if err := writer.Write(header); err != nil {
		return "", 0, fmt.Errorf("failed to write header: %w", err)
	}

	// Write data rows
	for _, task := range data.Tasks {
		row := []string{
			strconv.Itoa(task.ID),
			task.Title,
			task.Description,
			strconv.Itoa(task.Priority),
			strconv.Itoa(task.Status),
			task.CreatedAt.Format(time.RFC3339),
			task.UpdatedAt.Format(time.RFC3339),
		}

		// Add due date
		if task.DueDate != nil {
			row = append(row, task.DueDate.Format(time.RFC3339))
		} else {
			row = append(row, "")
		}

		// Add user ID
		if task.UserID != nil {
			row = append(row, strconv.Itoa(*task.UserID))
		} else {
			row = append(row, "")
		}

		// Add category ID
		if task.CategoryID != nil {
			row = append(row, strconv.Itoa(*task.CategoryID))
		} else {
			row = append(row, "")
		}

		// Add archived status
		row = append(row, strconv.FormatBool(task.IsArchived))

		if err := writer.Write(row); err != nil {
			file.Close()
			return "", 0, fmt.Errorf("failed to write row: %w", err)
		}
	}

	// Flush writer and close file
	writer.Flush()
	file.Close()

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return filePath, fileInfo.Size(), nil
}

// GetExportHistory returns a list of export files
func (es *ExportService) GetExportHistory() ([]ExportResult, error) {
	// This would typically read from a database or file system
	// For now, return empty list
	return []ExportResult{}, nil
}

// CleanupOldExports removes old export files
func (es *ExportService) CleanupOldExports(olderThan time.Duration) error {
	// This would typically clean up old files
	// For now, return nil
	return nil
}
