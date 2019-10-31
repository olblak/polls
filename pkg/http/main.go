package http

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/olblak/polls/pkg/ldap"
	"github.com/olblak/polls/pkg/polls"
	"log"
	"net/http"
	"strconv"
)

func participateHandler(w http.ResponseWriter, r *http.Request) {

	results := map[string]string{"confirmed": "false"}
	required_data := true

	// Read file
	email := r.FormValue("email")
	token := r.FormValue("token")
	poll := r.FormValue("poll")

	if email == "" {
		log.Printf("Missing 'email' Parameter")
		required_data = false
	}

	if token == "" {
		log.Printf("Missing 'token' Parameter for %v", email)
		required_data = false
	}

	if poll == "" {
		log.Printf("Missing 'poll' Parameter for %v", email)
		required_data = false
	}

	if required_data == true {
		var p polls.Poll

		p.OpenDatabaseConnection()
		defer p.CloseDatabaseConnection()

		isConfirmed := p.IsParticipantConfirmed(email, poll)

		if !isConfirmed {
			p.SetParticipation("true", email, token, poll)
		}
		// For now, we test a second time to be sure that the db is correct
		isConfirmed = p.IsParticipantConfirmed(email, poll)

		results["confirmed"] = strconv.FormatBool(isConfirmed)
	} else {
		results["confirmed"] = strconv.FormatBool(required_data)
	}

	js, err := json.Marshal(results)

	w.Write(js)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Participation Request received for %v\n", email)
}

func participantsGetHandler(w http.ResponseWriter, r *http.Request) {

	var p polls.Poll

	p.OpenDatabaseConnection()
	defer p.CloseDatabaseConnection()

	poll := r.FormValue("poll")

	if poll == "" {
		log.Println("Poll parameter is empty")
	} else {
		log.Printf("List All participants for poll: %v\n", poll)
		log.Println(p.ParticipantsAsCSV(poll))

		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment;filename=participants.csv")

		csv := csv.NewWriter(w)
		csv.WriteAll(p.ParticipantsAsCSV(poll))
		if err := csv.Error(); err != nil {
			log.Println("error writing csv:", err)
		}

	}

}

func participantsPostHandler(w http.ResponseWriter, r *http.Request) {

	var p polls.Poll

	p.OpenDatabaseConnection()
	defer p.CloseDatabaseConnection()

	poll := r.FormValue("poll")
	group := r.FormValue("group")
	seniorityCriteriaDate := r.FormValue("seniorityCriteriaDate")

	if group == "" {
		group = "all"
	}

	if seniorityCriteriaDate == "" {
		seniorityCriteriaDate = "2019/09/01 00:00:00"
	}

	if poll == "" {
		log.Println("Poll parameter is empty")

	} else {

		seniorAccounts := []map[string]string{{}}
		//newAccounts := []map[string]string{{}}

		seniorAccounts, _ = ldap.Users(seniorityCriteriaDate, group)
		p.CreateParticipants(poll, seniorAccounts)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func makeRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/participate", participateHandler).Methods("GET")
	router.HandleFunc("/api/participants", participantsGetHandler).Methods("GET")
	router.HandleFunc("/api/participants", participantsPostHandler).Methods("POST")
	router.HandleFunc("/api/status", statusHandler)
	return router
}

func StartHttp() {
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", makeRouter()))
}
