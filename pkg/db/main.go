package db

import (
	"fmt"
	"log"
	"os"

	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	migrationDir string = "./migrations"
)

var (
	Database_url string = os.Getenv("DATABASE_URL")
)

func RunDatabaseMigration() {
	fmt.Println(Database_url)
	db, err := sql.Open("postgres", Database_url)
	defer db.Close()

	if err != nil {
		log.Printf("could not connect to the Postgresql database... %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("could not ping DB... %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"postgres", driver)

	if err != nil {
		log.Printf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("An error occurred while syncing the database.. %v", err)
	}
}
