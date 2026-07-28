package service

import (
	"fmt"

	"sample-wails-react-golang/internal/pkg/model"
	"sample-wails-react-golang/internal/pkg/repository"
)

// NoteService handles business logic for note operations.
type NoteService struct {
	repo *repository.NoteRepo // underlying data access layer for notes.
}

// NewNoteService creates a new NoteService with the given repository.
func NewNoteService(repo *repository.NoteRepo) *NoteService {
	return &NoteService{repo: repo}
}

// List returns all notes ordered by most recently updated.
func (s *NoteService) List() ([]model.Note, error) {
	return s.repo.List()
}

// Get returns a single note by its ID.
func (s *NoteService) Get(id int) (*model.Note, error) {
	return s.repo.Get(id)
}

// Create creates a new note after validating that the content is not empty.
func (s *NoteService) Create(input model.NoteCreateInput) (*model.Note, error) {
	if input.Content == "" {
		return nil, fmt.Errorf("content is required")
	}
	return s.repo.Create(input)
}

// Update modifies an existing note with the provided fields.
func (s *NoteService) Update(id int, input model.NoteUpdateInput) (*model.Note, error) {
	return s.repo.Update(id, input)
}

// Delete removes a note by its ID.
func (s *NoteService) Delete(id int) error {
	return s.repo.Delete(id)
}

// Search finds notes whose title or content match the query string.
// An empty query returns an empty slice without hitting the database.
func (s *NoteService) Search(query string) ([]model.Note, error) {
	if query == "" {
		return []model.Note{}, nil
	}
	return s.repo.Search(query)
}
