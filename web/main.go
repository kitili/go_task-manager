package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"learn-go-capstone/internal/task"
)

var taskManager = task.NewTaskManager()

func main() {
	// Add some demo tasks
	taskManager.AddTask("Learn Go", "Study Go programming", task.High, nil)
	taskManager.AddTask("Build Web App", "Create a web application", task.Medium, nil)
	taskManager.AddTask("Deploy to Cloud", "Deploy the application", task.Low, nil)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// API routes
	http.HandleFunc("/api/tasks", handleTasks)
	http.HandleFunc("/api/tasks/", handleTaskByID)
	http.HandleFunc("/api/stats", handleStats)

	// Web routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/add", handleAddTask)
	http.HandleFunc("/update", handleUpdateTask)
	http.HandleFunc("/delete", handleDeleteTask)

	fmt.Println("🚀 Web Server starting on http://localhost:8080")
	fmt.Println("📱 Open your browser and go to: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Task Manager</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto; 
            background: white; 
            border-radius: 15px; 
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white; 
            padding: 30px; 
            text-align: center; 
        }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .header p { font-size: 1.2em; opacity: 0.9; }
        .content { padding: 30px; }
        .add-task { 
            background: #f8f9fa; 
            padding: 25px; 
            border-radius: 10px; 
            margin-bottom: 30px;
            border-left: 5px solid #667eea;
        }
        .form-group { margin-bottom: 15px; }
        .form-group label { 
            display: block; 
            margin-bottom: 5px; 
            font-weight: 600; 
            color: #333; 
        }
        .form-group input, .form-group select, .form-group textarea { 
            width: 100%; 
            padding: 12px; 
            border: 2px solid #e1e5e9; 
            border-radius: 8px; 
            font-size: 16px;
            transition: border-color 0.3s;
        }
        .form-group input:focus, .form-group select:focus, .form-group textarea:focus { 
            outline: none; 
            border-color: #667eea; 
        }
        .btn { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white; 
            padding: 12px 25px; 
            border: none; 
            border-radius: 8px; 
            cursor: pointer; 
            font-size: 16px;
            font-weight: 600;
            transition: transform 0.2s;
        }
        .btn:hover { transform: translateY(-2px); }
        .btn-danger { background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%); }
        .btn-success { background: linear-gradient(135deg, #51cf66 0%, #40c057 100%); }
        .tasks-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); 
            gap: 20px; 
        }
        .task-card { 
            background: white; 
            border: 2px solid #e1e5e9; 
            border-radius: 12px; 
            padding: 20px; 
            transition: all 0.3s;
            position: relative;
        }
        .task-card:hover { 
            transform: translateY(-5px); 
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
        }
        .task-title { 
            font-size: 1.3em; 
            font-weight: 700; 
            margin-bottom: 10px; 
            color: #333; 
        }
        .task-description { 
            color: #666; 
            margin-bottom: 15px; 
            line-height: 1.5; 
        }
        .task-meta { 
            display: flex; 
            justify-content: space-between; 
            align-items: center; 
            margin-bottom: 15px; 
        }
        .priority, .status { 
            padding: 5px 12px; 
            border-radius: 20px; 
            font-size: 0.9em; 
            font-weight: 600; 
        }
        .priority-low { background: #e3f2fd; color: #1976d2; }
        .priority-medium { background: #fff3e0; color: #f57c00; }
        .priority-high { background: #fce4ec; color: #c2185b; }
        .priority-urgent { background: #ffebee; color: #d32f2f; }
        .status-pending { background: #fff3e0; color: #f57c00; }
        .status-in-progress { background: #e3f2fd; color: #1976d2; }
        .status-completed { background: #e8f5e8; color: #2e7d32; }
        .status-cancelled { background: #ffebee; color: #d32f2f; }
        .task-actions { 
            display: flex; 
            gap: 10px; 
        }
        .stats { 
            background: #f8f9fa; 
            padding: 20px; 
            border-radius: 10px; 
            margin-bottom: 30px;
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
        }
        .stat-item { 
            text-align: center; 
            padding: 15px; 
            background: white; 
            border-radius: 8px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.05);
        }
        .stat-number { 
            font-size: 2em; 
            font-weight: 700; 
            color: #667eea; 
        }
        .stat-label { 
            color: #666; 
            margin-top: 5px; 
        }
        .loading { 
            text-align: center; 
            padding: 50px; 
            color: #666; 
        }
        .empty-state { 
            text-align: center; 
            padding: 50px; 
            color: #666; 
        }
        .empty-state h3 { margin-bottom: 10px; color: #333; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Go Task Manager</h1>
            <p>A comprehensive Go learning project with web interface</p>
        </div>
        
        <div class="content">
            <!-- Add Task Form -->
            <div class="add-task">
                <h2>➕ Add New Task</h2>
                <form id="addTaskForm">
                    <div class="form-group">
                        <label for="title">Task Title</label>
                        <input type="text" id="title" name="title" required>
                    </div>
                    <div class="form-group">
                        <label for="description">Description</label>
                        <textarea id="description" name="description" rows="3"></textarea>
                    </div>
                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 15px;">
                        <div class="form-group">
                            <label for="priority">Priority</label>
                            <select id="priority" name="priority">
                                <option value="0">Low</option>
                                <option value="1">Medium</option>
                                <option value="2" selected>High</option>
                                <option value="3">Urgent</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="dueDate">Due Date</label>
                            <input type="date" id="dueDate" name="dueDate">
                        </div>
                    </div>
                    <button type="submit" class="btn">Add Task</button>
                </form>
            </div>

            <!-- Statistics -->
            <div class="stats" id="stats">
                <div class="loading">Loading statistics...</div>
            </div>

            <!-- Tasks -->
            <div id="tasks">
                <div class="loading">Loading tasks...</div>
            </div>
        </div>
    </div>

    <script>
        // Load data on page load
        document.addEventListener('DOMContentLoaded', function() {
            loadStats();
            loadTasks();
        });

        // Add task form handler
        document.getElementById('addTaskForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const formData = new FormData(this);
            const taskData = {
                title: formData.get('title'),
                description: formData.get('description'),
                priority: parseInt(formData.get('priority')),
                dueDate: formData.get('dueDate') || null
            };

            fetch('/api/tasks', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(taskData)
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    alert('Error: ' + data.error);
                } else {
                    this.reset();
                    loadStats();
                    loadTasks();
                }
            });
        });

        // Load statistics
        function loadStats() {
            fetch('/api/stats')
            .then(response => response.json())
            .then(data => {
                document.getElementById('stats').innerHTML = 
                    '<div class="stat-item">' +
                        '<div class="stat-number">' + data.total + '</div>' +
                        '<div class="stat-label">Total Tasks</div>' +
                    '</div>' +
                    '<div class="stat-item">' +
                        '<div class="stat-number">' + data.pending + '</div>' +
                        '<div class="stat-label">Pending</div>' +
                    '</div>' +
                    '<div class="stat-item">' +
                        '<div class="stat-number">' + data.inProgress + '</div>' +
                        '<div class="stat-label">In Progress</div>' +
                    '</div>' +
                    '<div class="stat-item">' +
                        '<div class="stat-number">' + data.completed + '</div>' +
                        '<div class="stat-label">Completed</div>' +
                    '</div>';
            });
        }

        // Load tasks
        function loadTasks() {
            fetch('/api/tasks')
            .then(response => response.json())
            .then(data => {
                if (data.length === 0) {
                    document.getElementById('tasks').innerHTML = `
                        <div class="empty-state">
                            <h3>No tasks yet</h3>
                            <p>Add your first task using the form above!</p>
                        </div>
                    `;
                    return;
                }

                const tasksHtml = data.map(function(task) {
                    return '<div class="task-card">' +
                        '<div class="task-title">' + task.title + '</div>' +
                        '<div class="task-description">' + (task.description || 'No description') + '</div>' +
                        '<div class="task-meta">' +
                            '<span class="priority priority-' + getPriorityClass(task.priority) + '">' + getPriorityText(task.priority) + '</span>' +
                            '<span class="status status-' + getStatusClass(task.status) + '">' + getStatusText(task.status) + '</span>' +
                        '</div>' +
                        '<div class="task-actions">' +
                            '<button class="btn btn-success" onclick="updateTask(' + task.id + ', 2)">Complete</button>' +
                            '<button class="btn" onclick="updateTask(' + task.id + ', 1)">In Progress</button>' +
                            '<button class="btn btn-danger" onclick="deleteTask(' + task.id + ')">Delete</button>' +
                        '</div>' +
                    '</div>';
                }).join('');

                document.getElementById('tasks').innerHTML = 
                    '<h2>📋 Your Tasks</h2>' +
                    '<div class="tasks-grid">' + tasksHtml + '</div>';
            });
        }

        // Update task status
        function updateTask(id, status) {
            fetch(`/api/tasks/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ status: status })
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    alert('Error: ' + data.error);
                } else {
                    loadStats();
                    loadTasks();
                }
            });
        }

        // Delete task
        function deleteTask(id) {
            if (confirm('Are you sure you want to delete this task?')) {
                fetch(`/api/tasks/${id}`, {
                    method: 'DELETE'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        alert('Error: ' + data.error);
                    } else {
                        loadStats();
                        loadTasks();
                    }
                });
            }
        }

        // Helper functions
        function getPriorityText(priority) {
            const priorities = ['Low', 'Medium', 'High', 'Urgent'];
            return priorities[priority] || 'Unknown';
        }

        function getPriorityClass(priority) {
            const classes = ['low', 'medium', 'high', 'urgent'];
            return classes[priority] || 'low';
        }

        function getStatusText(status) {
            const statuses = ['Pending', 'In Progress', 'Completed', 'Cancelled'];
            return statuses[status] || 'Unknown';
        }

        function getStatusClass(status) {
            const classes = ['pending', 'in-progress', 'completed', 'cancelled'];
            return classes[status] || 'pending';
        }
    </script>
</body>
</html>
`
	
	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
		DueDate     string `json:"dueDate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&taskData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var dueDate *time.Time
	if taskData.DueDate != "" {
		if parsed, err := time.Parse("2006-01-02", taskData.DueDate); err == nil {
			dueDate = &parsed
		}
	}

	newTask := taskManager.AddTask(
		taskData.Title,
		taskData.Description,
		task.Priority(taskData.Priority),
		dueDate,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tasks := taskManager.GetAllTasks()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := path[len("/api/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		var updateData struct {
			Status int `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err := taskManager.UpdateTaskStatus(id, task.Status(updateData.Status))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Task updated successfully"})

	case "DELETE":
		err := taskManager.DeleteTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	allTasks := taskManager.GetAllTasks()
	
	stats := map[string]int{
		"total":      len(allTasks),
		"pending":    len(taskManager.GetTasksByStatus(task.Pending)),
		"inProgress": len(taskManager.GetTasksByStatus(task.InProgress)),
		"completed":  len(taskManager.GetTasksByStatus(task.Completed)),
		"cancelled":  len(taskManager.GetTasksByStatus(task.Cancelled)),
		"overdue":    len(taskManager.GetOverdueTasks()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
