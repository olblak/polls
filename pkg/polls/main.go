package polls

import (
	"fmt"
	"log"

	"database/sql"

	"strconv"

	"github.com/olblak/polls/pkg/db"
)

// isTokenValid check if the token provide by a http request match the one defined in the database for a specific email adress
// then return true or false
func isTokenValid(mail, token, poll string) bool {
	db, err := sql.Open("postgres", db.Database_url)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Printf("unable to use data source name %v\n", err)
	}
	query := fmt.Sprintf("SELECT token FROM participate WHERE mail = '%v' AND poll = '%v'", mail, poll)

	rows := db.QueryRow(query)

	if err != nil {
		log.Println(err)
	}

	var dbToken string

	if err := rows.Scan(&dbToken); err != nil {
		log.Println(err)
	}

	if token == dbToken {
		return true
	}
	log.Printf("Wrong Token provided for %v\n", mail)

	return false
}

// isParticipant check if a specific email address has been registered for a specific poll
func isParticipant(mail, poll string) bool {
	db, err := sql.Open("postgres", db.Database_url)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Printf("unable to use data source name %v\n", err)
	}
	query := fmt.Sprintf("SELECT id FROM participate WHERE mail = '%v' AND poll = '%v'", mail, poll)

	rows := db.QueryRow(query)

	if err != nil {
		log.Println(err)
	}

	var id string

	if err := rows.Scan(&id); err != nil {
		log.Println(err)
	}

	if i, _ := strconv.Atoi(id); i > 0 {
		return true
	}

	return false
}

// Participants return a list of every participant who registered for a specific poll
func Participants(poll string) []map[string]string {

	var voters []map[string]string

	db, err := sql.Open("postgres", db.Database_url)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Printf("unable to use data source name : %v\n", err)
	}
	query := fmt.Sprintf("SELECT mail,token, participate FROM participate WHERE poll = '%v'", poll)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var mail string
		var token string
		var participate string

		if err := rows.Scan(&mail, &token, &participate); err != nil {
			log.Println(err)
		}
		voters = append(voters, map[string]string{"mail": mail, "token": token, "participate": participate})
	}

	return voters
}

// CreateParticipants generate a list of all participants from a specific ldap group and then insert it
// inside the database
func CreateParticipants(poll string, participants []map[string]string) {
	log.Printf("Create %v participants for poll number: $%v\n", len(participants), poll)

	db, err := sql.Open("postgres", db.Database_url)
	defer db.Close()
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Printf("unable to use data source name: %v\n", err)
	}

	for _, participant := range participants {

		if isParticipant(participant["mail"], poll) {
			log.Printf("%v is already in the participants invitation list", participant["mail"])
			continue
		}

		log.Println(participant["name"])
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
		}

		stmt, err := tx.Prepare("INSERT INTO participate (account, mail, poll, participate) VALUES ($1, $2,$3,'false')")

		defer stmt.Close()

		result, err := stmt.Exec(
			participant["ldap_account"],
			participant["mail"],
			poll)

		if err != nil {
			log.Println(err)
		}

		if _, err = result.RowsAffected(); err != nil {
			log.Println(err)
		}

		tx.Commit()
		log.Printf("%v was added to the participants invitation list", participant["mail"])
	}
}

// SetParticipation will confirm user participation for a specific poll
func SetParticipation(value, email, token, poll string) bool {
	db, err := sql.Open("postgres", db.Database_url)
	defer db.Close()
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Printf("unable to use data source name: %v\n", err)
	}

	if !isTokenValid(email, token, poll) {
		return false
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	stmt, err := tx.Prepare("UPDATE participate SET participate=$1 WHERE mail=$2 AND poll=$3")

	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(value, email, poll)

	if err != nil {
		log.Println(err)
	}

	if _, err = result.RowsAffected(); err != nil {
		log.Println(err)
	}

	tx.Commit()

	return true
}
