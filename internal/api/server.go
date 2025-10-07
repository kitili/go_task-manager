package api

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"learn-go-capstone/internal/auth"
	"learn-go-capstone/internal/task"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	handler *Handler
}

// NewServer creates a new API server
func NewServer(
	taskManager task.TaskManagerInterface,
	userManager *task.UserManager,
	categoryManager *task.CategoryManager,
	dependencyManager *task.DependencyManager,
	searchManager *task.SearchManager,
	exportManager *task.ExportManager,
	notificationManager *task.NotificationManager,
	authService *auth.AuthService,
) *Server {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(CORSMiddleware())
	router.Use(LoggingMiddleware())
	router.Use(ErrorHandlerMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(RateLimitMiddleware())

	// Create handler
	handler := NewHandler(
		taskManager,
		userManager,
		categoryManager,
		dependencyManager,
		searchManager,
		exportManager,
		notificationManager,
		authService,
	)

	return &Server{
		router:  router,
		handler: handler,
	}
}

// SetupRoutes sets up all the API routes
func (s *Server) SetupRoutes() {
	// Health check (no auth required)
	s.router.GET("/health", s.handler.HealthCheck)

	// API documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Public routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", s.handler.Register)
			auth.POST("/login", s.handler.Login)
		}

		// Protected routes (auth required)
		protected := v1.Group("")
		protected.Use(AuthMiddleware(s.handler.authService))
		{
			// Task routes
			tasks := protected.Group("/tasks")
			{
				tasks.POST("", s.handler.CreateTask)
				tasks.GET("", s.handler.GetTasks)
				tasks.GET("/:id", s.handler.GetTask)
				tasks.PUT("/:id/status", s.handler.UpdateTaskStatus)
				tasks.DELETE("/:id", s.handler.DeleteTask)
				tasks.POST("/search", s.handler.SearchTasks)
			}

			// Category routes
			categories := protected.Group("/categories")
			{
				categories.POST("", s.handler.CreateCategory)
				categories.GET("", s.handler.GetCategories)
				categories.GET("/:id", s.handler.GetCategory)
				categories.PUT("/:id", s.handler.UpdateCategory)
				categories.DELETE("/:id", s.handler.DeleteCategory)
			}

			// Tag routes
			tags := protected.Group("/tags")
			{
				tags.POST("", s.handler.CreateTag)
				tags.GET("", s.handler.GetTags)
				tags.GET("/:id", s.handler.GetTag)
				tags.PUT("/:id", s.handler.UpdateTag)
				tags.DELETE("/:id", s.handler.DeleteTag)
			}

			// Dependency routes
			dependencies := protected.Group("/dependencies")
			{
				dependencies.POST("", s.handler.CreateDependency)
				dependencies.GET("/task/:id", s.handler.GetTaskDependencies)
				dependencies.DELETE("/:id", s.handler.DeleteDependency)
			}

			// Export/Import routes
			export := protected.Group("/export")
			{
				export.POST("/tasks", s.handler.ExportTasks)
				export.POST("/import", s.handler.ImportTasks)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", s.handler.GetNotifications)
				notifications.POST("", s.handler.CreateNotification)
				notifications.PUT("/:id/read", s.handler.MarkNotificationAsRead)
				notifications.GET("/stats", s.handler.GetNotificationStats)
			}

			// Statistics routes
			stats := protected.Group("/statistics")
			{
				stats.GET("", s.handler.GetStatistics)
			}
		}
	}
}

// Run starts the server
func (s *Server) Run(port string) {
	log.Printf("Starting API server on port %s", port)
	log.Printf("API documentation available at: http://localhost:%s/swagger/index.html", port)
	
	if err := s.router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// GetRouter returns the Gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
