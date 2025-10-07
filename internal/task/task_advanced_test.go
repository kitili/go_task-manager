package task

import (
	"fmt"
	"testing"
	"time"

	testutils "learn-go-capstone/internal/testing"
)

func TestTaskCreation(t *testing.T) {
	// Test basic task creation
	task := Task{
		ID:          1,
		Title:       "Test Task",
		Description: "A test task",
		Priority:    High,
		Status:      Pending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	testutils.AssertEqual(t, 1, task.ID, "Task ID should be 1")
	testutils.AssertEqual(t, "Test Task", task.Title, "Task title should match")
	testutils.AssertEqual(t, "A test task", task.Description, "Task description should match")
	testutils.AssertEqual(t, High, task.Priority, "Task priority should be High")
	testutils.AssertEqual(t, Pending, task.Status, "Task status should be Pending")
}

func TestPriorityEnum(t *testing.T) {
	// Test priority enum values
	testutils.AssertEqual(t, 1, int(Low), "Low priority should be 1")
	testutils.AssertEqual(t, 2, int(Medium), "Medium priority should be 2")
	testutils.AssertEqual(t, 3, int(High), "High priority should be 3")
	testutils.AssertEqual(t, 4, int(Urgent), "Urgent priority should be 4")
}

func TestStatusEnum(t *testing.T) {
	// Test status enum values
	testutils.AssertEqual(t, 0, int(Pending), "Pending status should be 0")
	testutils.AssertEqual(t, 1, int(InProgress), "InProgress status should be 1")
	testutils.AssertEqual(t, 2, int(Completed), "Completed status should be 2")
	testutils.AssertEqual(t, 3, int(Cancelled), "Cancelled status should be 3")
}

func TestAddTaskWithDueDate(t *testing.T) {
	tm := NewTaskManager()
	
	dueDate := time.Now().Add(24 * time.Hour)
	task := tm.AddTask("Task with Due Date", "A task with due date", Medium, &dueDate)
	
	testutils.AssertNotNil(t, task, "Added task should not be nil")
	testutils.AssertNotNil(t, task.DueDate, "Due date should not be nil")
	testutils.AssertEqual(t, dueDate.Unix(), task.DueDate.Unix(), "Due date should match")
}

func TestUpdateTaskStatus(t *testing.T) {
	tm := NewTaskManager()
	
	// Add a task
	task := tm.AddTask("Test Task", "A test task", High, nil)
	
	// Update status
	err := tm.UpdateTaskStatus(task.ID, InProgress)
	testutils.AssertNoError(t, err, "Should not error when updating status")
	
	// Verify status was updated
	updatedTask, err := tm.GetTask(task.ID)
	testutils.AssertNoError(t, err, "Should not error when getting updated task")
	testutils.AssertEqual(t, InProgress, updatedTask.Status, "Status should be updated to InProgress")
}

func TestUpdateNonExistentTaskStatus(t *testing.T) {
	tm := NewTaskManager()
	
	// Try to update non-existent task
	err := tm.UpdateTaskStatus(999, InProgress)
	testutils.AssertError(t, err, "Should error when updating non-existent task")
}

func TestGetTasksByStatus(t *testing.T) {
	tm := NewTaskManager()
	
	// Add tasks with different statuses
	task1 := tm.AddTask("Task 1", "First task", Low, nil)
	task2 := tm.AddTask("Task 2", "Second task", Medium, nil)
	
	// Update one task to InProgress
	tm.UpdateTaskStatus(task1.ID, InProgress)
	
	// Get pending tasks
	pendingTasks := tm.GetTasksByStatus(Pending)
	testutils.AssertEqual(t, 1, len(pendingTasks), "Should have 1 pending task")
	testutils.AssertEqual(t, task2.ID, pendingTasks[0].ID, "Pending task should be task2")
	
	// Get in-progress tasks
	inProgressTasks := tm.GetTasksByStatus(InProgress)
	testutils.AssertEqual(t, 1, len(inProgressTasks), "Should have 1 in-progress task")
	testutils.AssertEqual(t, task1.ID, inProgressTasks[0].ID, "In-progress task should be task1")
}

func TestGetTasksByPriority(t *testing.T) {
	tm := NewTaskManager()
	
	// Add tasks with different priorities
	tm.AddTask("Task 1", "Low priority task", Low, nil)
	tm.AddTask("Task 2", "High priority task", High, nil)
	tm.AddTask("Task 3", "Medium priority task", Medium, nil)
	
	// Get high priority tasks
	highPriorityTasks := tm.GetTasksByPriority(High)
	testutils.AssertEqual(t, 1, len(highPriorityTasks), "Should have 1 high priority task")
	testutils.AssertEqual(t, "Task 2", highPriorityTasks[0].Title, "High priority task should be Task 2")
}

func TestSortTasksByPriority(t *testing.T) {
	tm := NewTaskManager()
	
	// Add tasks with different priorities
	tm.AddTask("Low Task", "Low priority", Low, nil)
	tm.AddTask("Urgent Task", "Urgent priority", Urgent, nil)
	tm.AddTask("Medium Task", "Medium priority", Medium, nil)
	tm.AddTask("High Task", "High priority", High, nil)
	
	// Sort by priority
	sortedTasks := tm.SortTasksByPriority()
	testutils.AssertEqual(t, 4, len(sortedTasks), "Should have 4 tasks")
	
	// Check that tasks are sorted by priority (Urgent, High, Medium, Low)
	testutils.AssertEqual(t, Urgent, sortedTasks[0].Priority, "First task should be Urgent")
	testutils.AssertEqual(t, High, sortedTasks[1].Priority, "Second task should be High")
	testutils.AssertEqual(t, Medium, sortedTasks[2].Priority, "Third task should be Medium")
	testutils.AssertEqual(t, Low, sortedTasks[3].Priority, "Fourth task should be Low")
}

func TestGetOverdueTasks(t *testing.T) {
	tm := NewTaskManager()
	
	// Add tasks with different due dates
	pastDate := time.Now().Add(-24 * time.Hour) // Yesterday
	futureDate := time.Now().Add(24 * time.Hour) // Tomorrow
	
	tm.AddTask("Overdue Task", "This task is overdue", High, &pastDate)
	tm.AddTask("Future Task", "This task is due in the future", Medium, &futureDate)
	tm.AddTask("No Due Date Task", "This task has no due date", Low, nil)
	
	// Get overdue tasks
	overdueTasks := tm.GetOverdueTasks()
	testutils.AssertEqual(t, 1, len(overdueTasks), "Should have 1 overdue task")
	testutils.AssertEqual(t, "Overdue Task", overdueTasks[0].Title, "Overdue task should be 'Overdue Task'")
}


func TestTaskManagerConcurrency(t *testing.T) {
	tm := NewTaskManager()
	
	// Test concurrent task creation
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(i int) {
			task := tm.AddTask(fmt.Sprintf("Concurrent Task %d", i), "A concurrent task", Medium, nil)
			testutils.AssertNotNil(t, task, "Concurrent task should not be nil")
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all tasks were created
	tasks := tm.GetAllTasks()
	testutils.AssertEqual(t, 10, len(tasks), "Should have 10 concurrent tasks")
}
