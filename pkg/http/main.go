package http

import (
	"github.com/gorilla/mux"
	"github.com/olblak/polls/pkg/ldap"
	"github.com/olblak/polls/pkg/polls"
	"io/ioutil"
	"log"
	"net/http"
)

func participateHandler(w http.ResponseWriter, r *http.Request) {
	// Read file
	email := r.FormValue("email")
	token := r.FormValue("token")
	poll := r.FormValue("poll")

	if email == "" {
		log.Printf("Missing 'email' Parameter")
		return
	}

	if token == "" {
		log.Printf("Missing 'token' Parameter for %v", email)
		return
	}

	if poll == "" {
		log.Printf("Missing 'poll' Parameter for %v", email)
		return
	}

	polls.SetParticipation("true", email, token, poll)
	log.Printf("Participation Request received for %v\n", email)
}

func participantsGetHandler(w http.ResponseWriter, r *http.Request) {

	poll := r.FormValue("poll")

	if poll == "" {
		log.Println("Poll parameter is empty")
	} else {
		log.Printf("List All participants for poll: %v\n", poll)
		log.Println(polls.Participants(poll))
	}

}

func participantsPostHandler(w http.ResponseWriter, r *http.Request) {

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
		polls.CreateParticipants(poll, seniorAccounts)
	}
}

func accountDeletionHandler(w http.ResponseWriter, r *http.Request) {
	// Read file
	log.Println("Request received")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Println(body)

	// Do Something
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func makeRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/participate", participateHandler).Methods("GET")
	router.HandleFunc("/participants", participantsGetHandler).Methods("GET")
	router.HandleFunc("/participants", participantsPostHandler).Methods("POST")
	router.HandleFunc("/account/delete", accountDeletionHandler).Methods("GET")
	router.HandleFunc("/status", statusHandler)
	return router
}

func StartHttp() {
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", makeRouter()))
}
