package main

import (
	"fmt"
	"time"
)

// This file demonstrates fundamental Go concepts through practical examples

// 1. Variables and Constants
func demonstrateVariables() {
	fmt.Println("=== Variables and Constants ===")
	
	// Different ways to declare variables
	var name string = "Go Learner"
	var age int = 25
	var isLearning bool = true
	
	// Short declaration (most common)
	language := "Go"
	version := 1.21
	
	// Multiple variable declaration
	var (
		firstName = "John"
		lastName  = "Doe"
		email     = "john@example.com"
	)
	
	// Constants
	const pi = 3.14159
	const company = "Google"
	
	fmt.Printf("Name: %s, Age: %d, Learning: %t\n", name, age, isLearning)
	fmt.Printf("Language: %s, Version: %.1f\n", language, version)
	fmt.Printf("Full Name: %s %s, Email: %s\n", firstName, lastName, email)
	fmt.Printf("Pi: %.5f, Company: %s\n", pi, company)
	fmt.Println()
}

// 2. Functions
func demonstrateFunctions() {
	fmt.Println("=== Functions ===")
	
	// Basic function
	result := add(10, 20)
	fmt.Printf("10 + 20 = %d\n", result)
	
	// Function with multiple return values
	sum, product := calculate(5, 3)
	fmt.Printf("Sum: %d, Product: %d\n", sum, product)
	
	// Variadic function
	total := sumAll(1, 2, 3, 4, 5)
	fmt.Printf("Sum of 1,2,3,4,5 = %d\n", total)
	
	// Anonymous function
	multiply := func(x, y int) int {
		return x * y
	}
	fmt.Printf("5 * 6 = %d\n", multiply(5, 6))
	
	fmt.Println()
}

func add(a, b int) int {
	return a + b
}

func calculate(x, y int) (int, int) {
	return x + y, x * y
}

func sumAll(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// 3. Structs and Methods
func demonstrateStructs() {
	fmt.Println("=== Structs and Methods ===")
	
	// Create a person
	person := Person{
		Name:    "Alice",
		Age:     30,
		Email:   "alice@example.com",
		IsActive: true,
	}
	
	// Call methods
	fmt.Printf("Person: %s\n", person.String())
	fmt.Printf("Is adult: %t\n", person.IsAdult())
	
	// Update person
	person.UpdateEmail("alice.new@example.com")
	fmt.Printf("Updated email: %s\n", person.Email)
	
	// Create person with constructor
	person2 := NewPerson("Bob", 25, "bob@example.com")
	fmt.Printf("New person: %s\n", person2.String())
	
	fmt.Println()
}

type Person struct {
	Name     string
	Age      int
	Email    string
	IsActive bool
}

// Method with value receiver
func (p Person) String() string {
	status := "inactive"
	if p.IsActive {
		status = "active"
	}
	return fmt.Sprintf("%s (%d years old, %s) - %s", p.Name, p.Age, p.Email, status)
}

// Method with pointer receiver (can modify the struct)
func (p *Person) UpdateEmail(newEmail string) {
	p.Email = newEmail
}

// Method with value receiver
func (p Person) IsAdult() bool {
	return p.Age >= 18
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

// 4. Interfaces
func demonstrateInterfaces() {
	fmt.Println("=== Interfaces ===")
	
	// Create different shapes
	circle := Circle{Radius: 5.0}
	rectangle := Rectangle{Width: 4.0, Height: 6.0}
	
	// Use interface
	shapes := []Shape{circle, rectangle}
	
	for i, shape := range shapes {
		fmt.Printf("Shape %d: Area = %.2f, Perimeter = %.2f\n", 
			i+1, shape.Area(), shape.Perimeter())
	}
	
	// Interface with multiple methods
	var writer Writer = &ConsoleWriter{}
	writer.Write([]byte("Hello from interface!"))
	
	fmt.Println()
}

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

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Writer interface {
	Write([]byte) (int, error)
}

type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(data []byte) (int, error) {
	n, err := fmt.Print(string(data))
	return n, err
}

// 5. Error Handling
func demonstrateErrorHandling() {
	fmt.Println("=== Error Handling ===")
	
	// Function that can return an error
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", result)
	}
	
	// Function that returns an error
	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 0 = %.2f\n", result)
	}
	
	// Custom error type
	_, err = processUser("")
	if err != nil {
		fmt.Printf("Custom error: %v\n", err)
	}
	
	_, err = processUser("valid@email.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("User processed successfully!")
	}
	
	fmt.Println()
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

func processUser(email string) (string, error) {
	if email == "" {
		return "", ValidationError{
			Field:   "email",
			Message: "email cannot be empty",
		}
	}
	return fmt.Sprintf("Processing user with email: %s", email), nil
}

// 6. Goroutines and Channels
func demonstrateConcurrency() {
	fmt.Println("=== Goroutines and Channels ===")
	
	// Basic goroutine
	go func() {
		fmt.Println("This runs in a goroutine!")
	}()
	
	// Channel for communication
	ch := make(chan string, 2)
	
	// Send data to channel
	go func() {
		ch <- "Hello"
		ch <- "World"
		close(ch)
	}()
	
	// Receive data from channel
	for msg := range ch {
		fmt.Printf("Received: %s\n", msg)
	}
	
	// Worker pool pattern
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	
	// Start workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	
	// Send jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Collect results
	for a := 1; a <= 5; a++ {
		result := <-results
		fmt.Printf("Job result: %d\n", result)
	}
	
	fmt.Println()
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(100 * time.Millisecond) // Simulate work
		results <- j * 2
	}
}

// 7. Slices and Maps
func demonstrateCollections() {
	fmt.Println("=== Slices and Maps ===")
	
	// Slices
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("Numbers: %v\n", numbers)
	
	// Append to slice
	numbers = append(numbers, 6, 7, 8)
	fmt.Printf("After append: %v\n", numbers)
	
	// Slice operations
	fmt.Printf("First 3: %v\n", numbers[:3])
	fmt.Printf("Last 3: %v\n", numbers[len(numbers)-3:])
	fmt.Printf("Middle: %v\n", numbers[2:5])
	
	// Maps
	ages := map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Carol": 35,
	}
	
	fmt.Printf("Ages: %v\n", ages)
	
	// Add to map
	ages["David"] = 28
	fmt.Printf("After adding David: %v\n", ages)
	
	// Check if key exists
	if age, exists := ages["Alice"]; exists {
		fmt.Printf("Alice is %d years old\n", age)
	}
	
	// Delete from map
	delete(ages, "Bob")
	fmt.Printf("After removing Bob: %v\n", ages)
	
	fmt.Println()
}

func main() {
	fmt.Println("🚀 Go Learning Examples - Fundamental Concepts")
	fmt.Println("=============================================")
	fmt.Println()
	
	demonstrateVariables()
	demonstrateFunctions()
	demonstrateStructs()
	demonstrateInterfaces()
	demonstrateErrorHandling()
	demonstrateConcurrency()
	demonstrateCollections()
	
	fmt.Println("🎉 All examples completed!")
}
