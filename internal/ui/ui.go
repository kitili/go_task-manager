package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"learn-go-capstone/internal/task"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// DisplayWelcome shows the welcome message
func DisplayWelcome() {
	color.Cyan("🚀 Welcome to the Go Task Manager!")
	color.Yellow("A comprehensive Go learning project that teaches:")
	color.White("  • Structs and methods")
	color.White("  • Interfaces and error handling")
	color.White("  • Concurrency with goroutines")
	color.White("  • Package organization")
	color.White("  • CLI application development")
	fmt.Println()
}

// RunInteractiveMode starts the interactive CLI
func RunInteractiveMode(tm task.TaskManagerInterface) error {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		displayMenu()
		fmt.Print("Enter your choice: ")
		
		if !scanner.Scan() {
			break
		}
		
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			addTaskInteractive(tm, scanner)
		case "2":
			listTasks(tm)
		case "3":
			updateTaskStatus(tm, scanner)
		case "4":
			deleteTask(tm, scanner)
		case "5":
			listTasksByStatus(tm, scanner)
		case "6":
			listTasksByPriority(tm, scanner)
		case "7":
			listOverdueTasks(tm)
		case "8":
			showStatistics(tm)
		case "9":
			runConcurrencyDemo(tm)
		case "0":
			color.Green("👋 Thanks for using Go Task Manager!")
			return nil
		default:
			color.Red("❌ Invalid choice. Please try again.")
		}
		
		fmt.Println()
	}
	
	return nil
}

func displayMenu() {
	color.Cyan("\n📋 Go Task Manager Menu:")
	fmt.Println("1. Add Task")
	fmt.Println("2. List All Tasks")
	fmt.Println("3. Update Task Status")
	fmt.Println("4. Delete Task")
	fmt.Println("5. Filter by Status")
	fmt.Println("6. Filter by Priority")
	fmt.Println("7. Show Overdue Tasks")
	fmt.Println("8. Show Statistics")
	fmt.Println("9. Concurrency Demo")
	fmt.Println("0. Exit")
}

func addTaskInteractive(tm task.TaskManagerInterface, scanner *bufio.Scanner) {
	fmt.Print("Enter task title: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	
	if title == "" {
		color.Red("❌ Title cannot be empty")
		return
	}
	
	fmt.Print("Enter task description: ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())
	
	priority := getPriorityFromUser(scanner)
	
	var dueDate *time.Time
	fmt.Print("Enter due date (YYYY-MM-DD) or press Enter to skip: ")
	scanner.Scan()
	dueDateStr := strings.TrimSpace(scanner.Text())
	
	if dueDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dueDateStr); err == nil {
			dueDate = &parsed
		} else {
			color.Red("❌ Invalid date format. Skipping due date.")
		}
	}
	
	newTask := tm.AddTask(title, description, priority, dueDate)
	color.Green("✅ Task added successfully!")
	color.White("Task ID: %d", newTask.ID)
}

func getPriorityFromUser(scanner *bufio.Scanner) task.Priority {
	for {
		fmt.Print("Enter priority (1=Low, 2=Medium, 3=High, 4=Urgent): ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		
		switch input {
		case "1":
			return task.Low
		case "2":
			return task.Medium
		case "3":
			return task.High
		case "4":
			return task.Urgent
		default:
			color.Red("❌ Invalid priority. Please enter 1, 2, 3, or 4.")
		}
	}
}

func listTasks(tm task.TaskManagerInterface) {
	tasks := tm.GetAllTasks()
	if len(tasks) == 0 {
		color.Yellow("📝 No tasks found.")
		return
	}
	
	displayTasksTable(tasks, "All Tasks")
}

func displayTasksTable(tasks []task.Task, title string) {
	color.Cyan("\n📋 %s", title)
	
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Priority", "Status", "Created", "Due Date"})
	table.SetBorder(true)
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	
	for _, t := range tasks {
		dueDate := "N/A"
		if t.DueDate != nil {
			dueDate = t.DueDate.Format("2006-01-02")
		}
		
		priorityColor := getPriorityColor(t.Priority)
		statusColor := getStatusColor(t.Status)
		
		table.Append([]string{
			fmt.Sprintf("%d", t.ID),
			t.Title,
			priorityColor(t.Priority.String()),
			statusColor(t.Status.String()),
			t.CreatedAt.Format("2006-01-02"),
			dueDate,
		})
	}
	
	table.Render()
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

func updateTaskStatus(tm task.TaskManagerInterface, scanner *bufio.Scanner) {
	fmt.Print("Enter task ID: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		color.Red("❌ Invalid task ID")
		return
	}
	
	_, err = tm.GetTask(id)
	if err != nil {
		color.Red("❌ %v", err)
		return
	}
	
	fmt.Print("Enter new status (1=Pending, 2=In Progress, 3=Completed, 4=Cancelled): ")
	scanner.Scan()
	statusStr := strings.TrimSpace(scanner.Text())
	
	var status task.Status
	switch statusStr {
	case "1":
		status = task.Pending
	case "2":
		status = task.InProgress
	case "3":
		status = task.Completed
	case "4":
		status = task.Cancelled
	default:
		color.Red("❌ Invalid status")
		return
	}
	
	err = tm.UpdateTaskStatus(id, status)
	if err != nil {
		color.Red("❌ %v", err)
		return
	}
	
	color.Green("✅ Task status updated successfully!")
}

func deleteTask(tm task.TaskManagerInterface, scanner *bufio.Scanner) {
	fmt.Print("Enter task ID to delete: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())
	
	id, err := strconv.Atoi(idStr)
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

func listTasksByStatus(tm task.TaskManagerInterface, scanner *bufio.Scanner) {
	fmt.Print("Enter status (1=Pending, 2=In Progress, 3=Completed, 4=Cancelled): ")
	scanner.Scan()
	statusStr := strings.TrimSpace(scanner.Text())
	
	var status task.Status
	switch statusStr {
	case "1":
		status = task.Pending
	case "2":
		status = task.InProgress
	case "3":
		status = task.Completed
	case "4":
		status = task.Cancelled
	default:
		color.Red("❌ Invalid status")
		return
	}
	
	tasks := tm.GetTasksByStatus(status)
	displayTasksTable(tasks, fmt.Sprintf("Tasks with Status: %s", status.String()))
}

func listTasksByPriority(tm task.TaskManagerInterface, scanner *bufio.Scanner) {
	fmt.Print("Enter priority (1=Low, 2=Medium, 3=High, 4=Urgent): ")
	scanner.Scan()
	priorityStr := strings.TrimSpace(scanner.Text())
	
	var priority task.Priority
	switch priorityStr {
	case "1":
		priority = task.Low
	case "2":
		priority = task.Medium
	case "3":
		priority = task.High
	case "4":
		priority = task.Urgent
	default:
		color.Red("❌ Invalid priority")
		return
	}
	
	tasks := tm.GetTasksByPriority(priority)
	displayTasksTable(tasks, fmt.Sprintf("Tasks with Priority: %s", priority.String()))
}

func listOverdueTasks(tm task.TaskManagerInterface) {
	tasks := tm.GetOverdueTasks()
	displayTasksTable(tasks, "Overdue Tasks")
}

func showStatistics(tm task.TaskManagerInterface) {
	allTasks := tm.GetAllTasks()
	
	if len(allTasks) == 0 {
		color.Yellow("📊 No tasks to show statistics for.")
		return
	}
	
	color.Cyan("\n📊 Task Manager Statistics")
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

// runConcurrencyDemo demonstrates Go's concurrency features
func runConcurrencyDemo(tm task.TaskManagerInterface) {
	color.Cyan("\n🔄 Go Concurrency Demo")
	fmt.Println("This demo shows goroutines and channels in action...")
	
	// Create some demo tasks
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
	var addedTasks []*task.Task
	for i := 0; i < len(demoTasks); i++ {
		addedTask := <-taskChan
		addedTasks = append(addedTasks, addedTask)
		color.Green("✅ Added task: %s (ID: %d)", addedTask.Title, addedTask.ID)
		time.Sleep(200 * time.Millisecond) // Simulate processing time
	}
	
	color.Cyan("\n🎉 Concurrency demo completed!")
	color.White("Added %d tasks using goroutines and channels", len(addedTasks))
}
