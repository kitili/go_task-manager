# 📚 Go Learning Guide

This comprehensive guide will take you from Go beginner to confident Go developer through hands-on learning with our Task Manager project.

## 🎯 Learning Objectives

By the end of this guide, you will:
- Understand Go's fundamental concepts and syntax
- Be able to write concurrent Go programs
- Know how to structure Go projects properly
- Have built a real-world application
- Be ready to tackle more complex Go projects

---

## 📋 Table of Contents

1. [Getting Started](#getting-started)
2. [Basic Concepts](#basic-concepts)
3. [Intermediate Concepts](#intermediate-concepts)
4. [Advanced Concepts](#advanced-concepts)
5. [Project Structure](#project-structure)
6. [Hands-on Exercises](#hands-on-exercises)
7. [Next Steps](#next-steps)

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+ installed
- Basic programming knowledge (any language)
- Terminal/command prompt access

### Setup
```bash
# Clone and setup
git clone <your-repo-url>
cd learn-go-capstone-main
go mod tidy

# Verify installation
go version
```

### Your First Go Program
```bash
# Run the basic examples
go run examples/basic_concepts.go
```

---

## 🔤 Basic Concepts

### 1. Variables and Types

Go is statically typed, meaning you must declare variable types.

```go
// Different ways to declare variables
var name string = "Go Learner"        // Explicit type
var age int = 25                      // Integer
var isLearning bool = true            // Boolean

// Short declaration (most common)
language := "Go"                      // Type inferred
version := 1.21                       // Float

// Multiple variables
var (
    firstName = "John"
    lastName  = "Doe"
    email     = "john@example.com"
)

// Constants
const pi = 3.14159
const company = "Google"
```

**Exercise**: Modify the examples to add your own variables and constants.

### 2. Functions

Functions are first-class citizens in Go.

```go
// Basic function
func add(a, b int) int {
    return a + b
}

// Multiple return values
func calculate(x, y int) (int, int) {
    return x + y, x * y
}

// Named return values
func divide(a, b float64) (result float64, err error) {
    if b == 0 {
        err = fmt.Errorf("division by zero")
        return
    }
    result = a / b
    return
}

// Variadic functions
func sumAll(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}
```

**Exercise**: Create a function that calculates the area of different shapes.

### 3. Control Structures

```go
// If statements
if age >= 18 {
    fmt.Println("Adult")
} else if age >= 13 {
    fmt.Println("Teenager")
} else {
    fmt.Println("Child")
}

// For loops
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// Range loops
numbers := []int{1, 2, 3, 4, 5}
for index, value := range numbers {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}

// Switch statements
switch day {
case "Monday":
    fmt.Println("Start of work week")
case "Friday":
    fmt.Println("TGIF!")
default:
    fmt.Println("Regular day")
}
```

**Exercise**: Create a function that categorizes tasks by priority using switch statements.

---

## 🏗️ Intermediate Concepts

### 1. Structs and Methods

Structs are Go's way of creating custom types.

```go
type Person struct {
    Name     string
    Age      int
    Email    string
    IsActive bool
}

// Method with value receiver
func (p Person) String() string {
    return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

// Method with pointer receiver (can modify struct)
func (p *Person) UpdateEmail(newEmail string) {
    p.Email = newEmail
}

// Constructor function
func NewPerson(name string, age int, email string) *Person {
    return &Person{
        Name:     name,
        Age:      age,
        Email:    email,
        IsActive: true,
    }
}
```

**Exercise**: Add a `Task` struct with methods to calculate completion percentage.

### 2. Interfaces

Interfaces define behavior, not implementation.

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

// Interface composition
type ReadWriter interface {
    Reader
    Writer
}
```

**Exercise**: Create a `TaskProcessor` interface with methods for different task types.

### 3. Error Handling

Go uses explicit error handling, not exceptions.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Usage
result, err := divide(10, 2)
if err != nil {
    log.Printf("Error: %v", err)
    return
}
fmt.Printf("Result: %.2f\n", result)

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}
```

**Exercise**: Add validation to the Task struct with custom error types.

---

## 🚀 Advanced Concepts

### 1. Goroutines

Goroutines are lightweight threads managed by the Go runtime.

```go
// Basic goroutine
go func() {
    fmt.Println("This runs concurrently!")
}()

// Goroutine with parameters
go processTask(taskID, taskData)

// Anonymous goroutine
go func(id int) {
    fmt.Printf("Processing task %d\n", id)
}(taskID)
```

**Exercise**: Modify the task manager to process tasks concurrently.

### 2. Channels

Channels are Go's way of communicating between goroutines.

```go
// Create a channel
ch := make(chan string)

// Send data
go func() {
    ch <- "Hello from goroutine!"
}()

// Receive data
msg := <-ch
fmt.Println(msg)

// Buffered channel
ch := make(chan int, 3)
ch <- 1
ch <- 2
ch <- 3

// Range over channel
for msg := range ch {
    fmt.Println(msg)
}
```

**Exercise**: Create a worker pool that processes tasks using channels.

### 3. Select Statements

Select allows you to work with multiple channels.

```go
select {
case msg1 := <-ch1:
    fmt.Printf("Received from ch1: %s\n", msg1)
case msg2 := <-ch2:
    fmt.Printf("Received from ch2: %s\n", msg2)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout!")
default:
    fmt.Println("No message ready")
}
```

**Exercise**: Add timeout handling to task processing.

### 4. Sync Package

```go
// WaitGroup - wait for goroutines to complete
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        processTask(id)
    }(i)
}
wg.Wait()

// Mutex - protect shared data
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

**Exercise**: Add thread-safe task counting to the task manager.

---

## 🏗️ Project Structure

Understanding Go project structure is crucial for building maintainable applications.

```
project/
├── main.go              # Entry point
├── go.mod              # Module definition
├── internal/           # Private application code
│   ├── task/          # Task management
│   └── ui/            # User interface
├── cmd/               # Command-line tools
├── examples/          # Learning examples
└── docs/             # Documentation
```

### Package Organization Principles

1. **`internal/`** - Private code, not importable by other modules
2. **`cmd/`** - Command-line applications
3. **`examples/`** - Example code and learning materials
4. **`docs/`** - Documentation and guides

**Exercise**: Create a new package for task persistence.

---

## 🎯 Hands-on Exercises

### Beginner Exercises

1. **Add Task Categories**
   ```go
   type Category struct {
       ID   int
       Name string
   }
   
   type Task struct {
       // ... existing fields
       Category Category
   }
   ```

2. **Create Task Validation**
   ```go
   func (t Task) Validate() error {
       if t.Title == "" {
           return errors.New("title cannot be empty")
       }
       // Add more validation rules
   }
   ```

3. **Add Task Search**
   ```go
   func (tm *TaskManager) SearchTasks(query string) []Task {
       // Implement search functionality
   }
   ```

### Intermediate Exercises

1. **JSON Persistence**
   ```go
   func (tm *TaskManager) SaveToFile(filename string) error {
       // Save tasks to JSON file
   }
   
   func (tm *TaskManager) LoadFromFile(filename string) error {
       // Load tasks from JSON file
   }
   ```

2. **Task Dependencies**
   ```go
   type Task struct {
       // ... existing fields
       Dependencies []int // IDs of tasks this depends on
   }
   ```

3. **Task Scheduling**
   ```go
   func (tm *TaskManager) ScheduleTask(taskID int, scheduleTime time.Time) error {
       // Schedule task for later execution
   }
   ```

### Advanced Exercises

1. **Web API**
   ```go
   func handleGetTasks(w http.ResponseWriter, r *http.Request) {
       // HTTP handler for getting tasks
   }
   ```

2. **Database Integration**
   ```go
   type TaskRepository interface {
       Save(task Task) error
       FindByID(id int) (Task, error)
       FindAll() ([]Task, error)
   }
   ```

3. **Plugin System**
   ```go
   type TaskProcessor interface {
       Process(task Task) error
       Name() string
   }
   ```

---

## 🧪 Testing Your Code

Go has excellent built-in testing support.

```go
// task_test.go
func TestAddTask(t *testing.T) {
    tm := NewTaskManager()
    task := tm.AddTask("Test", "Test task", High, nil)
    
    if task == nil {
        t.Error("Expected task to be created")
    }
    
    if task.Title != "Test" {
        t.Errorf("Expected title 'Test', got '%s'", task.Title)
    }
}

func TestGetTask(t *testing.T) {
    tm := NewTaskManager()
    tm.AddTask("Test", "Test task", High, nil)
    
    task, err := tm.GetTask(1)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if task.Title != "Test" {
        t.Errorf("Expected title 'Test', got '%s'", task.Title)
    }
}
```

Run tests with:
```bash
go test ./...
go test -v ./...  # Verbose output
go test -cover ./...  # Coverage report
```

---

## 📚 Additional Resources

### Official Go Resources
- [Go Documentation](https://golang.org/doc/)
- [Go Tour](https://tour.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)

### Books
- "The Go Programming Language" by Alan Donovan and Brian Kernighan
- "Go in Action" by William Kennedy
- "Concurrency in Go" by Katherine Cox-Buday

### Online Courses
- [Go: The Complete Developer's Guide](https://www.udemy.com/course/go-the-complete-developers-guide/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)

---

## 🎉 Next Steps

After completing this learning guide:

1. **Build Your Own Project** - Create something unique using Go
2. **Contribute to Open Source** - Find Go projects on GitHub
3. **Learn Go Web Development** - Build web APIs and services
4. **Explore Go Ecosystem** - Discover popular Go libraries and frameworks
5. **Join the Community** - Participate in Go forums and meetups

---

## 🤝 Getting Help

- **Go Forum**: https://forum.golang.org/
- **Reddit**: r/golang
- **Stack Overflow**: Tag your questions with `go`
- **GitHub Issues**: Open an issue in this repository

---

**Happy Learning! 🚀**

*Remember: The best way to learn Go is by writing Go code. Start with the examples, modify them, and build something amazing!*
