package migrations

import "embed"

// FS holds the migration SQL files for use with golang-migrate.
//go:embed *.sql
var FS embed.FS
