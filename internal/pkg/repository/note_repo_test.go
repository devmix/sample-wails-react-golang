package repository

import (
	"strings"
	"testing"
	"time"

	"sample-wails-react-golang/internal/pkg/model"
)

func TestNoteRepo_CreateAndGet(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	input := model.NoteCreateInput{
		Title:   "Test Note",
		Content: "Some content here",
	}

	note, err := repo.Create(input)
	requireNoError(t, err)
	if note.Title != "Test Note" {
		t.Errorf("expected title 'Test Note', got %q", note.Title)
	}

	fetched, err := repo.Get(note.ID)
	requireNoError(t, err)
	if fetched.Content != "Some content here" {
		t.Errorf("expected matching content, got %q", fetched.Content)
	}
}

func TestNoteRepo_GetNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	_, err := repo.Get(999)
	requireError(t, err, "Get non-existent note")
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestNoteRepo_CreateWithTags(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	input := model.NoteCreateInput{
		Title:    "Tagged Note",
		Content:  "Important note",
		TagNames: []string{"idea", "work"},
	}

	note, err := repo.Create(input)
	requireNoError(t, err)
	if len(note.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(note.Tags))
	}
}

func TestNoteRepo_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	note, err := repo.Create(model.NoteCreateInput{Title: "Old Title", Content: "Old content"})
	requireNoError(t, err)

	newContent := "Updated content"
	updated, err := repo.Update(note.ID, model.NoteUpdateInput{
		Content: &newContent,
	})
	requireNoError(t, err)
	if updated.Content != "Updated content" {
		t.Errorf("expected 'Updated content', got %q", updated.Content)
	}
}

func TestNoteRepo_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	note, err := repo.Create(model.NoteCreateInput{Title: "Delete Me", Content: "content"})
	requireNoError(t, err)

	err = repo.Delete(note.ID)
	requireNoError(t, err)

	_, err = repo.Get(note.ID)
	requireError(t, err, "Get deleted note")
}

func TestNoteRepo_DeleteNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	err := repo.Delete(999)
	requireError(t, err, "Delete non-existent note")
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestNoteRepo_ListOrdering(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	repo.Create(model.NoteCreateInput{Title: "First", Content: "content1"})
	note2, err := repo.Create(model.NoteCreateInput{Title: "Second", Content: "content2"})
	requireNoError(t, err)

	time.Sleep(1500 * time.Millisecond)

	_, err = repo.Update(note2.ID, model.NoteUpdateInput{
		Title: strPtr("Updated Second"),
	})
	requireNoError(t, err)

	notes, err := repo.List()
	requireNoError(t, err)
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
	if notes[0].Title != "Updated Second" {
		t.Errorf("expected most recently updated note first, got %q", notes[0].Title)
	}
}

func TestNoteRepo_Search(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	repo.Create(model.NoteCreateInput{Title: "Meeting Notes", Content: "Discuss project"})
	repo.Create(model.NoteCreateInput{Title: "Shopping List", Content: "Buy milk"})

	results, err := repo.Search("meeting")
	requireNoError(t, err)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Meeting Notes" {
		t.Errorf("expected 'Meeting Notes', got %q", results[0].Title)
	}
}

func TestNoteRepo_SearchByContent(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepo(db)

	repo.Create(model.NoteCreateInput{Title: "Random Title", Content: "Remember to buy eggs"})

	results, err := repo.Search("eggs")
	requireNoError(t, err)
	if len(results) != 1 {
		t.Fatalf("expected 1 result searching content, got %d", len(results))
	}
}

func strPtr(s string) *string {
	return &s
}
