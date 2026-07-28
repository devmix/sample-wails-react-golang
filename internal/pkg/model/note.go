// Package model defines the data structures used across the application layers.
package model

// Note represents a free-form text note with optional tags.
type Note struct {
	ID        int    `json:"id"`        // unique database identifier.
	Title     string `json:"title"`     // short title of the note.
	Content   string `json:"content"`   // full body content of the note.
	Tags      []Tag  `json:"tags"`      // associated tags for categorization.
	CreatedAt string `json:"createdAt"` // creation timestamp in ISO-8601 format.
	UpdatedAt string `json:"updatedAt"` // last modification timestamp in ISO-8601 format.
}

// NoteCreateInput carries the fields required to create a new note.
type NoteCreateInput struct {
	Title    string   `json:"title"`    // title of the new note.
	Content  string   `json:"content"`  // body content (required).
	TagNames []string `json:"tagNames"` // tag names to associate with the note.
}

// NoteUpdateInput carries partial fields for updating an existing note.
// Only non-nil fields are applied during the update.
type NoteUpdateInput struct {
	Title    *string  `json:"title,omitempty"`    // new title value.
	Content  *string  `json:"content,omitempty"`  // new body content.
	TagNames []string `json:"tagNames,omitempty"` // replaces existing tags if non-empty.
}
