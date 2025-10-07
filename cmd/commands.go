package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"learn-go-capstone/internal/task"
	"github.com/fatih/color"
)

// HandleCommand processes command line arguments
func HandleCommand(args []string, tm task.TaskManagerInterface) {
	if len(args) == 0 {
		return
	}
	
	command := args[0]
	
	switch command {
	case "add":
		handleAddCommand(args[1:], tm)
	case "list":
		handleListCommand(args[1:], tm)
	case "update":
		handleUpdateCommand(args[1:], tm)
	case "delete":
		handleDeleteCommand(args[1:], tm)
	case "stats":
		handleStatsCommand(tm)
	case "demo":
		handleDemoCommand(tm)
	case "help":
		handleHelpCommand()
	default:
		color.Red("❌ Unknown command: %s", command)
		handleHelpCommand()
	}
}

func handleAddCommand(args []string, tm task.TaskManagerInterface) {
	if len(args) < 2 {
		color.Red("❌ Usage: go run main.go add <title> <description> [priority] [due_date]")
		color.White("Priority: 1=Low, 2=Medium, 3=High, 4=Urgent")
		color.White("Due date format: YYYY-MM-DD")
		return
	}
	
	title := args[0]
	description := args[1]
	
	priority := task.Medium // Default priority
	if len(args) > 2 {
		if p, err := strconv.Atoi(args[2]); err == nil && p >= 1 && p <= 4 {
			priority = task.Priority(p - 1)
		}
	}
	
	var dueDate *time.Time
	if len(args) > 3 {
		if parsed, err := time.Parse("2006-01-02", args[3]); err == nil {
			dueDate = &parsed
		}
	}
	
	newTask := tm.AddTask(title, description, priority, dueDate)
	color.Green("✅ Task added successfully!")
	color.White("ID: %d | Title: %s | Priority: %s", 
		newTask.ID, newTask.Title, newTask.Priority.String())
}

func handleListCommand(args []string, tm task.TaskManagerInterface) {
	tasks := tm.GetAllTasks()
	
	if len(tasks) == 0 {
		color.Yellow("📝 No tasks found.")
		return
	}
	
	// Check for filters
	if len(args) > 0 {
		filter := args[0]
		switch filter {
		case "pending":
			tasks = tm.GetTasksByStatus(task.Pending)
		case "in-progress":
			tasks = tm.GetTasksByStatus(task.InProgress)
		case "completed":
			tasks = tm.GetTasksByStatus(task.Completed)
		case "cancelled":
			tasks = tm.GetTasksByStatus(task.Cancelled)
		case "overdue":
			tasks = tm.GetOverdueTasks()
		case "priority":
			if len(args) > 1 {
				if p, err := strconv.Atoi(args[1]); err == nil && p >= 1 && p <= 4 {
					tasks = tm.GetTasksByPriority(task.Priority(p - 1))
				}
			}
		}
	}
	
	color.Cyan("📋 Task List:")
	for _, t := range tasks {
		dueDate := "N/A"
		if t.DueDate != nil {
			dueDate = t.DueDate.Format("2006-01-02")
		}
		
		priorityColor := getPriorityColor(t.Priority)
		statusColor := getStatusColor(t.Status)
		
		fmt.Printf("ID: %d | %s | %s | %s | Due: %s\n",
			t.ID, t.Title,
			priorityColor(t.Priority.String()),
			statusColor(t.Status.String()),
			dueDate)
	}
}

func handleUpdateCommand(args []string, tm task.TaskManagerInterface) {
	if len(args) < 2 {
		color.Red("❌ Usage: go run main.go update <task_id> <status>")
		color.White("Status: pending, in-progress, completed, cancelled")
		return
	}
	
	id, err := strconv.Atoi(args[0])
	if err != nil {
		color.Red("❌ Invalid task ID")
		return
	}
	
	statusStr := strings.ToLower(args[1])
	var status task.Status
	
	switch statusStr {
	case "pending":
		status = task.Pending
	case "in-progress":
		status = task.InProgress
	case "completed":
		status = task.Completed
	case "cancelled":
		status = task.Cancelled
	default:
		color.Red("❌ Invalid status. Use: pending, in-progress, completed, cancelled")
		return
	}
	
	err = tm.UpdateTaskStatus(id, status)
	if err != nil {
		color.Red("❌ %v", err)
		return
	}
	
	color.Green("✅ Task status updated successfully!")
}

func handleDeleteCommand(args []string, tm task.TaskManagerInterface) {
	if len(args) < 1 {
		color.Red("❌ Usage: go run main.go delete <task_id>")
		return
	}
	
	id, err := strconv.Atoi(args[0])
	if err != nil {
		color.Red("❌ Invalid task ID")
		return
	}
	
	err = tm.DeleteTask(id)
	if err != nil {
		color.Red("❌ %v", err)
		return
	}
	
	color.Green("✅ Task deleted successfully!")
}

func handleStatsCommand(tm task.TaskManagerInterface) {
	allTasks := tm.GetAllTasks()
	
	if len(allTasks) == 0 {
		color.Yellow("📊 No tasks to show statistics for.")
		return
	}
	
	color.Cyan("📊 Task Manager Statistics")
	fmt.Println("========================")
	
	// Total tasks
	color.White("Total Tasks: %d", len(allTasks))
	
	// Status breakdown
	pending := len(tm.GetTasksByStatus(task.Pending))
	inProgress := len(tm.GetTasksByStatus(task.InProgress))
	completed := len(tm.GetTasksByStatus(task.Completed))
	cancelled := len(tm.GetTasksByStatus(task.Cancelled))
	
	fmt.Printf("Pending: %d\n", pending)
	fmt.Printf("In Progress: %d\n", inProgress)
	fmt.Printf("Completed: %d\n", completed)
	fmt.Printf("Cancelled: %d\n", cancelled)
	
	// Priority breakdown
	low := len(tm.GetTasksByPriority(task.Low))
	medium := len(tm.GetTasksByPriority(task.Medium))
	high := len(tm.GetTasksByPriority(task.High))
	urgent := len(tm.GetTasksByPriority(task.Urgent))
	
	fmt.Printf("\nPriority Breakdown:\n")
	fmt.Printf("Low: %d\n", low)
	fmt.Printf("Medium: %d\n", medium)
	fmt.Printf("High: %d\n", high)
	fmt.Printf("Urgent: %d\n", urgent)
	
	// Overdue tasks
	overdue := len(tm.GetOverdueTasks())
	fmt.Printf("\nOverdue Tasks: %d\n", overdue)
}

func handleDemoCommand(tm task.TaskManagerInterface) {
	color.Cyan("🔄 Go Concurrency Demo")
	fmt.Println("Adding demo tasks using goroutines and channels...")
	
	// Create demo tasks
	demoTasks := []struct {
		title       string
		description string
		priority    task.Priority
	}{
		{"Learn Goroutines", "Understanding concurrent programming", task.High},
		{"Master Channels", "Communication between goroutines", task.High},
		{"Practice Interfaces", "Go's interface system", task.Medium},
		{"Build CLI Apps", "Command line applications", task.Medium},
	}
	
	// Use channels to add tasks concurrently
	taskChan := make(chan *task.Task, len(demoTasks))
	
	// Start goroutines to add tasks
	for _, demoTask := range demoTasks {
		go func(title, desc string, priority task.Priority) {
			newTask := tm.AddTask(title, desc, priority, nil)
			taskChan <- newTask
		}(demoTask.title, demoTask.description, demoTask.priority)
	}
	
	// Collect results
	for i := 0; i < len(demoTasks); i++ {
		addedTask := <-taskChan
		color.Green("✅ Added task: %s (ID: %d)", addedTask.Title, addedTask.ID)
		time.Sleep(200 * time.Millisecond) // Simulate processing time
	}
	
	color.Cyan("🎉 Concurrency demo completed!")
}

func handleHelpCommand() {
	color.Cyan("🚀 Go Task Manager - Command Line Interface")
	fmt.Println()
	color.White("Available commands:")
	fmt.Println()
	
	color.Yellow("Interactive Mode:")
	color.White("  go run main.go                    # Start interactive mode")
	fmt.Println()
	
	color.Yellow("Command Line Mode:")
	color.White("  go run main.go add <title> <description> [priority] [due_date]")
	color.White("  go run main.go list [filter]")
	color.White("  go run main.go update <id> <status>")
	color.White("  go run main.go delete <id>")
	color.White("  go run main.go stats")
	color.White("  go run main.go demo")
	color.White("  go run main.go help")
	fmt.Println()
	
	color.Yellow("Examples:")
	color.White("  go run main.go add \"Learn Go\" \"Study Go programming\" 3 2024-12-31")
	color.White("  go run main.go list pending")
	color.White("  go run main.go update 1 completed")
	color.White("  go run main.go delete 1")
	fmt.Println()
	
	color.Yellow("Filters for list command:")
	color.White("  pending, in-progress, completed, cancelled, overdue, priority <1-4>")
}

func getPriorityColor(p task.Priority) func(string) string {
	switch p {
	case task.Low:
		return func(s string) string { return color.BlueString(s) }
	case task.Medium:
		return func(s string) string { return color.YellowString(s) }
	case task.High:
		return func(s string) string { return color.MagentaString(s) }
	case task.Urgent:
		return func(s string) string { return color.RedString(s) }
	default:
		return func(s string) string { return color.WhiteString(s) }
	}
}

func getStatusColor(s task.Status) func(string) string {
	switch s {
	case task.Pending:
		return func(s string) string { return color.YellowString(s) }
	case task.InProgress:
		return func(s string) string { return color.BlueString(s) }
	case task.Completed:
		return func(s string) string { return color.GreenString(s) }
	case task.Cancelled:
		return func(s string) string { return color.RedString(s) }
	default:
		return func(s string) string { return color.WhiteString(s) }
	}
}
