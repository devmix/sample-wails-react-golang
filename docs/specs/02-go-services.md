---
title: ToDo Notes — Go Services & Repository Layer
version: 2.0
date_created: 2025-07-28
tags: [go, services, repository, wails-bindings]
part_of: spec-app-architecture
---

# Go Services & Repository Layer

## Architecture

```
frontend/src/hooks/ → @bindings/* (auto-generated) → internal/app/service/ → internal/pkg/repository/ → SQLite
```

Services are registered in `main.go` as Wails services, making their methods available to the frontend through auto-generated TypeScript bindings.

## Service Registration (`main.go`)

```go
app := application.New(application.Options{
    Name:        "ToDo Notes",
    Description: "A to-do and notes application",
    Services: []application.Service{
        application.NewService(service.NewTodoService(todoRepo)),
        application.NewService(service.NewNoteService(noteRepo)),
        application.NewService(service.NewTagService(tagRepo)),
        application.NewService(service.NewAppService()),
    },
    // ...
})
```

## TodoService (`internal/app/service/todo_service.go`)

| Method | Signature | Description |
|--------|-----------|-------------|
| `List` | `(filter *TodoFilter) ([]Todo, error)` | Filtered list of todos |
| `Get` | `(id int) (*Todo, error)` | Single todo by ID |
| `Create` | `(input TodoCreateInput) (*Todo, error)` | Create with title validation |
| `Update` | `(id int, input TodoUpdateInput) (*Todo, error)` | Partial update (non-nil fields only) |
| `Delete` | `(id int) error` | Remove todo by ID |
| `Search` | `(query string) ([]Todo, error)` | FTS on title/description; empty query → empty slice |

**Validation**: `Create` returns error if title is empty.

## NoteService (`internal/app/service/note_service.go`)

| Method | Signature | Description |
|--------|-----------|-------------|
| `List` | `() ([]Note, error)` | All notes, ordered by most recently updated |
| `Get` | `(id int) (*Note, error)` | Single note by ID |
| `Create` | `(input NoteCreateInput) (*Note, error)` | Create with content validation |
| `Update` | `(id int, input NoteUpdateInput) (*Note, error)` | Partial update (non-nil fields only) |
| `Delete` | `(id int) error` | Remove note by ID |
| `Search` | `(query string) ([]Note, error)` | FTS on title/content; empty query → empty slice |

**Validation**: `Create` returns error if content is empty.

## TagService (`internal/app/service/tag_service.go`)

| Method | Signature | Description |
|--------|-----------|-------------|
| `ListAll` | `() ([]Tag, error)` | All tags, ordered alphabetically by name |

## AppService (`internal/app/service/app_service.go`)

| Method | Signature | Description |
|--------|-----------|-------------|
| `Version` | `() string` | Returns application version string |

Global version is set via `service.SetVersion(v string)` at startup.

## Repository Layer (`internal/pkg/repository/`)

### DB (`db.go`)

- `NewDB(dbPath string) (*DB, error)` — opens connection, runs migration
- `Conn() *sql.DB` — raw connection access
- `Close() error` — close connection
- `migrate()` — idempotent schema creation with `CREATE TABLE IF NOT EXISTS`

### TodoRepo (`todo_repo.go`)

Data access for todos. All queries use prepared statements. Tag association handled within Create/Update via transactional insert/delete on `todo_tags`.

### NoteRepo (`note_repo.go`)

Data access for notes. Same pattern as TodoRepo with `note_tags` join table.

### TagRepo (`tag_repo.go`)

Data access for tags. Tags are case-insensitive (enforced by `COLLATE NOCASE` in schema).

## See Also

- [01-data-model.md](./01-data-model.md) — SQLite schema, Go models
- [03-frontend.md](./03-frontend.md) — Frontend hooks that call these services
- [05-testing.md](./05-testing.md) — Repository tests with temp SQLite DBs
