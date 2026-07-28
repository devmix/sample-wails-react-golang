---
title: ToDo Notes — Frontend Architecture
version: 2.0
date_created: 2025-07-28
tags: [frontend, react, typescript, components]
part_of: spec-app-architecture
---

# Frontend Architecture

## Technology Stack

- React 19 + TypeScript (strict mode)
- Vite build tool with HMR dev server (port 9245)
- Tailwind CSS v4 for styling
- shadcn/ui components (Radix UI primitives)
- Lucide icons (`lucide-react`)
- State management: React hooks only, no external state library

## Entry Point

`frontend/src/main.tsx` → `<App />` in `App.tsx`. No router — page switching via `useState`.

## TypeScript Interfaces

Types are re-exported from auto-generated Wails bindings:

```typescript
// frontend/src/types/models.ts
export type { Todo, Note, Tag, TodoCreateInput, TodoUpdateInput, TodoFilter, NoteCreateInput, NoteUpdateInput }
  from '@bindings/sample-wails-react-golang/internal/pkg/model/models';
```

## Hooks (`frontend/src/hooks/`)

### useTodos

Full CRUD lifecycle for todos. Returns `{ todos, loading, error, fetchTodos, createTodo, updateTodo, deleteTodo, search }`. Fetches on mount, normalizes `tags` to empty array when null.

### useNotes

Full CRUD lifecycle for notes. Returns `{ notes, loading, error, fetchNotes, createNote, updateNote, deleteNote, search }`. New notes are prepended to local state.

### useTags

Fetches all tags on mount. Returns `{ tags, loading }`. No mutation methods — tags are managed through todo/note operations.

## Component Tree (`frontend/src/components/`)

```
ui/                    # shadcn/ui base components
├── button.tsx         # CVA variants: default, destructive, outline, etc.
├── input.tsx
├── textarea.tsx
├── select.tsx         # Radix Select primitive
├── badge.tsx          # CVA variants: default, secondary, destructive
└── card.tsx           # Card, CardHeader, CardContent, CardFooter

todo/                  # Todo-specific components
├── TodoList.tsx       # Renders list of TodoItem with filters
├── TodoItem.tsx       # Single todo row with status/priority/tags
├── TodoForm.tsx       # Create/edit form
└── TodoFilters.tsx    # Status, priority, tag filter controls

note/                  # Note-specific components
├── NoteList.tsx       # Renders grid of NoteCard
├── NoteCard.tsx       # Single note card with preview
└── NoteEditor.tsx     # Full-screen note editor

common/                # Shared utility components
├── TagChip.tsx        # Colored tag badge
├── SearchBar.tsx      # Debounced search input
└── PriorityBadge.tsx  # Color-coded priority indicator

layout/
└── Sidebar.tsx        # Navigation + dark mode toggle
```

## Pages (`frontend/src/pages/`)

| Page | File | Description |
|------|------|-------------|
| TodosPage | `TodosPage.tsx` | Todo list with filters, search, create form |
| NotesPage | `NotesPage.tsx` | Note grid with editor modal |
| SettingsPage | `SettingsPage.tsx` | App settings (theme toggle, version info) |

## Theme System (`frontend/src/index.css`)

CSS custom properties mapped to Tailwind via `@theme inline`. Dark mode toggled via `.dark` class on `<html>`, persisted to `localStorage`, respects `prefers-color-scheme` on first visit.

### Core tokens

| Token | Light | Dark |
|-------|-------|------|
| `--background` | oklch(1 0 0) | oklch(0.145 0 0) |
| `--foreground` | oklch(0.145 0 0) | oklch(0.985 0 0) |
| `--primary` | oklch(0.205 0 0) | oklch(0.985 0 0) |
| `--card` | oklch(1 0 0) | oklch(0.145 0 0) |
| `--muted` | oklch(0.97 0 0) | oklch(0.269 0 0) |

### Priority colors (both themes)

| Token | Value |
|-------|-------|
| `--priority-low` | oklch(0.723 0.11 149.64) — green |
| `--priority-medium` | oklch(0.705 0.185 56.06) — amber |
| `--priority-high` | oklch(0.646 0.222 31.97) — orange |
| `--priority-urgent` | oklch(0.577 0.245 27.325) — red |

### Tag colors

Light: bg oklch(0.93 0.03 264.5), fg oklch(0.55 0.15 264.5). Dark: bg oklch(0.3 0.08 264.5), fg oklch(0.8 0.15 264.5).

## Utilities (`frontend/src/utils/`)

- `debounce.ts` — Debounce utility for search input throttling

## See Also

- [03-frontend.md](./03-frontend.md) — This file
- [09-frontend-conventions.md](./09-frontend-conventions.md) — Tailwind, shadcn/ui setup, path aliases
- [02-go-services.md](./02-go-services.md) — Go services called by hooks
