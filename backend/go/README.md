# Meeting Cost Calculator - Go Backend

Backend API for the meeting cost calculator application. Built with [Fiber](https://gofiber.io/) and [GORM](https://gorm.io/).

## Project Structure

- `cmd/` - Application entry points (main, migrate)
- `internal/` - Private application code
  - `models/` - GORM data models
  - `repository/` - Repository interfaces and implementations
  - `service/` - Business logic layer
  - `handler/` - HTTP handlers
  - `middleware/` - HTTP middleware
  - `config/` - Configuration
  - `cache/` - Cache abstractions (Valkey/Redis)
  - `errors/` - Error definitions
  - `logger/` - Structured logging
- `migrations/` - Versioned SQL migrations
- `pkg/` - Public reusable packages (if any)

## Prerequisites

- Go 1.22+
- PostgreSQL
- Valkey or Redis (for caching)

## Local Development

1. Copy `.env.example` to `.env` and set values.
2. Run migrations: `go run ./cmd/migrate`
3. Start the API: `go run ./cmd/api`

## Testing

```bash
go test ./...
```

## Linting

```bash
golangci-lint run
```
