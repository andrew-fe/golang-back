package store

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

const migrationsTable = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version    TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)
`

// Migrate aplica, en orden y de forma idempotente, todas las migraciones SQL
// embebidas que aún no se hayan registrado en la tabla schema_migrations.
func Migrate(db *sql.DB) error {
	if _, err := db.Exec(migrationsTable); err != nil {
		return fmt.Errorf("store: creando tabla de migraciones: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("store: leyendo migraciones embebidas: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		applied, err := isMigrationApplied(db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("store: leyendo migración %s: %w", name, err)
		}

		if err := applyMigration(db, name, string(content)); err != nil {
			return fmt.Errorf("store: aplicando migración %s: %w", name, err)
		}
	}

	return nil
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("store: verificando migración %s: %w", version, err)
	}
	return exists, nil
}

func applyMigration(db *sql.DB, version, sqlContent string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(sqlContent); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return err
	}

	return tx.Commit()
}
