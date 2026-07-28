---
title: ToDo Notes — Development Conventions
version: 2.0
date_created: 2025-07-28
tags: [conventions, code-style, naming, git]
part_of: spec-app-architecture
---

# Development Conventions

## Go Conventions

### Package structure

All Go code lives under `internal/`. No external imports of internal packages.

```
internal/app/service/    → business logic, Wails bindings
internal/pkg/model/      → shared data structures
internal/pkg/repository/ → data access layer
```

### Naming

- **Packages**: lowercase, no underscores (`service`, `repository`, `model`)
- **Exported types**: PascalCase with descriptive suffix (`TodoService`, `TodoRepo`, `TodoFilter`)
- **Functions**: PascalCase for exported, camelCase for unexported (`NewDB`, `migrate`)
- **Files**: snake_case with type suffix (`todo_service.go`, `todo_repo.go`, `db_test.go`)

### Error handling

- Use `fmt.Errorf("context: %w", err)` for wrapping
- Service layer returns errors to frontend — no panic/recover in service methods
- Repository layer wraps DB errors with context (`"create todo: %w"`)

### Struct fields

- JSON tags on all exported struct fields in `model/` package
- Pointer types (`*string`) for optional fields; value types for required fields
- Struct field comments document purpose and constraints

## TypeScript Conventions

### Naming

- **Components**: PascalCase (`TodoList`, `TodoItem`, `SearchBar`)
- **Hooks**: camelCase with `use` prefix (`useTodos`, `useNotes`, `useTags`)
- **Files**: Match component/hook name (`TodoForm.tsx`, `useTodos.ts`)
- **Types**: PascalCase, re-exported from bindings via `types/models.ts`

### Imports

```typescript
// Order: React → bindings → local types → components → utils
import { useState } from 'react';
import * as TodoService from '@bindings/...';
import type { Todo } from '@/types/models';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
```

### Type safety

- `strict: true` in tsconfig, with `noUnusedLocals: true`, `noImplicitAny: false`
- All hook return values typed explicitly
- Event handlers use proper React types (`React.ChangeEvent<HTMLInputElement>`)

## Git Conventions

### Commit messages

Conventional Commits format recommended:

```
type(scope): description

feat(todos): add priority filter to TodoFilters
fix(notes): handle null tags in NoteCard
chore(deps): update wails to alpha2.118
```

Types: `feat`, `fix`, `chore`, `refactor`, `docs`, `test`

### Branching

- `main` — production-ready code
- Feature branches: `feature/<short-description>`
- Bug fixes: `fix/<short-description>`

## Code Review Checklist

- [ ] No hardcoded paths (use xdg for data, relative imports in frontend)
- [ ] Prepared statements used for all SQL queries
- [ ] Error wrapping with context (`%w`)
- [ ] TypeScript types match Go model structs
- [ ] Components use `cn()` utility for class merging
- [ ] No edits to `frontend/bindings/` (auto-generated)

## See Also

- [04-project-structure.md](./04-project-structure.md) — Directory layout conventions
- [09-frontend-conventions.md](./09-frontend-conventions.md) — Frontend-specific conventions
