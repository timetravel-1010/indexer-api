package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/timetravel-1010/indexer-api/internal/zinc"
)

type EmailHandler struct{}

var (
	c = zinc.Config{
		Host:     "localhost",
		Port:     "4080",
		Username: "admin",
		Password: "Complexpass#123",
	}
)

// SearchByTerm
func (eh EmailHandler) SearchByTerm(w http.ResponseWriter, r *http.Request) {
	query, err := zinc.BuildQuery(zinc.ZincQuery{
		Params:     r.URL.Query(),
		SearchType: zinc.MATCH_QUERY,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res, err := zinc.DoZincRequest(r, query, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res)
}

// GetEmails
func (eh EmailHandler) GetEmails(w http.ResponseWriter, r *http.Request) {
	q, err := zinc.BuildQuery(zinc.ZincQuery{
		Params:     r.URL.Query(),
		SearchType: zinc.MATCHALL_QUERY,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res, err := zinc.DoZincRequest(r, q, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res)
}
