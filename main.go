package main

import (
	"github.com/olblak/polls/pkg/db"
	"github.com/olblak/polls/pkg/http"
	"os"
	"string"
)

func main() {
	if migration_enabled := os.Getenv("DB_MIGRATION_ENABLED"); strings.ToLower(migration_enabled) == "true" {
		db.RunDatabaseMigration()
	}
	http.StartHttp()
}
