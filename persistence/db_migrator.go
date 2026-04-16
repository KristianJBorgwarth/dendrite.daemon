package persistence

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func InitializeIndex(vaultPath string) (*sql.DB, error) {
	indexDir, err := getIndexDBPath(vaultPath)
	if err != nil {
		slog.Error("failed to get index directory path", "error", err)
		return nil, err
	}

	if err = os.MkdirAll(indexDir, 0755); err != nil {
		slog.Error("failed to create index directory", "error", err)
		return nil, err
	}

	dbPath := filepath.Join(indexDir, "index.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, err
	}


	if err := applyMigrations(db); err != nil {
		slog.Error("failed to apply migrations", "error", err)
		return nil, err
	}

	return db, nil
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
