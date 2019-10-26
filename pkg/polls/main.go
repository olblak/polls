package polls

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/olblak/polls/pkg/db"
)

type Poll struct {
	db *sql.DB
}

func (p *Poll) OpenDatabaseConnection() {
	var err error

	p.db, err = sql.Open("postgres", db.Database_url)

	p.db.SetMaxIdleConns(100)
	p.db.SetMaxOpenConns(100)

	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Printf("unable to use data source name %v\n", err)
	}
}

func (p *Poll) CloseDatabaseConnection() {
	p.db.Close()
}

// isTokenValid check if the token provide by a http request match the one defined in the database for a specific email adress
// then return true or false
func (p *Poll) isTokenValid(mail, token, poll string) bool {
	query := fmt.Sprintf("SELECT token FROM participate WHERE mail = '%v' AND poll = '%v'", mail, poll)

	rows := p.db.QueryRow(query)

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
func (p *Poll) isParticipant(mail, poll string) bool {

	query := fmt.Sprintf("SELECT id FROM participate WHERE mail = '%v' AND poll = '%v'", mail, poll)

	rows := p.db.QueryRow(query)

	var id string

	err := rows.Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Println(err)
		return false
	default:
		return true
	}

}

// Participants return a list of every participant who registered for a specific poll
func (p *Poll) Participants(poll string) []map[string]string {

	var voters []map[string]string

	query := fmt.Sprintf("SELECT mail,token, participate FROM participate WHERE poll = '%v'", poll)

	rows, err := p.db.Query(query)
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

	// TODO export csv cn,mail,token,poll

	return voters
}

// Participants return a list of every participant who registered for a specific poll
func (p *Poll) ParticipantsAsCSV(poll string) [][]string {

	voters := [][]string{{"mail", "token", "participate"}}

	query := fmt.Sprintf("SELECT mail,token, participate FROM participate WHERE poll = '%v'", poll)

	rows, err := p.db.Query(query)
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
		voters = append(voters, []string{mail, token, participate})
	}

	return voters
}

// CreateParticipants generate a list of all participants from a specific ldap group and then insert it
// inside the database
func (p *Poll) CreateParticipants(poll string, participants []map[string]string) {
	//values := ""
	log.Printf("Create %v participants for poll number: %v\n", len(participants), poll)

	for _, participant := range participants {

		if len(participant["mail"]) > 50 {
			log.Printf("Wrong email: %v\n", participant["mail"])
			continue
		}

		if len(participant["account"]) > 50 {
			log.Printf("Wrong account: %v\n", participant["mail"])
			continue
		}

		if p.isParticipant(participant["mail"], poll) {
			log.Printf("%v is already in the participants invitation list", participant["mail"])
			continue
		}

		tx, err := p.db.Begin()
		defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

		if err != nil {
			log.Println(err)
		}

		stmt, err := tx.Prepare("INSERT INTO participate (account, mail, poll, participate) VALUES ($1,$2,$3,'false')")
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

	}
}

// SetParticipation will confirm user participation for a specific poll
func (p *Poll) SetParticipation(value, email, token, poll string) bool {

	if !p.isTokenValid(email, token, poll) {
		return false
	}

	tx, err := p.db.Begin()
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
