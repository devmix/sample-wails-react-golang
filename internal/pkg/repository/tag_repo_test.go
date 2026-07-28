package repository

import (
	"sort"
	"testing"

	"sample-wails-react-golang/internal/pkg/model"
)

func TestTagRepo_ListAll(t *testing.T) {
	db := newTestDB(t)
	tagRepo := NewTagRepo(db)
	todoRepo := NewTodoRepo(db)

	todoRepo.Create(model.TodoCreateInput{Title: "Task1", TagNames: []string{"beta", "alpha"}})

	tags, err := tagRepo.ListAll()
	requireNoError(t, err)
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}

	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	sort.Strings(names)

	if names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("expected [alpha, beta], got %v", names)
	}
}

func TestTagRepo_EmptyList(t *testing.T) {
	db := newTestDB(t)
	tagRepo := NewTagRepo(db)

	tags, err := tagRepo.ListAll()
	requireNoError(t, err)
	if len(tags) != 0 {
		t.Errorf("expected empty list, got %d tags", len(tags))
	}
}
