package main

import (
	"log"
	"os"

	"learn-go-capstone/cmd"
	"learn-go-capstone/internal/config"
	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/task"
	"learn-go-capstone/internal/ui"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Initialize task manager based on configuration
	var taskManager task.TaskManagerInterface
	
	if cfg.IsDatabaseEnabled() {
		// Connect to database
		dbConfig := database.DefaultConfig()
		if cfg.Database.Driver == "postgres" {
			dbConfig = database.PostgreSQLConfig(
				cfg.Database.Host,
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.DBName,
				cfg.Database.Port,
			)
		}
		
		db, err := database.Connect(dbConfig)
		if err != nil {
			log.Printf("Warning: Failed to connect to database (%v), falling back to memory storage", err)
			taskManager = task.NewTaskManager()
		} else {
			defer database.Close(db)
			
			// Run migrations
			migrationManager := database.NewMigrationManager(db)
			if err := migrationManager.Migrate(); err != nil {
				log.Printf("Warning: Failed to run migrations (%v), falling back to memory storage", err)
				taskManager = task.NewTaskManager()
			} else {
				// Create hybrid task manager
				repository := database.NewSQLiteRepository(db)
				storageType := task.MemoryStorage
				
				switch cfg.Database.StorageType {
				case "database":
					storageType = task.DatabaseStorage
				case "hybrid":
					storageType = task.HybridStorage
				default:
					storageType = task.MemoryStorage
				}
				
				taskManager = task.NewHybridTaskManager(repository, storageType)
				log.Printf("Using %s storage with database backend", cfg.Database.StorageType)
			}
		}
	} else {
		// Use memory storage
		taskManager = task.NewTaskManager()
		log.Println("Using memory storage")
	}
	
	// Display welcome message
	ui.DisplayWelcome()
	
	// Check if we should run in interactive mode or show help
	if len(os.Args) > 1 {
		cmd.HandleCommand(os.Args[1:], taskManager)
		return
	}
	
	// Run interactive mode
	if err := ui.RunInteractiveMode(taskManager); err != nil {
		log.Fatalf("Error running interactive mode: %v", err)
	}
}