// Package persistence handles database schema migrations
package persistence

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
)

func InitializeIndex(vaulPath string) error {
	dbPath := filepath.Join(vaulPath, "index.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	migrationsDir := filepath.Join("persistence", "migrations")

	if err := applyMigrations(db, migrationsDir); err != nil {
		return err
	}
	return nil
}

func applyMigrations(db *sql.DB, dir string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, e := range entries {
		if e.IsDir(){
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

		path := filepath.Join(dir, name)

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return err
		}

		if _, err := db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, name); err != nil {
			return err
		}
	}

	return nil
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var v string
	err := db.QueryRow(`SELECT version FROM schema_migrations WHERE version= ?`, version).Scan(&v)

	if err == sql.ErrNoRows {
		return false, nil
	}

	return err == nil, err
}
