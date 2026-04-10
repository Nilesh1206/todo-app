package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"todo-app/models"
	"todo-app/store"
)

const dueDateLayout = "2026-04-09"

type TodoHandler struct {
	store *store.TodoStore
}

// ----------------------------------------------------------------------------

// NewTodoHandler creates a TodoHandler wired to the given store.
func NewTodoHandler(s *store.TodoStore) *TodoHandler {
	return &TodoHandler{store: s}
}

// ----------------------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// ----------------------------------------------------------------------------

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// ----------------------------------------------------------------------------

// generateID creates a simple unique ID using the current nanosecond timestamp.
func generateID() string {
	return time.Now().Format("20060102150405.000000000")
}

// ----------------------------------------------------------------------------

// extractIDFromPath pulls the last path segment from a URL like /todos/{id}.
func extractIDFromPath(path string) string {
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	return parts[len(parts)-1]
}

// ----------------------------------------------------------------------------

// RegisterRoutes handles the routes.
func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/todos", h.handleCollection)
	mux.HandleFunc("/todos/", h.handleItem)
}

// ----------------------------------------------------------------------------

// handleCollection dispatches POST /todos and GET /todos.
func (h *TodoHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTodo(w, r)
	case http.MethodGet:
		h.ListTodos(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ----------------------------------------------------------------------------

// handleItem dispatches GET /todos/{id}, PUT /todos/{id}, DELETE /todos/{id}.
func (h *TodoHandler) handleItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTodo(w, r)
	case http.MethodPut:
		h.UpdateTodo(w, r)
	case http.MethodDelete:
		h.DeleteTodo(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ----------------------------------------------------------------------------

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(req.Task) == "" {
		writeError(w, http.StatusBadRequest, "task field is required")
		return
	}

	if strings.TrimSpace(req.DueDate) == "" {
		writeError(w, http.StatusBadRequest, "due_date field is required (format: YYYY-MM-DD)")
		return
	}

	dueDate, err := time.Parse(dueDateLayout, req.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid due_date format, expected YYYY-MM-DD")
		return
	}

	now := time.Now()
	item := &models.TodoItem{
		ID:        generateID(),
		Task:      strings.TrimSpace(req.Task),
		DueDate:   dueDate,
		Status:    models.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created := h.store.Create(item)
	writeJSON(w, http.StatusCreated, created)
}

// ----------------------------------------------------------------------------

// GetTodo handles GET /todos/{id}
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing todo id")
		return
	}

	item, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, item)
}

// ----------------------------------------------------------------------------

// ListTodos handles GET /todos
// Query param: include_completed=true  -> includes completed items
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	includeCompleted := r.URL.Query().Get("include_completed") == "true"

	items := h.store.List(store.ListOptions{IncludeCompleted: includeCompleted})

	// Return an empty JSON array instead of null when there are no items.
	if items == nil {
		items = []*models.TodoItem{}
	}

	writeJSON(w, http.StatusOK, items)
}

// UpdateTodo handles PUT /todos/{id}
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing todo id")
		return
	}

	existing, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Apply only the fields that were provided.
	if req.Task != nil {
		if strings.TrimSpace(*req.Task) == "" {
			writeError(w, http.StatusBadRequest, "task cannot be empty")
			return
		}
		existing.Task = strings.TrimSpace(*req.Task)
	}

	if req.DueDate != nil {
		d, err := time.Parse(dueDateLayout, *req.DueDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid due_date format, expected YYYY-MM-DD")
			return
		}
		existing.DueDate = d
	}

	if req.Status != nil {
		if *req.Status != models.StatusPending && *req.Status != models.StatusCompleted {
			writeError(w, http.StatusBadRequest, "status must be 'pending' or 'completed'")
			return
		}
		existing.Status = *req.Status
	}

	existing.UpdatedAt = time.Now()

	updated, err := h.store.Update(existing)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ----------------------------------------------------------------------------

// DeleteTodo handles DELETE /todos/{id}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromPath(r.URL.Path)
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing todo id")
		return
	}

	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "todo deleted successfully"})
}

// ----------------------------------------------------------------------------
