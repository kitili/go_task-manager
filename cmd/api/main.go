package main

import (
	"database/sql"
	"log"
	"os"

	"learn-go-capstone/internal/api"
	"learn-go-capstone/internal/auth"
	"learn-go-capstone/internal/config"
	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/notifications"
	"learn-go-capstone/internal/task"
	_ "learn-go-capstone/docs" // Import docs for Swagger
)

// @title Go Task Manager API
// @version 1.0.0
// @description A comprehensive task management API built with Go, featuring user authentication, task management, categories, tags, dependencies, search, export/import, and notifications.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	log.Println("🚀 Starting Go Task Manager API Server...")

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("📋 Configuration loaded - Database enabled: %v", cfg.IsDatabaseEnabled())

	// Connect to database
	var db *sql.DB
	var repository database.Repository
	var err error

	if cfg.IsDatabaseEnabled() {
		config := &database.Config{
			Driver: "sqlite3",
			DSN:    "task_manager.db",
		}
		db, err = database.Connect(config)
		if err != nil {
			log.Fatalf("❌ Failed to connect to database: %v", err)
		}
		defer database.Close(db)

		// Run migrations
		migrationManager := database.NewMigrationManager(db)
		if err := migrationManager.Migrate(); err != nil {
			log.Fatalf("❌ Failed to run migrations: %v", err)
		}
		log.Println("✅ Database migrations completed successfully")

		// Create repository
		repository = database.NewSQLiteRepository(db)
		log.Println("✅ Database repository initialized")
	} else {
		log.Println("⚠️  Database disabled, using in-memory storage")
	}

	// Create task manager
	var taskManager task.TaskManagerInterface
	if cfg.IsDatabaseEnabled() && repository != nil {
		taskManager = task.NewHybridTaskManager(repository, task.DatabaseStorage)
		log.Println("✅ Hybrid task manager initialized (Database + Memory)")
	} else {
		taskManager = task.NewTaskManager()
		log.Println("✅ In-memory task manager initialized")
	}

	// Create managers
	userManager := task.NewUserManager(repository)
	categoryManager := task.NewCategoryManager(repository)
	dependencyManager := task.NewDependencyManager(repository)
	searchManager := task.NewSearchManager(repository)
	exportManager := task.NewExportManager(repository)

	// Create auth service
	authService := auth.NewAuthService(repository)

	// Create notification service
	notificationService := notifications.NewNotificationService(repository, notifications.DefaultNotificationConfig())
	notificationManager := task.NewNotificationManager(repository, notificationService)

	// Create API server
	server := api.NewServer(
		taskManager,
		userManager,
		categoryManager,
		dependencyManager,
		searchManager,
		exportManager,
		notificationManager,
		authService,
	)

	// Setup routes
	server.SetupRoutes()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	server.Run(port)
}
