package search

import (
	"time"
)

// SearchQuery represents a search query with various filters
type SearchQuery struct {
	Query       string    `json:"query"`        // Text search query
	UserID      *int      `json:"user_id"`      // Filter by user ID
	Status      *int      `json:"status"`       // Filter by status
	Priority    *int      `json:"priority"`     // Filter by priority
	CategoryID  *int      `json:"category_id"`  // Filter by category ID
	TagNames    []string  `json:"tag_names"`    // Filter by tag names
	DateFrom    *time.Time `json:"date_from"`   // Filter by creation date from
	DateTo      *time.Time `json:"date_to"`     // Filter by creation date to
	DueDateFrom *time.Time `json:"due_date_from"` // Filter by due date from
	DueDateTo   *time.Time `json:"due_date_to"`   // Filter by due date to
	IsOverdue   *bool     `json:"is_overdue"`   // Filter by overdue status
	Limit       int       `json:"limit"`        // Limit number of results
	Offset      int       `json:"offset"`       // Offset for pagination
	SortBy      string    `json:"sort_by"`      // Sort field (title, created_at, due_date, priority)
	SortOrder   string    `json:"sort_order"`   // Sort order (asc, desc)
}

// SearchResult represents a search result with metadata
type SearchResult struct {
	Tasks      []TaskResult `json:"tasks"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
	TotalPages int          `json:"total_pages"`
	Query      SearchQuery  `json:"query"`
}

// TaskResult represents a task in search results
type TaskResult struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	UserID      *int      `json:"user_id,omitempty"`
	CategoryID  *int      `json:"category_id,omitempty"`
	IsArchived  bool      `json:"is_archived"`
	// Additional fields for search results
	Category    *CategoryResult `json:"category,omitempty"`
	Tags        []TagResult     `json:"tags,omitempty"`
	User        *UserResult     `json:"user,omitempty"`
	RelevanceScore float64      `json:"relevance_score,omitempty"`
}

// CategoryResult represents a category in search results
type CategoryResult struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// TagResult represents a tag in search results
type TagResult struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// UserResult represents a user in search results
type UserResult struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// FilterOptions represents available filter options
type FilterOptions struct {
	Statuses   []StatusOption   `json:"statuses"`
	Priorities []PriorityOption `json:"priorities"`
	Categories []CategoryOption `json:"categories"`
	Tags       []TagOption      `json:"tags"`
	Users      []UserOption     `json:"users"`
}

// StatusOption represents a status filter option
type StatusOption struct {
	Value int    `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// PriorityOption represents a priority filter option
type PriorityOption struct {
	Value int    `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// CategoryOption represents a category filter option
type CategoryOption struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Count int    `json:"count"`
}

// TagOption represents a tag filter option
type TagOption struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Count int    `json:"count"`
}

// UserOption represents a user filter option
type UserOption struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Count    int    `json:"count"`
}

// SearchStats represents search statistics
type SearchStats struct {
	TotalTasks     int `json:"total_tasks"`
	SearchResults  int `json:"search_results"`
	SearchTime     int `json:"search_time_ms"`
	FilteredBy     map[string]int `json:"filtered_by"`
}

// DefaultSearchQuery returns a default search query
func DefaultSearchQuery() SearchQuery {
	return SearchQuery{
		Limit:     20,
		Offset:    0,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// Validate validates the search query
func (sq *SearchQuery) Validate() error {
	// Set defaults
	if sq.Limit <= 0 {
		sq.Limit = 20
	}
	if sq.Offset < 0 {
		sq.Offset = 0
	}
	if sq.SortBy == "" {
		sq.SortBy = "created_at"
	}
	if sq.SortOrder == "" {
		sq.SortOrder = "desc"
	}
	
	// Validate sort fields
	validSortFields := map[string]bool{
		"title":      true,
		"created_at": true,
		"updated_at": true,
		"due_date":   true,
		"priority":   true,
		"status":     true,
	}
	if !validSortFields[sq.SortBy] {
		sq.SortBy = "created_at"
	}
	
	// Validate sort order
	if sq.SortOrder != "asc" && sq.SortOrder != "desc" {
		sq.SortOrder = "desc"
	}
	
	return nil
}
