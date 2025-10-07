# 🏋️ Go Learning Exercises

This document contains hands-on exercises to reinforce your Go learning. Complete these exercises in order, as they build upon each other.

---

## 🎯 Exercise Levels

- **🟢 Beginner** - Basic Go concepts
- **🟡 Intermediate** - Structs, interfaces, error handling
- **🔴 Advanced** - Concurrency, channels, advanced patterns

---

## 🟢 Beginner Exercises

### Exercise 1: Variables and Functions
**Goal**: Create a simple calculator

```go
// Create a file: exercises/calculator.go
package main

import "fmt"

// TODO: Create functions for basic math operations
// add, subtract, multiply, divide
// Each should take two float64 parameters and return float64

func main() {
    // TODO: Test your functions
    // Example: fmt.Println(add(5.5, 3.2))
}
```

**Requirements**:
- Create functions for +, -, *, /
- Handle division by zero with proper error handling
- Test all functions with different inputs

### Exercise 2: Control Structures
**Goal**: Build a number guessing game

```go
// Create a file: exercises/guessing_game.go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    // TODO: Generate a random number between 1-100
    // TODO: Create a loop that asks for user input
    // TODO: Provide hints (too high/too low)
    // TODO: Count number of guesses
    // TODO: Congratulate when correct
}
```

**Requirements**:
- Use `rand.Intn()` for random numbers
- Use `fmt.Scanln()` for user input
- Provide feedback on each guess
- Count and display total guesses

### Exercise 3: Slices and Maps
**Goal**: Create a student grade tracker

```go
// Create a file: exercises/grade_tracker.go
package main

import "fmt"

type Student struct {
    Name  string
    Grades []int
}

func main() {
    // TODO: Create a map of students
    // TODO: Add grades for each student
    // TODO: Calculate average grade for each student
    // TODO: Find the student with highest average
}
```

**Requirements**:
- Use a map with student names as keys
- Store multiple grades per student
- Calculate and display averages
- Find and display the top performer

---

## 🟡 Intermediate Exercises

### Exercise 4: Structs and Methods
**Goal**: Build a bank account system

```go
// Create a file: exercises/bank_account.go
package main

import (
    "fmt"
    "errors"
)

type Account struct {
    ID      int
    Name    string
    Balance float64
}

// TODO: Add methods to Account
// Deposit(amount float64) error
// Withdraw(amount float64) error
// GetBalance() float64
// String() string

func main() {
    // TODO: Create accounts and test all methods
}
```

**Requirements**:
- Prevent negative balances
- Return appropriate errors
- Implement String() method for display
- Test all edge cases

### Exercise 5: Interfaces
**Goal**: Create a shape calculator

```go
// Create a file: exercises/shapes.go
package main

import "fmt"

// TODO: Define Shape interface with Area() and Perimeter() methods

// TODO: Implement Circle, Rectangle, and Triangle structs
// Each should implement the Shape interface

// TODO: Create a function that calculates total area of multiple shapes

func main() {
    // TODO: Create different shapes and test calculations
}
```

**Requirements**:
- Define a proper Shape interface
- Implement at least 3 different shapes
- Create a function that works with any Shape
- Test with mixed shape types

### Exercise 6: Error Handling
**Goal**: Build a file reader with robust error handling

```go
// Create a file: exercises/file_reader.go
package main

import (
    "fmt"
    "os"
    "io/ioutil"
)

// TODO: Create custom error types
// FileNotFoundError, PermissionError, InvalidFileError

// TODO: Create a function that reads a file and returns content
// Handle different types of errors appropriately

func main() {
    // TODO: Test with existing and non-existing files
    // TODO: Test with files you don't have permission to read
}
```

**Requirements**:
- Create custom error types
- Handle file not found errors
- Handle permission errors
- Provide meaningful error messages

---

## 🔴 Advanced Exercises

### Exercise 7: Goroutines and Channels
**Goal**: Build a concurrent web scraper

```go
// Create a file: exercises/web_scraper.go
package main

import (
    "fmt"
    "net/http"
    "time"
)

type ScrapedData struct {
    URL     string
    Content string
    Error   error
}

// TODO: Create a function that scrapes a single URL
// TODO: Use goroutines to scrape multiple URLs concurrently
// TODO: Use channels to collect results
// TODO: Implement timeout handling

func main() {
    urls := []string{
        "https://httpbin.org/delay/1",
        "https://httpbin.org/delay/2",
        "https://httpbin.org/delay/3",
    }
    
    // TODO: Scrape all URLs concurrently
    // TODO: Display results as they come in
}
```

**Requirements**:
- Use goroutines for concurrent scraping
- Use channels to collect results
- Implement proper timeout handling
- Handle errors gracefully

### Exercise 8: Worker Pool Pattern
**Goal**: Build a task processor with worker pools

```go
// Create a file: exercises/worker_pool.go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Task struct {
    ID   int
    Data string
}

type Result struct {
    TaskID int
    Output string
    Error  error
}

// TODO: Create a worker pool that processes tasks
// TODO: Use WaitGroup to wait for all workers to complete
// TODO: Use channels for task distribution and result collection

func main() {
    // TODO: Create a batch of tasks
    // TODO: Process them using worker pool
    // TODO: Collect and display results
}
```

**Requirements**:
- Implement a configurable worker pool
- Use WaitGroup for synchronization
- Use channels for communication
- Handle worker errors

### Exercise 9: Context and Cancellation
**Goal**: Build a long-running process with cancellation

```go
// Create a file: exercises/context_demo.go
package main

import (
    "context"
    "fmt"
    "time"
)

// TODO: Create a function that does long-running work
// TODO: Use context for cancellation
// TODO: Implement timeout handling
// TODO: Handle context cancellation gracefully

func main() {
    // TODO: Create context with timeout
    // TODO: Start long-running process
    // TODO: Test cancellation scenarios
}
```

**Requirements**:
- Use context.Context for cancellation
- Implement timeout handling
- Handle context cancellation
- Clean up resources properly

---

## 🎯 Project-Based Exercises

### Exercise 10: Enhanced Task Manager
**Goal**: Extend the existing task manager with new features

```go
// Modify existing files to add these features:

// 1. Task Categories
type Category struct {
    ID   int
    Name string
}

// 2. Task Dependencies
type Task struct {
    // ... existing fields
    Dependencies []int // Task IDs this task depends on
    Dependents   []int // Task IDs that depend on this task
}

// 3. Task Scheduling
func (tm *TaskManager) ScheduleTask(taskID int, scheduleTime time.Time) error {
    // TODO: Implement task scheduling
}

// 4. Task Search
func (tm *TaskManager) SearchTasks(query string) []Task {
    // TODO: Implement full-text search
}

// 5. Task Export/Import
func (tm *TaskManager) ExportToJSON(filename string) error {
    // TODO: Export tasks to JSON file
}

func (tm *TaskManager) ImportFromJSON(filename string) error {
    // TODO: Import tasks from JSON file
}
```

### Exercise 11: Web API
**Goal**: Create HTTP endpoints for the task manager

```go
// Create a file: exercises/web_api.go
package main

import (
    "encoding/json"
    "net/http"
    "strconv"
)

// TODO: Create HTTP handlers for:
// GET /tasks - List all tasks
// POST /tasks - Create new task
// GET /tasks/{id} - Get specific task
// PUT /tasks/{id} - Update task
// DELETE /tasks/{id} - Delete task

func main() {
    // TODO: Set up HTTP routes
    // TODO: Start HTTP server
    // TODO: Test all endpoints
}
```

### Exercise 12: Database Integration
**Goal**: Add persistent storage to the task manager

```go
// Create a file: exercises/database.go
package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

// TODO: Create database schema for tasks
// TODO: Implement CRUD operations
// TODO: Add database connection pooling
// TODO: Handle database errors

func main() {
    // TODO: Initialize database
    // TODO: Test all database operations
}
```

---

## 🧪 Testing Exercises

### Exercise 13: Unit Testing
**Goal**: Write comprehensive tests for the task manager

```go
// Create a file: exercises/task_test.go
package main

import "testing"

// TODO: Write tests for all TaskManager methods
// TODO: Test error conditions
// TODO: Test concurrent access
// TODO: Achieve high test coverage

func TestAddTask(t *testing.T) {
    // TODO: Test task creation
}

func TestGetTask(t *testing.T) {
    // TODO: Test task retrieval
}

func TestUpdateTaskStatus(t *testing.T) {
    // TODO: Test status updates
}

func TestDeleteTask(t *testing.T) {
    // TODO: Test task deletion
}

func TestConcurrentAccess(t *testing.T) {
    // TODO: Test thread safety
}
```

---

## 📊 Progress Tracking

### Beginner Level (🟢)
- [ ] Exercise 1: Calculator
- [ ] Exercise 2: Guessing Game
- [ ] Exercise 3: Grade Tracker

### Intermediate Level (🟡)
- [ ] Exercise 4: Bank Account
- [ ] Exercise 5: Shapes
- [ ] Exercise 6: File Reader

### Advanced Level (🔴)
- [ ] Exercise 7: Web Scraper
- [ ] Exercise 8: Worker Pool
- [ ] Exercise 9: Context Demo

### Project Level
- [ ] Exercise 10: Enhanced Task Manager
- [ ] Exercise 11: Web API
- [ ] Exercise 12: Database Integration
- [ ] Exercise 13: Unit Testing

---

## 🎉 Completion Criteria

You've mastered Go when you can:
- [ ] Write idiomatic Go code
- [ ] Handle errors properly
- [ ] Use interfaces effectively
- [ ] Write concurrent programs
- [ ] Structure projects properly
- [ ] Write comprehensive tests
- [ ] Build real-world applications

---

## 🤝 Getting Help

If you get stuck on any exercise:

1. **Read the Go documentation** for the relevant packages
2. **Check the examples** in the `examples/` directory
3. **Look at the existing code** in the task manager
4. **Ask questions** in the Go community
5. **Experiment** with small code snippets

---

**Happy Coding! 🚀**

*Remember: The best way to learn is by doing. Don't just read the exercises - implement them!*
