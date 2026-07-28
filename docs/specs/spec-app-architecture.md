---
title: ToDo Notes & Tasks Application — Architecture Specification Index
version: 2.0
date_created: 2025-07-28
tags: [app, architecture, index]
---

# Architecture Specification Index

This specification has been split into focused documents for easier navigation and maintenance.

## Documents

| # | Document | Contents |
|---|----------|----------|
| 00 | [Overview & Requirements](./00-overview.md) | Purpose, definitions, functional/non-functional requirements, security, constraints, patterns, dependencies, rationale |
| 01 | [Data Model](./01-data-model.md) | SQLite schema, Go model structs (Todo, Note, Tag), DB connection settings |
| 02 | [Go Services & Repository Layer](./02-go-services.md) | Service API (Wails bindings), TodoService, NoteService, TagService, AppService, repository pattern |
| 03 | [Frontend Architecture](./03-frontend.md) | React hooks, component tree, pages, theme system with CSS tokens |
| 04 | [Project Structure](./04-project-structure.md) | Full directory layout: `internal/`, `frontend/`, `build/` |
| 05 | [Test Strategy](./05-testing.md) | Go tests only (repository layer), coverage targets, future frontend/E2E considerations |
| 06 | [Edge Cases & Acceptance Criteria](./06-edge-cases.md) | AC-001..006, validation criteria V-001..006, edge cases EC-001..005 |
| 07 | [Development Conventions](./07-dev-conventions.md) | Go naming, TypeScript conventions, git workflow, code review checklist |
| 08 | [Build & Deployment](./08-build-deployment.md) | Dev server, production build, platform targets, Taskfile.yml, embed FS, server mode |
| 09 | [Frontend Conventions](./09-frontend-conventions.md) | Tailwind v4 CSS-first config, shadcn/ui setup, path aliases, auto-generated bindings, .npmrc |

## Quick Start

```bash
go mod tidy
cd frontend && npm install && cd ..
wails3 dev          # hot-reload dev server (port 9245)
```

## Key Commands

| Action | Command |
|--------|---------|
| Dev server | `wails3 dev` or `task dev` |
| Production build | `wails3 build` or `task build` |
| Frontend-only build | `cd frontend && npm run build` |
| Frontend typecheck | `cd frontend && npx tsc --noEmit` |
| Go tests | `go test ./...` |

## Architecture at a Glance

```
frontend/src/hooks/  →  @bindings/* (auto-generated)  →  internal/app/service/  →  internal/pkg/repository/  →  SQLite
```

- **Entry point**: `main.go` wires DB → repos → services → Wails app
- **Frontend entry**: `frontend/src/main.tsx` → `<App />` (state-based page switching, no router)
- **Database**: `$XDG_DATA_HOME/sample-wails-react-golang/app.db` (via `github.com/adrg/xdg`)

## References

[Wails v3](https://v3.wails.io/) · [shadcn/ui](https://ui.shadcn.com/) · [Radix UI](https://www.radix-ui.com/) · [Tailwind CSS v4](https://tailwindcss.com/docs/v4) · [SQLite Docs](https://www.sqlite.org/docs.html)
