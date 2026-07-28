# ToDo Notes — Wails v3 + React Demo App

A desktop application built with **Wails v3**, **React 19** and **SQLite**. Demonstrates a full-stack architecture: Go backend services exposed to a TypeScript frontend via Wails bindings, backed by a local SQLite database.

This project is a reference example for a technical blog post about building cross-platform desktop apps with Wails v3.

> [!NOTE]
> https://www.tekblueprint.org/blog/dev/guide-wails-react-golang/.

## Features

- **Todos** — create, edit, delete tasks with status tracking (`pending`, `in_progress`, `completed`), priority levels (`low`, `medium`, `high`, `urgent`) and optional due dates
- **Notes** — plain text notes with title, content and tags
- **Tags** — categorize todos and notes with shared tag system
- **Search & Filter** — full-text search across todos/notes; filter by status, priority or tag
- **Local Storage** — all data persisted in a SQLite database under the XDG data directory

## Tech Stack

| Layer     | Technology                                         |
|-----------|----------------------------------------------------|
| Framework | [Wails v3](https://v3.wails.io/) (alpha)           |
| Backend   | Go 1.25, `modernc.org/sqlite`                      |
| Frontend  | React 19, TypeScript 5, Vite 8                     |
| UI        | Tailwind CSS v4, shadcn/ui, Radix UI, Lucide icons |
| Build     | Wails CLI (`wails3 dev`, `wails3 build`)           |

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Frontend (React)                   │
│  hooks/useTodos.ts ──→ Wails bindings ──→ Go API    │
│  hooks/useNotes.ts                                  │
│  hooks/useTags.ts                                   │
├─────────────────────────────────────────────────────┤
│              Service Layer (internal/app/service)    │
│  TodoService · NoteService · TagService · AppService│
├─────────────────────────────────────────────────────┤
│           Repository Layer (internal/pkg/repository) │
│     TodoRepo · NoteRepo · TagRepo · SQLite DB       │
└─────────────────────────────────────────────────────┘
```

Three-layer architecture:

- **Repository** — raw SQL queries, transactions, schema migrations
- **Service** — business logic, validation, exposed to frontend via Wails bindings
- **Frontend** — React hooks call Go services directly through auto-generated TypeScript bindings. Navigation is a simple `useState` + `switch`, no router library.

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) 20+ with npm
- [Wails CLI](https://v3.wails.io/docs/introduction) (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)

## Getting Started

```bash
# Clone and enter the project
git clone <repo-url>
cd sample-wails-react-golang

# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd frontend && npm install && cd ..

# Run in development mode with hot-reload
wails3 dev
```

## Build & Package

```bash
# Production build for current OS
wails3 build

# Cross-compile (Linux, macOS, Windows)
GOOS=linux   wails3 build
GOOS=darwin  wails3 build
GOOS=windows wails3 build
```

## Project Structure

```
├── main.go                          # Entry point — DB init, service wiring, window config
├── internal/
│   ├── app/service/                 # Business logic layer (exposed to frontend)
│   │   ├── todo_service.go
│   │   ├── note_service.go
│   │   ├── tag_service.go
│   │   └── app_service.go
│   └── pkg/
│       ├── model/                   # Shared data structures
│       │   ├── todo.go
│       │   ├── note.go
│       │   └── tag.go
│       └── repository/              # Data access layer (SQLite)
│           ├── db.go                # Connection, WAL mode, migrations
│           ├── todo_repo.go
│           ├── note_repo.go
│           └── tag_repo.go
├── frontend/
│   ├── src/
│   │   ├── main.tsx                   # React entry point — renders <App />
│   │   ├── App.tsx                    # Root component — sidebar + page switcher (no router)
│   │   ├── lib/utils.ts               # cn() className merge, formatDate helper
│   │   ├── components/
│   │   │   ├── ui/                  # shadcn/ui primitives (Button, Card, Badge...)
│   │   │   ├── layout/              # Sidebar navigation
│   │   │   ├── common/              # SearchBar, PriorityBadge, TagChip
│   │   │   ├── todo/                # TodoForm, TodoList, TodoItem, TodoFilters
│   │   │   └── note/                # NoteEditor, NoteList, NoteCard
│   │   ├── hooks/                   # useTodos, useNotes, useTags
│   │   ├── pages/                   # TodosPage, NotesPage, SettingsPage
│   │   ├── types/models.ts          # Auto-generated Go→TS type bindings
│   │   └── utils/debounce.ts        # Debounce utility for search input
│   └── package.json
└── build/                           # Wails config (config.yml) and platform assets
```

## Database Schema

SQLite database with 5 tables and foreign key constraints:

| Table       | Purpose                                    |
|-------------|--------------------------------------------|
| `todos`     | Tasks with status, priority, due date      |
| `notes`     | Notes with title and text content          |
| `tags`      | Unique tag names (auto-upserted on create) |
| `todo_tags` | Many-to-many: todos ↔ tags                 |
| `note_tags` | Many-to-many: notes ↔ tags                 |

Database file is stored at `$XDG_DATA_HOME/sample-wails-react-golang/app.db`.

## License

MIT
