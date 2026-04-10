package models

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

// TodoItem is the core data model for a to-do entry.
type TodoItem struct {
	ID        string    `json:"id"`
	Task      string    `json:"task"`
	DueDate   time.Time `json:"due_date"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ----------------------------------------------------------------------------

// CreateTodoRequest is the request body for creating a new to-do item.
type CreateTodoRequest struct {
	Task    string `json:"task"`
	DueDate string `json:"due_date"` // expected format: "2006-01-02"
}

// ----------------------------------------------------------------------------

// UpdateTodoRequest is the request body for updating an existing to-do item.
type UpdateTodoRequest struct {
	Task    *string `json:"task"`
	DueDate *string `json:"due_date"` // expected format: "2006-01-02"
	Status  *Status `json:"status"`
}

// ----------------------------------------------------------------------------
