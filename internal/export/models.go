package export

import (
	"time"
)

// ExportFormat represents the supported export formats
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// ExportOptions represents options for exporting data
type ExportOptions struct {
	Format      ExportFormat `json:"format"`
	IncludeTags bool         `json:"include_tags"`
	IncludeCategories bool   `json:"include_categories"`
	IncludeUsers bool        `json:"include_users"`
	DateFrom    *time.Time   `json:"date_from,omitempty"`
	DateTo      *time.Time   `json:"date_to,omitempty"`
	UserID      *int         `json:"user_id,omitempty"`
	Status      *int         `json:"status,omitempty"`
	Priority    *int         `json:"priority,omitempty"`
}

// ImportOptions represents options for importing data
type ImportOptions struct {
	Format      ExportFormat `json:"format"`
	UpdateExisting bool      `json:"update_existing"`
	SkipDuplicates bool      `json:"skip_duplicates"`
	ValidateData   bool      `json:"validate_data"`
	DryRun        bool       `json:"dry_run"`
}

// ExportResult represents the result of an export operation
type ExportResult struct {
	Format      ExportFormat `json:"format"`
	FileName    string       `json:"file_name"`
	FilePath    string       `json:"file_path"`
	RecordCount int          `json:"record_count"`
	FileSize    int64        `json:"file_size"`
	ExportedAt  time.Time    `json:"exported_at"`
	Options     ExportOptions `json:"options"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Format        ExportFormat `json:"format"`
	FileName      string       `json:"file_name"`
	TotalRecords  int          `json:"total_records"`
	Imported      int          `json:"imported"`
	Skipped       int          `json:"skipped"`
	Errors        int          `json:"errors"`
	ImportedAt    time.Time    `json:"imported_at"`
	Options       ImportOptions `json:"options"`
	ErrorDetails  []ImportError `json:"error_details,omitempty"`
}

// ImportError represents an error during import
type ImportError struct {
	Row     int    `json:"row"`
	Field   string `json:"field,omitempty"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

// TaskExport represents a task in export format
type TaskExport struct {
	ID          int                 `json:"id" csv:"id"`
	Title       string              `json:"title" csv:"title"`
	Description string              `json:"description" csv:"description"`
	Priority    int                 `json:"priority" csv:"priority"`
	Status      int                 `json:"status" csv:"status"`
	CreatedAt   time.Time           `json:"created_at" csv:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" csv:"updated_at"`
	DueDate     *time.Time          `json:"due_date,omitempty" csv:"due_date"`
	UserID      *int                `json:"user_id,omitempty" csv:"user_id"`
	CategoryID  *int                `json:"category_id,omitempty" csv:"category_id"`
	IsArchived  bool                `json:"is_archived" csv:"is_archived"`
	Category    *CategoryExport     `json:"category,omitempty" csv:"-"`
	Tags        []TagExport         `json:"tags,omitempty" csv:"-"`
	User        *UserExport         `json:"user,omitempty" csv:"-"`
	Dependencies []int              `json:"dependencies,omitempty" csv:"dependencies"`
}

// CategoryExport represents a category in export format
type CategoryExport struct {
	ID          int       `json:"id" csv:"id"`
	Name        string    `json:"name" csv:"name"`
	Description string    `json:"description" csv:"description"`
	Color       string    `json:"color" csv:"color"`
	CreatedAt   time.Time `json:"created_at" csv:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" csv:"updated_at"`
}

// TagExport represents a tag in export format
type TagExport struct {
	ID        int       `json:"id" csv:"id"`
	Name      string    `json:"name" csv:"name"`
	Color     string    `json:"color" csv:"color"`
	CreatedAt time.Time `json:"created_at" csv:"created_at"`
}

// UserExport represents a user in export format
type UserExport struct {
	ID        int       `json:"id" csv:"id"`
	Username  string    `json:"username" csv:"username"`
	Email     string    `json:"email" csv:"email"`
	IsActive  bool      `json:"is_active" csv:"is_active"`
	CreatedAt time.Time `json:"created_at" csv:"created_at"`
	UpdatedAt time.Time `json:"updated_at" csv:"updated_at"`
}

// ExportData represents the complete export data structure
type ExportData struct {
	Version     string           `json:"version"`
	ExportedAt  time.Time        `json:"exported_at"`
	Tasks       []TaskExport     `json:"tasks"`
	Categories  []CategoryExport `json:"categories,omitempty"`
	Tags        []TagExport      `json:"tags,omitempty"`
	Users       []UserExport     `json:"users,omitempty"`
	Metadata    ExportMetadata   `json:"metadata"`
}

// ExportMetadata represents metadata about the export
type ExportMetadata struct {
	TotalTasks     int `json:"total_tasks"`
	TotalCategories int `json:"total_categories"`
	TotalTags      int `json:"total_tags"`
	TotalUsers     int `json:"total_users"`
	ExportOptions  ExportOptions `json:"export_options"`
}

// DefaultExportOptions returns default export options
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format:           FormatJSON,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
}

// DefaultImportOptions returns default import options
func DefaultImportOptions() ImportOptions {
	return ImportOptions{
		Format:         FormatJSON,
		UpdateExisting: false,
		SkipDuplicates: true,
		ValidateData:   true,
		DryRun:         false,
	}
}
