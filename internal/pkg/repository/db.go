// Package repository provides data access layers backed by a SQLite database.
// It handles schema migrations and CRUD operations for todos, notes, and tags.
package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps a sql.DB connection and manages the application's SQLite database.
type DB struct {
	conn *sql.DB // underlying database connection.
}

// NewDB opens a SQLite database at the given path, creates the directory if needed,
// and runs schema migrations. It configures WAL journal mode and foreign keys.
func NewDB(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

// Conn returns the underlying sql.DB connection for direct query access.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Close closes the database connection. Callers should defer this after opening.
func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		description TEXT DEFAULT '',
		status      TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','in_progress','completed')),
		priority    TEXT NOT NULL DEFAULT 'medium' CHECK(priority IN ('low','medium','high','urgent')),
		due_date    TEXT DEFAULT NULL,
		completed_at TEXT DEFAULT NULL,
		created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f','now')),
		updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f','now'))
	);

	CREATE TABLE IF NOT EXISTS notes (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		title      TEXT NOT NULL DEFAULT '',
		content    TEXT NOT NULL DEFAULT '',
		created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f','now')),
		updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f','now'))
	);

	CREATE TABLE IF NOT EXISTS tags (
		id   INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE COLLATE NOCASE
	);

	CREATE TABLE IF NOT EXISTS todo_tags (
		todo_id INTEGER NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
		tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY(todo_id, tag_id)
	);

	CREATE TABLE IF NOT EXISTS note_tags (
		note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
		tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY(note_id, tag_id)
	);

	CREATE INDEX IF NOT EXISTS idx_todos_status   ON todos(status);
	CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
	CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
	`

	_, err := db.conn.Exec(schema)
	return err
}
