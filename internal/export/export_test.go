package export

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
)

func TestExportService(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_export.db"
	defer os.Remove(tempDB)
	
	// Create database configuration
	config := &database.Config{
		Driver: "sqlite3",
		DSN:    tempDB,
	}
	
	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	
	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Create repository and export service
	repository := database.NewSQLiteRepository(db)
	exportService := NewExportService(repository)
	
	// Create test data
	user := &database.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		IsActive: true,
	}
	err = repository.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	category := &database.Category{
		Name:        "Work",
		Description: "Work-related tasks",
		Color:       "#ff0000",
	}
	err = repository.CreateCategory(category)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	
	tag := &database.Tag{
		Name:  "urgent",
		Color: "#ff0000",
	}
	err = repository.CreateTag(tag)
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}
	
	// Create test tasks
	now := time.Now()
	task1 := &database.DatabaseTask{
		Title:       "Test Task 1",
		Description: "Test Description 1",
		Priority:    3,
		Status:      0,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserID:      &user.ID,
		CategoryID:  &category.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(task1)
	if err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}
	
	task2 := &database.DatabaseTask{
		Title:       "Test Task 2",
		Description: "Test Description 2",
		Priority:    2,
		Status:      1,
		CreatedAt:   now.Add(-24 * time.Hour),
		UpdatedAt:   now.Add(-24 * time.Hour),
		UserID:      &user.ID,
		CategoryID:  &category.ID,
		IsArchived:  false,
	}
	err = repository.CreateTask(task2)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}
	
	// Add tag to task
	err = repository.AddTagToTask(task1.ID, tag.ID)
	if err != nil {
		t.Fatalf("Failed to add tag to task: %v", err)
	}
	
	// Test JSON export
	jsonOptions := ExportOptions{
		Format:           FormatJSON,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
	
	jsonResult, err := exportService.ExportTasks(jsonOptions)
	if err != nil {
		t.Fatalf("Failed to export to JSON: %v", err)
	}
	
	if jsonResult.Format != FormatJSON {
		t.Errorf("Expected JSON format, got %s", jsonResult.Format)
	}
	
	if jsonResult.RecordCount != 2 {
		t.Errorf("Expected 2 records, got %d", jsonResult.RecordCount)
	}
	
	if jsonResult.FileSize == 0 {
		t.Errorf("Expected file size to be greater than 0, got %d. File path: %s", jsonResult.FileSize, jsonResult.FilePath)
	}
	
	// Test CSV export
	csvOptions := ExportOptions{
		Format:           FormatCSV,
		IncludeTags:      true,
		IncludeCategories: true,
		IncludeUsers:     true,
	}
	
	csvResult, err := exportService.ExportTasks(csvOptions)
	if err != nil {
		t.Fatalf("Failed to export to CSV: %v", err)
	}
	
	if csvResult.Format != FormatCSV {
		t.Errorf("Expected CSV format, got %s", csvResult.Format)
	}
	
	if csvResult.RecordCount != 2 {
		t.Errorf("Expected 2 records, got %d", csvResult.RecordCount)
	}
	
	if csvResult.FileSize == 0 {
		t.Errorf("Expected file size to be greater than 0, got %d. File path: %s", csvResult.FileSize, csvResult.FilePath)
	}
	
	// Test filtered export
	filteredOptions := ExportOptions{
		Format:    FormatJSON,
		UserID:    &user.ID,
		Status:    func() *int { s := 0; return &s }(),
		Priority:  func() *int { p := 3; return &p }(),
	}
	
	filteredResult, err := exportService.ExportTasks(filteredOptions)
	if err != nil {
		t.Fatalf("Failed to export filtered tasks: %v", err)
	}
	
	if filteredResult.RecordCount != 1 {
		t.Errorf("Expected 1 filtered record, got %d", filteredResult.RecordCount)
	}
	
	// Clean up export files
	os.Remove(jsonResult.FilePath)
	os.Remove(csvResult.FilePath)
	os.Remove(filteredResult.FilePath)
}

func TestImportService(t *testing.T) {
	// Create a temporary database for testing
	tempDB := "test_import.db"
	defer os.Remove(tempDB)
	
	// Create database configuration
	config := &database.Config{
		Driver: "sqlite3",
		DSN:    tempDB,
	}
	
	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	
	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Create repository and import service
	repository := database.NewSQLiteRepository(db)
	importService := NewImportService(repository)
	
	// Create a test JSON file
	exportData := ExportData{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Tasks: []TaskExport{
			{
				ID:          1,
				Title:       "Imported Task 1",
				Description: "Imported Description 1",
				Priority:    3,
				Status:      0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				IsArchived:  false,
			},
			{
				ID:          2,
				Title:       "Imported Task 2",
				Description: "Imported Description 2",
				Priority:    2,
				Status:      1,
				CreatedAt:   time.Now().Add(-24 * time.Hour),
				UpdatedAt:   time.Now().Add(-24 * time.Hour),
				IsArchived:  false,
			},
		},
		Categories: []CategoryExport{
			{
				ID:          1,
				Name:        "Imported Category",
				Description: "Imported Category Description",
				Color:       "#00ff00",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		Tags: []TagExport{
			{
				ID:        1,
				Name:      "imported-tag",
				Color:     "#0000ff",
				CreatedAt: time.Now(),
			},
		},
		Users: []UserExport{
			{
				ID:        1,
				Username:  "importeduser",
				Email:     "imported@example.com",
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Metadata: ExportMetadata{
			TotalTasks:     2,
			TotalCategories: 1,
			TotalTags:      1,
			TotalUsers:     1,
		},
	}
	
	// Write test JSON file
	jsonFile := "test_import.json"
	file, err := os.Create(jsonFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(jsonFile)
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportData); err != nil {
		t.Fatalf("Failed to write test JSON: %v", err)
	}
	file.Close()
	
	// Test JSON import
	importOptions := ImportOptions{
		Format:         FormatJSON,
		UpdateExisting: false,
		SkipDuplicates: true,
		ValidateData:   true,
		DryRun:         false,
	}
	
	importResult, err := importService.ImportTasks(jsonFile, importOptions)
	if err != nil {
		t.Fatalf("Failed to import from JSON: %v", err)
	}
	
	if importResult.Format != FormatJSON {
		t.Errorf("Expected JSON format, got %s", importResult.Format)
	}
	
	if importResult.TotalRecords != 2 {
		t.Errorf("Expected 2 total records, got %d", importResult.TotalRecords)
	}
	
	if importResult.Imported != 2 {
		t.Errorf("Expected 2 imported records, got %d", importResult.Imported)
	}
	
	if importResult.Errors != 0 {
		t.Errorf("Expected 0 errors, got %d", importResult.Errors)
	}
	
	// Verify tasks were imported
	allTasks, err := repository.GetAllTasks()
	if err != nil {
		t.Fatalf("Failed to get all tasks: %v", err)
	}
	
	if len(allTasks) != 2 {
		t.Errorf("Expected 2 tasks in database, got %d", len(allTasks))
	}
	
	// Verify categories were imported
	categories, err := repository.GetAllCategories()
	if err != nil {
		t.Fatalf("Failed to get all categories: %v", err)
	}
	
	if len(categories) != 1 {
		t.Errorf("Expected 1 category in database, got %d", len(categories))
	}
	
	// Verify tags were imported
	tags, err := repository.GetAllTags()
	if err != nil {
		t.Fatalf("Failed to get all tags: %v", err)
	}
	
	if len(tags) != 1 {
		t.Errorf("Expected 1 tag in database, got %d", len(tags))
	}
	
	// Verify users were imported
	users, err := repository.GetAllUsers()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}
	
	if len(users) != 1 {
		t.Errorf("Expected 1 user in database, got %d", len(users))
	}
}

func TestExportOptions(t *testing.T) {
	// Test default export options
	options := DefaultExportOptions()
	if options.Format != FormatJSON {
		t.Errorf("Expected default format JSON, got %s", options.Format)
	}
	
	if !options.IncludeTags {
		t.Error("Expected default to include tags")
	}
	
	if !options.IncludeCategories {
		t.Error("Expected default to include categories")
	}
	
	if !options.IncludeUsers {
		t.Error("Expected default to include users")
	}
	
	// Test default import options
	importOptions := DefaultImportOptions()
	if importOptions.Format != FormatJSON {
		t.Errorf("Expected default format JSON, got %s", importOptions.Format)
	}
	
	if importOptions.UpdateExisting {
		t.Error("Expected default to not update existing")
	}
	
	if !importOptions.SkipDuplicates {
		t.Error("Expected default to skip duplicates")
	}
	
	if !importOptions.ValidateData {
		t.Error("Expected default to validate data")
	}
	
	if importOptions.DryRun {
		t.Error("Expected default to not be dry run")
	}
}
