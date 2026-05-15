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

### 2. Configure environment (optinal), only to run on local machine

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
cd frontend
npm run start
```

## Make Targets

```bash
make build   # Compile binary
make run     # Build and run
make test    # Run tests
make clean   # Remove build artifacts
```
