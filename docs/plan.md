# laga - MVP Plan

## Project Goal
Build a mobile-first, offline-first shopping list application. The primary objective of the MVP is to establish the full end-to-end stack: a Svelte frontend, a Go backend, a PostgreSQL database, and the synchronization logic connecting them. The app must be fully usable without any internet connection after the first load.

## Decisions

1.  **Go Router:** Standard library `net/http` to minimize dependencies.
2.  **Go Module:** `github.com/herrnan/laga`.
3.  **Sync Strategy:** Client-overwrites-server (Last-Write-Wins) for the MVP. Multi-device conflict resolution is deferred to v2.
4.  **Backend Hosts Frontend:** The Go backend serves the compiled frontend (`frontend/dist/`) using `//go:embed` in production. During development, it reverse-proxies to the Vite dev server when `VITE_DEV_PROXY` is set. No CORS is needed because frontend and API share the same origin.
5.  **Styling:** Plain CSS for the MVP to keep dependencies minimal.
6.  **Database Access:** `sqlc` generates type-safe Go code from annotated SQL queries. We write SQL; sqlc writes the Go.
7.  **Migrations:** `goose` manages schema changes as versioned, single-file SQL migrations.
8.  **Schema Source:** `backend/schema.sql` is a generated file (not hand-written, not in git). A Go script (`scripts/generate-schema.go`) spins up a temporary PostgreSQL container, applies all migrations via goose, and dumps the resulting schema via `pg_dump --schema-only`. sqlc then reads `schema.sql` for type inference.
9.  **Task Runner:** `mise` manages all dev tools (Go, Node, sqlc, goose, pnpm) and defines project tasks inline in `mise.toml` (no dot prefix).
10. **Node Package Manager:** `pnpm`.
11. **Local PostgreSQL for Dev:** Runs as a Docker container managed by mise tasks.

## MVP Scope

### Features
-   **Add Item:** Free-text input (e.g., "2 kg potatoes").
-   **List Items:** Display items with a checkbox.
-   **Check / Uncheck:** Toggle item completion state.
-   **Offline Usage:** App loads and functions without any network access after the first visit.
-   **Sync:** Push local state to the backend and pull the server state when online.

### Defer to v2
-   Delete all checked items.
-   Recipes and cooking mode.
-   Item categorization and custom group sorting.
-   Advanced conflict resolution (CRDTs, OT, vector clocks).
-   Renaming items or changing amounts after creation.

## Architecture & Stack

| Layer | Tech | Rationale |
|-------|------|-----------|
| **Frontend** | Svelte + Vite | Lightweight, reactive, fast build. |
| **Offline Load** | Vite PWA Plugin (`vite-plugin-pwa`) | Caches the application shell, allowing the SPA to load without network requests. |
| **Local State** | IndexedDB (`idb` npm package) | Async, structured storage larger than `localStorage`, survives reloads. |
| **Backend** | Go 1.22+ (`net/http`) | Stable, zero-dependency routing. |
| **Database** | PostgreSQL | Robust, standard SQL, handles concurrency and future multi-user features better than SQLite. |
| **DB Access** | `sqlc` (engine: postgresql, driver: database/sql) | Type-safe queries generated from SQL. Avoids ORMs while eliminating hand-written boilerplate. |
| **Migrations** | `goose` (single-file up/down) | Simple, versioned schema changes. |
| **Dev Tools** | `mise` (`mise.jdx.dev`) | Manages `go`, `node`, `sqlc`, `goose`, `pnpm` versions and project tasks. |
| **Dev PostgreSQL** | Docker container | Isolated, reproducible, easy to reset. |

## Data Model

### Database Schema (`backend/schema.sql` — generated)
```sql
CREATE TABLE IF NOT EXISTS items (
    id TEXT PRIMARY KEY,                    -- UUID v4, generated client-side
    text TEXT NOT NULL,
    checked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Rationale for client-side UUIDs:** Because the `POST /api/sync` endpoint replaces the entire server table with the client's state, using stable client-generated IDs prevents the database from accumulating duplicate rows across multiple syncs.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/health` | Liveness probe. |
| `GET` | `/api/items` | Fetch current server state. |
| `POST` | `/api/items` | Add a single item (utility/testing endpoint). |
| `PUT` | `/api/items/{id}` | Update a single item (check/uncheck). |
| `DELETE` | `/api/items/{id}` | Delete a single item. |
| `POST` | `/api/sync` | **Primary sync flow.** Client sends full local array; server atomically replaces its state and returns the canonical list. |

## Project Structure

```
laga/
├── mise.toml                  -- Dev tool versions + project tasks
├── scripts/
│   └── generate-schema.go     -- Spins up temp postgres, applies migrations, dumps schema.sql
├── backend/
│   ├── migrations/
│   │   └── 00001_init_schema.sql
│   ├── queries/
│   │   └── items.sql          -- Annotated SQL for sqlc (e.g., -- name: ListItems :many)
│   ├── internal/
│   │   └── gen/
│   │       └── db/            -- Generated by sqlc (DO NOT EDIT)
│   │           ├── models.go
│   │           └── querier.go
│   ├── sqlc.yaml              -- sqlc config (engine: postgresql, emit_db_tags: true, output: internal/gen/db)
│   ├── schema.sql             -- GENERATED by scripts/generate-schema.go (not in git)
│   ├── dist/                  -- GENERATED by `pnpm run build` in frontend (embedded by backend via go:embed)
│   ├── main.go                -- Server boot, routes, static file serving, dev proxy
│   ├── handlers.go            -- REST handlers (use *db.Queries from sqlc)
│   └── ...
├── frontend/
│   ├── src/
│   │   ├── App.svelte
│   │   ├── lib/
│   │   │   ├── AddItem.svelte
│   │   │   └── ItemList.svelte
│   │   ├── stores.js          -- Svelte writable store for the items array
│   │   ├── db.js              -- IndexedDB get/set logic
│   │   └── api.js             -- Fetch wrappers + sync() function
│   ├── index.html
│   ├── vite.config.js         -- Vite + PWA plugin configuration
│   └── manifest.json          -- PWA manifest (name: laga, display: standalone)
├── docs/
│   ├── plan.md                -- This document
│   └── sync-research.md       -- Research note for v2 conflict resolution
└── README.md                  -- Updated project description
```

## mise Task Definitions (mise.toml)

```toml
[tools]
go = "latest"
node = "latest"
sqlc = "latest"
goose = "latest"
pnpm = "latest"

[env]
DATABASE_URL = "postgres://laga:laga@localhost:5432/laga?sslmode=disable"
VITE_DEV_PROXY = "http://localhost:5173"

[tasks.db-up]
description = "Start local PostgreSQL container for development"
run = "docker run -d --name laga-db -p 5432:5432 -e POSTGRES_USER=laga -e POSTGRES_PASSWORD=laga -e POSTGRES_DB=laga postgres:16-alpine || docker start laga-db"

[tasks.db-down]
description = "Stop local PostgreSQL container"
run = "docker stop laga-db"

[tasks.db-reset]
description = "Remove and recreate local PostgreSQL container"
depends = ["db-down"]
run = "docker rm laga-db && mise run db-up"

[tasks.generate-schema]
description = "Apply all migrations to temp postgres container and dump schema.sql"
run = "go run scripts/generate-schema.go"

[tasks.sqlc-generate]
description = "Generate type-safe Go DB code from schema.sql and queries/"
depends = ["generate-schema"]
run = "cd backend && sqlc generate"

[tasks.migrate]
description = "Run database migrations via goose"
run = "cd backend && goose postgres $DATABASE_URL up"

[tasks.migrate-status]
description = "Check migration status"
run = "cd backend && goose postgres $DATABASE_URL status"

[tasks.dev-fe]
description = "Run Vite dev server (backend proxies to this)"
run = "cd frontend && pnpm dev"

[tasks.dev-be]
description = "Run backend with Vite proxy (VITE_DEV_PROXY must be set)"
run = "cd backend && go run ."

[tasks.dev]
description = "Run backend, frontend, and database in parallel"
depends = ["db-up", "dev-fe", "dev-be"]
```

## Implementation Roadmap

### Phase 0: Scaffold & Tooling
1.  Write `mise.toml` pinning Go, Node, sqlc, goose, and pnpm versions. Define all project tasks.
2.  Initialize `backend/` as a Go module: `go mod init github.com/herrnan/laga`.
3.  Initialize `frontend/` with the official Vite Svelte template using `pnpm`.
4.  Install `vite-plugin-pwa` via `pnpm` and configure `vite.config.js` to generate a Service Worker that precaches the app shell.
5.  Write `scripts/generate-schema.go`: a Go script that runs a temporary PostgreSQL container via Docker, applies all goose migrations, then executes `pg_dump --schema-only` to write `backend/schema.sql`.
6.  Write the initial migration `backend/migrations/00001_init_schema.sql` containing the `items` table DDL.
7.  Run `mise run db-up` to start the dev PostgreSQL container.
8.  Run `mise run generate-schema` and verify `backend/schema.sql` is created.

### Phase 1: Backend Foundation (with sqlc + goose)
1.  Write `backend/sqlc.yaml` targeting `postgresql` and `database/sql`, outputting to `internal/gen/db/`.
2.  Write `backend/queries/items.sql` with annotated queries for ListItems, GetItem, CreateItem, UpdateItem, DeleteItem, DeleteAllItems, and SyncItems.
3.  Run `mise run sqlc-generate` to produce `internal/gen/db/models.go` and `querier.go`.
4.  Implement `db.go`: open a `*sql.DB` connection using `pgx/stdlib`, reading `DATABASE_URL` from env.
5.  Implement `handlers.go`:
    -   `GET /api/health`
    -   `GET /api/items` → `q.ListItems(ctx)`
    -   `POST /api/items` → `q.CreateItem(ctx, ...)`
    -   `PUT /api/items/{id}` → `q.UpdateItem(ctx, ...)`
    -   `DELETE /api/items/{id}` → `q.DeleteItem(ctx, ...)`
    -   `POST /api/sync` → Begin `sql.Tx`, call `q.WithTx(tx).DeleteAllItems(ctx)` then bulk insert via `q.WithTx(tx).SyncItems(ctx, ...)` (or equivalent), commit, return `q.WithTx(tx).ListItems(ctx)`.
6.  Implement `main.go`:
    -   Register API routes under `/api/*`.
    -   If `VITE_DEV_PROXY` is set, reverse-proxy all other requests to the Vite dev server using `httputil.ReverseProxy`.
    -   If not set, serve embedded static files from `frontend/dist/` via `//go:embed` and `http.FileServer`.
7.  Run `mise run migrate` to apply the initial migration to the dev PostgreSQL container.

### Phase 2: Frontend UI & State
1.  Build `AddItem.svelte`: an input field and a button that parses free text and appends a new item (with client-side UUID) to the store.
2.  Build `ItemList.svelte`: renders items from the store with checkboxes. Each toggle updates the store directly.
3.  Create `stores.js`: a Svelte writable store holding the array of items. This store is the single source of truth for the UI.

### Phase 3: Offline Persistence
1.  Create `db.js`: wrap IndexedDB (using the `idb` npm package for ergonomics) to read/write the items array under a single object store.
2.  On app mount (`App.svelte`), read from IndexedDB into the store.
3.  Subscribe to the store: any change writes the entire array to IndexedDB.
4.  **Verification:** Stop the backend, reload the browser, add and check items, reload again. The state must survive.

### Phase 4: The Glue (Sync)
1.  Create `api.js` with a `sync()` function that:
    -   `POST`s the current store contents to `/api/sync`.
    -   On success, updates the store with the server response.
    -   On network/server error, leaves local state untouched and logs a warning.
2.  Add a "Sync now" button to `App.svelte`.
3.  Enable/disable the sync button based on `navigator.onLine`.
4.  (Optional) Listen for the browser `online` event to trigger sync automatically.

### Phase 5: Polish & Notes
1.  Write `docs/sync-research.md` documenting the Last-Write-Wins strategy and suggesting research into CRDTs, Operational Transformation (OT), or vector clocks for future multi-device support.
2.  Ensure `manifest.json` has `display: standalone` and appropriate icons for a mobile-like experience.
3.  Implement plain CSS with a mobile-first mindset: large tap targets (`min-height: 48px`), simple flexbox layout, readable font sizes.
4.  Final verification: Use browser DevTools in mobile view, disable network, reload, add items, re-enable network, sync, and verify the server database reflects the local state.

## Research Notes (for v2)

-   **Sync Conflict Resolution:** The current MVP uses a naive "client overwrites server" strategy. Before implementing multi-user or multi-device support, research how other applications handle offline-to-online synchronization. Potential topics include:
    -   Conflict-free Replicated Data Types (CRDTs).
    -   Operational Transformation (OT).
    -   Per-item vector clocks or logical timestamps.
    -   Event sourcing with replay.
-   **Delete Checked:** Implement a "clear completed" action that removes all checked items in one go.
