package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	// Database configuration
	Database DatabaseConfig
	
	// Application configuration
	App AppConfig
	
	// Feature flags
	Features FeatureFlags
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Driver     string
	DSN        string
	FilePath   string // For SQLite
	Host       string // For PostgreSQL
	Port       int    // For PostgreSQL
	User       string // For PostgreSQL
	Password   string // For PostgreSQL
	DBName     string // For PostgreSQL
	SSLMode    string // For PostgreSQL
	StorageType string // memory, database, hybrid
}

// AppConfig holds application-related configuration
type AppConfig struct {
	Port        string
	Environment string
	LogLevel    string
}

// FeatureFlags holds feature toggle configuration
type FeatureFlags struct {
	DatabaseEnabled    bool
	CategoriesEnabled  bool
	TagsEnabled        bool
	UsersEnabled       bool
	SearchEnabled      bool
	NotificationsEnabled bool
	APIEnabled         bool
}

// LoadConfig loads configuration from environment variables and defaults
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Driver:      getEnv("DB_DRIVER", "sqlite3"),
			DSN:         getEnv("DB_DSN", ""),
			FilePath:    getEnv("DB_FILE_PATH", "data/tasks.db"),
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnvAsInt("DB_PORT", 5432),
			User:        getEnv("DB_USER", ""),
			Password:    getEnv("DB_PASSWORD", ""),
			DBName:      getEnv("DB_NAME", "tasks"),
			SSLMode:     getEnv("DB_SSL_MODE", "disable"),
			StorageType: getEnv("STORAGE_TYPE", "memory"),
		},
		App: AppConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Features: FeatureFlags{
			DatabaseEnabled:     getEnvAsBool("FEATURE_DATABASE", true),
			CategoriesEnabled:   getEnvAsBool("FEATURE_CATEGORIES", false),
			TagsEnabled:         getEnvAsBool("FEATURE_TAGS", false),
			UsersEnabled:        getEnvAsBool("FEATURE_USERS", false),
			SearchEnabled:       getEnvAsBool("FEATURE_SEARCH", false),
			NotificationsEnabled: getEnvAsBool("FEATURE_NOTIFICATIONS", false),
			APIEnabled:          getEnvAsBool("FEATURE_API", true),
		},
	}
}

// GetStorageType returns the storage type as an enum
func (c *Config) GetStorageType() string {
	return c.Database.StorageType
}

// IsDatabaseEnabled returns true if database features are enabled
func (c *Config) IsDatabaseEnabled() bool {
	return c.Features.DatabaseEnabled && c.Database.StorageType != "memory"
}

// IsFeatureEnabled returns true if a specific feature is enabled
func (c *Config) IsFeatureEnabled(feature string) bool {
	switch feature {
	case "categories":
		return c.Features.CategoriesEnabled
	case "tags":
		return c.Features.TagsEnabled
	case "users":
		return c.Features.UsersEnabled
	case "search":
		return c.Features.SearchEnabled
	case "notifications":
		return c.Features.NotificationsEnabled
	case "api":
		return c.Features.APIEnabled
	default:
		return false
	}
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetDSN returns the appropriate DSN based on the database driver
func (dc *DatabaseConfig) GetDSN() string {
	if dc.DSN != "" {
		return dc.DSN
	}
	
	switch dc.Driver {
	case "sqlite3":
		return dc.FilePath
	case "postgres":
		return "host=" + dc.Host + " port=" + strconv.Itoa(dc.Port) + " user=" + dc.User + " password=" + dc.Password + " dbname=" + dc.DBName + " sslmode=" + dc.SSLMode
	default:
		return dc.FilePath
	}
}

// IsProduction returns true if the environment is production
func (ac *AppConfig) IsProduction() bool {
	return strings.ToLower(ac.Environment) == "production"
}

// IsDevelopment returns true if the environment is development
func (ac *AppConfig) IsDevelopment() bool {
	return strings.ToLower(ac.Environment) == "development"
}
