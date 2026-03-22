package persistence

import (
	"database/sql"
	"embed"
	"log/slog"
	"path/filepath"
	"sort"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func InitializeIndex(vaultPath string) error {
	dbPath := filepath.Join(vaultPath, "index.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return err
	}

	defer db.Close()

	if err := applyMigrations(db); err != nil {
		slog.Error("failed to apply migrations", "error", err)
		return err
	}

	return nil
}

func applyMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY
		);
	`)
	if err != nil {
		return err
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		applied, err := isMigrationApplied(db, name)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		sqlBytes, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return err
		}

		if _, err := db.Exec(
			`INSERT INTO schema_migrations (version) VALUES (?)`,
			name,
		); err != nil {
			return err
		}

		slog.Info("applied migration", "file", name)
	}

	return nil
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var v string

	err := db.QueryRow(
		`SELECT version FROM schema_migrations WHERE version = ?`,
		version,
	).Scan(&v)

	if err == sql.ErrNoRows {
		return false, nil
	}

	return err == nil, err
}
