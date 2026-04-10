package store

import (
	"errors"
	"sort"
	"sync"

	"todo-app/models"
)

// Sentinel errors returned by store operations.
var (
	ErrNotFound = errors.New("todo item not found")
)

// TodoStore is a thread-safe in-memory store for to-do items.
type TodoStore struct {
	mu    sync.RWMutex
	items map[string]*models.TodoItem
}

// ----------------------------------------------------------------------------

// NewTodoStore creates and returns an empty TodoStore.
func NewTodoStore() *TodoStore {
	return &TodoStore{
		items: make(map[string]*models.TodoItem),
	}
}

// ----------------------------------------------------------------------------

// Create inserts a new TodoItem into the store.
func (s *TodoStore) Create(item *models.TodoItem) *models.TodoItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store a copy to prevent external mutation.
	stored := *item
	s.items[stored.ID] = &stored

	// Return a second copy so the caller cannot mutate the stored value.
	result := stored
	return &result
}

// ----------------------------------------------------------------------------

// GetByID retrieves a single TodoItem by its ID.
func (s *TodoStore) GetByID(id string) (*models.TodoItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}

	clone := *item
	return &clone, nil
}

// ----------------------------------------------------------------------------

// ListOptions controls filtering behaviour for the List call.
type ListOptions struct {
	IncludeCompleted bool
}

// ----------------------------------------------------------------------------

// List returns all to-do items sorted by DueDate (earliest first).
// By default it excludes completed items; pass IncludeCompleted: true to get them too.
func (s *TodoStore) List(opts ListOptions) []*models.TodoItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.TodoItem, 0, len(s.items))
	for _, item := range s.items {
		if !opts.IncludeCompleted && item.Status == models.StatusCompleted {
			continue
		}
		clone := *item
		result = append(result, &clone)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DueDate.Before(result[j].DueDate)
	})

	return result
}

// ----------------------------------------------------------------------------

// Update applies changes to an existing TodoItem and returns the updated copy.
func (s *TodoStore) Update(item *models.TodoItem) (*models.TodoItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[item.ID]; !ok {
		return nil, ErrNotFound
	}

	clone := *item
	s.items[clone.ID] = &clone
	return &clone, nil
}

// ----------------------------------------------------------------------------

// Delete removes a TodoItem from the store by ID.
func (s *TodoStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}

	delete(s.items, id)
	return nil
}

// ----------------------------------------------------------------------------
