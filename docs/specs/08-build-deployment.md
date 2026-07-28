---
title: ToDo Notes — Build & Deployment
version: 2.0
date_created: 2025-07-28
tags: [build, wails, taskfile, deployment]
part_of: spec-app-architecture
---

# Build & Deployment

## Development Server

```bash
wails3 dev                          # Full dev server with hot-reload (port 9245)
task dev                            # Same via Taskfile.yml
cd frontend && npm run dev          # Frontend-only Vite dev server
```

Dev mode configuration in `build/config.yml`:

- Watches `.go` and `.js/.ts` files (via `//wails:include` directive)
- Excludes `frontend/`, `.git`, `node_modules`, `bin` from Go watcher
- Frontend HMR handled separately by Vite
- Debounce: 1000ms on file changes

## Production Build

```bash
wails3 build                        # Full production build (embeds frontend/dist)
task build                          # Same via Taskfile.yml, dispatches to platform task
cd frontend && npm run build        # Frontend-only build → frontend/dist/
```

### Embed mechanism

`main.go` uses `//go:embed all:frontend/dist` to bundle the built frontend into the Go binary. The production binary is self-contained — no separate frontend files needed at runtime.

## Platform Targets

Taskfile.yml dispatches to platform-specific Taskfiles via `GOOS` variable:

| Platform | Command | Output |
|----------|---------|--------|
| Linux | `task linux:build` or `GOOS=linux task build` | `bin/` directory |
| macOS | `task darwin:build` or `GOOS=darwin task build` | `.app` bundle |
| Windows | `task windows:build` or `GOOS=windows task build` | `.exe` binary |
| iOS | `task ios:build` (commented out by default) | Xcode project |
| Android | `task android:build` (commented out by default) | APK/AAB |

Cross-compilation via Docker:
```bash
task setup:docker                   # Build Docker image (~800MB download)
GOOS=linux task build               # Cross-compile for Linux
```

## Server Mode

Wails v3 supports a headless server mode (no GUI, HTTP only):

```bash
task build:server                   # Build server-mode binary
task run:server                     # Run in server mode
task build:docker                   # Docker image for server deployment
task run:docker                     # Build and run Docker container
```

## Frontend Build Pipeline

1. `npm run build` → Vite bundles React app to `frontend/dist/`
2. `wails3 build` detects `frontend/dist/`, embeds via Go `embed.FS`
3. Final binary contains both Go runtime and frontend assets

### Vite configuration

- Wails vite plugin (`@wailsio/runtime/plugins/vite`) generates bindings to `frontend/bindings/`
- Bindings regenerated on each dev/build cycle — never edit manually
- Dev server port: 9245 (configurable via `WAILS_VITE_PORT` env var)

## Build Configuration (`build/config.yml`)

```yaml
version: '3'
info:
  companyName: "My Company"
  productName: "My Product"
  productIdentifier: "com.mycompany.myproduct"
  version: "0.0.1"

dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir: [.git, node_modules, frontend, bin]
    file: [.DS_Store, .gitignore, .gitkeep, "*_test.go"]
    watched_extension: ["*.go", "*.js", "*.ts"]
```

To update build assets after changing `info` or `fileAssociations`:
```bash
wails3 task common:update:build-assets
```

## Quick Reference

| Action | Command |
|--------|---------|
| Dev server | `wails3 dev` or `task dev` |
| Production build | `wails3 build` or `task build` |
| Frontend-only build | `cd frontend && npm run build` |
| Frontend typecheck | `cd frontend && npx tsc --noEmit` |
| Go tests | `go test ./...` |
| Cross-compile Linux | `GOOS=linux task build` (requires Docker) |

## See Also

- [04-project-structure.md](./04-project-structure.md) — Project directory layout
- [07-dev-conventions.md](./07-dev-conventions.md) — Development conventions
