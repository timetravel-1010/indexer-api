package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/timetravel-1010/indexer-api/internal/models"
)

type EmailHandler struct{}

type ZincResponse struct {
	Hits struct {
		Hits []models.Hit `json:"hits"`
	} `json:"hits"`
}

type Config struct {
	host     string
	port     string
	username string
	password string
}

type ZincQuery struct {
	query      []byte
	params     map[string][]string
	searchType ZincSearchType
}

type ZincSearchType string

const (
	match    ZincSearchType = "match"
	matchall ZincSearchType = "matchall"

	DEFAULT_PAGE_SIZE = 10
)

var (
	requiredParams = []string{
		"page",
	}

	matchQueryParams = append([]string{
		"term",
	}, requiredParams...)

	matchAllQueryParams = append([]string{}, requiredParams...)

	c Config = Config{
		host:     "localhost",
		port:     "4080",
		username: "admin",
		password: "Complexpass#123",
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
	query, err := buildQuery(ZincQuery{
		query:      queryTemp,
		params:     r.URL.Query(),
		searchType: match,
	})

	fmt.Println("url:", query)

	res, err := doZincRequest(r, query)
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
	q, err := buildQuery(ZincQuery{
		query:      query,
		params:     r.URL.Query(),
		searchType: matchall,
	})
	res, err := doZincRequest(r, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res)
}

func buildQuery(zq ZincQuery) (q string, err error) {
	/*
		if validPage() != nil {
			numResults := DEFAULT_PAGE_SIZE
			n, err := strconv.Atoi(page)
			if err != nil {
				return "", err
			}
			numResults = n

		}
	*/
	switch zq.searchType {
	case match:
		return buildMatchQuery(zq.query, zq.params)
	case matchall:
		return buildMatchAllQuery(zq.query, zq.params)
	}

	return "", nil
}

func getParams(params map[string][]string, toGet []string) (p []any, err error) {
	for _, tg := range toGet {
		urlParam := params[tg][0]
		if urlParam == "" {
			return nil, errors.New(fmt.Sprintf("missing param %s", p))
		}
		p = append(p, urlParam)
	}
	return p, nil
}

func buildMatchQuery(queryTemp []byte, params map[string][]string) (q string, err error) {
	urlParams, err := getParams(params, matchQueryParams)
	if err != nil {
		return "", err
	}
	q = fmt.Sprintf(string(queryTemp), urlParams...)
	return q, err
}

func buildMatchAllQuery(queryTemp []byte, params map[string][]string) (q string, err error) {
	urlParams, err := getParams(params, matchAllQueryParams)
	if err != nil {
		return "", err
	}
	q = fmt.Sprintf(string(queryTemp), urlParams...)
	return q, nil
}

func doZincRequest(r *http.Request, query string) (*ZincResponse, error) {
	log.Println("url query:", r.URL.Query())

	index := r.URL.Query().Get("index")
	url := fmt.Sprintf("http://%s:%s/api/%s/_search", c.host, c.port, index)
	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(query),
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	log.Println(resp.StatusCode)

	zr := &ZincResponse{}
	err = json.NewDecoder(resp.Body).Decode(&zr)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return zr, nil
}
