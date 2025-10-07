package testing

import (
	"fmt"
	"testing"
	"time"

	"learn-go-capstone/internal/database"
	"learn-go-capstone/internal/task"
)

// BenchmarkTaskManagerAddTask benchmarks task creation
func BenchmarkTaskManagerAddTask(b *testing.B) {
	tm := task.NewTaskManager()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.AddTask(fmt.Sprintf("Task %d", i), "Benchmark task", task.Medium, nil)
	}
}

// BenchmarkTaskManagerGetTask benchmarks task retrieval
func BenchmarkTaskManagerGetTask(b *testing.B) {
	tm := task.NewTaskManager()
	
	// Pre-populate with tasks
	for i := 0; i < 1000; i++ {
		tm.AddTask(fmt.Sprintf("Task %d", i), "Benchmark task", task.Medium, nil)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.GetTask(i%1000 + 1)
	}
}

// BenchmarkTaskManagerGetAllTasks benchmarks getting all tasks
func BenchmarkTaskManagerGetAllTasks(b *testing.B) {
	tm := task.NewTaskManager()
	
	// Pre-populate with tasks
	for i := 0; i < 1000; i++ {
		tm.AddTask(fmt.Sprintf("Task %d", i), "Benchmark task", task.Medium, nil)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.GetAllTasks()
	}
}

// BenchmarkDatabaseCreateTask benchmarks database task creation
func BenchmarkDatabaseCreateTask(b *testing.B) {
	// Create a temporary test database
	cfg := &database.Config{
		Driver: "sqlite3",
		DSN:    "file:benchmark_test.db?mode=memory&cache=shared",
	}
	
	db, err := database.Connect(cfg)
	if err != nil {
		b.Fatalf("Failed to connect to test database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		b.Fatalf("Failed to run test migrations: %v", err)
	}

	repository := database.NewSQLiteRepository(db)

	// Create a test user
	user := &database.User{
		Username:  "benchuser",
		Email:     "bench@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repository.CreateUser(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &database.DatabaseTask{
			Title:       fmt.Sprintf("Benchmark Task %d", i),
			Description: "Benchmark task",
			Priority:    3,
			Status:      0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserID:      &user.ID,
			IsArchived:  false,
		}
		repository.CreateTask(task)
	}
}

// BenchmarkDatabaseGetTask benchmarks database task retrieval
func BenchmarkDatabaseGetTask(b *testing.B) {
	// Create a temporary test database
	cfg := &database.Config{
		Driver: "sqlite3",
		DSN:    "file:benchmark_test.db?mode=memory&cache=shared",
	}
	
	db, err := database.Connect(cfg)
	if err != nil {
		b.Fatalf("Failed to connect to test database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		b.Fatalf("Failed to run test migrations: %v", err)
	}

	repository := database.NewSQLiteRepository(db)

	// Create a test user
	user := &database.User{
		Username:  "benchuser",
		Email:     "bench@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repository.CreateUser(user)

	// Pre-populate with tasks
	var taskIDs []int
	for i := 0; i < 1000; i++ {
		task := &database.DatabaseTask{
			Title:       fmt.Sprintf("Benchmark Task %d", i),
			Description: "Benchmark task",
			Priority:    3,
			Status:      0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserID:      &user.ID,
			IsArchived:  false,
		}
		repository.CreateTask(task)
		taskIDs = append(taskIDs, task.ID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repository.GetTask(taskIDs[i%len(taskIDs)])
	}
}

// BenchmarkDatabaseSearchTasks benchmarks database task search
func BenchmarkDatabaseSearchTasks(b *testing.B) {
	// Create a temporary test database
	cfg := &database.Config{
		Driver: "sqlite3",
		DSN:    "file:benchmark_test.db?mode=memory&cache=shared",
	}
	
	db, err := database.Connect(cfg)
	if err != nil {
		b.Fatalf("Failed to connect to test database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		b.Fatalf("Failed to run test migrations: %v", err)
	}

	repository := database.NewSQLiteRepository(db)

	// Create a test user
	user := &database.User{
		Username:  "benchuser",
		Email:     "bench@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repository.CreateUser(user)

	// Pre-populate with tasks
	for i := 0; i < 1000; i++ {
		task := &database.DatabaseTask{
			Title:       fmt.Sprintf("Benchmark Task %d", i),
			Description: "Benchmark task description",
			Priority:    3,
			Status:      0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserID:      &user.ID,
			IsArchived:  false,
		}
		repository.CreateTask(task)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repository.SearchTasks("Benchmark")
	}
}

// BenchmarkConcurrentTaskCreation benchmarks concurrent task creation
func BenchmarkConcurrentTaskCreation(b *testing.B) {
	tm := task.NewTaskManager()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			tm.AddTask(fmt.Sprintf("Concurrent Task %d", i), "Benchmark task", task.Medium, nil)
			i++
		}
	})
}

// BenchmarkConcurrentDatabaseOperations benchmarks concurrent database operations
func BenchmarkConcurrentDatabaseOperations(b *testing.B) {
	// Create a temporary test database
	cfg := &database.Config{
		Driver: "sqlite3",
		DSN:    "file:benchmark_test.db?mode=memory&cache=shared",
	}
	
	db, err := database.Connect(cfg)
	if err != nil {
		b.Fatalf("Failed to connect to test database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	migrationManager := database.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		b.Fatalf("Failed to run test migrations: %v", err)
	}

	repository := database.NewSQLiteRepository(db)

	// Create a test user
	user := &database.User{
		Username:  "benchuser",
		Email:     "bench@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repository.CreateUser(user)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			task := &database.DatabaseTask{
				Title:       fmt.Sprintf("Concurrent Task %d", i),
				Description: "Benchmark task",
				Priority:    3,
				Status:      0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UserID:      &user.ID,
				IsArchived:  false,
			}
			repository.CreateTask(task)
			i++
		}
	})
}

// TestMemoryUsage tests memory usage with large datasets
func TestMemoryUsage(t *testing.T) {
	tm := task.NewTaskManager()
	
	// Create a large number of tasks
	numTasks := 10000
	start := time.Now()
	
	for i := 0; i < numTasks; i++ {
		tm.AddTask(fmt.Sprintf("Memory Test Task %d", i), "Memory test task", task.Medium, nil)
	}
	
	creationTime := time.Since(start)
	t.Logf("Created %d tasks in %v", numTasks, creationTime)
	
	// Verify all tasks were created
	tasks := tm.GetAllTasks()
	AssertEqual(t, numTasks, len(tasks), "Should have created all tasks")
	
	// Test retrieval performance
	start = time.Now()
	for i := 1; i <= numTasks; i++ {
		_, err := tm.GetTask(i)
		AssertNoError(t, err, "Should retrieve task successfully")
	}
	retrievalTime := time.Since(start)
	t.Logf("Retrieved %d tasks in %v", numTasks, retrievalTime)
	
	// Test search performance
	start = time.Now()
	pendingTasks := tm.GetTasksByStatus(task.Pending)
	searchTime := time.Since(start)
	t.Logf("Found %d pending tasks in %v", len(pendingTasks), searchTime)
}

// TestDatabaseMemoryUsage tests database memory usage with large datasets
func TestDatabaseMemoryUsage(t *testing.T) {
	db, repository, cleanup := SetupTestDB(t)
	defer cleanup()
	_ = db

	// Create a test user
	user := &database.User{
		Username:  "memuser",
		Email:     "mem@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repository.CreateUser(user)

	// Create a large number of tasks
	numTasks := 1000
	start := time.Now()
	
	for i := 0; i < numTasks; i++ {
		task := &database.DatabaseTask{
			Title:       fmt.Sprintf("Memory Test Task %d", i),
			Description: "Memory test task",
			Priority:    3,
			Status:      0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserID:      &user.ID,
			IsArchived:  false,
		}
		repository.CreateTask(task)
	}
	
	creationTime := time.Since(start)
	t.Logf("Created %d database tasks in %v", numTasks, creationTime)
	
	// Verify all tasks were created
	tasks, err := repository.GetAllTasks()
	AssertNoError(t, err, "Should retrieve all tasks")
	AssertEqual(t, numTasks, len(tasks), "Should have created all tasks")
	
	// Test search performance
	start = time.Now()
	searchResults, err := repository.SearchTasks("Memory")
	AssertNoError(t, err, "Should search tasks successfully")
	searchTime := time.Since(start)
	t.Logf("Found %d tasks matching 'Memory' in %v", len(searchResults), searchTime)
}

// TestConcurrentAccess tests concurrent access to task manager
func TestConcurrentAccess(t *testing.T) {
	tm := task.NewTaskManager()
	
	// Test concurrent reads and writes
	numGoroutines := 10
	numOperations := 100
	
	done := make(chan bool, numGoroutines)
	
	start := time.Now()
	
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < numOperations; j++ {
				// Add task
				taskItem := tm.AddTask(fmt.Sprintf("Concurrent Task %d-%d", goroutineID, j), "Concurrent task", task.Medium, nil)
				
				// Read task
				_, err := tm.GetTask(taskItem.ID)
				AssertNoError(t, err, "Should retrieve task successfully")
				
				// Update task status
				err = tm.UpdateTaskStatus(taskItem.ID, task.InProgress)
				AssertNoError(t, err, "Should update task status successfully")
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	totalTime := time.Since(start)
	t.Logf("Completed %d concurrent operations in %v", numGoroutines*numOperations*3, totalTime)
	
	// Verify final state
	tasks := tm.GetAllTasks()
	AssertEqual(t, numGoroutines*numOperations, len(tasks), "Should have created all tasks")
}

// TestDatabaseConcurrentAccess tests concurrent access to database
func TestDatabaseConcurrentAccess(t *testing.T) {
	db, repository, cleanup := SetupTestDB(t)
	defer cleanup()
	_ = db

	// Create a test user
	user := &database.User{
		Username:  "concurrentuser",
		Email:     "concurrent@example.com",
		Password:  "hashedpassword",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repository.CreateUser(user)

	// Test concurrent reads and writes
	numGoroutines := 5
	numOperations := 50
	
	done := make(chan bool, numGoroutines)
	
	start := time.Now()
	
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < numOperations; j++ {
				// Add task
				task := &database.DatabaseTask{
					Title:       fmt.Sprintf("Concurrent Task %d-%d", goroutineID, j),
					Description: "Concurrent task",
					Priority:    3,
					Status:      0,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					UserID:      &user.ID,
					IsArchived:  false,
				}
				err := repository.CreateTask(task)
				AssertNoError(t, err, "Should create task successfully")
				
				// Read task
				_, err = repository.GetTask(task.ID)
				AssertNoError(t, err, "Should retrieve task successfully")
				
				// Update task status
				task.Status = 1 // InProgress
				task.UpdatedAt = time.Now()
				err = repository.UpdateTask(task)
				AssertNoError(t, err, "Should update task successfully")
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	totalTime := time.Since(start)
	t.Logf("Completed %d concurrent database operations in %v", numGoroutines*numOperations*3, totalTime)
	
	// Verify final state
	tasks, err := repository.GetAllTasks()
	AssertNoError(t, err, "Should retrieve all tasks")
	AssertEqual(t, numGoroutines*numOperations, len(tasks), "Should have created all tasks")
}
