package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// This file demonstrates advanced Go concepts

// 1. Context Package
func demonstrateContext() {
	fmt.Println("=== Context Package ===")
	
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// Simulate work with context
	result := doWorkWithContext(ctx, "Task 1")
	fmt.Printf("Result: %s\n", result)
	
	// Context with cancellation
	ctx2, cancel2 := context.WithCancel(context.Background())
	
	go func() {
		time.Sleep(1 * time.Second)
		cancel2() // Cancel after 1 second
	}()
	
	result2 := doWorkWithContext(ctx2, "Task 2")
	fmt.Printf("Result: %s\n", result2)
	
	fmt.Println()
}

func doWorkWithContext(ctx context.Context, taskName string) string {
	select {
	case <-time.After(3 * time.Second):
		return fmt.Sprintf("%s completed", taskName)
	case <-ctx.Done():
		return fmt.Sprintf("%s cancelled: %v", taskName, ctx.Err())
	}
}

// 2. WaitGroup
func demonstrateWaitGroup() {
	fmt.Println("=== WaitGroup ===")
	
	var wg sync.WaitGroup
	
	// Start multiple goroutines
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Worker %d finished\n", id)
		}(i)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Println("All workers completed!")
	fmt.Println()
}

// 3. Mutex for Synchronization
func demonstrateMutex() {
	fmt.Println("=== Mutex Synchronization ===")
	
	counter := &Counter{}
	var wg sync.WaitGroup
	
	// Start multiple goroutines that increment counter
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter.Value())
	fmt.Println()
}

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// 4. Select Statement
func demonstrateSelect() {
	fmt.Println("=== Select Statement ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	// Send data to channels
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "Message from channel 1"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Message from channel 2"
	}()
	
	// Select on multiple channels
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received: %s\n", msg2)
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout!")
		}
	}
	
	fmt.Println()
}

// 5. Buffered Channels
func demonstrateBufferedChannels() {
	fmt.Println("=== Buffered Channels ===")
	
	// Buffered channel with capacity 3
	ch := make(chan int, 3)
	
	// Send values (won't block until buffer is full)
	ch <- 1
	ch <- 2
	ch <- 3
	
	fmt.Printf("Channel length: %d, capacity: %d\n", len(ch), cap(ch))
	
	// Receive values
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("Received: %d\n", <-ch)
	
	fmt.Println()
}

// 6. Channel Direction
func demonstrateChannelDirection() {
	fmt.Println("=== Channel Direction ===")
	
	ch := make(chan string)
	
	// Send-only channel
	go sendOnly(ch)
	
	// Receive-only channel
	receiveOnly(ch)
	
	fmt.Println()
}

func sendOnly(ch chan<- string) {
	ch <- "Hello from send-only channel"
	close(ch)
}

func receiveOnly(ch <-chan string) {
	for msg := range ch {
		fmt.Printf("Received: %s\n", msg)
	}
}

// 7. Interface Composition
func demonstrateInterfaceComposition() {
	fmt.Println("=== Interface Composition ===")
	
	// Create a file that implements both Reader and Writer
	file := &File{name: "example.txt", content: "Hello, World!"}
	
	// Use as Reader
	var reader Reader = file
	data := make([]byte, 100)
	n, err := reader.Read(data)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
	} else {
		fmt.Printf("Read %d bytes: %s\n", n, string(data[:n]))
	}
	
	// Use as Writer
	var writer Writer = file
	n, err = writer.Write([]byte("New content"))
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
	} else {
		fmt.Printf("Wrote %d bytes\n", n)
	}
	
	// Use as ReadWriter (composed interface)
	var readWriter ReadWriter = file
	fmt.Printf("File content: %s\n", readWriter.GetContent())
	
	fmt.Println()
}

type Reader interface {
	Read([]byte) (int, error)
}

type Writer interface {
	Write([]byte) (int, error)
}

type ReadWriter interface {
	Reader
	Writer
	GetContent() string
}

type File struct {
	name    string
	content string
}

func (f *File) Read(data []byte) (int, error) {
	copy(data, f.content)
	return len(f.content), nil
}

func (f *File) Write(data []byte) (int, error) {
	f.content = string(data)
	return len(data), nil
}

func (f *File) GetContent() string {
	return f.content
}

// 8. Type Assertions and Type Switches
func demonstrateTypeAssertions() {
	fmt.Println("=== Type Assertions and Type Switches ===")
	
	var i interface{} = "Hello, World!"
	
	// Type assertion
	if str, ok := i.(string); ok {
		fmt.Printf("String value: %s\n", str)
	}
	
	// Type switch
	processValue(42)
	processValue("Hello")
	processValue(3.14)
	processValue(true)
	
	fmt.Println()
}

func processValue(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	case float64:
		fmt.Printf("Float: %.2f\n", v)
	case bool:
		fmt.Printf("Boolean: %t\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

// 9. Reflection-like behavior with interfaces
func demonstrateEmptyInterface() {
	fmt.Println("=== Empty Interface ===")
	
	// Empty interface can hold any type
	var values []interface{}
	values = append(values, 42, "Hello", 3.14, true, []int{1, 2, 3})
	
	for i, v := range values {
		fmt.Printf("Index %d: %v (type: %T)\n", i, v, v)
	}
	
	fmt.Println()
}

// 10. Panic and Recover
func demonstratePanicRecover() {
	fmt.Println("=== Panic and Recover ===")
	
	// Function that might panic
	safeDivide := func(a, b int) (result int, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v", r)
			}
		}()
		
		if b == 0 {
			panic("division by zero")
		}
		
		result = a / b
		return result, nil
	}
	
	// Test safe division
	result, err := safeDivide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %d\n", result)
	}
	
	result, err = safeDivide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 0 = %d\n", result)
	}
	
	fmt.Println()
}

func mainAdvanced() {
	fmt.Println("🚀 Go Learning Examples - Advanced Concepts")
	fmt.Println("===========================================")
	fmt.Println()
	
	demonstrateContext()
	demonstrateWaitGroup()
	demonstrateMutex()
	demonstrateSelect()
	demonstrateBufferedChannels()
	demonstrateChannelDirection()
	demonstrateInterfaceComposition()
	demonstrateTypeAssertions()
	demonstrateEmptyInterface()
	demonstratePanicRecover()
	
	fmt.Println("🎉 All advanced examples completed!")
}
