---
title: ToDo Notes — Overview & Requirements
version: 2.0
date_created: 2025-07-28
tags: [app, architecture, requirements]
part_of: spec-app-architecture
---

# Overview

Cross-platform desktop app for managing to-do items and notes. Go backend + SQLite persistence, Wails v3 wrapper, TypeScript/React frontend with shadcn/ui (Radix UI primitives) and Vite.

## Purpose & Scope

Defines architecture, data model, interfaces, and implementation guidelines. Targets Windows, macOS, Linux via Wails v3.

**Assumptions**: Wails v3 alpha stable enough for dev; SQLite adequate for single-user desktop; shadcn/ui + Radix UI API stable during development.

## Definitions

| Term | Definition |
|------|-----------|
| Wails v3 | Go framework for native desktop apps with web frontends (v3 alpha) |
| shadcn/ui | Accessible, composable React component library built on Radix UI primitives |
| Radix UI | Headless, accessible UI primitives (Select, Dialog, Popover, etc.) |
| Todo | Task item: title, description, status, priority, due date, completion state |
| Note | Free-form text note: title, content, tags, timestamps |
| Go-JS Bindings | Auto-generated type-safe bindings (Go services → TypeScript) |

## Functional Requirements

- **REQ-001**: CRUD for Todo items
- **REQ-002**: CRUD for Note items
- **REQ-003**: Todo status: `pending`, `in_progress`, `completed`
- **REQ-004**: Todo priority: `low`, `medium`, `high`, `urgent`
- **REQ-005**: Tags on Todos and Notes (many-to-many)
- **REQ-006**: Filter by status, priority, tags
- **REQ-007**: Full-text search across todos and notes
- **REQ-008**: Optional due date on todos
- **REQ-009**: Overdue visual highlight (red accent)
- **REQ-010**: All data persisted to SQLite

## Non-Functional Requirements

- **PER-001**: Startup < 1s; CRUD response < 200ms
- **PER-002**: Binary size < 25MB
- **PER-003**: Prepared statements only (no SQL injection)

## Security

- **SEC-001**: Parameterized queries only
- **SEC-002**: SQLite file: owner-only permissions
- **SEC-003**: Data path via `github.com/adrg/xdg` (`$XDG_DATA_HOME/sample-wails-react-golang/app.db`)

## Constraints

- **CON-001**: Go 1.25+
- **CON-002**: Wails v3 (`github.com/wailsapp/wails/v3`)
- **CON-003**: TypeScript strict mode
- **CON-004**: shadcn/ui + Radix UI exclusive UI library
- **CON-005**: Vite build tool
- **CON-006**: SQLite via `modernc.org/sqlite` (pure Go, no CGO)

## Patterns

- **PAT-001**: Repository pattern — data access in `internal/pkg/repository/`
- **PAT-002**: Service layer — business logic + Wails bindings in `internal/app/service/`
- **PAT-003**: Custom React hooks — API calls + local state in `frontend/src/hooks/`

## Dependencies

| ID | Dependency | Purpose |
|----|-----------|---------|
| PLT-001 | Go 1.25+ | Wails v3 compatibility |
| PLT-002 | Node.js 20+ | Vite + TypeScript tooling |
| INF-001 | `modernc.org/sqlite` | Pure Go SQLite driver |
| INF-002 | `github.com/adrg/xdg` | XDG data directory resolution |
| SVC-001 | shadcn/ui (Radix UI) | Accessible UI components |
| SVC-002 | Tailwind CSS v4 | Utility-first styling + theme variables |
| SVC-003 | Vite 8+ | Frontend build + HMR dev server |
| SVC-004 | `class-variance-authority` | Component variant system (CVA) |
| SVC-005 | `lucide-react` | Icon library |

## Rationale

- **Wails v3 over Electron**: ~15MB binary vs 150MB+, ~10MB memory vs 100MB+, <0.5s startup
- **SQLite over JSON**: ACID, indexed queries, concurrent reads — critical as data grows
- **`modernc.org/sqlite`**: Pure Go, no CGO → simpler cross-platform builds
- **shadcn/ui + Radix UI**: Accessible components out of the box (ARIA, keyboard nav), composable design, Tailwind-native styling, dark mode support via CSS variables

## See Also

- [01-data-model.md](./01-data-model.md) — SQLite schema, Go models
- [02-go-services.md](./02-go-services.md) — Service API, repository layer
- [03-frontend.md](./03-frontend.md) — Frontend architecture, components, theme
- [04-project-structure.md](./04-project-structure.md) — Directory layout
- [05-testing.md](./05-testing.md) — Test strategy
- [06-edge-cases.md](./06-edge-cases.md) — Edge cases, acceptance criteria
- [07-dev-conventions.md](./07-dev-conventions.md) — Code style, naming, git workflow
- [08-build-deployment.md](./08-build-deployment.md) — Build process, platforms, Taskfile
- [09-frontend-conventions.md](./09-frontend-conventions.md) — Tailwind, shadcn/ui, aliases
