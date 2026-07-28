package repository

import (
	"strings"
	"testing"

	"sample-wails-react-golang/internal/pkg/model"
)

func TestTodoRepo_CreateAndGet(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	title := "Test Todo"
	input := model.TodoCreateInput{
		Title:       title,
		Description: "Test description",
		Priority:    "high",
	}

	todo, err := repo.Create(input)
	requireNoError(t, err)
	if todo.Title != title {
		t.Errorf("expected title %q, got %q", title, todo.Title)
	}
	if todo.Status != "pending" {
		t.Errorf("expected status pending, got %q", todo.Status)
	}
	if todo.Priority != "high" {
		t.Errorf("expected priority high, got %q", todo.Priority)
	}

	fetched, err := repo.Get(todo.ID)
	requireNoError(t, err)
	if fetched.ID != todo.ID {
		t.Errorf("expected id %d, got %d", todo.ID, fetched.ID)
	}
}

func TestTodoRepo_GetNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	_, err := repo.Get(999)
	requireError(t, err, "Get non-existent todo")
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestTodoRepo_CreateDefaultPriority(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	input := model.TodoCreateInput{
		Title: "Default Priority Todo",
	}

	todo, err := repo.Create(input)
	requireNoError(t, err)
	if todo.Priority != "medium" {
		t.Errorf("expected default priority medium, got %q", todo.Priority)
	}
}

func TestTodoRepo_CreateWithTags(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	input := model.TodoCreateInput{
		Title:    "Tagged Todo",
		TagNames: []string{"work", "urgent"},
	}

	todo, err := repo.Create(input)
	requireNoError(t, err)
	if len(todo.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(todo.Tags))
	}

	tagNames := make([]string, len(todo.Tags))
	for i, tag := range todo.Tags {
		tagNames[i] = tag.Name
	}
	if !contains(tagNames, "work") || !contains(tagNames, "urgent") {
		t.Errorf("expected tags [work, urgent], got %v", tagNames)
	}
}

func TestTodoRepo_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	todo, err := repo.Create(model.TodoCreateInput{Title: "Original"})
	requireNoError(t, err)

	newStatus := "completed"
	updated, err := repo.Update(todo.ID, model.TodoUpdateInput{
		Status: &newStatus,
	})
	requireNoError(t, err)
	if updated.Status != "completed" {
		t.Errorf("expected status completed, got %q", updated.Status)
	}
	if updated.CompletedAt == nil {
		t.Error("expected CompletedAt to be set when marking as completed")
	}
}

func TestTodoRepo_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	todo, err := repo.Create(model.TodoCreateInput{Title: "To Delete"})
	requireNoError(t, err)

	err = repo.Delete(todo.ID)
	requireNoError(t, err)

	_, err = repo.Get(todo.ID)
	requireError(t, err, "Get deleted todo")
}

func TestTodoRepo_DeleteNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	err := repo.Delete(999)
	requireError(t, err, "Delete non-existent todo")
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestTodoRepo_ListWithFilter(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	repo.Create(model.TodoCreateInput{Title: "High Priority", Priority: "high"})
	repo.Create(model.TodoCreateInput{Title: "Low Priority", Priority: "low"})

	highPriority := "high"
	todos, err := repo.List(&model.TodoFilter{Priority: &highPriority})
	requireNoError(t, err)

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo with high priority, got %d", len(todos))
	}
	if todos[0].Title != "High Priority" {
		t.Errorf("expected 'High Priority', got %q", todos[0].Title)
	}
}

func TestTodoRepo_Search(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	repo.Create(model.TodoCreateInput{Title: "Buy groceries"})
	repo.Create(model.TodoCreateInput{Title: "Walk the dog"})

	results, err := repo.Search("groceries")
	requireNoError(t, err)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Buy groceries" {
		t.Errorf("expected 'Buy groceries', got %q", results[0].Title)
	}
}

func TestTodoRepo_SearchEmpty(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	results, err := repo.Search("")
	requireNoError(t, err)
	if len(results) != 0 {
		t.Errorf("expected empty results for empty query, got %d", len(results))
	}
}

func TestTodoRepo_ListWithTagsFilter(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepo(db)

	repo.Create(model.TodoCreateInput{Title: "Work Task", TagNames: []string{"work"}})
	repo.Create(model.TodoCreateInput{Title: "Personal Task", TagNames: []string{"personal"}})

	tagName := "work"
	todos, err := repo.List(&model.TodoFilter{TagName: &tagName})
	requireNoError(t, err)
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo with tag 'work', got %d", len(todos))
	}
	if todos[0].Title != "Work Task" {
		t.Errorf("expected 'Work Task', got %q", todos[0].Title)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
