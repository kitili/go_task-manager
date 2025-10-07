package main

import (
	"log"
	"os"

	"learn-go-capstone/cmd"
	"learn-go-capstone/internal/task"
	"learn-go-capstone/internal/ui"
)

func main() {
	// Initialize the task manager
	taskManager := task.NewTaskManager()
	
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