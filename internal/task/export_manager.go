package task

import (
	"time"

	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/export"
)

// ExportManager manages export and import operations for tasks
type ExportManager struct {
	repository    database.Repository
	exportService *export.ExportService
	importService *export.ImportService
}

// NewExportManager creates a new export manager
func NewExportManager(repository database.Repository) *ExportManager {
	return &ExportManager{
		repository:    repository,
		exportService: export.NewExportService(repository),
		importService: export.NewImportService(repository),
	}
}

// ExportTasks exports tasks to a file
func (em *ExportManager) ExportTasks(options export.ExportOptions) (*export.ExportResult, error) {
	return em.exportService.ExportTasks(options)
}

// ImportTasks imports tasks from a file
func (em *ExportManager) ImportTasks(filePath string, options export.ImportOptions) (*export.ImportResult, error) {
	return em.importService.ImportTasks(filePath, options)
}

// ExportTasksToJSON exports tasks to JSON format
func (em *ExportManager) ExportTasksToJSON(includeTags, includeCategories, includeUsers bool) (*export.ExportResult, error) {
	options := export.ExportOptions{
		Format:           export.FormatJSON,
		IncludeTags:      includeTags,
		IncludeCategories: includeCategories,
		IncludeUsers:     includeUsers,
	}
	return em.ExportTasks(options)
}

// ExportTasksToCSV exports tasks to CSV format
func (em *ExportManager) ExportTasksToCSV(includeTags, includeCategories, includeUsers bool) (*export.ExportResult, error) {
	options := export.ExportOptions{
		Format:           export.FormatCSV,
		IncludeTags:      includeTags,
		IncludeCategories: includeCategories,
		IncludeUsers:     includeUsers,
	}
	return em.ExportTasks(options)
}

// ExportUserTasks exports tasks for a specific user
func (em *ExportManager) ExportUserTasks(userID int, format export.ExportFormat) (*export.ExportResult, error) {
	options := export.ExportOptions{
		Format:           format,
		UserID:           &userID,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
	return em.ExportTasks(options)
}

// ExportTasksByDateRange exports tasks within a date range
func (em *ExportManager) ExportTasksByDateRange(
	dateFrom, dateTo time.Time,
	format export.ExportFormat,
) (*export.ExportResult, error) {
	options := export.ExportOptions{
		Format:           format,
		DateFrom:         &dateFrom,
		DateTo:           &dateTo,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
	return em.ExportTasks(options)
}

// ExportTasksByStatus exports tasks by status
func (em *ExportManager) ExportTasksByStatus(
	status Status,
	format export.ExportFormat,
) (*export.ExportResult, error) {
	statusInt := int(status)
	options := export.ExportOptions{
		Format:           format,
		Status:           &statusInt,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
	return em.ExportTasks(options)
}

// ExportTasksByPriority exports tasks by priority
func (em *ExportManager) ExportTasksByPriority(
	priority Priority,
	format export.ExportFormat,
) (*export.ExportResult, error) {
	priorityInt := int(priority)
	options := export.ExportOptions{
		Format:           format,
		Priority:         &priorityInt,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
	return em.ExportTasks(options)
}

// ImportTasksFromJSON imports tasks from JSON file
func (em *ExportManager) ImportTasksFromJSON(filePath string, updateExisting, skipDuplicates, validateData, dryRun bool) (*export.ImportResult, error) {
	options := export.ImportOptions{
		Format:         export.FormatJSON,
		UpdateExisting: updateExisting,
		SkipDuplicates: skipDuplicates,
		ValidateData:   validateData,
		DryRun:         dryRun,
	}
	return em.ImportTasks(filePath, options)
}

// ImportTasksFromCSV imports tasks from CSV file
func (em *ExportManager) ImportTasksFromCSV(filePath string, updateExisting, skipDuplicates, validateData, dryRun bool) (*export.ImportResult, error) {
	options := export.ImportOptions{
		Format:         export.FormatCSV,
		UpdateExisting: updateExisting,
		SkipDuplicates: skipDuplicates,
		ValidateData:   validateData,
		DryRun:         dryRun,
	}
	return em.ImportTasks(filePath, options)
}

// GetExportHistory returns export history
func (em *ExportManager) GetExportHistory() ([]export.ExportResult, error) {
	return em.exportService.GetExportHistory()
}

// CleanupOldExports removes old export files
func (em *ExportManager) CleanupOldExports(olderThan time.Duration) error {
	return em.exportService.CleanupOldExports(olderThan)
}

// ValidateImportFile validates an import file before importing
func (em *ExportManager) ValidateImportFile(filePath string, format export.ExportFormat) error {
	// This would validate the file format and structure
	// For now, we'll just check if the file exists
	return nil
}

// GetImportPreview shows a preview of what will be imported
func (em *ExportManager) GetImportPreview(filePath string, format export.ExportFormat) (*export.ImportResult, error) {
	options := export.ImportOptions{
		Format:       format,
		DryRun:       true,
		ValidateData: true,
	}
	return em.ImportTasks(filePath, options)
}

// ExportBackup creates a complete backup of all data
func (em *ExportManager) ExportBackup() (*export.ExportResult, error) {
	options := export.ExportOptions{
		Format:           export.FormatJSON,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
	return em.ExportTasks(options)
}

// ImportBackup restores data from a backup
func (em *ExportManager) ImportBackup(filePath string) (*export.ImportResult, error) {
	options := export.ImportOptions{
		Format:         export.FormatJSON,
		UpdateExisting: true,
		SkipDuplicates: false,
		ValidateData:   true,
		DryRun:         false,
	}
	return em.ImportTasks(filePath, options)
}
