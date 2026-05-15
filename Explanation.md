# Backend Architecture & Design Decisions

## High-Level Overview

The backend is a Go HTTP server that provides a photo upload and sharing service with real-time updates via WebSocket. It follows a layered architecture with clear separation between transport (HTTP/WS), business logic, and infrastructure (database, file storage).

```
┌─────────────────────────────────────────────────┐
│                   Clients                        │
│          (React SPA / WebSocket)                 │
└────────────┬──────────────────────┬─────────────┘
             │ HTTP (REST)          │ WS
┌────────────▼──────────────────────▼─────────────┐
│              Transport Layer                      │
│   gorilla/mux router + middleware (CORS, logs)   │
└────────────┬──────────────────────┬─────────────┘
             │                      │
┌────────────▼──────────┐  ┌───────▼──────────────┐
│    API Handlers        │  │   WebSocket Hub       │
│  (upload, list, health)│  │  (broadcast to all)   │
└────────────┬──────────┘  └───────────────────────┘
             │
     ┌───────┴────────┐
     │                 │
┌────▼─────┐   ┌──────▼──────┐
│  Storage  │   │  Database   │
│ (local fs)│   │ (PostgreSQL)│
└───────────┘   └─────────────┘
```

## Project Layout

```
backend/
├── cmd/main.go                    # Entry point, wiring, graceful shutdown
├── internal/
│   ├── api/                       # HTTP handlers and middleware (not importable outside module)
│   └── websocket/                 # WebSocket hub for real-time broadcast
├── pkg/
│   ├── database/                  # Repository interface + Postgres implementation
│   ├── models/                    # Domain models (Photo)
│   ├── storage/                   # Storage interface + local filesystem implementation
│   └── image_normalization/       # Optional image processing (grayscale)
├── migrations/                    # SQL migration files
├── docs/openapi.yaml              # API contract
└── docker-compose.yml             # Local dev environment (Postgres + backend)
```

**Decision**: Use Go's standard `internal/` vs `pkg/` convention.
- `internal/` — code that is private to this module (handlers, websocket hub). Cannot be imported by other Go modules.
- `pkg/` — code that is conceptually reusable and could be imported externally (interfaces, models, storage).

## Key Design Decisions

### 1. Interface-Based Dependency Injection

All infrastructure dependencies are defined as interfaces and injected via constructor functions:

```go
type StorageService interface {
    UploadFoto(fileName string, fileBytes []byte) (string, error)
    GetFoto(filePath string) ([]byte, error)
}

type PhotoRepository interface {
    InsertPhoto(ctx context.Context, photo models.Photo) error
    ListPhotos(ctx context.Context, params ListPhotosParams) ([]models.Photo, error)
}
```

**Why**: Decouples handlers from concrete implementations. Makes it straightforward to swap local storage for S3, or Postgres for another database, without touching handler code. Also enables unit testing with mocks.

### 2. Cursor-Based Pagination

The `/photos` endpoint uses `created_at` timestamp as a cursor instead of offset-based pagination.

**Why**:
- Stable results when new photos are uploaded between page fetches (no skipped/duplicated items).
- Performs well with an index on `created_at DESC` — the DB seeks directly to the cursor position.
- Offset pagination degrades on large datasets because the DB must skip N rows.

### 3. Real-Time Updates via WebSocket Hub

A simple hub pattern manages connected clients:
- New connections are registered in a map.
- On photo upload, the hub broadcasts the new photo JSON to all connected clients.
- A `readPump` per connection handles ping/pong keepalive and detects disconnects.

**Why**: Provides instant feed updates without polling. The hub pattern is lightweight and sufficient for a single-instance deployment. No external message broker needed at this scale.

### 4. Optional Image Normalization

Image normalization (grayscale conversion) is toggled by the `NORMALIZATION` environment variable. When enabled, uploaded images are converted to grayscale PNG before storage.

**Why**: Feature-flagged via env var so it can be enabled/disabled per environment without code changes. Keeps the processing pipeline optional and non-breaking.

### 5. Graceful Shutdown

The server listens for `os.Interrupt` and performs a graceful shutdown with a configurable timeout (default 15s). In-flight requests are allowed to complete before the process exits.

**Why**: Prevents dropped connections during deployments or restarts. Standard practice for production Go services.

### 6. UUID-Based File Naming

Uploaded files are renamed to `<uuid>.<original-extension>` before storage.

**Why**: Eliminates filename collisions and prevents path traversal attacks from crafted filenames. The original filename is never used as a storage path.

## What's Not Included (and Why)

- **Rate Limiting** — Acceptable for a single-user or small-group deployment.
- **Structured Logging** — Current `log.Println` is sufficient for development. A structured logger (zerolog, slog) would be appropriate for production.
- **Database Migrations Runner** — SQL files exist but are applied manually. A tool like golang-migrate could automate this.
- **Tests for Handlers** — The interface-based design supports it; test coverage can be added incrementally.
