# Local Docker Development

This directory contains the Docker Compose configuration for running the Go backend, PostgreSQL, and Valkey for local development.

## Prerequisites

- Docker
- Docker Compose (v2)

## Configuration

Copy `.env.example` to `.env` and adjust values if needed:

```bash
cp infrastructure/docker/.env.example infrastructure/docker/.env
```

Key variables:

- `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`: PostgreSQL connection
- `CACHE_PORT`: Valkey port
- `PORT`: Backend HTTP port

## Running the Stack

From the repo root:

```bash
docker compose -f infrastructure/docker/docker-compose.yml up --build
```

This starts:

- `postgres` on `localhost:${DB_PORT}` (default 5432)
- `valkey` on `localhost:${CACHE_PORT}` (default 6379)
- `backend` API on `http://localhost:${PORT}` (default 8080)

Health endpoints:

- `GET /health`
- `GET /api/v1/health`

## Running Database Migrations

To apply migrations inside the backend container:

```bash
docker compose -f infrastructure/docker/docker-compose.yml run --rm backend /bin/migrate up
```

To roll back one step:

```bash
docker compose -f infrastructure/docker/docker-compose.yml run --rm backend /bin/migrate down
```

## Stopping and Cleaning Up

```bash
# Stop containers
docker compose -f infrastructure/docker/docker-compose.yml down

# Remove volumes (including database data)
docker compose -f infrastructure/docker/docker-compose.yml down -v
```

# Docker – Local Development

Run the Go backend, PostgreSQL, and Valkey in containers. No local Go or database install required.

## Prerequisites

- Docker and Docker Compose
- No need for Go, PostgreSQL, or Valkey installed on the host

## Quick start

From the **repository root**:

```bash
# Optional: copy and edit env
cp infrastructure/docker/.env.example infrastructure/docker/.env

# Build and start all services
docker compose -f infrastructure/docker/docker-compose.yml up --build -d

# API: http://localhost:8080
# Health: http://localhost:8080/health
```

## Run migrations

After Postgres is up, run migrations once (or after schema changes):

```bash
docker compose -f infrastructure/docker/docker-compose.yml run --rm backend /bin/migrate up
```

Rollback one version:

```bash
docker compose -f infrastructure/docker/docker-compose.yml run --rm backend /bin/migrate down
```

## Backend image only

Build and run the Go app image only (expects DB and cache elsewhere):

```bash
cd backend/go
docker build -t meeting-cost-backend .
docker run --rm -e DB_HOST=host.docker.internal -e DB_PASSWORD=postgres -p 8080:8080 meeting-cost-backend
```

To run migrations with the same image:

```bash
docker run --rm -e DB_HOST=host.docker.internal -e DB_PASSWORD=postgres --entrypoint /bin/migrate meeting-cost-backend up
```

## Stop and clean

```bash
docker compose -f infrastructure/docker/docker-compose.yml down
# With volume removal:
docker compose -f infrastructure/docker/docker-compose.yml down -v
```
