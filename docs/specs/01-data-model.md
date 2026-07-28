---
title: ToDo Notes — Data Model
version: 2.0
date_created: 2025-07-28
tags: [data, sqlite, schema, models]
part_of: spec-app-architecture
---

# Data Model

## SQLite Schema

Schema is created via migration in `internal/pkg/repository/db.go` on `NewDB()` call. All tables use `CREATE TABLE IF NOT EXISTS` for idempotent runs.

```sql
-- Todos table
CREATE TABLE todos (
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

-- Notes table
CREATE TABLE notes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT NOT NULL DEFAULT '',
    content    TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f','now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f','now'))
);

-- Tags table (case-insensitive unique names)
CREATE TABLE tags (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE COLLATE NOCASE
);

-- Many-to-many: todo ↔ tag
CREATE TABLE todo_tags (
    todo_id INTEGER NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY(todo_id, tag_id)
);

-- Many-to-many: note ↔ tag
CREATE TABLE note_tags (
    note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY(note_id, tag_id)
);

-- Indexes
CREATE INDEX idx_todos_status   ON todos(status);
CREATE INDEX idx_todos_priority ON todos(priority);
CREATE INDEX idx_todos_due_date ON todos(due_date);
```

### Database connection settings

- Journal mode: `WAL` (Write-Ahead Logging) for concurrent reads
- Foreign keys: enabled (`_foreign_keys=on`)
- Driver: `modernc.org/sqlite` (pure Go, no CGO)
- Path: `$XDG_DATA_HOME/sample-wails-react-golang/app.db` (resolved via `github.com/adrg/xdg`)

## Go Model Structs

### Todo (`internal/pkg/model/todo.go`)

```go
type Todo struct {
    ID          int      `json:"id"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Status      string   `json:"status"`        // pending, in_progress, completed
    Priority    string   `json:"priority"`      // low, medium, high, urgent
    DueDate     *string  `json:"dueDate,omitempty"`
    CompletedAt *string  `json:"completedAt,omitempty"`
    Tags        []Tag    `json:"tags"`
    CreatedAt   string   `json:"createdAt"`
    UpdatedAt   string   `json:"updatedAt"`
}

type TodoCreateInput struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Priority    string   `json:"priority"`
    DueDate     *string  `json:"dueDate,omitempty"`
    TagNames    []string `json:"tagNames"`
}

type TodoUpdateInput struct {
    Title       *string  `json:"title,omitempty"`
    Description *string  `json:"description,omitempty"`
    Status      *string  `json:"status,omitempty"`
    Priority    *string  `json:"priority,omitempty"`
    DueDate     *string  `json:"dueDate,omitempty"`
    TagNames    []string `json:"tagNames,omitempty"`
}

type TodoFilter struct {
    Status   *string `json:"status,omitempty"`
    Priority *string `json:"priority,omitempty"`
    TagName  *string `json:"tagName,omitempty"`
}
```

### Note (`internal/pkg/model/note.go`)

```go
type Note struct {
    ID        int     `json:"id"`
    Title     string  `json:"title"`
    Content   string  `json:"content"`
    Tags      []Tag   `json:"tags"`
    CreatedAt string  `json:"createdAt"`
    UpdatedAt string  `json:"updatedAt"`
}

type NoteCreateInput struct {
    Title    string   `json:"title"`
    Content  string   `json:"content"`
    TagNames []string `json:"tagNames"`
}

type NoteUpdateInput struct {
    Title    *string  `json:"title,omitempty"`
    Content  *string  `json:"content,omitempty"`
    TagNames []string `json:"tagNames,omitempty"`
}
```

### Tag (`internal/pkg/model/tag.go`)

```go
type Tag struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
```

## See Also

- [00-overview.md](./00-overview.md) — Overview & requirements
- [02-go-services.md](./02-go-services.md) — Service API using these models
