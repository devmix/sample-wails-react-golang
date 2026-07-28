package service

import (
	"fmt"

	"sample-wails-react-golang/internal/pkg/model"
	"sample-wails-react-golang/internal/pkg/repository"
)

// TodoService handles business logic for todo operations.
type TodoService struct {
	repo *repository.TodoRepo // underlying data access layer for todos.
}

// NewTodoService creates a new TodoService with the given repository.
func NewTodoService(repo *repository.TodoRepo) *TodoService {
	return &TodoService{repo: repo}
}

// List returns todos matching the optional filter criteria.
func (s *TodoService) List(filter *model.TodoFilter) ([]model.Todo, error) {
	return s.repo.List(filter)
}

// Get returns a single todo by its ID.
func (s *TodoService) Get(id int) (*model.Todo, error) {
	return s.repo.Get(id)
}

// Create creates a new todo after validating that the title is not empty.
func (s *TodoService) Create(input model.TodoCreateInput) (*model.Todo, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	return s.repo.Create(input)
}

// Update modifies an existing todo with the provided fields.
func (s *TodoService) Update(id int, input model.TodoUpdateInput) (*model.Todo, error) {
	return s.repo.Update(id, input)
}

// Delete removes a todo by its ID.
func (s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}

// Search finds todos whose title or description match the query string.
// An empty query returns an empty slice without hitting the database.
func (s *TodoService) Search(query string) ([]model.Todo, error) {
	if query == "" {
		return []model.Todo{}, nil
	}
	return s.repo.Search(query)
}
