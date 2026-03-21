package main

import (
	"database/sql"
	"log"
	"path/filepath"

	migrate "github.com/KristianJBorgwarth/dendrite.daemon/persistence"
	_ "modernc.org/sqlite"
)

func main() {

}

func runMigration() {
	dbPath := filepath.Join("persistence", "index.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrationsDir := filepath.Join("persistence", "migrations")

	if err := migrate.Run(db, migrationsDir); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations applied successfully")
}
