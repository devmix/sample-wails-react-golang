---
title: ToDo Notes — Test Strategy
version: 2.0
date_created: 2025-07-28
tags: [testing, go-test, repository]
part_of: spec-app-architecture
---

# Test Strategy

## Current State

Only Go tests exist. No frontend test framework is configured.

| Layer | Framework | Approach | Status |
|-------|-----------|----------|--------|
| Go Repository Unit | `testing` (stdlib) | Temp SQLite DB via `t.TempDir()` per test | ✅ Implemented |
| Frontend Unit | — | — | ❌ Not implemented |
| E2E | — | — | ❌ Not implemented |

## Go Tests (`internal/pkg/repository/*_test.go`)

### Test files

- `db_test.go` — DB connection and migration tests
- `todo_repo_test.go` — TodoRepo CRUD + tag association
- `note_repo_test.go` — NoteRepo CRUD + tag association
- `tag_repo_test.go` — TagRepo list operations

### Pattern

Each test creates an isolated SQLite database in a temporary directory:

```go
func TestTodoRepo_Create(t *testing.T) {
    tmpDir := t.TempDir()
    db, err := NewDB(filepath.Join(tmpDir, "test.db"))
    require.NoError(t, err)
    defer db.Close()

    repo := NewTodoRepo(db)
    // ... assertions
}
```

No shared state between test runs. No external dependencies required.

### Running tests

```bash
go test ./...
```

From repository root. Tests run in parallel by default (`-parallel` flag supported).

## Coverage Target

80%+ on `internal/pkg/repository/` and `internal/app/service/` packages.

Check coverage:
```bash
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## Future Considerations

### Frontend tests (not yet implemented)

If frontend testing is added, the recommended approach would be:
- Vitest + React Testing Library for component/unit tests
- Mock Wails bindings (`@bindings/*`) with vi.mock()
- Focus on hooks that normalize data and handle error states

### E2E tests (not yet implemented)

For critical user flows:
- Playwright in Wails test mode, or
- Manual smoke tests documented in a checklist

## See Also

- [02-go-services.md](./02-go-services.md) — Service layer under test
- [06-edge-cases.md](./06-edge-cases.md) — Edge cases to cover in tests
