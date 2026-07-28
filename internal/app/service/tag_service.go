package service

import (
	"sample-wails-react-golang/internal/pkg/model"
	"sample-wails-react-golang/internal/pkg/repository"
)

// TagService handles business logic for tag operations.
type TagService struct {
	repo *repository.TagRepo // underlying data access layer for tags.
}

// NewTagService creates a new TagService with the given repository.
func NewTagService(repo *repository.TagRepo) *TagService {
	return &TagService{repo: repo}
}

// ListAll returns every tag in the database, ordered alphabetically by name.
func (s *TagService) ListAll() ([]model.Tag, error) {
	return s.repo.ListAll()
}
