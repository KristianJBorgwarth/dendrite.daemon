package main

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/KristianJBorgwarth/dendrite.daemon/persistence"
  "github.com/KristianJBorgwarth/dendrite.daemon/core"
	_ "modernc.org/sqlite"
)

func main() {
	server := core.NewServer()
	server.Run()
}

func runMigration() {
	dbPath := filepath.Join("persistence", "index.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrationsDir := filepath.Join("persistence", "migrations")

	if err := persistence.ApplyMigrations(db, migrationsDir); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations applied successfully")
}
