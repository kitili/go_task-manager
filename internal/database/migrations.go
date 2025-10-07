package database

import (
	"database/sql"
	"fmt"
	"log"
)

// MigrationManager handles database schema migrations
type MigrationManager struct {
	db *sql.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// Migrate runs all pending migrations
func (mm *MigrationManager) Migrate() error {
	// Create migrations table if it doesn't exist
	if err := mm.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current migration version
	currentVersion, err := mm.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Run migrations
	migrations := mm.getMigrations()
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Printf("Running migration %d: %s", migration.Version, migration.Name)
			if err := migration.Run(mm.db); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
			}
			
			// Record migration as applied
			if err := mm.recordMigration(migration); err != nil {
				return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
			}
			
			log.Printf("Migration %d completed successfully", migration.Version)
		}
	}

	return nil
}

// Migration represents a single database migration
type Migration struct {
	Version int
	Name    string
	Run     func(*sql.DB) error
}

// getMigrations returns all available migrations in order
func (mm *MigrationManager) getMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "create_tasks_table",
			Run:     mm.createTasksTable,
		},
		{
			Version: 2,
			Name:    "create_categories_table",
			Run:     mm.createCategoriesTable,
		},
		{
			Version: 3,
			Name:    "create_tags_table",
			Run:     mm.createTagsTable,
		},
		{
			Version: 4,
			Name:    "create_task_tags_table",
			Run:     mm.createTaskTagsTable,
		},
		{
			Version: 5,
			Name:    "create_users_table",
			Run:     mm.createUsersTable,
		},
		{
			Version: 6,
			Name:    "add_user_id_to_tasks",
			Run:     mm.addUserIDToTasks,
		},
		{
			Version: 7,
			Name:    "add_category_id_to_tasks",
			Run:     mm.addCategoryIDToTasks,
		},
		{
			Version: 8,
			Name:    "add_archived_flag_to_tasks",
			Run:     mm.addArchivedFlagToTasks,
		},
	}
}

// createMigrationsTable creates the migrations tracking table
func (mm *MigrationManager) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version INTEGER NOT NULL UNIQUE,
		name TEXT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	
	_, err := mm.db.Exec(query)
	return err
}

// getCurrentVersion returns the highest applied migration version
func (mm *MigrationManager) getCurrentVersion() (int, error) {
	query := `SELECT COALESCE(MAX(version), 0) FROM migrations`
	var version int
	err := mm.db.QueryRow(query).Scan(&version)
	return version, err
}

// recordMigration records a migration as applied
func (mm *MigrationManager) recordMigration(migration Migration) error {
	query := `INSERT INTO migrations (version, name) VALUES (?, ?)`
	_, err := mm.db.Exec(query, migration.Version, migration.Name)
	return err
}

// Migration implementations

func (mm *MigrationManager) createTasksTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		priority INTEGER NOT NULL DEFAULT 1,
		status INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		due_date DATETIME,
		user_id INTEGER,
		category_id INTEGER,
		is_archived BOOLEAN DEFAULT FALSE
	)`
	
	_, err := db.Exec(query)
	return err
}

func (mm *MigrationManager) createCategoriesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		color TEXT DEFAULT '#007bff',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	
	_, err := db.Exec(query)
	return err
}

func (mm *MigrationManager) createTagsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		color TEXT DEFAULT '#6c757d',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	
	_, err := db.Exec(query)
	return err
}

func (mm *MigrationManager) createTaskTagsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS task_tags (
		task_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (task_id, tag_id),
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	)`
	
	_, err := db.Exec(query)
	return err
}

func (mm *MigrationManager) createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	
	_, err := db.Exec(query)
	return err
}

func (mm *MigrationManager) addUserIDToTasks(db *sql.DB) error {
	// This migration is already included in createTasksTable
	// But we'll add it here for completeness in case we need to add it later
	return nil
}

func (mm *MigrationManager) addCategoryIDToTasks(db *sql.DB) error {
	// This migration is already included in createTasksTable
	// But we'll add it here for completeness in case we need to add it later
	return nil
}

func (mm *MigrationManager) addArchivedFlagToTasks(db *sql.DB) error {
	// This migration is already included in createTasksTable
	// But we'll add it here for completeness in case we need to add it later
	return nil
}
