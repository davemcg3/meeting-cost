# Database Migrations

Versioned SQL migrations for the meeting cost calculator. Use [golang-migrate](https://github.com/golang-migrate/migrate) for production.

## Running Migrations

### Using golang-migrate CLI

```bash
# Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Up (apply all)
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/meetingcost?sslmode=disable" up

# Down (rollback one)
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/meetingcost?sslmode=disable" down 1
```

### Using the migrate command (cmd/migrate)

When `cmd/migrate` is implemented, it will run migrations on startup using the same SQL files.

## Migration Naming

- `NNN_description.up.sql` - Apply migration
- `NNN_description.down.sql` - Rollback migration

## Dependencies

Migrations must be applied in order. The initial schema (001) creates all tables; later migrations add or change objects.
