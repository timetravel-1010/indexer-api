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
	queryTemp := []byte(`{
    "search_type": "match",
    "query": {
        "term": "%v",
        "field": "_all"
    },
    "sort_fields": ["-@timestamp"],
    "from": 0,
    "max_results": %v,
    "_source": [ ]
}`)
	//term := r.URL.Query().Get("term")
	query, err := zinc.BuildQuery(zinc.ZincQuery{
		Query:      queryTemp,
		Params:     r.URL.Query(),
		SearchType: zinc.MATCH_QUERY,
	})

	res, err := zinc.DoZincRequest(r, query, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res)
}

// GetEmails
func (eh EmailHandler) GetEmails(w http.ResponseWriter, r *http.Request) {

	query := []byte(`{
	    "search_type": "matchall",
	    "from": 0,
	    "max_results": %v,
	    "_source": []
	}`)
	q, err := zinc.BuildQuery(zinc.ZincQuery{
		Query:      query,
		Params:     r.URL.Query(),
		SearchType: zinc.MATCHALL_QUERY,
	})
	res, err := zinc.DoZincRequest(r, q, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res)
}
