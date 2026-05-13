# Insta-Like

A photo upload and sharing service with real-time WebSocket notifications.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose
- PostgreSQL client (for running migrations)

## Getting Started

### 1. Start the database

```bash
cd backend
docker compose up -d
```

### 2. Configure environment

Create a `.env` file in the `backend/` directory (see `.env.example`):

```env
DATABASE_URL=postgres://postgresql:postgresql@localhost:5732/insta?sslmode=disable
UPLOAD_PATH=./uploads
```

### 3. Run the migration

```bash
psql "$DATABASE_URL" -f migrations/0001-up-init-photo.sql
```

### 4. Run the server

```bash
cd backend
make run
```

The server starts at `http://127.0.0.1:8000`.

## API Endpoints

| Method | Path      | Description                          |
|--------|-----------|--------------------------------------|
| GET    | /health   | Health check                         |
| POST   | /upload   | Upload a photo with tags             |
| GET    | /photos   | List photos (with tag filter & cursor pagination) |
| GET    | /ws       | WebSocket for real-time upload notifications |

## Usage Examples

**Upload a photo:**

```bash
curl -X POST http://localhost:8000/upload \
  -F "photo=@myimage.jpg" \
  -F "tags=sunset, beach"
```

**List photos with tag filter and pagination:**

```bash
curl "http://localhost:8000/photos?tags=sunset&limit=10"
curl "http://localhost:8000/photos?cursor=2026-05-13T14:30:00Z&limit=10"
```

**Connect to WebSocket:**

```bash
npx wscat -c ws://127.0.0.1:8000/ws
```

## Project Structure

```
backend/
├── cmd/main.go              # Entry point
├── internal/
│   ├── api/                 # HTTP handlers & middleware
│   └── websocket/           # WebSocket hub
├── pkg/
│   ├── database/            # Repository layer
│   ├── models/              # Data models
│   └── storage/             # File storage interface
├── migrations/              # SQL migrations
├── docs/openapi.yaml        # API documentation
├── Makefile
└── docker-compose.yml
```

## Make Targets

```bash
make build   # Compile binary
make run     # Build and run
make test    # Run tests
make fmt     # Format code
make vet     # Static analysis
make clean   # Remove build artifacts
```
