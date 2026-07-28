---
title: ToDo Notes — Edge Cases & Acceptance Criteria
version: 2.0
date_created: 2025-07-28
tags: [edge-cases, acceptance-criteria, validation]
part_of: spec-app-architecture
---

# Edge Cases & Acceptance Criteria

## Acceptance Criteria

| # | Criterion | Verification |
|---|-----------|--------------|
| AC-001 | "New Todo" opens form → submit → todo appears with status `pending` | Manual UI test |
| AC-002 | Status → `completed` sets `completedAt` to current datetime | Repository test |
| AC-003 | Past due + not completed → red visual indicator | PriorityBadge component |
| AC-004 | Filter by priority/status/tag shows only matching items | TodoFilters + useTodos hook |
| AC-005 | Global search matches across todo titles/descriptions and note content | Service Search method |
| AC-006 | Data persists across app restart from SQLite | Integration test: write → close → reopen → verify |

## Validation Criteria

| # | Criterion | Method |
|---|-----------|--------|
| V-001 | Go builds without CGO | `go build` on all platforms |
| V-002 | Wails bindings generate | `wails3 dev` → TS declarations in `frontend/bindings/` |
| V-003 | TypeScript strict compiles | `cd frontend && npx tsc --noEmit` zero errors |
| V-004 | SQLite schema creates clean | Migration test passes (`db_test.go`) |
| V-005 | CRUD end-to-end works | Repository tests pass (`go test ./...`) |
| V-006 | Data persists on restart | Integration test: write → close → reopen → verify |

## Edge Cases

### EC-001: Empty todo title

`TodoService.Create` returns `fmt.Errorf("title is required")` when title is empty string. Frontend should display validation error.

### EC-002: Tag name case-insensitive

Schema enforces `UNIQUE COLLATE NOCASE`. "Work" and "work" resolve to the same tag. `TagRepo.ListAll()` returns canonical casing from DB.

### EC-003: Delete tag cascades to associations, not items

`ON DELETE CASCADE` on `todo_tags` and `note_tags` removes associations but preserves todos/notes. Tag itself is not auto-deleted (no direct delete API exposed).

### EC-004: Due date comparison

Due date stored as `YYYY-MM-DD`. Overdue check compares at midnight boundary — a todo due 2025-07-28 is overdue on 2025-07-29 00:00.

### EC-005: DB corruption detection

If SQLite returns an integrity error, the application should detect and offer restore or fresh start. Currently handled by `NewDB` returning error — `main.go` calls `log.Fatalf`. Future improvement: graceful recovery UI.

## See Also

- [01-data-model.md](./01-data-model.md) — Schema constraints (CHECK, UNIQUE COLLATE NOCASE)
- [02-go-services.md](./02-go-services.md) — Service validation logic
- [05-testing.md](./05-testing.md) — Test coverage for edge cases
