# AGENTS.md — sample-wails-react-golang

## Quick Start

```bash
go mod tidy
cd frontend && npm install && cd ..
wails3 dev          # hot-reload dev server (port 9245)
```

Prerequisites: Go 1.25+, Node.js 20+, `wails3` CLI (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`).

## Architecture

Three-layer desktop app (Wails v3 alpha + React 19 + SQLite):

```
frontend/src/hooks/          → calls Go services via auto-generated bindings
internal/app/service/        → business logic, registered as Wails services in main.go
internal/pkg/repository/     → raw SQL, migrations, DB connection
```

- **Entry point:** `main.go` wires DB → repos → services → Wails app.
- **Frontend entry:** `frontend/src/main.tsx` → `<App />` (no router; page switching via `useState`).
- **Navigation:** simple state-based switch in `App.tsx`, not a router library.

## Key Commands

| Action              | Command                                       |
|---------------------|-----------------------------------------------|
| Dev server          | `wails3 dev`                                  |
| Production build    | `wails3 build` (embeds `frontend/dist`)       |
| Frontend-only build | `cd frontend && npm run build`                |
| Frontend typecheck  | `cd frontend && npx tsc --noEmit`             |
| Go tests            | `go test ./...`                               |
| Task runner         | `task dev`, `task build`, etc. (Taskfile.yml) |

## Critical Gotchas

### Auto-generated bindings

The `frontend/bindings/` directory is **auto-generated** by the Wails vite plugin (`@wailsio/runtime/plugins/vite`). Never edit files there manually — they are regenerated on each dev/build cycle. Frontend hooks import from `@bindings/...`.

### Path aliases

- `@/*` → `frontend/src/*`
- `@bindings/*` → `frontend/bindings/*`
  Both configured in `tsconfig.json` and `vite.config.ts`.

### Production build embeds frontend

`main.go` uses `//go:embed all:frontend/dist`. The production binary bundles the built frontend. Always run `cd frontend && npm run build` before `wails3 build`, or let `wails3 build` handle it automatically.

### Database location

SQLite DB lives at `$XDG_DATA_HOME/sample-wails-react-golang/app.db` (resolved via `github.com/adrg/xdg`). Not in the project directory.

### Tests

- Go tests only, under `internal/pkg/repository/*_test.go`. No frontend tests exist.
- Tests use `t.TempDir()` for isolated SQLite DBs — no shared state between runs.
- Run: `go test ./...` from repo root.

### Dev mode file watching

Wails dev watches `.go` and `.js/.ts` files (via `//wails:include`). The `frontend/` directory is excluded from Go watcher since Vite handles frontend HMR separately. Config in `build/config.yml`.

### Taskfile

Root `Taskfile.yml` dispatches platform-specific tasks via `GOOS` var. Includes for linux/darwin/windows/ios/android live under `build/`. Use `task dev` as shorthand for `wails3 dev -config ./build/config.yml`.

## Frontend Conventions

- UI: Tailwind CSS v4 + shadcn/ui (Radix primitives) + Lucide icons
- State management: React hooks only, no external state library
- TypeScript: strict mode, `noUnusedLocals: true`, `noImplicitAny: false`
- `.npmrc`: enforces `minimum-release-age=10080` (7 days) on package installs
