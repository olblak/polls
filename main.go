package main

import (
	"github.com/olblak/polls/pkg/db"
	"github.com/olblak/polls/pkg/http"
)

func main() {
	db.RunDatabaseMigration()
	http.StartHttp()
}
