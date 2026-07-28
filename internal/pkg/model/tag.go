// Package model defines the data structures used across the application layers.
package model

// Tag represents a named label that can be attached to todos and notes.
type Tag struct {
	ID   int    `json:"id"`   // unique database identifier.
	Name string `json:"name"` // human-readable tag name (case-insensitive).
}
