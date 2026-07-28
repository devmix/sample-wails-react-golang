// Package model defines the data structures used across the application layers.
package model

// Todo represents a task item with status, priority, and optional tags.
type Todo struct {
	ID          int     `json:"id"`                    // unique database identifier.
	Title       string  `json:"title"`                 // short title of the todo.
	Description string  `json:"description"`           // detailed description text.
	Status      string  `json:"status"`                // one of pending, in_progress, completed.
	Priority    string  `json:"priority"`              // one of low, medium, high, urgent.
	DueDate     *string `json:"dueDate,omitempty"`     // optional due date as ISO-8601 string.
	CompletedAt *string `json:"completedAt,omitempty"` // timestamp when status became completed.
	Tags        []Tag   `json:"tags"`                  // associated tags for categorization.
	CreatedAt   string  `json:"createdAt"`             // creation timestamp in ISO-8601 format.
	UpdatedAt   string  `json:"updatedAt"`             // last modification timestamp in ISO-8601 format.
}

// TodoCreateInput carries the fields required to create a new todo.
type TodoCreateInput struct {
	Title       string   `json:"title"`             // title of the new todo (required).
	Description string   `json:"description"`       // optional description text.
	Priority    string   `json:"priority"`          // priority level; defaults to medium if empty.
	DueDate     *string  `json:"dueDate,omitempty"` // optional due date as ISO-8601 string.
	TagNames    []string `json:"tagNames"`          // tag names to associate with the todo.
}

// TodoUpdateInput carries partial fields for updating an existing todo.
// Only non-nil fields are applied during the update.
type TodoUpdateInput struct {
	Title       *string  `json:"title,omitempty"`       // new title value.
	Description *string  `json:"description,omitempty"` // new description text.
	Status      *string  `json:"status,omitempty"`      // new status value.
	Priority    *string  `json:"priority,omitempty"`    // new priority level.
	DueDate     *string  `json:"dueDate,omitempty"`     // new due date as ISO-8601 string.
	TagNames    []string `json:"tagNames,omitempty"`    // replaces existing tags if non-empty.
}

// TodoFilter specifies optional criteria for filtering todo lists.
type TodoFilter struct {
	Status   *string `json:"status,omitempty"`   // filter by status value.
	Priority *string `json:"priority,omitempty"` // filter by priority level.
	TagName  *string `json:"tagName,omitempty"`  // filter by associated tag name (case-insensitive).
}
