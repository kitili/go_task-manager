package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Config holds database configuration
type Config struct {
	Driver   string
	DSN      string
	FilePath string // For SQLite
	Host     string // For PostgreSQL
	Port     int    // For PostgreSQL
	User     string // For PostgreSQL
	Password string // For PostgreSQL
	DBName   string // For PostgreSQL
	SSLMode  string // For PostgreSQL
}

// DefaultConfig returns a default SQLite configuration
func DefaultConfig() *Config {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Warning: Could not create data directory: %v", err)
	}

	return &Config{
		Driver:   "sqlite3",
		FilePath: filepath.Join(dataDir, "tasks.db"),
		DSN:      filepath.Join(dataDir, "tasks.db"),
	}
}

// PostgreSQLConfig returns a PostgreSQL configuration
func PostgreSQLConfig(host, user, password, dbname string, port int) *Config {
	return &Config{
		Driver:   "postgres",
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
		SSLMode:  "disable", // Change to "require" in production
		DSN:      fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, "disable"),
	}
}

// Connect establishes a database connection
func Connect(config *Config) (*sql.DB, error) {
	var dsn string
	
	switch config.Driver {
	case "sqlite3":
		dsn = config.DSN
	case "postgres":
		dsn = config.DSN
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Printf("Successfully connected to %s database", config.Driver)
	return db, nil
}

// Close closes the database connection
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetDriver returns the database driver name
func (c *Config) GetDriver() string {
	return c.Driver
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return c.DSN
}
