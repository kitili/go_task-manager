package task

import (
	"testing"
)

func TestNewTaskManager(t *testing.T) {
	tm := NewTaskManager()
	if tm == nil {
		t.Fatal("Expected TaskManager to be created")
	}
	if tm.nextID != 1 {
		t.Errorf("Expected nextID to be 1, got %d", tm.nextID)
	}
	if len(tm.tasks) != 0 {
		t.Errorf("Expected empty tasks slice, got %d tasks", len(tm.tasks))
	}
}

func TestAddTask(t *testing.T) {
	tm := NewTaskManager()
	
	task := tm.AddTask("Test Task", "A test task", Medium, nil)
	
	if task == nil {
		t.Fatal("Expected task to be created")
	}
	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", task.Title)
	}
}

func TestGetTask(t *testing.T) {
	tm := NewTaskManager()
	tm.AddTask("Test Task", "Description", Medium, nil)
	
	task, err := tm.GetTask(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", task.Title)
	}
	
	// Test non-existent task
	_, err = tm.GetTask(999)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func BenchmarkAddTask(b *testing.B) {
	tm := NewTaskManager()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.AddTask("Benchmark Task", "Description", Medium, nil)
	}
}
